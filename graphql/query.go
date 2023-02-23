package graphql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"

	// "github.com/digimakergo/digimaker/dmeditor"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// Digimaker scalar type
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
			cond, err := generateQueryCondition(valueAST, db.EmptyCond(), def.FieldMap)
			if err != nil {
				return err
			}
			return cond
		},
	})

	return filterInput
}

func generateQueryCondition(valueAST ast.Value, cond db.Condition, fieldMap map[string]fieldtype.FieldDef) (db.Condition, error) {
	switch v := valueAST.(type) {
	case *ast.ListValue:
		for _, item := range v.Values {
			objectCond, err := generateQueryCondition(item, db.EmptyCond(), fieldMap)
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
			if key == "id" {
				cond = cond.And("c.id", value)
			} else if _, ok := fieldMap[key]; ok {
				cond = cond.And(key, value)
			} else if strings.HasPrefix(key, "_location_") {
				name := strings.TrimPrefix(key, "_location_")
				if util.Contains(definition.LocationColumns, name) {
					return cond.And("l."+name, value), nil
				} else {
					return db.FalseCond(), fmt.Errorf("%v not found in location", name)
				}
			} else if strings.HasPrefix(key, "_metadata_") {
				name := strings.TrimPrefix(key, "_metadata_")
				if util.Contains(definition.MetaColumns, name) {
					return cond.And("c._"+name, value), nil
				} else {
					return db.FalseCond(), fmt.Errorf("%v not found in metadata", name)
				}
			} else {
				return db.FalseCond(), fmt.Errorf("Field %v not found", key)
			}
		}
		return cond, nil
	default:
		return db.FalseCond(), errors.New("Unknown type in filter")
	}
}

func queryType() *graphql.Object {
	//content types
	gqlContentTypes := graphql.Fields{}
	for cType, def := range definition.GetDefinitionList()["default"] {
		contentGQLType := getContentGQLType(def)
		gqlContentTypes[cType] = &graphql.Field{
			Name: def.Name,
			Type: graphql.NewList(contentGQLType),
			Args: buildQueryArgs(cType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolveQuery(p.Context, p)
			},
		}
	}
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: gqlContentTypes,
	})
}

// build arguments on a content type
func buildQueryArgs(cType string) graphql.FieldConfigArgument {
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

// execute quer
func resolveQuery(ctx context.Context, p graphql.ResolveParams) (interface{}, error) {
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
