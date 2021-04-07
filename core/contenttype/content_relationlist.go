//Author xc, Created on 2019-04-25 21:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"
	"sort"

	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/pkg/errors"
)

//ContentRelations as a struct which is linked into a content.
//The purpose is for binding & access.
type ContentRelationList map[string]fieldtype.RelationList

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

	var relationList []fieldtype.Relation
	err := json.Unmarshal(source, &relationList)
	if err != nil {
		return errors.Wrap(err, "Can not convert to Relation. Relation data is not correct: "+string(source))
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
func (relations ContentRelationList) groupRelations(list []fieldtype.Relation) ContentRelationList {
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
				relationMap[identifier] = fieldtype.RelationList{relation}
			}
		}
	}

	return relationMap
}

//Get relationlist field, create if it's not in map.
//nb: does't checking if it's a relationlist type.
func (relations ContentRelationList) GetField(identifier string) fieldtype.RelationList {
	if len(relations) == 0 {
		return nil
	} else {
		if rl, ok := relations[identifier]; ok {
			return rl
		} else {
			return nil
		}
	}
}
