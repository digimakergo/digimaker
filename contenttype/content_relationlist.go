//Author xc, Created on 2019-04-25 21:00
//{COPYRIGHTS}

package contenttype

import (
	"dm/fieldtype"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

//ContentRelations as a struct which is linked into a content.
//The purpose is for binding & access.
type ContentRelationList struct {
	Value map[string][]fieldtype.RelationField `json:"value"`
}

func (relations *ContentRelationList) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Unknow scan value.")
	}

	relationData := "[" + source + "]"

	fmt.Println(relationData)

	var mapObject []map[string]interface{}
	err := json.Unmarshal([]byte(relationData), &mapObject)
	if err != nil {
		return errors.Wrap(err, "Can not convert to Relation. Relation data is not correct: "+source)
	}

	//todo: validate keys and make sure it's pre defined
	value := map[string][]fieldtype.RelationField{}
	for _, item := range mapObject {
		if item["identifier"] != nil {
			identifier := item["identifier"].(string)
			if _, ok := value[identifier]; ok {
				value[identifier] = append(value[identifier], item)
			} else {
				var fieldValue fieldtype.RelationField = item
				fieldValue.Restructure()
				value[identifier] = []fieldtype.RelationField{fieldValue}
			}
		}
	}
	relations.Value = value

	return nil
}
