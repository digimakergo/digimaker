package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/rest"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"io"
	"net/http"
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
							if len(commonArgs) > 0 {
								for key, commonArg := range commonArgs {
									args[key] = commonArg
								}
							}

							queryType.AddFieldConfig(cType, &graphql.Field{
								Type: graphql.NewList(cFieldOfType),
								Args: genRootArgs(args), // return a config
								Resolve: func(p graphql.ResolveParams) (interface{}, error) {
									// todo : generate query params
									var list interface{}
									var lerr error
									condition := db.Condition{}
									// filters params
									fmt.Println("args:", p.Args)
									if filter, ok := p.Args["filter"]; ok {
										fmt.Println("filter params:", buildJsonResult(filter))
										if filterMap, mapOk := filter.(map[string]interface{}); mapOk {
											if andMap, andOk := filterMap["and"].(map[string]interface{}); andOk {
												fmt.Println("andMap params:", buildJsonResult(andMap))
												if andMap != nil && len(andMap) > 0 {
													for k, v := range andMap {
														if k == "cid" {
															k = "id"
														}
														condition = condition.And(k, v)
													}
												}
											}
											if orMap, orOk := filterMap["or"].(map[string]interface{}); orOk {
												fmt.Println("orMap params:", buildJsonResult(orMap))
												if orMap != nil && len(orMap) > 0 {
													for k, v := range orMap {
														condition = condition.Or(k, v)
													}
												}
											}
										}
									}

									// sort params
									//if sorts, ok := p.Args["sort"].(string); ok {
									//	condition = condition.Sortby(sorts)
									//}

									// page params
									//condition = condition.Limit(p.Args["offset"].(int), p.Args["limit"].(int))

									conditionRet, a := db.BuildCondition(condition)
									fmt.Println("condition:", conditionRet, "|", a)
									list, _, lerr = query.List(ctx, cType, condition)
									return list, lerr
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

func genRootArgs(content graphql.FieldConfigArgument) (ret graphql.FieldConfigArgument) {
	inputArgs := make(graphql.InputObjectConfigFieldMap, 0)
	sortArgs := make(graphql.InputObjectConfigFieldMap, 0)
	if len(content) > 0 {
		for key, config := range content {
			inputArgs[key] = &graphql.InputObjectFieldConfig{
				Type:        config.Type,
				Description: "input" + key,
			}
			sortArgs[key] = &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "sort" + key,
			}
		}
	}

	// conditon args
	conditionArgs := graphql.InputObjectConfigFieldMap{
		"and": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "and args",
				Fields:      inputArgs,
				Description: "and args",
			}),
		},
		"or": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "or args",
				Fields:      inputArgs,
				Description: "or args",
			}),
		},
	}

	// root args
	rootArgs := graphql.FieldConfigArgument{
		"filter": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "filter args",
				Fields:      conditionArgs,
				Description: "filter args",
			}),
		},
		"sort": &graphql.ArgumentConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "sort args",
				Fields:      sortArgs,
				Description: "sort args",
			}),
		},
		"limit": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 10,
			Description:  "pageSize",
		},
		"offset": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 1,
			Description:  "page",
		},
	}
	return rootArgs
}

func buildJsonResult(input interface{}) string {
	ret, err := json.Marshal(input)
	if err != nil {
		return "{\"data\":" + err.Error() + "}"
	}
	return string(ret)
}

func init() {
	// try to diff method
	rest.RegisterRoute("/graphql", QueryGraphql)
}
