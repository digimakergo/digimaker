package graphql

import (
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// Digimaker scalar type
var DMScalarType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "DMScalarType",
	Description: "Digimaker scalar type.",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case fieldtypes.Json:
			// result, _ := dmeditor.ProceedData(context.Background(), value)
			return value
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
	gqlFields := graphql.Fields{"id": {
		Type: DMScalarType,
		Name: "ID",
	}}
	for fieldIdentifier, fieldDef := range def.FieldMap {
		//set fields to gqlFields
		gqlField := graphql.Field{}
		gqlField.Type = DMScalarType
		gqlField.Name = fieldDef.Name
		gqlFields[fieldIdentifier] = &gqlField
	}

	// Metadata.
	metadataFields := graphql.Fields{}
	for _, metaField := range definition.MetaColumns {
		metadataFields[metaField] = &graphql.Field{
			Type: DMScalarType,
		}
	}
	gqlFields["metadata"] = &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   "MetadataFields",
			Fields: metadataFields,
		}),
	}

	// Location.
	locationFields := graphql.Fields{}
	for _, locationField := range definition.LocationColumns {
		locationFields[locationField] = &graphql.Field{
			Type: DMScalarType,
		}
	}
	gqlFields["location"] = &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   "LocationFields",
			Fields: locationFields,
		}),
	}

	//customized type
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   def.Name,
		Fields: gqlFields,
	})
}
