//Author xc, Created on 2019-04-23 16:50
//{COPYRIGHTS}

package entity

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type RelationSetting map[string]string

//A relationtype is to cache relation group information in an article,
//in db it can be a serilized string or json
//The real reation is in dm_relation
//So this is a low level field type(datatype), which is simiar to string
type ContentRelationList struct {
	data string
	//Value is a group of relations.
	//eg. { "relation_articles": [ {"name":"Break news",
	//                               "uid":"999y32gghh",
	//                                "id": 1123 } ],
	//       "related_links": ... }
	Value map[string][]RelationSetting
}

func (relations *ContentRelationList) GetRelationList(identifier string) []RelationSetting {
	return relations.Value[identifier]
}

func (relations *ContentRelationList) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		errors.New("Unknow scan value.")
	}

	relations.data = source

	//var objmap map[string]*json.RawMessage
	err := json.Unmarshal(src.([]byte), &relations.Value)
	if err != nil {
		return errors.Wrap(err, "Can not convert to ContentRelationList. Relation data is not correct: "+source)
	}

	return nil
}
