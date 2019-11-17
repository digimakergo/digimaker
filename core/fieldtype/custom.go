//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

//Package fieldtype implements build-in field types(value and fieldtype handler).
package fieldtype

import (
	"encoding/json"
)

//TextField is a field for normal text line. It implements Datatyper
type CustomField struct {
	FieldValue
}

func (t *CustomField) Scan(src interface{}) error {
	err := t.SetData(src, "custom")
	return err
}

//convert data to view data.
func (t CustomField) ViewValue() string {
	return t.Raw
}

//implement FieldtypeHandler
type CustomHandler struct{}

func (t CustomHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t CustomHandler) NewValueFromInput(input interface{}) interface{} {
	r := CustomField{}
	str, _ := json.Marshal(input)
	r.Scan(string(str))
	return r
}

func (t CustomHandler) IsEmpty(input interface{}) bool {
	return false
}

func init() {
	RegisterHandler("custom", CustomHandler{})
}
