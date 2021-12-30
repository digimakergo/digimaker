//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"
	"strings"
)

//Content to json, used for internal content storing(eg. version data, draft data )
func ContentToJson(content ContentTyper) (string, error) {
	//todo: use a new tag instead of json(eg. version: 'summary', version: '-' to ignore that.)
	result, err := json.Marshal(content)
	return string(result), err
}

func MarchallToOutput(content ContentTyper) ([]byte, error) {
	contentMap := content.ToMap()
	//todo: use a new tag instead of json(eg. version: 'summary', version: '-' to ignore that.)
	result, err := json.Marshal(contentMap)
	return result, err
}

//Json to Content, used for internal content recoving. (eg. versioning, draft)
func JsonToContent(contentJson string, content ContentTyper) error {
	err := json.Unmarshal([]byte(contentJson), content)
	return err
}

//Convert content to map
func ContentToMap(content ContentTyper) (ContentMap, error) {
	jsonData, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	result := ContentMap{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func IsUnderLocation(subLocation Location, location Location) bool {
	underHierarchy := subLocation.Hierarchy
	hierarchy := location.Hierarchy

	if strings.HasPrefix(hierarchy, underHierarchy) {
		return true
	} else {
		return false
	}
}
