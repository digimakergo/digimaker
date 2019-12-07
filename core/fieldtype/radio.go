//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

//Package fieldtype implements build-in field types(value and fieldtype handler).
package fieldtype

import (
	"strings"
)

//TextField is a field for normal text line. It implements Datatyper
type RadioField struct {
	FieldValue
}

func (t *RadioField) Scan(src interface{}) error {
	err := t.SetData(src, "radio")
	return err
}

//convert data to view data.
func (t RadioField) ViewValue() string {
	return t.Raw
}

//implement FieldtypeHandler
type RadioHandler struct{}

func (t RadioHandler) Validate(input interface{}) (bool, string) {
	//todo: support int
	s := input.(string)
	if s != "" {
		if s != "-1" && s != "1" && s != "0" {
			return false, "Invalid radio value."
		}
	}
	return true, ""
}

func (t RadioHandler) NewValueFromInput(input interface{}) interface{} {
	r := TextField{}
	r.Scan(input.(string))
	return r
}

func (t RadioHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("radio", RadioHandler{})
}
