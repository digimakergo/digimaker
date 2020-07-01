//Author xc, Created on 2019-04-25 21:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/util"
)

type RelationList []Relation

//LoadFromInput load data from input before validation
func (rl *RelationList) LoadFromInput(input interface{}) error {
	s := fmt.Sprint(input)

	arr := strings.Split(s, ";")
	if len(arr) != 2 {
		return errors.New("wrong format.")
	}
	arrInt, err := util.ArrayStrToInt(strings.Split(arr[0], ","))

	if err != nil {
		return errors.New("Not int(s)")
	}
	arrType := strings.Split(arr[1], ",")
	if len(arrType) != len(arrInt) {
		return errors.New("id and type are not same length")
	}

	relationlist := []Relation{}
	for i, v := range arrInt {
		r := Relation{}
		r.FromContentID = v
		r.FromType = arrType[i]
		r.Priority = len(arrInt) - i
		relationlist = append(relationlist, r)
	}
	*rl = relationlist
	return nil
}

func (rl RelationList) FieldValue() interface{} {
	return rl
}

func (rl RelationList) Validate(rule fieldtype.VaidationRule) (bool, string) {
	//todo: check if the id exist or has permission
	return true, ""
}

func (rl RelationList) IsEmpty() bool {
	return false
}

func (rl RelationList) Type() string {
	return "relationlist"
}

//ContentRelations as a struct which is linked into a content.
//The purpose is for binding & access.
type ContentRelationList struct {
	Map      map[string]*RelationList `json:"-"`
	List     []Relation               `json:"-"`
	Existing []Relation               `json:"-"`
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

	var relationList []Relation
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

//Convert list into Map
func (relations *ContentRelationList) groupRelations() {
	//todo: validate keys and make sure it's pre defined
	groupedList := map[string]*RelationList{}
	for _, relation := range relations.List {
		if relation.Identifier != "" {
			identifier := relation.Identifier
			fmt.Println(identifier)
			if _, ok := groupedList[identifier]; ok {
				rl := groupedList[identifier]
				*rl = append(*rl, relation)
				groupedList[identifier] = rl
			} else {
				groupedList[identifier] = &RelationList{relation}
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
func (relations *ContentRelationList) GetField(identifier string) *RelationList {
	if relations.Map == nil {
		relations.Map = map[string]*RelationList{}
	}
	if rl, ok := relations.Map[identifier]; ok {
		return rl
	} else {
		relations.Map[identifier] = &RelationList{}
		return relations.Map[identifier]
	}
}

//LoadFromInput load data from input before validation
//This method will be invoked several times. For uninvoked identifier, it will use exising data.
func (relations *ContentRelationList) LoadFromInput(identifier string, input interface{}) error {
	if relationlist, ok := relations.Map[identifier]; ok {
		err := relationlist.LoadFromInput(input)
		if err != nil {
			return err
		}
		relations.Map[identifier] = relationlist
	} else {
		relationlist := RelationList{}
		relationlist.LoadFromInput(input)
		relations.Map[identifier] = &relationlist
	}
	return nil
}
