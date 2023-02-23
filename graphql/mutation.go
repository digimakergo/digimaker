package graphql

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// Returns the following schema:
//
//	mutation update {
//	  article(updateData: [{data: {title: "", summary: ""}, id: 1}]) {
//	    id
//	    title
//	  }
//	}
func mutationType() *graphql.Object {
	gqlContentTypes := graphql.Fields{}
	for cType, def := range definition.GetDefinitionList()["default"] {
		contentGQLType := getContentGQLType(def)
		gqlContentTypes[cType] = &graphql.Field{
			Name: def.Name,
			Type: graphql.NewList(contentGQLType),
			Args: buildMutationArgs(cType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolveMutation(p.Context, p)
			},
		}
	}
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: gqlContentTypes,
	})
}

func buildMutationArgs(cType string) graphql.FieldConfigArgument {
	result := make(graphql.FieldConfigArgument)
	result["updateData"] = &graphql.ArgumentConfig{
		Type: graphql.NewScalar(graphql.ScalarConfig{
			Name:        "DataInput",
			Description: "Data input type.",
			Serialize: func(value interface{}) interface{} {
				return value
			},
			ParseValue: func(value interface{}) interface{} {
				return value
			},
			ParseLiteral: func(valueAST ast.Value) interface{} {
				switch v := valueAST.(type) {
				case *ast.ListValue:
					def, _ := definition.GetDefinition(cType)
					result := []map[string]interface{}{}
					for _, v := range v.Values {
						if object, ok := v.(*ast.ObjectValue); ok {
							updateData := map[string]interface{}{
								"data": make(handler.InputMap),
							}
							if err := parseInput(object.Fields, updateData, def.FieldMap); err != nil {
								return fmt.Errorf("failed to parse input data: %w", err)
							}
							result = append(result, updateData)
						}
					}
					return result
				default:
					return errors.New("unknown updateData")
				}
			},
		}),
		Description: "Data to be updated.",
	}
	return result
}

func resolveMutation(ctx context.Context, p graphql.ResolveParams) (interface{}, error) {
	if inputs, ok := p.Args["updateData"]; ok {
		switch v := inputs.(type) {
		case []map[string]interface{}:
			for _, item := range v {
				userId := util.CurrentUserID(ctx)
				if userId == 0 {
					return nil, errors.New("need to login")
				}
				idInt, err := strconv.Atoi(item["id"].(string))
				if err != nil {
					return nil, fmt.Errorf("could not convert id string to int: %w", err)
				}
				data := item["data"].(handler.InputMap)
				_, err = handler.UpdateByContentID(ctx, p.Info.FieldName, idInt, data, userId)
				if err != nil {
					return nil, fmt.Errorf("could not update content by id: %w", err)
				}
			}
		case error:
			return nil, v
		default:
			return nil, errors.New("could not resolve, unknown updateData")
		}
	}
	return nil, nil
}

func parseInput(fields []*ast.ObjectField, inputData map[string]interface{}, fieldMap map[string]fieldtype.FieldDef) error {
	for _, field := range fields {
		key := field.Name.Value
		value := field.Value.GetValue()
		if key == "id" {
			if v, ok := value.(string); ok {
				inputData["id"] = v
			} else {
				return errors.New("id field should be string")
			}
		}
		if key == "data" {
			if v, ok := value.([]*ast.ObjectField); ok {
				if err := parseInput(v, inputData, fieldMap); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("data field should be an object: %#v", value)
			}
		}
		if key == "title" {
			if _, ok := fieldMap[key]; !ok {
				return fmt.Errorf("invalid query field: %s", key)
			}
			if v, ok := value.(string); ok {
				inputData["data"].(handler.InputMap)["title"] = v
			} else {
				return errors.New("title field should be string")
			}
		}
		if key == "summary" {
			if _, ok := fieldMap[key]; !ok {
				return fmt.Errorf("invalid query field: %s", key)
			}
			if v, ok := value.(string); ok {
				inputData["data"].(handler.InputMap)["summary"] = v
			} else {
				return errors.New("summary field should be string")
			}
		}
	}
	return nil
}
