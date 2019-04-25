//Author xc, Created on 2019-04-24 20:00
//{COPYRIGHTS}
package entity

import "dm/db"

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

type Relations struct {
	CID          int
	RelationList *[]Relation
}

//todo: add relation, remove relation, save relation
func (r Relations) Add(values map[string]interface{}) error {
	db := db.DBHanlder()
	_, err := db.Insert("dm_relation", values)
	return err
}

func (r Relations) Update() {

}

func (r Relations) Delete() {

}
