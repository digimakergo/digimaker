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
				Type:        config.Type,
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
	conditionArgs := buildConditionArgs(inputArgs)

	// root args add children args
	// filter query
	rootArgs["filter"] = &graphql.ArgumentConfig{
		Type: graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        "FILTER",
			Fields:      conditionArgs,
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

func buildConditionArgs(inputArgs graphql.InputObjectConfigFieldMap) (conditionArgs graphql.InputObjectConfigFieldMap) {
	conditionArgs = graphql.InputObjectConfigFieldMap{
		"and": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "AND",
				Fields:      inputArgs,
				Description: "and args",
			}),
		},
		"or": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "OR",
				Fields:      inputArgs,
				Description: "or args",
			}),
		},
		"gt": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "GT",
				Fields:      inputArgs,
				Description: "gt args (>)",
			}),
		},
		"ge": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "GE",
				Fields:      inputArgs,
				Description: "gt args (>=)",
			}),
		},
		"lt": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "LT",
				Fields:      inputArgs,
				Description: "lt args (<)",
			}),
		},
		"le": &graphql.InputObjectFieldConfig{
			Type: graphql.NewInputObject(graphql.InputObjectConfig{
				Name:        "LE",
				Fields:      inputArgs,
				Description: "le args (<=)",
			}),
		},
	}
	return conditionArgs
}

func isConditionKey(key string) bool {
	params := map[string]bool{
		"and":    true,
		"or":     true,
		"gt":     true,
		"ge":     true,
		"lt":     true,
		"le":     true,
		"limit":  true,
		"offset": true,
		"filter": true,
	}
	if v, ok := params[key]; ok {
		return v
	}
	return false
}

func verifyKey(key string) bool {
	switch key {
	case "and", "or", "gt", "ge", "lt", "le":
		return true
	case "filter":
		return true
	case "limit", "offset":
		return true
	default:
		return false
	}
}

func generateFilterArgs(filterMap map[string]interface{}, condition db.Condition) db.Condition {
	if andMap, ok := filterMap["and"].(map[string]interface{}); ok {
		if andMap != nil && len(andMap) > 0 {
			for k, v := range andMap {
				//condition = condition.And(k, v)
				condition = condition.And(k, v)
			}
		}
	}
	if orMap, ok := filterMap["or"].(map[string]interface{}); ok {
		if orMap != nil && len(orMap) > 0 {
			for k, v := range orMap {
				//condition = condition.Or(k, v)
				condition = condition.Or(k, v)
			}
		}
	}
	if gtMap, ok := filterMap["gt"].(map[string]interface{}); ok {
		if gtMap != nil && len(gtMap) > 0 {
			for k, v := range gtMap {
				condition = condition.Cond(k+" >", v)
			}
		}
	}

	if geMap, ok := filterMap["ge"].(map[string]interface{}); ok {
		if geMap != nil && len(geMap) > 0 {
			for k, v := range geMap {
				condition = condition.Cond(k+" >=", v)
			}
		}
	}

	if ltMap, ok := filterMap["lt"].(map[string]interface{}); ok {
		if ltMap != nil && len(ltMap) > 0 {
			for k, v := range ltMap {
				//condition = condition.Or(k, v)
				condition = condition.Cond(k+" <", v)
			}
		}
	}

	if leMap, ok := filterMap["le"].(map[string]interface{}); ok {
		if leMap != nil && len(leMap) > 0 {
			for k, v := range leMap {
				//condition = condition.Or(k, v)
				condition = condition.Cond(k+" <=", v)
			}
		}
	}
	return condition
}

func parseSolveParams(ctx context.Context, cType string, p graphql.ResolveParams) (list interface{}, err error) {

	if p.Args == nil || len(p.Args) == 0 {
		return list, errors.New("args is nil")
	}

	condition := db.Condition{}
	// condition key query
	for k, v := range p.Args {
		// verify key
		if isConditionKey(k) {
			continue
		}
		condition = condition.Cond(k, v)
		list, _, err = query.List(ctx, cType, condition)
		return list, err
	}
	if filter, ok := p.Args["filter"]; ok {
		if filterMap, ok := filter.(map[string]interface{}); ok {
			condition = generateFilterArgs(filterMap, condition)
		}
	}

	// sort params
	if sorts, ok := p.Args["sort"].([]interface{}); ok {
		sortPs := make([]string, 0)
		for _, i := range sorts {
			sortPs = append(sortPs, fmt.Sprint(i))
		}
		condition = condition.Sortby(sortPs...)
	}

	// page params
	condition = condition.Limit(p.Args["offset"].(int), p.Args["limit"].(int))
	conRet, conParams := db.BuildCondition(condition)
	log.Info(fmt.Sprintln("build condition:", conRet, "|", conParams, "|", condition))
	list, _, err = query.List(ctx, cType, condition)
	return list, err
}

func buildJsonResult(input interface{}) string {
	ret, err := json.Marshal(input)
	if err != nil {
		return "{\"data\": \"" + err.Error() + "\"}"
	}
	return string(ret)
}

func init() {
	// try to diff method
	rest.RegisterRoute("/graphql", QueryGraphql)
}
