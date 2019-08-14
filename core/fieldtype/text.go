//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

//Package fieldtype implements build-in field types(value and fieldtype handler).
package fieldtype

//TextField is a field for normal text line. It implements Datatyper
type TextField struct {
	FieldtypeValue
}

func (t *TextField) Scan(src interface{}) error {
	err := t.SetData(src, "text")
	return err
}

//convert data to view data.
func (t TextField) ViewValue() string {
	return t.Raw
}

//implement FieldtypeHandler
type TextFieldHandler struct {
	*FieldtypeHandler
}

func init() {
	RegisterHandler("text", TextFieldHandler{})
}
