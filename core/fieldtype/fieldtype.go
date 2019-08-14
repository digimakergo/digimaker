//Author xc, Created on 2019-03-26 20:44
//{COPYRIGHTS}

package fieldtype

import (
	"database/sql/driver"
	"errors"
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
