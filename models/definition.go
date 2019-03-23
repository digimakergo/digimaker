package models

import "dm/utils"

type ContentTypeSetting struct {
	TableName  string            `json:"table_name"`
	Versioning bool              `json:"versioning"`
	Fields     map[string]string `json:"fields"`
}

type DatatypeSetting struct {
	Identifier   string            `json:"identifier"`
	Name         string            `json:"name"`
	Searchable   bool              `json:"searchable"`
	Translations map[string]string `json:"translations"`
}

//todo: use dynamic way
var DMPath = "/Users/xc/go/caf-prototype/src/dm"

//LoadDefinition Load all setting in file into memory.
func LoadDefinition() error {

	//Load contenttype.json into ContentTypeDefinition
	var contentDef map[string]ContentTypeSetting
	err := utils.UnmarshalData(DMPath+"/configs/"+"contenttype.json", &contentDef)
	if err != nil {
		return err
	}
	ContentTypeDefinition = contentDef

	//Load datatype.json into DatatypeDefinition
	var datatypeDef map[string]DatatypeSetting
	err = utils.UnmarshalData(DMPath+"/configs/"+"datatype.json", &datatypeDef)
	if err != nil {
		return err
	}
	DatatypeDefinition = datatypeDef

	return nil
}

var ContentTypeDefinition map[string]ContentTypeSetting
var DatatypeDefinition map[string]DatatypeSetting
