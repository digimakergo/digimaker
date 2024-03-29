package fieldtypes

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query/querier"
)

//Map defines a key - value map
type Map map[string]interface{}

//List defines an array Map
type MapList []Map

var jsonOutputerMap map[string]querier.Outputer = map[string]querier.Outputer{}

func RegisterJSONOutputer(identifier string, outputer querier.Outputer) {
	jsonOutputerMap[identifier] = outputer
}

func (a *Map) Scan(value interface{}) error {
	obj := Map{}
	if value != nil {
		if string(value.([]byte)) != "" {
			err := json.Unmarshal(value.([]byte), &obj)
			if err != nil {
				return err
			}
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
		if string(value.([]byte)) != "" {
			err := json.Unmarshal(value.([]byte), &obj)
			if err != nil {
				return err
			}
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
	fieldtype.FieldDef
}

func (handler MapHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
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
	fieldtype.FieldDef
}

func (handler MapListHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
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

//generaic json

type Json struct {
	Content []byte
}

func (j *Json) Scan(value interface{}) error {
	obj := Json{}
	if value != nil {
		data := value.([]byte)
		if string(data) != "" {
			if json.Valid(data) {
				var jsonData = make([]byte, len(data))
				copy(jsonData, data)
				obj.Content = jsonData
			} else {
				return errors.New("Not a valid json")
			}
		}
	}
	*j = obj
	return nil
}

func (j Json) MarshalJSON() ([]byte, error) {

	if j.Content == nil {
		return json.Marshal(nil)
	}

	//slice
	v := []*json.RawMessage{}
	err := json.Unmarshal(j.Content, &v)
	if err == nil {
		return json.Marshal(v)
	}

	//map
	m := map[string]*json.RawMessage{}
	err = json.Unmarshal(j.Content, &m)
	return json.Marshal(m)
}

//insert as string
func (j Json) Value() (driver.Value, error) {
	if string(j.Content) == "" {
		return nil, nil
	}
	return string(j.Content), nil
}

func (j *Json) String() string {
	return string(j.Content)
}

type JSONParameters struct {
	Format   string `mapstructure:"format"`
	Settings string `mapstructure:"settings"`
}

//JSON Handler
type JSONHandler struct {
	fieldtype.FieldDef
	Params JSONParameters
}

//support string, []byte, object
func (handler JSONHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	obj := Json{}
	if input == nil {
		return obj, nil
	}

	var data []byte
	switch input.(type) {
	case string:
		data = []byte(input.(string))
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
	obj.Content = data

	return obj, nil
}

func (handler JSONHandler) Output(ctx context.Context, querier querier.Querier, value interface{}) interface{} {
	jsonFormat := handler.Params.Format
	if jsonFormat != "" {
		outputer, exist := jsonOutputerMap[jsonFormat]
		if exist {
			return outputer.Output(ctx, querier, value)
		} else {
			log.Warning("Outputer for "+jsonFormat+" doesn't exist. Return raw", "")
			return value
		}
	}
	return value
}

func (handler JSONHandler) DBField() string {
	return "JSON"
}

func init() {
	fieldtype.Register(
		fieldtype.Definition{Name: "map",
			DataType: "fieldtypes.Map",
			Package:  "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return MapHandler{FieldDef: def}
			}})
	fieldtype.Register(fieldtype.Definition{Name: "maplist",
		DataType: "fieldtypes.MapList",
		Package:  "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes",
		NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
			return MapListHandler{FieldDef: def}
		}})
	fieldtype.Register(fieldtype.Definition{Name: "json",
		DataType: "fieldtypes.Json",
		Package:  "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes",
		NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
			params := JSONParameters{}
			err := ConvertParameters(def.Parameters, &params)
			if err != nil {
				log.Error("Definition error on json, parameters ignored: "+err.Error(), "")
			}
			return JSONHandler{FieldDef: def, Params: params}
		}})
}
