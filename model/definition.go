//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package model

import (
	"dm/util"
)

type ContentTypeSetting struct {
	TableName  string                  `json:"table_name"`
	Versioning bool                    `json:"versioning"`
	Fields     map[string]ContentField `json:"fields"`
}

type ContentField struct {
	FieldType string `json:"type"`
	Required  bool   `json:"required"`
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
//
// It will not load anything unless all json' format matches the struct definition.
//
func LoadDefinition() error {

	//Load contenttype.json into ContentTypeDefinition
	var contentDef map[string]ContentTypeSetting
	err := util.UnmarshalData(DMPath+"/configs/"+"contenttype.json", &contentDef)
	if err != nil {
		return err
	}

	//Load datatype.json into DatatypeDefinition
	var datatypeDef map[string]DatatypeSetting
	err = util.UnmarshalData(DMPath+"/configs/"+"datatype.json", &datatypeDef)
	if err != nil {
		return err
	}

	ContentTypeDefinition = contentDef
	DatatypeDefinition = datatypeDef

	return nil
}

func GetContentDefinition(contentType string) ContentTypeSetting {
	definition := ContentTypeDefinition
	result, ok := definition[contentType]
	if ok {
		return result
	} else {
		return ContentTypeSetting{}
	}

}

//ContentTypeDefinition Content types which defined in contenttype.json
var ContentTypeDefinition map[string]ContentTypeSetting

//DatatypeDefinition Datatypes which defined in datatype.json
var DatatypeDefinition map[string]DatatypeSetting

//this is predefined, internal use
var LocationFields = map[string]string{"id": "int",
	"parent_id":    "int",
	"main_id":      "int",
	"content_type": "string",
	"content_id":   "int",
	"language":     "string",
	"name":         "string",
	"section":      "string",
	"remote_id":    "string",
	"p":            "string"}
