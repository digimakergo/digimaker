package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/rest"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/spf13/viper"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

// todo dynamic settings,can set global param
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{},
	})

var commonStruct = contenttype.ContentCommon{}

var commonArgs = graphql.BindArg(commonStruct, commonStruct.IdentifierList()...)

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

var staticKeys = map[string]string{
	"and": "&&",
	"or":  "||",
	"gt":  ">",
	"ge":  ">=",
	"lt":  "<",
	"le":  "<=",
	"eq":  "==",
	"ne":  "!=",
}

var DMScalarType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "DMScalarType",
	Description: "Digimaker scalar type.",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case fieldtypes.Json:
			return value.String()
		default:
			return nil
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

	// parse querying
	astDocument, err := parser.Parse(parser.ParseParams{
		Source: queryParams,
		Options: parser.ParseOptions{
			NoSource:   true,
			NoLocation: true,
		},
	})
	if err != nil {
		rest.HandleError(err, w)
		return
	}

	if len(astDocument.Definitions) == 0 {
		rest.HandleError(errors.New("definitions length is 0"), w)
		return
	}

	for _, adDefinition := range astDocument.Definitions {
		if def, ok := adDefinition.(*ast.OperationDefinition); ok {
			// todo verify def.Name.Value => "content"
			if len(def.SelectionSet.Selections) > 0 {
				for _, selection := range def.SelectionSet.Selections {
					if sel, isOk := selection.(*ast.Field); isOk {
						cType := sel.Name.Value
						if cType != "" {
							_, err = definition.GetDefinition(cType)
							if err != nil {
								rest.HandleError(err, w)
								return
							}
							contentInstance := contenttype.NewInstance(cType)
							def, _ := definition.GetDefinition(cType)

							// default build in types
							contentFields := graphql.BindFields(contentInstance)
							// customized field types
							for _, identifier := range contentInstance.IdentifierList() {
								fieldDef := def.FieldMap[identifier]
								switch fieldDef.FieldType {
								case "json":
									gqlField := graphql.Field{}
									gqlField.Type = DMScalarType
									contentFields[identifier] = &gqlField
								default:
								}
							}

							cFieldOfType := graphql.NewObject(graphql.ObjectConfig{Name: cType, Fields: contentFields})

							args := graphql.BindArg(contentInstance)
							for _, arg := range args {
								//								fmt.Println("config original:", k+"|"+config.Type.Name()+"|"+config.Type.String())
								arg.Type = graphql.NewList(arg.Type)
								//								fmt.Println("config now:", k+"|"+config.Type.Name()+"|"+config.Type.String())
							}

							if len(commonArgs) > 0 {
								for key, commonArg := range commonArgs {
									args[key] = commonArg
								}
							}

							queryType.AddFieldConfig(cType, &graphql.Field{
								Type: graphql.NewList(cFieldOfType),
								Args: buildRootArgs(args), // return a config
								Resolve: func(p graphql.ResolveParams) (interface{}, error) {
									log.Info(p)
									return parseSolveParams(ctx, cType, p)
								},
							})
						}
					}
				}
			}
		}
	}

	// setting in the end
	var schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		})

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: queryParams,
	})
	rest.WriteResponse(result, w)
}

