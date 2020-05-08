//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

import (
	"fmt"
	"strings"

	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

type funcNewField = func() FieldTyper

//global variable for registering handlers
//A handler is always singleton
var handlerRegistry = map[string]interface{}{}

func RegisterFieldType(newType funcNewField) {
	emptytype := newType()
	fieldtype := emptytype.(FieldTyper).Type()
	log.Info("Registering field type " + fieldtype)
	handlerRegistry[fieldtype] = newType
}

// func RegisterHandler(fieldtype string, handler Handler) {
// 	log.Info("Registering handler for field type " + fieldtype)
// 	handlerRegistry[fieldtype] = handler
// }

func NewField(fieldtype string) FieldTyper {
	newtype := handlerRegistry[fieldtype]
	field := newtype.(funcNewField)()
	return field
}

type RelationSetting struct {
	DataFields  string `json:"data_fields"`
	DataPattern string `json:"data_pattern"`
}

//ValidationRule defines rule for a field's validation. eg. max length
type VaidationRule map[string]interface{}

type FieldtypeDef struct {
	Identifier       string            `json:"identifier"`
	Name             string            `json:"name"`
	HasVariable      bool              `json:"has_variable"`
	Searchable       bool              `json:"searchable"`
	Value            string            `json:"value"`
	Translations     map[string]string `json:"translations"`
	IsRelation       bool              `json:"is_relation"`
	IsContainer      bool              `json:"is_container"`
	RelationSettings RelationSetting   `json:"relation_settings"`
}

// Datatypes which defined in datatype.json
var fieldtypeDefinition map[string]FieldtypeDef

func LoadDefinition() error {
	//Load datatype.json into DatatypeDefinition
	var defMap map[string]FieldtypeDef
	err := util.UnmarshalData(util.ConfigPath()+"/fieldtype.json", &defMap)
	if err != nil {
		return err
	}
	for identifier, setting := range defMap {
		setting.Identifier = identifier
		defMap[identifier] = setting
	}
	fieldtypeDefinition = defMap
	fmt.Println(fieldtypeDefinition["text"])
	return nil
}

func GetDef(fieldtype string) FieldtypeDef {
	return fieldtypeDefinition[fieldtype]
}

func GetAllDefinition() map[string]FieldtypeDef {
	return fieldtypeDefinition
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

//Convert input to nil when input is empty string "".
func EmtpyToNull(input interface{}) interface{} {
	if input == nil {
		return nil
	}
	s, ok := input.(string)
	if ok && s == "" {
		return nil
	}
	return input
}

//Relation field handler can convert relations into RelationField
type RelationFieldHandler interface {
	ToStorage(contents interface{}) interface{}
	UpdateOne(toContent interface{}, identifier string, from interface{})
}
