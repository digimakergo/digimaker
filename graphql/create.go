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

// For the following schema:
//
//	mutation {
//		createArticle(createData: [{data: {title: "", body: ""}, parent: 3}]) {
//		  id
//		  title
//		}
//	}
func buildCreateType(gqlContentTypes graphql.Fields, cType string, def definition.ContentType, contentGQLType *graphql.Object) {
	cType = "create" + util.UpperName(cType)
	gqlContentTypes[cType] = &graphql.Field{
		Name: def.Name,
		Type: graphql.NewList(contentGQLType),
		Args: buildCreateArgs(cType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return resolveCreateMutation(p.Context, p)
		},
	}
}

func buildCreateArgs(cType string) graphql.FieldConfigArgument {
	result := make(graphql.FieldConfigArgument)
	result["createData"] = &graphql.ArgumentConfig{
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
					result := make([]createInputData, len(v.Values))
					for i, v := range v.Values {
						if object, ok := v.(*ast.ObjectValue); ok {
							createData := createInputData{data: make(handler.InputMap)}
							if err := parseCreateInput(object.Fields, &createData, def.FieldMap); err != nil {
								return fmt.Errorf("failed to parse input data: %w", err)
							}
							result[i] = createData
						}
					}
					return result
				default:
					return errors.New("unknown createData")
				}
			},
		}),
		Description: "Data to be Created.",
	}
	return result
}

func resolveCreateMutation(ctx context.Context, p graphql.ResolveParams) (interface{}, error) {
	var result []contenttype.ContentTyper
	if inputs, ok := p.Args["createData"]; ok {
		switch v := inputs.(type) {
		case []createInputData:
			result = make([]contenttype.ContentTyper, len(v))
			for i, item := range v {
				userId := util.CurrentUserID(ctx)
				if userId == 0 {
					return nil, errors.New("need to login")
				}
				content, err := handler.Create(ctx, userId, p.Info.FieldName, item.data, item.parentID)
				if err != nil {
					return nil, fmt.Errorf("could not update content by id: %w", err)
				}
				result[i] = content
			}
		case error:
			return nil, v
		default:
			return nil, errors.New("could not resolve, unknown createData")
		}
	}
	return result, nil
}

func parseCreateInput(fields []*ast.ObjectField, input *createInputData, fieldMap map[string]fieldtype.FieldDef) error {
	for _, field := range fields {
		key := field.Name.Value
		value := field.Value.GetValue()
		if key == "parent" {
			parentID, err := strconv.Atoi(value.(string))
			if err != nil {
				return fmt.Errorf("could not convert parent id string to int: %w", err)
			}
			input.parentID = parentID
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

type createInputData struct {
	parentID int
	data     handler.InputMap
}
