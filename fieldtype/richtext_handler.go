package fieldtype

import (
	"strings"
)

//implement FieldtypeHandler
type RichTextFieldHandler struct {
}

func (t RichTextFieldHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t RichTextFieldHandler) ToStorage(input interface{}) interface{} {
	return RichTextField{Data: input.(string)}
}

func (t RichTextFieldHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHanlder("richtext", RichTextFieldHandler{})
}
