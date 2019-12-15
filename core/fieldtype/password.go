//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

//Package fieldtype implements build-in field types(value and fieldtype handler).
package fieldtype

import (
	"dm/core/util"
	"strings"
)

//TextField is a field for normal text line. It implements Datatyper
type PasswordField struct {
	FieldValue
}

func (t *PasswordField) Scan(src interface{}) error {
	err := t.SetData(src, "password")
	return err
}

//convert data to view data.
func (t PasswordField) ViewValue() string {
	return t.Raw
}

//implement FieldtypeHandler
type PasswordHandler struct{}

func (t PasswordHandler) Validate(input interface{}) (bool, string) {
	str := input.(string)
	if len(str) < 8 {
		return false, "Password needs to be more than 8 characters." //todo: more rule based on configuration.
	}
	return true, ""
}

func (t PasswordHandler) NewValueFromInput(input interface{}) interface{} {
	r := PasswordField{}
	str := input.(string)
	hash, _ := util.HashPassword(str)
	r.Scan(hash)
	return r
}

func (t PasswordHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("password", PasswordHandler{})
}
