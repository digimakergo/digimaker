//Author xc, Created on 2019-04-24 20:00
//{COPYRIGHTS}
package entity

import (
	"dm/contenttype"
	"dm/fieldtype"
	"encoding/json"

	"github.com/pkg/errors"
)

type Relation struct {
	RID             int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	FromContentID   int    `boil:"from_content_id" json:"from_content_id" toml:"from_content_id" yaml:"from_content_id"`
	FromContentType string `boil:"from_type" json:"from_type" toml:"from_type" yaml:"from_type"`
	RelationType    string `boil:"relation_type" json:"relation_type" toml:"relation_type" yaml:"relation_type"`
	Priority        int    `boil:"priority" json:"priority" toml:"priority" yaml:"priority"`
	Identifier      string `boil:"identifier" json:"identifier" toml:"identifier" yaml:"identifier"`
	Description     string `boil:"description" json:"description" toml:"description" yaml:"description"`
	Data            string `boil:"data" json:"data" toml:"data" yaml:"data"`
	RemoteID        string `boil:"remote_id" json:"remote_id" toml:"remote_id" yaml:"remote_id"`
}

//Generate data based on data_pattern in the relation type.
//eg. typical relation list is name+url
//   cover image can be name+image id+image alias(eg. medium/coverimage)
// slideshows can be same as cover image, but will be more images, which is supported automatically
func (r *Relation) GenerateData(toContent contenttype.ContentTyper, identifier string, fromContent contenttype.ContentTyper) {

}

type ContentRelations struct {
	Value map[string][]fieldtype.RelationField
}

func (relations *ContentRelations) Scan(src interface{}) error {
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

/*
//todo: add relation, remove relation, save relation
func (r ContentRelations) Add(values map[string]interface{}) error {
	db := db.DBHanlder()
	_, err := db.Insert("dm_relation", values)
	return err
}

func (r ContentRelations) Update() {

}

func (r ContentRelations) Delete() {

}
*/
