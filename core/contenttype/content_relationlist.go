//Author xc, Created on 2019-04-25 21:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"
	"fmt"
	"sort"

	"errors"
)

//ContentRelations as a struct which is linked into a content.
//The purpose is for binding & access.
type ContentRelationList map[string]RelationList

func (relations *ContentRelationList) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Unknow scan value.")
	}

	var relationList []Relation
	err := json.Unmarshal(source, &relationList)
	if err != nil {
		return fmt.Errorf("Can not convert to Relation. Relation data is not correct - %v: %w ", string(source), err)
	}

	//Sort by priority (highest will be on top)
	sort.Slice(relationList, func(i int, j int) bool {
		return relationList[i].Priority > relationList[j].Priority
	})

	groupedRelations := relations.groupRelations(relationList)
	*relations = groupedRelations

	return nil
}

//Convert list into Map when scaning
func (relations ContentRelationList) groupRelations(list []Relation) ContentRelationList {
	//todo: validate keys and make sure it's pre defined
	relationMap := ContentRelationList{}
	for _, relation := range list {
		if relation.Identifier != "" {
			identifier := relation.Identifier
			if _, ok := relationMap[identifier]; ok {
				rl := relationMap[identifier]
				rl = append(rl, relation)
				relationMap[identifier] = rl
			} else {
				relationMap[identifier] = RelationList{relation}
			}
		}
	}

	return relationMap
}
