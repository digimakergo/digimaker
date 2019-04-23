//Author xc, Created on 2019-04-23 22:20
//{COPYRIGHTS}

package fieldtype

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type RelationList struct {
	data string
	//Value is a group of relations.
	//eg. { "relation_articles": [ {"name":"Break news",
	//                               "uid":"999y32gghh",
	//                                "id": 1123 } ],
	//       "related_links": ... }
	Value []map[string]string
}

func (relations *RelationList) Scan(src interface{}) error {
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
