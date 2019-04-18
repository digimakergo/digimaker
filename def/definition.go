//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package def

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

//todo: use dynamic way
var DMPath = "/Users/xc/go/caf-prototype/src/dm"

//ContentTypeDefinition Content types which defined in contenttype.json
var contentTypeDefinition map[string]ContentTypeSetting

//LoadDefinition Load all setting in file into memory.
//
// It will not load anything unless all json' format matches the struct definition.
//
func LoadDefinition(configPath string) error {

	//Load contenttype.json into ContentTypeDefinition
	var contentDef map[string]ContentTypeSetting
	err := util.UnmarshalData(configPath+"/contenttype.json", &contentDef)
	if err != nil {
		return err
	}

	contentTypeDefinition = contentDef

	return nil
}

func GetContentDefinition(contentType string) ContentTypeSetting {
	definition := contentTypeDefinition
	result, ok := definition[contentType]
	if ok {
		return result
	} else {
		return ContentTypeSetting{}
	}

}

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
