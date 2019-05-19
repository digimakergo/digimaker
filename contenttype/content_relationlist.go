//Author xc, Created on 2019-04-25 21:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"

	"github.com/pkg/errors"
)

//ContentRelations as a struct which is linked into a content.
//The purpose is for binding & access.
type ContentRelationList struct {
	Map  map[string][]Relation `json:"-"`
	List []Relation            `json:"list"`
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

	//relationData := "[" + source + "]"

	var relationList []Relation
	err := json.Unmarshal([]byte(source), &relationList)
	if err != nil {
		return errors.Wrap(err, "Can not convert to Relation. Relation data is not correct: "+source)
	}

	relations.List = relationList
	relations.groupRelations()
	return nil
}

//Convert list into Map
func (relations *ContentRelationList) groupRelations() {
	//todo: validate keys and make sure it's pre defined
	groupedList := map[string][]Relation{}
	for _, relation := range relations.List {
		if relation.Identifier != "" {
			identifier := relation.Identifier
			if _, ok := groupedList[identifier]; ok {
				groupedList[identifier] = append(groupedList[identifier], relation)
			} else {
				groupedList[identifier] = []Relation{relation}
			}
		}
	}
	relations.Map = groupedList
}
