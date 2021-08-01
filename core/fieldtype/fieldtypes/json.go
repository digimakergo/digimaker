package fieldtypes

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
)

type Map map[string]interface{}

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

func init() {
	fieldtype.Register(
		fieldtype.Definition{Name: "map",
			DataType: "fieldtypes.Map",
			Package:  "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes",
			NewHandler: func(def definition.FieldDef) fieldtype.Handler {
				return MapHandler{FieldDef: def}
			}})
}
