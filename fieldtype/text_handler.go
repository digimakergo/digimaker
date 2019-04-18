package fieldtype

import (
	"strings"
)

//implement FieldtypeHandler
type TextFieldHandler struct {
}

func (t TextFieldHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t TextFieldHandler) ToStorage(input interface{}) interface{} {
	return input.(string)
}

func (t TextFieldHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHanlder("text", TextFieldHandler{})
}
