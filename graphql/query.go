package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/dmeditor"
	"github.com/digimakergo/digimaker/rest"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/spf13/viper"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

type ArrModel struct {
	ArrInt    []int
	ArrString []string
}

//set graphql/api_key in dm.yaml
func AuthAPIKey(r *http.Request) error {
	apiKey := viper.GetString("graphql.api_key")
	if apiKey == "" {
		log.Error("Not api key set up", "")
		return errors.New("Set up issue")
	}
	rApiKey := r.Header.Get("apiKey")
	if rApiKey == "" {
		return errors.New("Need authorization")
	}
	if rApiKey == apiKey {
		return nil
	} else {
		return errors.New("Wrong api key")
	}
}

//Digimaker scalar type
func getFilterType(cType string) *graphql.Scalar {
	def, _ := definition.GetDefinition(cType)

	filterInput := graphql.NewScalar(graphql.ScalarConfig{
		Name:        "FilterInput",
		Description: "Filter input type.",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			cond, err := generateCondition(valueAST, db.EmptyCond(), def.FieldMap)
			if err != nil {
				return err
			}
			return cond
		},
	})

	return filterInput
}

func generateCondition(valueAST ast.Value, cond db.Condition, fieldMap map[string]definition.FieldDef) (db.Condition, error) {
	switch v := valueAST.(type) {
	case *ast.ListValue:
		for _, item := range v.Values {
			objectCond, err := generateCondition(item, db.EmptyCond(), fieldMap)
			if err != nil {
				return db.FalseCond(), err
			}
			cond = cond.Or(objectCond)
		}
		return cond, nil
	case *ast.ObjectValue:
		for _, gqlField := range v.Fields {
			key := gqlField.Name.Value
			//todo: convert graphql value to our value or nested condition
			value := gqlField.Value.GetValue()
			//todo: use meta from definition here
			common := []string{"id", "published", "modified"}
			if util.Contains(common, key) {
				if key == "id" {
					key = "l.id"
				}
				cond = cond.And(key, value)
			} else if _, ok := fieldMap[key]; ok {
				cond = cond.And(key, value)
			} else {
				return db.FalseCond(), fmt.Errorf("Field %v not found", key)
			}
		}
		return cond, nil
	default:
		return db.FalseCond(), errors.New("Unknown type in filter")
	}
}

//Digimaker scalar type
var DMScalarType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "DMScalarType",
	Description: "Digimaker scalar type.",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case fieldtypes.Json:
			result, _ := dmeditor.ProceedData(context.Background(), value)
			return result
		default:
			return value
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			j := fieldtypes.Json{}
			j.Content = []byte(value)
			return j
		default:
			return nil
		}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return fieldtypes.Json{Content: []byte(valueAST.Value)}
		default:
			return nil
		}
	},
})

// var typeMapping map[string]graphql.Type = map[string]graphql.Type{
// 	"text":     graphql.String,
// 	"richtext": graphql.String,
// 	"int":      graphql.Int,
// }

func getContentGQLType(def definition.ContentType) *graphql.Object {
	//fields on a content type
	gqlFields := graphql.Fields{}
	for fieldIdentifier, fieldDef := range def.FieldMap {
		//set fields to gqlFields
		gqlField := graphql.Field{}
		gqlField.Type = DMScalarType
		gqlField.Name = fieldDef.Name
		gqlFields[fieldIdentifier] = &gqlField
	}

	//customized type
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   def.Name,
		Fields: gqlFields,
	})
}

var querySchema graphql.Schema

func initQuerySchema() {
	//content types
	gqlContentTypes := graphql.Fields{}
	for cType, def := range definition.GetDefinitionList()["default"] {
		contentGQLType := getContentGQLType(def)
		gqlContentTypes[cType] = &graphql.Field{
			Name: def.Name,
			Type: graphql.NewList(contentGQLType),
			Args: buildContentArgs(cType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return executeQuery(p.Context, p)
			},
		}
	}

	//query
	var gqlQuery = graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: gqlContentTypes,
	})

	querySchema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: gqlQuery,
		})
}

func QueryGraphql(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authErr := AuthAPIKey(r)
	if authErr != nil {
		rest.HandleError(authErr, w)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			rest.HandleError(err.(error), w)
			return
		}
	}()

	queryParams := ""
	if http.MethodGet == r.Method {
		queryParams = r.URL.Query().Get("query")
		if queryParams == "" {
			rest.HandleError(errors.New("query is nil"), w)
			return
		}
	} else {
		// todo : parse post request
		postParams := postData{}
		requestStr, err := io.ReadAll(r.Body)
		if err != nil {
			rest.HandleError(err, w)
			return
		}
		err = json.Unmarshal(requestStr, &postParams)
		if err != nil {
			rest.HandleError(errors.New("post unmarshal is nil"), w)
			return
		}
		// set params
		queryParams = postParams.Query
	}

	result := graphql.Do(graphql.Params{
		Context:       ctx,
		Schema:        querySchema,
		RequestString: queryParams,
	})

	resultStr, err := json.Marshal(result)
	if err != nil {
		rest.HandleError(err, w)
	}
	w.Write(resultStr)
}

//build arguments on a content type
func buildContentArgs(cType string) graphql.FieldConfigArgument {
	result := make(graphql.FieldConfigArgument, 0)

	result["filter"] = &graphql.ArgumentConfig{
		Type: getFilterType(cType),
	}

	//todo: check key
	// sort query : sort:{cid:"desc"}
	result["sort"] = &graphql.ArgumentConfig{
		Type:        graphql.NewList(graphql.String),
		Description: "sort args",
	}
	// page limit query :
	result["limit"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 10,
		Description:  "pageSize",
	}
	// page offset query :
	result["offset"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 0,
		Description:  "page",
	}
	return result
}

//execute quer
func executeQuery(ctx context.Context, p graphql.ResolveParams) (interface{}, error) {
	cType := p.Info.FieldName
	condition := db.EmptyCond()
	// condition key query
	args := p.Args

	// filter params
	if filters, ok := args["filter"]; ok {
		switch filters.(type) {
		case error:
			return nil, filters.(error)
		case db.Condition:
			condition = filters.(db.Condition)
		default:
			return nil, errors.New("Unknown filter")
		}
	}

	if sort, ok := args["sort"]; ok {
		sortPs := make([]string, 0)
		sortsArr := sort.([]interface{})
		for _, i := range sortsArr {
			sortPs = append(sortPs, fmt.Sprint(i))
		}
		condition = condition.Sortby(sortPs...)
	}

	// page params (default value -> offset:0,limit:10)
	condition = condition.Limit(p.Args["offset"].(int), p.Args["limit"].(int))

	list, _, err := query.List(ctx, cType, condition)
	return list, err
}

func init() {
	initQuerySchema()

	// try to diff method
	rest.RegisterRoute("/graphql", QueryGraphql, "POST")
}
