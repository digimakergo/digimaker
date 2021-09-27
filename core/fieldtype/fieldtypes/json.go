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

func (a Map) Value() (driver.Value, error) {
	value, err := json.Marshal(a)
	return value, err
}

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

type MapHandler struct {
	definition.FieldDef
}

//Only allow 1/0 or "1"/"0"
func (handler MapHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	if _, ok := input.(Map); ok {
		return input, nil
	}
	data, _ := json.Marshal(input)
	m := Map{}
	err := json.Unmarshal(data, &m)
	return m, err
}

func (handler MapHandler) DBField() string {
	return "JSON"
}

type MapListHandler struct {
	definition.FieldDef
}

func (handler MapListHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	if _, ok := input.(MapList); ok {
		return input, nil
	}
	data, _ := json.Marshal(input)
	list := MapList{}
	err := json.Unmarshal(data, &list)
	return list, err
}

func (handler MapListHandler) DBField() string {
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
		DataType: "fieldtypes.List",
		Package:  "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes",
		NewHandler: func(def definition.FieldDef) fieldtype.Handler {
			return MapListHandler{FieldDef: def}
		}})
}
