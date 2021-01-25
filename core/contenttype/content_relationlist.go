//Author xc, Created on 2019-04-25 21:00
//{COPYRIGHTS}

package contenttype

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/pkg/errors"
)

//RelationList is field type which is in relations in ContentCommon
type RelationList []Relation

//LoadFromInput load data from input before validation
func (rl *RelationList) LoadFromInput(input interface{}, params fieldtype.FieldParameters) error {
	s := fmt.Sprint(input)
	if s == "" {
		*rl = []Relation{}
		return nil
	}
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

	dbHandler := db.DBHanlder()
	relationlist := []Relation{}
	for i, v := range arrInt {
		fromType := arrType[i]
		fromCid := v
		def, err := GetDefinition(fromType)
		if err != nil {
			return err
		}

		r := Relation{}

		r.FromContentID = fromCid
		r.FromType = fromType
		r.Priority = len(arrInt) - i

		relationDataFields := def.RelationData
		if len(relationDataFields) > 0 {
			//get content
			contents := db.DatamapList{}
			dbHandler.GetEntity(def.TableName, db.Cond("id", fromCid), nil, nil, &contents)
			if len(contents) == 0 {
				return errors.New("No content found on " + strconv.Itoa(fromCid))
			}

			//If there is one field, use it on data, otherwise use json map as data
			if len(relationDataFields) == 1 {
				r.Data = fmt.Sprint(contents[0][relationDataFields[0]])
			} else {
				datamap := map[string]interface{}{}
				for _, field := range relationDataFields {
					datamap[field] = contents[0][field]
				}
				data, _ := json.Marshal(datamap)
				r.Data = string(data)
			}
		}

		relationlist = append(relationlist, r)
	}
	*rl = relationlist
	return nil
	//todo: update data after updating from_content
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

//Convert list into Map when scaning
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

func init() {
	fieldtype.RegisterFieldType(
		fieldtype.FieldtypeDef{Type: "relationlist", Value: "contenttype.RelationList", IsRelation: true},
		func() fieldtype.FieldTyper { return &RelationList{} })
}
