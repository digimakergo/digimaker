//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package type_default

import (
	//content "dm/type_default/content"
	"dm/model"
	field "dm/type_default/field"
)

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
			//result = content.Article{}
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

func NewFieldType(fieldType string) model.Fielder {
	var result model.Fielder
	switch fieldType {
	case "text":
		result = field.TextField{}
	case "richtext":
		result = field.RichTextField{}
	}
	return result
}

func NewHandler(fieldType string) model.FieldtypeHandler {
	var result model.FieldtypeHandler
	switch fieldType {
	case "text":
		result = field.TextFieldHandler{}
	default:
		result = field.TextFieldHandler{}
	}
	return result
}
