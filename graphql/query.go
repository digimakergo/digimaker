package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/rest"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"io"
	"net/http"
	"reflect"
	"strings"
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

func QueryGraphql(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	for _, definition := range astDocument.Definitions {
		if def, ok := definition.(*ast.OperationDefinition); ok {
			// todo verify def.Name.Value => "content"
			if len(def.SelectionSet.Selections) > 0 {
				for _, selection := range def.SelectionSet.Selections {
					if sel, isOk := selection.(*ast.Field); isOk {
						cType := sel.Name.Value
						if cType != "" {
							cTypeField := contenttype.NewInstance(cType)
							cFieldOfType := graphql.NewObject(graphql.ObjectConfig{Name: cType, Fields: graphql.BindFields(cTypeField)})

							args := graphql.BindArg(cTypeField, cTypeField.IdentifierList()...)
							for k, config := range args {
								fmt.Println("config original:", k+"|"+config.Type.Name()+"|"+config.Type.String())
								config.Type = graphql.NewList(config.Type)
								fmt.Println("config now:", k+"|"+config.Type.Name()+"|"+config.Type.String())
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
				Description: "input" + key,
			}
			sortArgs[key] = &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "sort" + key,
			}
			// root query
			rootArgs[key] = config
		}
	}

	// build conditon args
	childArgs := buildGraphqlMap(getKeysMap(), inputArgs, inputArgs)

	parentArgs := buildGraphqlMap(getKeysMap(), childArgs, inputArgs)
	fmt.Println(parentArgs)

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
					fmt.Println("and :", reflect.TypeOf(v).Kind())
					switch child := v.(type) {
					case map[string]interface{}:
						fmt.Println("and map[string]interface{}:")
						for ck, cv := range child {
							condition = condition.Cond(ck, cv)
						}
					default:
						fmt.Println("and default")
						condition = condition.Cond(k, v)
					}
				}
			case "or":
				for k, v := range childMap {
					fmt.Println("or :", reflect.TypeOf(v).Kind())
					switch child := v.(type) {
					case map[string]interface{}:
						fmt.Println("or map[string]interface{}:")
						for ck, cv := range child {
							condition = condition.Or(ck, cv)
						}
					default:
						fmt.Println("or default")
						condition = condition.Or(k, v)
					}
				}
			case "gt", "ge", "lt", "le":
				for k, v := range childMap {
					condition = condition.Cond(k+operatorByKey(key), v)
				}
			}
		}
	}
	return condition
}

func parseSolveParams(ctx context.Context, cType string, p graphql.ResolveParams) (list interface{}, err error) {

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
			fmt.Println("is not condition key")
			fmt.Println("kind:", reflect.ValueOf(v))

			// valid slice
			condition = isSlice(k, v, condition)
		}
	}

	// page params (default value -> offset:0,limit:10)
	condition = condition.Limit(p.Args["offset"].(int), p.Args["limit"].(int))

	conRet, conParams := db.BuildCondition(condition)
	fmt.Println("root build condition:", conRet, "|", conParams, "|", condition)
	list, _, err = query.List(ctx, cType, condition)
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

func getKeysMap() map[string]string {
	return map[string]string{
		"and": "&&",
		"or":  "||",
		"gt":  ">",
		"ge":  ">=",
		"lt":  "<",
		"le":  "<=",
		"eq":  "==",
		"ne":  "!=",
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
			switch c := v1.(type) {
			case string:
				childStringArr = append(childStringArr, c)
			case int:
				childIntArr = append(childIntArr, c)
			}
		}
		if len(childStringArr) > 0 {
			condition = condition.Cond(key, childStringArr)
		}
		if len(childIntArr) > 0 {
			condition = condition.Cond(key, childIntArr)
		}
	} else {
		condition = condition.Cond(key, value)
	}
	return condition
}

func ifDoFetchQuery(key string) bool {
	switch key {
	case "cid":
		return true
	default:
		return false
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
	rest.RegisterRoute("/graphql", QueryGraphql)
}
