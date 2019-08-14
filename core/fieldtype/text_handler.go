package fieldtype

//implement FieldtypeHandler
type TextFieldHandler struct {
	*FieldtypeHandler
}

func init() {
	RegisterHandler("text", TextFieldHandler{})
}
