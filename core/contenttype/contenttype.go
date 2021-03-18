//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"

	"github.com/digimakergo/digimaker/core/fieldtype"
)

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
