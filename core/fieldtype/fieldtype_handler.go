//Author xc, Created on 2019-08-13 23:23
//{COPYRIGHTS}

package fieldtype

import "strings"

type FieldtypeHandler struct {
	Fieldtype string
}

func (t FieldtypeHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t FieldtypeHandler) ToStorage(input interface{}) interface{} {
	r := RichTextField{}
	r.Raw = input.(string)
	return r
}

func (t FieldtypeHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func (t *FieldtypeHandler) SetIdentifier(identifier string) {
	t.Fieldtype = identifier
}

func (t FieldtypeHandler) Definition() FieldtypeSetting {
	return GetFieldTypeDef(t.Fieldtype)
}
