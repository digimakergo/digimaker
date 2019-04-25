//Author xc, Created on 2019-04-24 20:00
//{COPYRIGHTS}
package entity

import (
	"dm/contenttype"
)

type ContentRelation struct {
	RID             int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	FromContentID   int    `boil:"from_content_id" json:"from_content_id" toml:"from_content_id" yaml:"from_content_id"`
	FromContentType string `boil:"from_type" json:"from_type" toml:"from_type" yaml:"from_type"`
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
func (r *ContentRelation) GenerateData(toContent contenttype.ContentTyper, identifier string, fromContent contenttype.ContentTyper) {

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
