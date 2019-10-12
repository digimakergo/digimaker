//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package contenttype

import (
	"dm/core/fieldtype"
	"dm/core/util"
	"encoding/json"
	"errors"
	"strings"
)

type ContentTypeList map[string]ContentType

type ContentType struct {
	Name         string            `json:"name"`
	TableName    string            `json:"table_name"`
	HasVersion   bool              `json:"has_version"`
	HasLocation  bool              `json:"has_location"`
	AllowedTypes []string          `json:"allowed_types"`
	Fields       ContentFieldArray `json:"fields"`
	//All fields where identifier is the key.
	FieldMap map[string]ContentField `json:"-"`
}

func (c *ContentType) Init() {
	//set all fields into FieldMap
	fieldMap := map[string]ContentField{}
	for _, field := range c.Fields {
		identifier := field.Identifier
		fieldMap[identifier] = field
		//get sub fields
		subFields := field.GetSubFields()
		for subIdentifier, subField := range subFields {
			fieldMap[subIdentifier] = subField
		}
	}
	c.FieldMap = fieldMap
}

//Custom type for ContentField Array, to have more functions from the list
type ContentFieldArray []ContentField

func (cfArray ContentFieldArray) GetField(identifier string) (ContentField, bool) {
	for _, field := range cfArray {
		if field.Identifier == identifier {
			return field, true
		}
	}
	return ContentField{}, false
}

//Content field definition
type ContentField struct {
	Identifier  string                 `json:"identifier"`
	Name        string                 `json:"name"`
	FieldType   string                 `json:"type"`
	Required    bool                   `json:"required"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Children    ContentFieldArray      `json:"children"`
}

func (cf *ContentField) GetSubFields() map[string]ContentField {
	return getSubFields(cf)
}

func getSubFields(cf *ContentField) map[string]ContentField {
	if cf.Children == nil {
		return nil
	} else {
		result := map[string]ContentField{}
		for _, field := range cf.Children {
			identifier := field.Identifier
			result[identifier] = field
			//get children under child
			children := getSubFields(&field)
			for _, item := range children {
				identifer := item.Identifier
				result[identifer] = item
			}
		}
		return result
	}
}

func (f *ContentField) GetDefinition() fieldtype.FieldtypeSetting {
	return fieldtype.GetDefinition(f.FieldType)
}

//ContentTypeDefinition Content types which defined in contenttype.json
var contentTypeDefinition ContentTypeList

//LoadDefinition Load all setting in file into memory.
func LoadDefinition() error {
	//Load contenttype.json into ContentTypeDefinition
	var def map[string]ContentType
	err := util.UnmarshalData(util.ConfigPath()+"/contenttype.json", &def)
	if err != nil {
		return err
	}

	for identifier, _ := range def {
		cDef := def[identifier]
		cDef.Init()
		def[identifier] = cDef
	}
	contentTypeDefinition = def

	return nil
}

func GetDefinitionList() ContentTypeList {
	return contentTypeDefinition
}

//Get a definition of a contenttype
func GetDefinition(contentType string) (ContentType, error) {
	definition := contentTypeDefinition
	result, ok := definition[contentType]
	if ok {
		return result, nil
	} else {
		return ContentType{}, errors.New("doesn't exist.")
	}
}

//Get fields based on path pattern including container,
//separated by /
//. eg. article/relations, report/step1
func GetFields(typePath string) (map[string]ContentField, error) {
	arr := strings.Split(typePath, "/")
	def, err := GetDefinition(arr[0])
	if err != nil {
		return nil, err
	}
	if len(arr) == 1 {
		return def.FieldMap, nil
	} else {
		//get first level field
		name := arr[1]
		var currentField ContentField
		for _, field := range def.Fields {
			if field.Identifier == name {
				currentField = field
			}
		}
		if currentField.Identifier == "" {
			return nil, errors.New(name + " doesn't exist.")
		}

		//get end level field
		for i := 2; i < len(arr); i++ {
			name = arr[i]
			field, ok := currentField.Children.GetField(name)
			if !ok {
				return nil, errors.New(name + " doesn't exist.")
			}
			currentField = field
		}

		if currentField.FieldType != "container" {
			return nil, errors.New("End field is not a container")
		}

		//get subfields of end level
		return currentField.GetSubFields(), nil
	}
}

//Content to json, used for internal content storing(eg. version data, draft data )
func ContentToJson(content ContentTyper) (string, error) {
	//todo: use a new tag instead of json(eg. version: 'summary', version: '-' to ignore that.)
	result, err := json.Marshal(content)
	return string(result), err
}

//Json to Content, used for internal content recoving. (eg. versioning, draft)
func JsonToContent(contentJson string, content ContentTyper) error {
	err := json.Unmarshal([]byte(contentJson), content)
	return err
}
