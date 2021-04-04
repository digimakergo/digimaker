package fieldtype

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
)

//All the definition of fieldtypes

//Definition includes a fieldtype basic information
type Definition struct {
	Name       string                            //eg. text
	DataType   string                            //eg. string or eg "fieldtype.CustomString"
	Package    string                            //empty if there is no additional package, otherwise it's like 'mycompany.fieldtype'. Used for generating entity's import.
	NewHandler func(definition.FieldDef) Handler //callback to create new handler for this fieldtype
}

//Fieldtyper is a implementation of a fieldtype, including main logic
type Handler interface {
	//Load from input, should return the value of BaseType, (eg. int), return error or validation error or empty error
	LoadInput(input interface{}, mode string) (interface{}, error)

	//output database field. todo: can support this to generate database field automatically
	DBField() string
}

type ValidationError struct {
	Message string
}

//Validation error
func (err ValidationError) Error() string {
	return err.Message
}

func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

//Empty Error
type EmptyError struct {
}

func (err EmptyError) Error() string {
	return "Field is empty"
}

//BeforeSaving is implemented when fieldtype has event before saving and transaction starts.
//todo: add a function to rollback if failed.
type Event interface {
	BeforeStoring(value interface{}, existing interface{}, mode string) (interface{}, error)
}

//Ouput is implemented when fieldtype needs converting when outputting
type Outputer interface {
	Ouput(value interface{}, params map[string]interface{}) interface{}
}

var fieldtypeMap map[string]Definition = map[string]Definition{}

//Register registers a fieldtype
func Register(definition Definition) {
	name := definition.Name
	if _, ok := fieldtypeMap[name]; ok {
		log.Warning("Fieldtype has been previous registered: "+name, "system")
	}
	log.Info("Registering fieldtype: " + name)
	fieldtypeMap[definition.Name] = definition
}

//GetFieldtype return a fieldtype
func GetFieldtype(fieldtype string) Definition {
	result, ok := fieldtypeMap[fieldtype]
	if ok {
		return result
	} else {
		log.Error("Field type doesn't exist: "+fieldtype, "system")
		return Definition{}
	}
}

//GetAllFieldtype get all fieldtype
func GetAllFieldtype() map[string]Definition {
	return fieldtypeMap
}

func GethHandler(def definition.FieldDef) Handler {
	fieldtypeStr := def.FieldType
	fieldtypeDef := GetFieldtype(fieldtypeStr)
	if fieldtypeDef.NewHandler == nil {
		return nil
	}
	handler := fieldtypeDef.NewHandler(def)
	return handler
}

// IsEmptyInput returns if it's an empty input.
// The input can be not string(eg. int - definately not empty)
func IsEmptyInput(input interface{}) bool {
	if input == nil {
		return true
	}
	s := fmt.Sprint(input)
	s = strings.TrimSpace(s)
	return s == ""
}

type RelationParameters struct {
	Type      string                 `json:"type"`
	Value     string                 `json:"value"`
	Condition map[string]interface{} `json:"condition"`
}

//Convert to parameters obj
func ConvertRelationParams(params definition.FieldParameters) (RelationParameters, error) {
	paramsData, _ := json.Marshal(params)
	rParams := RelationParameters{}
	err := json.Unmarshal(paramsData, &rParams)
	if err != nil {
		return rParams, errors.New("Wrong definition of parameters:" + err.Error())
	}
	return rParams, nil
}