func buildRootArgs(content graphql.FieldConfigArgument) (ret graphql.FieldConfigArgument) {
	inputArgs := make(graphql.InputObjectConfigFieldMap, 0)
	sortArgs := make(graphql.InputObjectConfigFieldMap, 0)
	rootArgs := make(graphql.FieldConfigArgument, 0)

	if len(content) > 0 {
		for key, config := range content {
			inputArgs[key] = &graphql.InputObjectFieldConfig{
				Type:        graphql.NewList(config.Type),
				Description: "input args " + key,
			}
			sortArgs[key] = &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "sort args " + key,
			}
			// root query
			rootArgs[key] = config
		}
	}

	// build conditon args
	childArgs := buildGraphqlMap(staticKeys, inputArgs, inputArgs)

	parentArgs := buildGraphqlMap(staticKeys, childArgs, inputArgs)

	// root args add children args
	// filter query
	rootArgs["filter"] = &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        "FILTER",
			Fields:      parentArgs,
			Description: "filter args",
		}),
	}
	// sort query : sort:{cid:"desc"}
	rootArgs["sort"] = &graphql.ArgumentConfig{
		Type:        graphql.NewList(graphql.String),
		Description: "sort args",
	}
	// page limit query :
	rootArgs["limit"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 10,
		Description:  "pageSize",
	}
	// page offset query :
	rootArgs["offset"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		DefaultValue: 0,
		Description:  "page",
	}
	return rootArgs
}

// replace generateFilterArgs
func generateFilterArgs(filterMap map[string]interface{}, condition db.Condition) db.Condition {

	if filterMap == nil || len(filterMap) == 0 {
		return condition
	}

	for key, value := range filterMap {
		if childMap, ok := value.(map[string]interface{}); ok {
			switch key {
			case "and":
				for k, v := range childMap {
					switch child := v.(type) {
					case map[string]interface{}:
						for ck, cv := range child {
							condition = condition.Cond(ck, cv)
						}
					default:
						condition = isSlice(k, v, condition)
						//condition = condition.Cond(k, v)
					}
				}
			case "or":
				for k, v := range childMap {
					switch child := v.(type) {
					case map[string]interface{}:
						for ck, cv := range child {
							condition = condition.Or(ck, cv)
						}
					default:
						condition = condition.Or(k, v)
						//arrStr := make([]string,0)
						//arrInt := make([]int,0)
						//condition = isSliceV2(ifDoFetchQuery(k), v, condition,arrStr,arrInt)
					}
				}
			case "gt", "ge", "lt", "le":
				for k, v := range childMap {
					arrStr := make([]string, 0)
					arrInt := make([]int, 0)
					condition = isSliceV2(ifDoFetchQuery(k)+operatorByKey(key), v, condition, arrStr, arrInt)
				}
			}
		}
	}
	return condition
}

func parseSolveParams(ctx context.Context, cType string, p graphql.ResolveParams) (interface{}, error) {

	condition := db.Condition{}
	// condition key query
	for k, v := range p.Args {
		// verify key
		if isSuitKey(k) {
			// filter params
			if k == "filter" {
				if filterMap, ok := v.(map[string]interface{}); ok {
					condition = generateFilterArgs(filterMap, condition)
				}
			}

			if k == "sort" {
				fmt.Println("do sort flag:", v)
				sortPs := make([]string, 0)
				sortsArr := v.([]interface{})
				for _, i := range sortsArr {
					sortPs = append(sortPs, fmt.Sprint(i))
				}
				condition = condition.Sortby(sortPs...)
			}

		} else {
			fmt.Println("is not condition key:" + k)
			fmt.Println("kind:", reflect.ValueOf(v), "|", reflect.TypeOf(v).Kind())

			// valid slice
			condition = isSlice(k, v, condition)
		}
	}

	// page params (default value -> offset:0,limit:10)
	condition = condition.Limit(p.Args["offset"].(int), p.Args["limit"].(int))

	conRet, conParams := db.BuildCondition(condition)
	fmt.Println("root build condition:", conRet, "|", conParams, "|", condition)
	list, _, err := query.List(ctx, cType, condition)
	return list, err
}

func operatorByKey(input string) string {
	// [lt:less than,le:less than or equal to,"]
	switch input {
	case "gt":
		return ">"
	case "ge":
		return ">="
	case "lt":
		return "<"
	case "le":
		return "<="
	case "eq":
		return "=="
	case "ne":
		return "!="
	default:
		log.Info("operator by key : " + input)
		return ""
	}
}

