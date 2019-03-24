package types

import (
	content "dm/type_default/content"
	field "dm/type_default/field"
)

type TypeLoader interface {
	Instance(extendedType string, identifier string) interface{}
	FieldTypeList() []string
	ContentTypeList() []string
}

//TypeLoaderDefault implements FieldInstancer and ContentTypeInstancer
type TypeLoaderDefault struct{}

func (TypeLoaderDefault) Instance(extendedType string, identifier string) interface{} {
	var result interface{}
	if extendedType == "field" {
		switch identifier {
		case "text":
			result = new(field.TextField)
		case "richtext":
			result = new(field.RichTextField)
		default:
		}
	} else if extendedType == "contenttype" {
		switch identifier {
		case "article":
			result = content.Article{}
		default:
		}
	}

	return result
}

func (TypeLoaderDefault) FieldTypeList() []string {
	return []string{"text", "richtext"}
}

func (TypeLoaderDefault) ContentTypeList() []string {
	return []string{"article", "folder"}
}
