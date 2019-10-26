//Author xc, Created on 2019-08-13 23:23
//{COPYRIGHTS}

package fieldtype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
)

type FieldValue struct {
	Raw        string
	Output     string //string
	Definition FieldtypeSetting
}

//when update value to db
func (t FieldValue) Value() (driver.Value, error) {
	fmt.Println("GOOOOd----------------")
	fmt.Println(t.Definition.Identifier)
	fmt.Println(t.Raw)
	return t.Raw, nil
}

//when binding data from db
func (t *FieldValue) SetData(src interface{}, fieldtype string) error {
	if t != nil {
		t.Definition = GetDefinition(fieldtype)
		var data string
		switch src.(type) {
		case string:
			data = src.(string)
		case []byte:
			data = string(src.([]byte))
		case int64:
			data = strconv.FormatInt(src.(int64), 10)
		case int:
			data = strconv.Itoa(src.(int))
		default:
			return errors.New("Incompatible type")
		}
		t.Raw = data
	}
	return nil
}

//return data from db(raw)
func (t FieldValue) Data() string {
	return t.Raw
}