func buildGraphqlMap(params map[string]string, input graphql.InputObjectConfigFieldMap, extra ...graphql.InputObjectConfigFieldMap) graphql.InputObjectConfigFieldMap {
	result := make(graphql.InputObjectConfigFieldMap, 0)

	for key, value := range params {
		result[key] = &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        strings.ToUpper(key),
				Fields:      input,
				Description: key + ":" + value,
			}),
		}
	}

	if len(extra) > 0 {
		for k, v := range extra[0] {
			result[k] = v
		}
	}

	return result
}

func isSuitKey(key string) bool {
	params := map[string]bool{
		"and":    true,
		"or":     true,
		"gt":     true,
		"ge":     true,
		"lt":     true,
		"le":     true,
		"sort":   true,
		"limit":  true,
		"offset": true,
		"filter": true,
	}
	if _, ok := params[key]; ok {
		return true
	}
	return false
}

func verifyKey(key string) bool {
	switch key {
	case "and", "or":
		return true
	case "gt", "ge", "lt", "le":
		return true
	case "filter":
		return true
	case "sort":
		return true
	case "limit", "offset":
		return true
	default:
		return false
	}
}

func buildJsonResult(input interface{}) string {
	ret, err := json.Marshal(input)
	if err != nil {
		return "{\"data\": \"" + err.Error() + "\"}"
	}
	return string(ret)
}

func isSlice(key string, value interface{}, condition db.Condition) db.Condition {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		arr := value.([]interface{})
		childStringArr := make([]string, 0)
		childIntArr := make([]int, 0)
		for _, v1 := range arr {
			fmt.Println("v1:", v1)
			switch c := v1.(type) {
			case []interface{}:
				for _, cv1 := range c {
					switch cV := cv1.(type) {
					case []interface{}:
						fmt.Println("s")
					case string:
						childStringArr = append(childStringArr, cV)
					case int:
						childIntArr = append(childIntArr, cV)
					}
				}
			case string:
				childStringArr = append(childStringArr, c)
			case int:
				childIntArr = append(childIntArr, c)
			}
		}
		if len(childStringArr) > 0 {
			condition = condition.Cond(ifDoFetchQuery(key), childStringArr)
		}
		if len(childIntArr) > 0 {
			condition = condition.Cond(ifDoFetchQuery(key), childIntArr)
		}
	} else {
		condition = condition.Cond(ifDoFetchQuery(key), value)
	}

	return condition
}

func isSliceV2(key string, value interface{}, condition db.Condition, arrStr []string, arrInt []int) db.Condition {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		arr := value.([]interface{})
		for _, item := range arr {
			switch info := item.(type) {
			case []interface{}:
				condition = isSliceV2(key, info, condition, arrStr, arrInt)
			case string:
				arrStr = append(arrStr, info)
			case int:
				arrInt = append(arrInt, info)
			}
		}
		if len(arrStr) > 0 {
			condition = condition.Cond(ifDoFetchQuery(key), arrStr)
		}
		if len(arrInt) > 0 {
			condition = condition.Cond(ifDoFetchQuery(key), arrInt)
		}
	} else {
		condition = condition.Cond(ifDoFetchQuery(key), value)
	}

	return condition
}

func isSliceV3() {

}

func ifDoFetchQuery(key string) string {
	switch key {
	case "id":
		return "l.id"
	case "version":
		return "l.name"
	case "published":
		return "l.published"
	case "modified":
		return "l.modified"
	case "author":
		return "l.author"
	case "author_name":
		return "l.author_name"
	case "cuid":
		return "l.cuid"
	case "status":
		return "l.status"
	default:
		return key
	}
}

// func exec one time,you can run this func in init
func settingOnce() {
	// 1. common args need convert graphql type -> list(type)
	for _, commonArg := range commonArgs {
		commonArg.Type = graphql.NewList(commonArg.Type)
	}
}

func init() {

	settingOnce()

	// try to diff method
	rest.RegisterRoute("/graphql", QueryGraphql, "POST")
}
