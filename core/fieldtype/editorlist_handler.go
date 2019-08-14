package fieldtype

//implement FieldtypeHandler
type EditorListHandler struct {
	*FieldtypeHandler
}

func init() {
	RegisterHandler("editorlist", EditorListHandler{})
}
