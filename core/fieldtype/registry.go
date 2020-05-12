//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

import (
	"fmt"
	"strings"

	"github.com/xc/digimaker/core/log"
)

type funcNewField = func() FieldTyper

var handlerRegistry = map[string]interface{}{}
var defMap = map[string]FieldtypeDef{}

//ValidationRule defines rule for a field's validation. eg. max length
type VaidationRule map[string]interface{}

type FieldtypeDef struct {
	Type             string            `json:"type"`
	Name             string            `json:"name"`
	HasVariable      bool              `json:"has_variable"`
	Value            string            `json:"value"`
	Import           string            `json:"import"`
	Translations     map[string]string `json:"translations"`
	IsRelation       bool              `json:"is_relation"`
	RelationSettings RelationSetting   `json:"relation_settings"`
}

type RelationSetting struct {
	DataFields  string `json:"data_fields"`
	DataPattern string `json:"data_pattern"`
}

//Relation field handler can convert relations into RelationField
type RelationFieldHandler interface {
	ToStorage(contents interface{}) interface{}
	UpdateOne(toContent interface{}, identifier string, from interface{})
}

func NewField(fieldtype string) FieldTyper {
	newtype := handlerRegistry[fieldtype]
	field := newtype.(funcNewField)()
	return field
}

func RegisterFieldType(def FieldtypeDef, newType funcNewField) {
	fieldtype := def.Type
	log.Info("Registering field type " + fieldtype)
	if _, ok := defMap[fieldtype]; ok {
		log.Warning("Field type "+fieldtype+" exists already. It will be replaced.", "")
	}
	handlerRegistry[fieldtype] = newType
	defMap[fieldtype] = def
}

func GetDef(fieldtype string) FieldtypeDef {
	return defMap[fieldtype]
}

func GetAllDefinition() map[string]FieldtypeDef {
	return defMap
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
