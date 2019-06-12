package fieldtype

import (
	"strings"
)

//implement FieldtypeHandler
type EditorListHandler struct {
}

func (t EditorListHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t EditorListHandler) ToStorage(input interface{}) interface{} {
	return EditorList{Data: input.(string)}
}

func (t EditorListHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHanlder("editorlist", EditorListHandler{})
}
