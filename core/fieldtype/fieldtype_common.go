//Author xc, Created on 2019-08-13 23:23
//{COPYRIGHTS}

package fieldtype

import (
	"database/sql/driver"
	"errors"
	"strings"
)

//All of the fields will implements this interface
type Fieldtyper interface {
	//Get value of
	//Value() string

	//Create()
	//Validate()
	//SetStoreData()
}

type FieldtypeHandlerI interface {
	ToStorage(input interface{}) interface{}
	Validate(input interface{}) (bool, string)
	IsEmpty(input interface{}) bool
	SetIdentifier(identifier string)
	Definition() FieldtypeSetting
}

//Relation field handler can convert relations into RelationField
type RelationFieldHandler interface {
	ToStorage(contents interface{}) interface{}
	UpdateOne(toContent interface{}, identifier string, from interface{})
}

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
		t.Definition = GetDefinition(fieldtype)
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
	return GetDefinition(t.Fieldtype)
}
