//Author xc, Created on 2019-04-25 21:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/pkg/errors"
)

//ContentRelations as a struct which is linked into a content.
//The purpose is for binding & access.
type ContentRelationList struct {
	Map      map[string]fieldtype.RelationList `json:"-"`
	List     []fieldtype.Relation              `json:"-"`
	Existing []fieldtype.Relation              `json:"-"`
}

func (cr *ContentRelationList) SetValue(identifier string, value interface{}) error {
	if cr.Map == nil {
		cr.Map = map[string]fieldtype.RelationList{}
	}
	if list, ok := value.(fieldtype.RelationList); ok {
		cr.Map[identifier] = list
	} else {
		return fmt.Errorf("Only RelationList is supported on field %v", identifier)
	}
	return nil
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

	var relationList []fieldtype.Relation
	err := json.Unmarshal([]byte(source), &relationList)
	if err != nil {
		return errors.Wrap(err, "Can not convert to Relation. Relation data is not correct: "+source)
	}

	//Sort by priority (highest will be on top)
	sort.Slice(relationList, func(i int, j int) bool {
		return relationList[i].Priority > relationList[j].Priority
	})

	//If not empty relation(only one record by everything is null/"")
	// if !(len(relationList) == 1 && relationList[0].Identifier == "") {
	relations.List = relationList
	relations.Existing = relationList
	relations.groupRelations()
	// }
	return nil
}

//Convert list into Map when scaning
func (relations *ContentRelationList) groupRelations() {
	//todo: validate keys and make sure it's pre defined
	groupedList := map[string]fieldtype.RelationList{}
	for _, relation := range relations.List {
		if relation.Identifier != "" {
			identifier := relation.Identifier
			if _, ok := groupedList[identifier]; ok {
				rl := groupedList[identifier]
				rl = append(rl, relation)
				groupedList[identifier] = rl
			} else {
				groupedList[identifier] = fieldtype.RelationList{relation}
			}
		}
	}

	relations.Map = groupedList
}

func (relations ContentRelationList) MarshalJSON() ([]byte, error) {
	return json.Marshal(relations.Map)
}

//Get relationlist field, create if it's not in map.
//nb: does't checking if it's a relationlist type.
func (relations ContentRelationList) GetField(identifier string) fieldtype.RelationList {
	if relations.Map == nil {
		return fieldtype.RelationList{}
	} else {
		if rl, ok := relations.Map[identifier]; ok {
			return rl
		} else {
			return fieldtype.RelationList{}
		}
	}
}
