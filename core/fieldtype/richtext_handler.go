package fieldtype

import (
	"strings"
)

//implement FieldtypeHandler
type RichTextFieldHandler struct {
	*FieldtypeHandler
}

func (t RichTextFieldHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t RichTextFieldHandler) ToStorage(input interface{}) interface{} {
	r := RichTextField{}
	r.Raw = input.(string)
	return r
}

func (t RichTextFieldHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("richtext", RichTextFieldHandler{})
}
