//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package contenttype

import (
	"dm/core/fieldtype"
	"dm/core/util"
	"errors"
	"strings"
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
	allFields     map[string]ContentField
}

func (c *ContentTypeSetting) GetAllFields() map[string]ContentField {
	if c.allFields == nil {
		result := map[string]ContentField{}
		for identifier, field := range c.Fields {
			result[identifier] = field
			//get sub fields
			_, subFields := field.GetSubFields()
			for identifier, subField := range subFields {
				result[identifier] = subField
			}
		}
		c.allFields = result
	}

	return c.allFields
}

type ContentField struct {
	Name          string                  `json:"name"`
	FieldType     string                  `json:"type"`
	Required      bool                    `json:"required"`
	Parameters    map[string]interface{}  `json:"parameters"`
	Description   string                  `json:"description"`
	ChildrenOrder []string                `json:"children_order"`
	Children      map[string]ContentField `json:"children"`
}

func (cf *ContentField) GetSubFields() ([]string, map[string]ContentField) {
	return getSubFields(cf)
}

func getSubFields(cf *ContentField) ([]string, map[string]ContentField) {
	if cf.Children == nil {
		return []string{}, nil
	} else {
		orderResult := []string{}
		result := map[string]ContentField{}
		for identifier, field := range cf.Children {
			result[identifier] = field
			//get children under child
			order, children := getSubFields(&field)
			for _, item := range order {
				orderResult = append(orderResult, item)
				result[item] = children[item]
			}
		}
		return orderResult, result
	}
}

func (f *ContentField) GetDefinition() fieldtype.FieldtypeSetting {
	return fieldtype.GetDefinition(f.FieldType)
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
func GetContentDefinition(contentType string) (ContentTypeSetting, error) {
	definition := contentTypeDefinition
	result, ok := definition[contentType]
	if ok {
		return result, nil
	} else {
		return ContentTypeSetting{}, errors.New("doesn't exist.")
	}
}

//Get fields based on path pattern including container, separated by /. eg. article/relations, report/step1
func GetFields(typePath string) ([]string, error) {
	arr := strings.Split(typePath, "/")
	def, err := GetContentDefinition(arr[0])
	if err != nil {
		return []string{}, err
	}
	fieldNames := def.FieldsDisplay
	fields := def.Fields
	if len(arr) == 1 {
		return fieldNames, nil
	} else {
		for i := 1; i < len(arr); i++ {
			fieldIdentifier := arr[i]
			if !util.Contains(fieldNames, fieldIdentifier) {
				return []string{}, errors.New(arr[i] + "is not in sub fields.")
			}
			if field, ok := fields[fieldIdentifier]; ok {
				if field.FieldType == "container" {
					fieldsResult := field.Parameters["fields"].([]interface{})
					result := []string{}
					for i := range fieldsResult {
						result = append(result, fieldsResult[i].(string))
					}
					return result, nil
				} else {
					return []string{}, errors.New(fieldIdentifier + "is not a container")
				}
			} else {
				return []string{}, errors.New(fieldIdentifier + "doesn't exist.")
			}
		}
		return fieldNames, nil
	}
}
