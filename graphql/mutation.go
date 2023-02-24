package graphql

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/digimakergo/digimaker/core/contenttype"
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
					result := make([]inputData, len(v.Values))
					for i, v := range v.Values {
						if object, ok := v.(*ast.ObjectValue); ok {
							updateData := inputData{data: make(handler.InputMap)}
							if err := parseInput(object.Fields, &updateData, def.FieldMap); err != nil {
								return fmt.Errorf("failed to parse input data: %w", err)
							}
							result[i] = updateData
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
	result := []contenttype.ContentTyper{}
	if inputs, ok := p.Args["updateData"]; ok {
		switch v := inputs.(type) {
		case []inputData:
			for _, item := range v {
				userId := util.CurrentUserID(ctx)
				if userId == 0 {
					return nil, errors.New("need to login")
				}
				content, err := handler.UpdateByContentID(ctx, p.Info.FieldName, item.id, item.data, userId)
				if err != nil {
					return nil, fmt.Errorf("could not update content by id: %w", err)
				}
				result = append(result, content)
			}
		case error:
			return nil, v
		default:
			return nil, errors.New("could not resolve, unknown updateData")
		}
	}
	return result, nil
}

func parseInput(fields []*ast.ObjectField, input *inputData, fieldMap map[string]fieldtype.FieldDef) error {
	for _, field := range fields {
		key := field.Name.Value
		value := field.Value.GetValue()
		if key == "id" {
			id, err := strconv.Atoi(value.(string))
			if err != nil {
				return fmt.Errorf("could not convert id string to int: %w", err)
			}
			input.id = id
		}
		if key == "data" {
			if fields, ok := value.([]*ast.ObjectField); ok {
				for _, field := range fields {
					key := field.Name.Value
					value := field.Value.GetValue()
					if _, ok := fieldMap[key]; !ok {
						return fmt.Errorf("invalid query field: %s", key)
					}
					input.data[key] = value
				}
			} else {
				return fmt.Errorf("data field should be an object: %#v", value)
			}
		}
	}
	return nil
}

// inputData is user input data.
type inputData struct {
	id   int
	data handler.InputMap
}
