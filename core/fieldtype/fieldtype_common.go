//Author xc, Created on 2019-08-13 23:23
//{COPYRIGHTS}

package fieldtype

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type FieldtypeValue struct {
	Raw        string
	Output     string           //string
	Definition FieldtypeSetting `identifier:"text"`
}

//when update value to db
func (t FieldtypeValue) Value() (driver.Value, error) {
	return t.Raw, nil
}

//when binding data from db
func (t *FieldtypeValue) SetData(src interface{}, fieldtype string) error {
	if t != nil {
		t.Definition = GetFieldTypeDef(fieldtype)
		var data string
		switch src.(type) {
		case string:
			data = src.(string)
		case []byte:
			data = string(src.([]byte))
		default:
			return errors.New("Incompatible type")
		}
		t.Raw = data
	}
	return nil
}

//return data from db(raw)
func (t FieldtypeValue) Data() string {
	return t.Raw
}

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
