package contenttype

import (
	"database/sql"
	"encoding/json"

	"github.com/digimakergo/digimaker/core/db"
)

type ContentCommon struct {
	CID        int                 `boil:"cid" json:"cid" toml:"cid" yaml:"cid"`
	Version    int                 `boil:"version" json:"version" toml:"version" yaml:"version"`
	Published  int                 `boil:"published" json:"published" toml:"published" yaml:"published"`
	Modified   int                 `boil:"modified" json:"modified" toml:"modified" yaml:"modified"`
	CUID       string              `boil:"cuid" json:"cuid" toml:"cuid" yaml:"cuid"`
	Status     int                 `boil:"status" json:"status" toml:"status" yaml:"status"`
	Author     int                 `boil:"author" json:"author" toml:"author" yaml:"author"`
	AuthorName string              `boil:"author_name" json:"author_name" toml:"author_name" yaml:"author_name"`
	Relations  ContentRelationList `boil:"relations" json:"relations" toml:"relations" yaml:"relations"`
}

//IdentifierList return list of all field names
func (c ContentCommon) IdentifierList() []string {
	return []string{"cid", "version", "published", "modified", "author", "author_name", "cuid", "status"}
}

//Values return values for insert/update DB. todo: rename to ToDBValues()
func (c ContentCommon) ToDBValues() map[string]interface{} {
	result := map[string]interface{}{
		"id":        c.CID,
		"version":   c.Version,
		"published": c.Published,
		"modified":  c.Modified,
		"status":    c.Status,
		"author":    c.Author,
		"cuid":      c.CUID,
	}
	return result
}

func (c *ContentCommon) Value(identifier string) interface{} {
	var result interface{}
	switch identifier {
	case "cid":
		result = c.CID
	case "version":
		result = c.Version
	case "modified":
		result = c.Modified
	case "published":
		result = c.Published
	case "author":
		result = c.Author
	case "author_name":
		result = c.AuthorName
	case "status":
		result = c.Status
	case "cuid":
		result = c.CUID
	case "relations":
		result = c.Relations
	}
	return result
}

func (c *ContentCommon) SetValue(identifier string, value interface{}) error {
	switch identifier {
	case "cid":
		c.CID = value.(int)
	case "version":
		c.Version = value.(int)
	case "published":
		c.Published = value.(int)
	case "modified":
		c.Modified = value.(int)
	case "status":
		c.Status = value.(int)
	case "author":
		c.Author = value.(int)
	case "author_name":
		c.AuthorName = value.(string)
	case "cuid":
		c.CUID = value.(string)
	}
	return nil
}

func (c *ContentCommon) GetCID() int {
	return c.CID
}

func (c *ContentCommon) GetRelations() *ContentRelationList {
	return &c.Relations
}

func (c *ContentCommon) StoreRelations(thisContenttype string, transaction ...*sql.Tx) error {
	dbHandler := db.DBHanlder()

	//delete
	err := dbHandler.Delete("dm_relation", db.Cond("to_content_id", c.CID).Cond("to_type", thisContenttype), transaction...)
	if err != nil {
		return nil
	}

	//insert
	for identifier, list := range c.Relations.Map {
		for _, relation := range *list {
			data, _ := json.Marshal(relation)
			dataMap := map[string]interface{}{}
			json.Unmarshal(data, &dataMap)
			dataMap["identifier"] = identifier
			dataMap["to_content_id"] = c.CID
			dataMap["to_type"] = thisContenttype
			dataMap["priority"] = relation.Priority
			_, err := dbHandler.Insert("dm_relation", dataMap, transaction...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//TODO: add more common methods related to content here.
