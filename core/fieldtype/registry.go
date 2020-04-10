//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

import (
	"dm/core/util"
	"dm/core/log"
	"fmt"
)

//global variable for registering handlers
//A handler is always singleton
var handlerRegistry = map[string]FieldHandler{}

func RegisterHandler(identifier string, handler FieldtypeHandler) {
	log.Info("Registering handler for field type " + identifier)
	f := FieldHandler{}
	f.Fieldtype = identifier
	f.handler = handler
	handlerRegistry[identifier] = f
}

func GetHandler(fieldType string) FieldHandler {
	return handlerRegistry[fieldType]
}

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
	IsContainer      bool              `json:"is_container"`
	RelationSettings RelationSetting   `json:"relation_settings"`
}

// Datatypes which defined in datatype.json
var fieldtypeDefinition map[string]FieldtypeSetting

func LoadDefinition() error {
	//Load datatype.json into DatatypeDefinition
	var defMap map[string]FieldtypeSetting
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

func GetDefinition(identifier string) FieldtypeSetting {
	return fieldtypeDefinition[identifier]
}

func GetAllDefinition() map[string]FieldtypeSetting {
	return fieldtypeDefinition
}
