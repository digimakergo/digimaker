//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

//Package fieldtype implements build-in field types(value and fieldtype handler).
package fieldtype

import "strings"

//TextField is a field for normal text line. It implements Datatyper
type TextField struct {
	FieldValue
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
type TextHandler struct{}

func (t TextHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t TextHandler) NewValueFromInput(input interface{}) interface{} {
	r := TextField{}
	r.Scan(input.(string))
	return r
}

func (t TextHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("text", TextHandler{})
        RegisterHandler("eth_indicator", TextHandler{})
}
