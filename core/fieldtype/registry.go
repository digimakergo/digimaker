//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

import "dm/core/util"

//global variable for registering handlers
//A handler is always singleton
var handlerRegistry = map[string]FieldtypeHandlerI{}

func RegisterHandler(identifier string, handler FieldtypeHandlerI) {
	util.Log("system", "Registering handler for field type "+identifier)
	handlerRegistry[identifier] = handler
}

func GetHandler(fieldType string) FieldtypeHandlerI {
	return handlerRegistry[fieldType]
}

//Global variable for registering fieldtypes
//Use call back to make sure it's not the same instance( the receiver can still singleton it )
var fieldtypeRegistry = map[string]func() Fieldtyper{}

func RegisterField(fieldType string, newFieldType func() Fieldtyper) {
	fieldtypeRegistry[fieldType] = newFieldType
}

func NewFieldType(fieldType string) Fieldtyper {
	return fieldtypeRegistry[fieldType]()
}

type FieldTypeSettings map[string]FieldtypeSetting

type RelationSetting struct {
	DataFields  string `json:"data_fields"`
	DataPattern string `json:"data_pattern"`
}

type FieldtypeSetting struct {
	Identifier       string            `json:"identifier"`
	Name             string            `json:"name"`
	Searchable       bool              `json:"searchable"`
	Value            string            `json:"value"`
	Translations     map[string]string `json:"translations"`
	IsRelation       bool              `json:"is_relation"`
	RelationSettings RelationSetting   `json:"relation_settings"`
}

// Datatypes which defined in datatype.json
var fieldtypeDefinition FieldTypeSettings

func LoadDefinition() error {
	//Load datatype.json into DatatypeDefinition
	var def map[string]FieldtypeSetting
	err := util.UnmarshalData(util.ConfigPath()+"/fieldtype.json", &def)
	if err != nil {
		return err
	}
	fieldtypeDefinition = def
	return nil
}

func GetDefinition() FieldTypeSettings {
	return fieldtypeDefinition
}

func GetFieldTypeDef(identifier string) FieldtypeSetting {
	return fieldtypeDefinition[identifier]
}
