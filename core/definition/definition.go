//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package definition

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/digimakergo/digimaker/core/config"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var LocationColumns []string = []string{
	"id",
	"parent_id",
	"main_id",
	"hierarchy",
	"depth",
	"identifier_path",
	"content_type",
	"content_id",
	"identifier_path",
	"name",
	"is_hidden",
	"is_invisible",
	"priority",
	"uid",
	"section",
	"p",
}

type ContentTypeList map[string]map[string]ContentType

//ValidationRule defines rule for a field's validation. eg. max length
type VaidationRule map[string]interface{}

type FieldParameters map[string]interface{}

type ContentType struct {
	Name          string      `json:"name"`
	TableName     string      `json:"table_name"`
	RelationData  []string    `json:"relation_data"`
	NamePattern   string      `json:"name_pattern"`
	HasVersion    bool        `json:"has_version"`
	HasLocation   bool        `json:"has_location"`
	HasLocationID bool        `json:"has_location_id"` //for non-location content
	Fields        []FieldDef  `json:"fields"`
	DataFields    []DataField `json:"data_fields"`
	//All fields where identifier is the key.
	FieldMap            map[string]FieldDef `json:"-"`
	FieldIdentifierList []string            `json:"-"`

	hasRelationlist int //cache of hasRelationlist, 1(yes), -1(no), 0(not set)
}

func (c *ContentType) HasDataField(identifier string) bool {
	result := false
	for _, item := range c.DataFields {
		if item.Identifier == identifier {
			result = true
			break
		}
	}
	return result
}

func (c ContentType) HasRelationlist() bool {
	if c.hasRelationlist == 0 {
		hasRelationList := false
		for _, field := range c.FieldMap {
			if field.FieldType == "relationlist" {
				hasRelationList = true
				break
			}
		}
		if hasRelationList {
			c.hasRelationlist = 1
		} else {
			c.hasRelationlist = -1
		}
	}

	if c.hasRelationlist == 1 {
		return true
	} else {
		return false
	}
}

func (c *ContentType) Init(fieldCallback ...func(*FieldDef)) {
	//set all fields into FieldMap
	fieldMap := map[string]FieldDef{}
	identifierList := []string{}
	for i, field := range c.Fields {
		identifier := field.Identifier
		if len(fieldCallback) > 0 {
			fieldCallback[0](&field)
		}

		fieldMap[identifier] = field
		identifierList = append(identifierList, identifier)
		//get sub fields
		subFields := field.GetSubFields(fieldCallback...)
		c.Fields[i] = field
		for subIdentifier, subField := range subFields {
			fieldMap[subIdentifier] = subField
			identifierList = append(identifierList, subIdentifier)
		}
	}
	c.FieldMap = fieldMap
	c.FieldIdentifierList = identifierList
}

//Content field definition
type FieldDef struct {
	Identifier   string          `json:"identifier"`
	Name         string          `json:"name"`
	FieldType    string          `json:"type"`
	DefaultValue interface{}     `json:"default_value"` //eg. checkbox 1 means checked
	Required     bool            `json:"required"`
	Description  string          `json:"description"`
	IsOutput     bool            `json:"is_output"`
	Parameters   FieldParameters `json:"parameters"`
	Children     []FieldDef      `json:"children"`
}

type DataField struct {
	Identifier string `json:"identifier"`
	FieldType  string `json:"fieldtype"`
	Name       string `json:"name"`
}

func (cf *FieldDef) GetSubFields(callback ...func(*FieldDef)) map[string]FieldDef {
	return getSubFields(cf, callback...)
}

func getSubFields(cf *FieldDef, callback ...func(*FieldDef)) map[string]FieldDef {
	if cf.Children == nil {
		return nil
	} else {
		result := map[string]FieldDef{}
		for i, field := range cf.Children {
			identifier := field.Identifier
			if len(callback) > 0 {
				callback[0](&field)
			}

			//get children under child
			children := getSubFields(&field, callback...)
			cf.Children[i] = field
			for _, item := range children {
				identifier2 := item.Identifier
				result[identifier2] = item
			}
			result[identifier] = field
		}
		return result
	}
}

//ContentTypeDefinition Content types which defined in contenttype.json
var contentTypeDefinition ContentTypeList

