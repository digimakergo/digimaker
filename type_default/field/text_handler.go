package field

//implement FieldtypeHandler
type TextFieldHandler struct {
	Input interface{}
}

func (t TextFieldHandler) Validate() (bool, string) {
	return true, ""
}

func (t TextFieldHandler) ToStorage() interface{} {
	return t.Input.(string)
}
