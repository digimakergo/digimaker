//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package contenttype

import (
	"dm/dm/fieldtype"
	"dm/dm/util"
)

type ContentTypeSettings map[string]ContentTypeSetting

type ContentTypeSetting struct {
	Name          string                  `json:"name"`
	TableName     string                  `json:"table_name"`
	HasVersion    bool                    `json:"has_version"`
	HasLocation   bool                    `json:"has_location"`
	FieldsDisplay []string                `json:"fields_display"`
	AllowedTypes  []string                `json:"allowed_types"`
	Fields        map[string]ContentField `json:"fields"`
}

type ContentField struct {
	Name       string            `json:"name"`
	FieldType  string            `json:"type"`
	Required   bool              `json:"required"`
	Parameters map[string]string `json:"parameters"`
}

func (f *ContentField) GetDefinition() fieldtype.FieldtypeSetting {
	return fieldtype.GetDefinition()[f.FieldType]
}

//ContentTypeDefinition Content types which defined in contenttype.json
var contentTypeDefinition ContentTypeSettings

//LoadDefinition Load all setting in file into memory.
//
// It will not load anything unless all json' format matches the struct definition.
//
func LoadDefinition() error {

	//Load contenttype.json into ContentTypeDefinition
	var contentDef map[string]ContentTypeSetting
	err := util.UnmarshalData(util.ConfigPath()+"/contenttype.json", &contentDef)
	if err != nil {
		return err
	}

	contentTypeDefinition = contentDef

	return nil
}

func GetDefinition() ContentTypeSettings {
	return contentTypeDefinition
}

//todo: Use a better name
func GetContentDefinition(contentType string) ContentTypeSetting {
	definition := contentTypeDefinition
	result, ok := definition[contentType]
	if ok {
		return result
	} else {
		return ContentTypeSetting{}
	}
}
