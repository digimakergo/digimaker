package fieldtypes

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
)

//Map defines a key - value map
type Map map[string]interface{}

//List defines an array Map
type MapList []Map

func (a *Map) Scan(value interface{}) error {
	obj := Map{}
	if value != nil {
		err := json.Unmarshal(value.([]byte), &obj)
		if err != nil {
			return err
		}
		*a = obj
	} else {
		*a = nil
	}
	return nil
}

func (a Map) Value() (driver.Value, error) {
	value, err := json.Marshal(a)
	return value, err
}

func (a *MapList) Scan(value interface{}) error {
	obj := MapList{}
	if value != nil {
		err := json.Unmarshal(value.([]byte), &obj)
		if err != nil {
			return err
		}
		*a = obj
	} else {
		*a = nil
	}
	return nil
}

func (a MapList) Value() (driver.Value, error) {
	data, err := json.Marshal(a)
	return data, err
}

//MapHandler
type MapHandler struct {
	definition.FieldDef
}

func (handler MapHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	if _, ok := input.(Map); ok {
		return input, nil
	}
	m := Map{}
	err := unmarshalInput(input, &m)
	return m, err
}

func (handler MapHandler) DBField() string {
	return "JSON"
}

//MapListHandler
type MapListHandler struct {
	definition.FieldDef
}

func (handler MapListHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	if _, ok := input.(MapList); ok {
		return input, nil
	}
	m := MapList{}
	err := unmarshalInput(input, &m)
	return m, err
}

func (handler MapListHandler) DBField() string {
	return "JSON"
}

func unmarshalInput(input interface{}, target interface{}) error {
	var data []byte
	switch input.(type) {
	case string:
		dataStr := input.(string)
		if dataStr == "" {
			return nil
		}
		data = []byte(dataStr)
		break
	case []byte:
		data = input.([]byte)
		break
	default:
		var err error
		data, err = json.Marshal(input)
		if err != nil {
			return fieldtype.NewValidationError("Not a valid json: " + err.Error())
		}
	}
	err := json.Unmarshal(data, target)
	return err
}

//JSON Handler
type JSONHandler struct {
	definition.FieldDef
}

//support string, []byte, object
func (handler JSONHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	if input == nil {
		return []byte{}, nil
	}

	var data []byte
	switch input.(type) {
	case string:
		dataStr := input.(string)
		if dataStr == "" {
			return []byte{}, nil
		}
		data = []byte(dataStr)
		break
	case []byte:
		data = input.([]byte)
		break
	default:
		var err error
		data, err = json.Marshal(input)
		if err != nil {
			return nil, fieldtype.NewValidationError("Not a valid json: " + err.Error())
		}
	}

	isValid := json.Valid(data)
	if !isValid {
		return "", fieldtype.NewValidationError("Not a valid json")
	}

	return string(data), nil
}

func (handler JSONHandler) DBField() string {
	return "JSON"
}

func init() {
	fieldtype.Register(
		fieldtype.Definition{Name: "map",
			DataType: "fieldtypes.Map",
			Package:  "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes",
			NewHandler: func(def definition.FieldDef) fieldtype.Handler {
				return MapHandler{FieldDef: def}
			}})
	fieldtype.Register(fieldtype.Definition{Name: "maplist",
		DataType: "fieldtypes.MapList",
		Package:  "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes",
		NewHandler: func(def definition.FieldDef) fieldtype.Handler {
			return MapListHandler{FieldDef: def}
		}})
	fieldtype.Register(fieldtype.Definition{Name: "json",
		DataType: "string",
		NewHandler: func(def definition.FieldDef) fieldtype.Handler {
			return JSONHandler{FieldDef: def}
		}})
}
