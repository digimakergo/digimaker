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
	filterArgs := make(graphql.InputObjectConfigFieldMap, 0)
	// sortArgs := make(graphql.InputObjectConfigFieldMap, 0)
	result := make(graphql.FieldConfigArgument, 0)

	def, _ := definition.GetDefinition(cType)

	for fieldIdentifier, _ := range def.FieldMap {
		//todo: customize this
		filterArgs[fieldIdentifier] = &graphql.InputObjectFieldConfig{
			Type:        graphql.String,
			Description: "input args " + fieldIdentifier,
		}
		// sortArgs[fieldIdentifier] = &graphql.InputObjectFieldConfig{
		// 	Type:        graphql.String,
		// 	Description: "sort args " + fieldIdentifier,
		// }
	}

	// root args add children args
	// filter query
	result["filter"] = &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        "FILTER",
			Fields:      filterArgs,
			Description: "filter args",
		}),
	}
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
		for identifier, value := range filters.(map[string]interface{}) {
			condition = condition.Cond(identifier, value)
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
