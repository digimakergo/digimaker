//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

//TypeLoaderDefault implements FieldInstancer and ContentTypeInstancer
type TypeLoaderDefault struct{}

func (TypeLoaderDefault) Instance(extendedType string, identifier string) interface{} {
	var result interface{}
	if extendedType == "field" {
		switch identifier {
		case "text":
			result = new(TextField)
		case "richtext":
			result = new(RichTextField)
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

func NewFieldType(fieldType string) Fielder {
	var result Fielder
	switch fieldType {
	case "text":
		result = TextField{}
	case "richtext":
		result = RichTextField{}
	}
	return result
}

func NewHandler(fieldType string) FieldtypeHandler {
	var result FieldtypeHandler
	switch fieldType {
	case "text":
		result = TextFieldHandler{}
	default:
		result = TextFieldHandler{}
	}
	return result
}
