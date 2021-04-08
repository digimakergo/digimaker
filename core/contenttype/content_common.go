package contenttype

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/friendsofgo/errors"
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
	Relations  ContentRelationList `boil:"relations" json:"-" toml:"relations" yaml:"relations"`
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

func (c *ContentCommon) GetRelations() ContentRelationList {
	return c.Relations
}

func (c *ContentCommon) SetRelation(identifier string, value interface{}) error {
	if len(c.Relations) == 0 {
		c.Relations = ContentRelationList{}
	}
	if list, ok := value.(fieldtype.RelationList); ok {
		rTemp := c.Relations
		rTemp[identifier] = list
		c.Relations = rTemp
	} else {
		return fmt.Errorf("Only RelationList is supported on field %v", identifier)
	}
	return nil
}

func (c *ContentCommon) StoreRelations(ctx context.Context, thisContenttype string, transaction ...*sql.Tx) error {
	existingList := []fieldtype.Relation{}
	//get existing
	if c.CID > 0 {
		cond := db.Cond("to_type", thisContenttype).Cond("to_content_id", c.CID)
		_, err := db.BindEntity(ctx, &existingList, "dm_relation", cond)
		if err != nil {
			return err
		}
	}

	//get to be deleted
	deleteCond := db.EmptyCond()
	for _, existing := range existingList {
		willDelete := true
		for _, list := range c.Relations {
			for _, relation := range list {
				if existing.FromContentID == relation.FromContentID && existing.FromType == relation.FromType {
					willDelete = false
				}
			}
		}
		if willDelete {
			deleteCond = deleteCond.Or(db.Cond("from_content_id", existing.FromContentID).Cond("from_type", existing.FromType))
		}
	}

	//get to be added
	toBeAdded := []fieldtype.Relation{}
	for _, list := range c.Relations {
		for _, relation := range list {
			exists := false
			for _, existing := range existingList {
				if existing.FromContentID == relation.FromContentID && existing.FromType == relation.FromType {
					exists = true
				}
			}
			if !exists {
				toBeAdded = append(toBeAdded, relation)
			}
		}
	}

	//execute delete
	if c.CID > 0 && !db.IsEmptyCond(deleteCond) {
		err := db.Delete(ctx, "dm_relation", db.Cond("to_content_id", c.CID).Cond("to_type", thisContenttype).And(deleteCond), transaction...)
		if err != nil {
			return err
		}
	}

	//execute insert
	for _, relation := range toBeAdded {
		dataMap := map[string]interface{}{}
		dataMap["identifier"] = relation.Identifier
		dataMap["to_content_id"] = c.CID
		dataMap["to_type"] = thisContenttype
		dataMap["from_content_id"] = relation.FromContentID
		dataMap["from_type"] = relation.FromType
		dataMap["data"] = relation.Data
		dataMap["priority"] = relation.Priority
		_, err := db.Insert(ctx, "dm_relation", dataMap, transaction...)
		if err != nil {
			return err
		}
	}

	return nil
}

//FinishBind sets related data after data binding. It will be better if SQLBoiler support interface for customized  binding for struct.
func FinishBind(content ContentTyper) error {
	contentType := content.ContentType()
	def, _ := definition.GetDefinition(contentType)
	if def.HasRelationlist() {
		relationMap := content.GetRelations()
		for identifier, fieldDef := range def.FieldMap {
			if fieldDef.FieldType == "relationlist" {
				if value, ok := relationMap[identifier]; ok {
					err := content.SetValue(identifier, value)
					if err != nil {
						return errors.Wrap(err, "Error when binding relationlist "+identifier)
					}
				}
			}
		}
	}
	return nil
}
