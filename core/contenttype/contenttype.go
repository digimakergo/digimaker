//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/util"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ContentTypeList map[string]map[string]ContentType

type ContentType struct {
	Name         string            `json:"name"`
	TableName    string            `json:"table_name"`
	NamePattern  string            `json:"name_pattern"`
	HasVersion   bool              `json:"has_version"`
	HasLocation  bool              `json:"has_location"`
	AllowedTypes []string          `json:"allowed_types"`
	Fields       ContentFieldArray `json:"fields"`
	DataFields   []DataField       `json:"data_fields"`
	//All fields where identifier is the key.
	FieldMap map[string]FieldDef `json:"-"`
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

func (c *ContentType) Init(fieldCallback ...func(*FieldDef)) {
	//set all fields into FieldMap
	fieldMap := map[string]FieldDef{}
	for i, field := range c.Fields {
		identifier := field.Identifier
		if len(fieldCallback) > 0 {
			fieldCallback[0](&field)
		}

		fieldMap[identifier] = field
		//get sub fields
		subFields := field.GetSubFields(fieldCallback...)
		c.Fields[i] = field
		for subIdentifier, subField := range subFields {
			fieldMap[subIdentifier] = subField
		}
	}
	c.FieldMap = fieldMap
}

//Custom type for FieldDef Array, to have more functions from the list
type ContentFieldArray []FieldDef

func (cfArray ContentFieldArray) GetField(identifier string) (FieldDef, bool) {
	for _, field := range cfArray {
		if field.Identifier == identifier {
			return field, true
		}
	}
	return FieldDef{}, false
}

//Content field definition
type FieldDef struct {
	Identifier  string                    `json:"identifier"`
	Name        string                    `json:"name"`
	FieldType   string                    `json:"type"`
	Required    bool                      `json:"required"`
	Validation  fieldtype.VaidationRule   `json:"validation"`
	Description string                    `json:"description"`
	IsOutput    bool                      `json:"is_output"`
	Parameters  fieldtype.FieldParameters `json:"parameters"`
	Children    ContentFieldArray         `json:"children"`
}

type DataField struct {
	Identifier string `json:"identifier"`
	FieldType  string `json:"fieldtype"`
	Type       string `json:"type"`
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

func (f *FieldDef) GetDefinition() fieldtype.FieldtypeDef {
	return fieldtype.GetDef(f.FieldType)
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
		util.UnmarshalData(util.ConfigPath()+"/contenttype.json", &def)

		//todo: formalize this: use folder, and loop through language
		viper := viper.New()
		viper.AddConfigPath(util.ConfigPath())
		filename := "contenttype_" + language
		viper.SetConfigName(filename)
		viper.ReadInConfig()
		viper.WatchConfig()

		translationObj := viper.AllSettings()
		str, _ := json.Marshal(translationObj)
		var translation = map[string][]map[string]string{}
		json.Unmarshal(str, &translation)
		viper.OnConfigChange(func(e fsnotify.Event) {
			if e.Name == util.ConfigPath()+"/"+filename+".json" {
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
		return ContentType{}, errors.New("doesn't exist.")
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

func MarchallToOutput(content ContentTyper) ([]byte, error) {
	contentMap := content.ToMap()
	for identifier, field := range contentMap {
		switch field.(type) {
		case fieldtype.FieldTyper:
			value := field.(fieldtype.FieldTyper)
			contentMap[identifier] = value.FieldValue()
		}
	}
	//todo: use a new tag instead of json(eg. version: 'summary', version: '-' to ignore that.)
	result, err := json.Marshal(contentMap)
	return result, err
}

//If field has variables, replace variables with real value
func OutputField(field fieldtype.FieldTyper) interface{} {
	value := field.FieldValue()
	def := fieldtype.GetDef(field.Type())
	if def.HasVariable {
		//todo: implement washing variable
	}
	return value
}

//Json to Content, used for internal content recoving. (eg. versioning, draft)
func JsonToContent(contentJson string, content ContentTyper) error {
	err := json.Unmarshal([]byte(contentJson), content)
	return err
}

//Convert content to map
func ContentToMap(content ContentTyper) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
