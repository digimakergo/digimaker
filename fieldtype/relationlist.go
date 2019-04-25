//Author xc, Created on 2019-04-23 22:20
//{COPYRIGHTS}

package fieldtype

import (
	"github.com/pkg/errors"
)

type RelationList struct {
	data string
	//Value is a group of relations.
	//eg. { "relation_articles": [ {"name":"Break news",
	//                               "uid":"999y32gghh",
	//                                "id": 1123 } ],
	//       "related_links": ... }
	Value map[string][]map[string]interface{}
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

	relations.data = "[" + source + "]"

	// //	var objmap map[string]*json.RawMessage
	// var objmap []map[string]interface{}
	//
	// err := json.Unmarshal([]byte(relations.data), &objmap)
	// if err != nil {
	// 	return errors.Wrap(err, "Can not convert to Relation. Relation data is not correct: "+source)
	// }
	//
	// if len(objmap) > 1 {
	// 	for _, item := range objmap {
	// 		identifier := item["identifier"].(string)
	// 		if _, ok := relations.Value[identifier]; ok {
	// 			relations.Value[identifier] = append(relations.Value[identifier], item)
	// 		} else {
	// 			relations.Value[identifier] = []map[string]interface{}{item}
	// 		}
	// 	}
	// }

	return nil
}