//LoadDefinition Load all setting in file into memory.
func LoadDefinition() error {
	//Load contenttype.json into ContentTypeDefinition
	var def map[string]ContentType
	err := util.UnmarshalData(config.ConfigPath()+"/contenttype.json", &def)
	if err != nil {
		return err
	}

	for identifier, _ := range def {
		cDef := def[identifier]
		cDef.Init()
		def[identifier] = cDef
	}
	contentTypeDefinition = map[string]map[string]ContentType{"default": def}

	//todo: use config or scan folder.
	//todo nb: the translation can be add later but listener should be there
	loadTranslations([]string{"nor-NO", "eng-GB"})

	return nil
}

//load translation based on existing definition
//todo: use locale folder
func loadTranslations(languages []string) {
	for _, language := range languages {
		var def map[string]ContentType
		util.UnmarshalData(config.ConfigPath()+"/contenttype.json", &def)

		//todo: formalize this: use folder, and loop through language
		viper := viper.New()
		viper.AddConfigPath(config.ConfigPath())
		filename := "contenttype_" + language
		viper.SetConfigName(filename)
		viper.ReadInConfig()
		viper.WatchConfig()

		translationObj := viper.AllSettings()
		str, _ := json.Marshal(translationObj)
		var translation = map[string][]map[string]string{}
		json.Unmarshal(str, &translation)
		viper.OnConfigChange(func(e fsnotify.Event) {
			if e.Name == config.ConfigPath()+"/"+filename+".json" {
				log.Println("Translation changed. file: " + filename)
				translationObj := viper.AllSettings()
				str, _ = json.Marshal(translationObj)
				json.Unmarshal(str, &translation)
				loadTranslation(def, translation)
			}
		})
		loadTranslation(def, translation)

		//set language related definition
		contentTypeDefinition[language] = def
	}

}

func loadTranslation(def map[string]ContentType, translation map[string][]map[string]string) {
	for contenttype, contenttypeDef := range def {
		translist, ok := translation[contenttype]
		if !ok {
			continue
		}
		origFields := contentTypeDefinition["default"][contenttype].FieldMap
		contenttypeDef.Init(func(field *FieldDef) {
			//translate name
			context := "field/" + field.Identifier + "/name"
			value := getTranslation(context, translist)
			origField := origFields[field.Identifier]

			if value != "" {
				field.Name = value
			} else {
				field.Name = origField.Name
			}

			//translate description
			context = "field/" + field.Identifier + "/description"
			value = getTranslation(context, translist)
			if value != "" {
				field.Description = value
			} else {
				field.Description = origField.Description
			}

			//translate parameters
			for key, param := range field.Parameters {
				switch param.(type) {
				case string:
					value = getTranslation("field/"+field.Identifier+"/parameters/"+key, translist)
					if value != "" {
						field.Parameters[key] = value
					} else {
						field.Parameters[key] = origField.Parameters[key]
					}
					break
				}
			}

		})
		def[contenttype] = contenttypeDef
	}
}

func getTranslation(context string, translist []map[string]string) string {
	value := ""
	for i := range translist {
		if translist[i]["context"] == context {
			value = translist[i]["translation"]
			break
		}
	}
	return value
}

func GetDefinitionList() ContentTypeList {
	return contentTypeDefinition
}

//Get a definition of a contenttype
func GetDefinition(contentType string, language ...string) (ContentType, error) {
	languageStr := "default"
	if len(language) > 0 {
		languageStr = language[0]
	}
	definition, ok := contentTypeDefinition[languageStr]
	if !ok {
		log.Println("Language " + languageStr + " doesn't exist. use default.")
		definition = contentTypeDefinition["default"]
	}
	result, ok := definition[contentType]
	if ok {
		return result, nil
	} else {
		return ContentType{}, errors.New("Content type doesn't exist: " + contentType)
	}
}

//Get fields based on path pattern including container,
//separated by /
//. eg. article/relations, report/step1
func GetFields(typePath string) (map[string]FieldDef, error) {
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
		var currentField FieldDef
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
			found := false
			for _, field := range currentField.Children {
				if field.Identifier == name {
					currentField = field
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New(name + " doesn't exist.")
			}
		}

		if currentField.FieldType != "container" {
			return nil, errors.New("End field is not a container")
		}

		//get subfields of end level
		return currentField.GetSubFields(), nil
	}
}

func init() {
	err := LoadDefinition()
	if err != nil {
		log.Fatal(err.Error())
	}
}
