//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
	"database/sql"
	"dm/contenttype"
	"dm/db"
	"dm/fieldtype"
	. "dm/query"
	"dm/util"
)

type Image struct {
	ContentCommon `boil:",bind"`

	AttachedLocation int `boil:"attached_location" json:"attached_location" toml:"attached_location" yaml:"attached_location"`

	Path fieldtype.TextField `boil:"path" json:"path" toml:"path" yaml:"path"`

	Title fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`

	Type fieldtype.TextField `boil:"type" json:"type" toml:"type" yaml:"type"`

	contenttype.Location `boil:"location,bind"`
}

func (*Image) TableName() string {
	return "dm_image"
}

func (*Image) ContentType() string {
	return "image"
}

func (c *Image) GetLocation() *contenttype.Location {
	return &c.Location
}

//todo: cache this? (then you need a reload?)
func (c *Image) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	result["attached_location"] = c.AttachedLocation

	result["path"] = c.Path

	result["title"] = c.Title

	result["type"] = c.Type

	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Image) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(), []string{"attached_location", "path", "title", "type"}...)
}

func (c *Image) DisplayIdentifierList() []string {
	return []string{"firstname", "lastname", "login", "password"}
}

func (c *Image) Value(identifier string) interface{} {
	if util.Contains(c.Location.IdentifierList(), identifier) {
		return c.Location.Field(identifier)
	}
	var result interface{}
	switch identifier {

	case "attached_location":

		result = c.AttachedLocation

	case "path":

		result = c.Path

	case "title":

		result = c.Title

	case "type":

		result = c.Type

	case "cid":
		result = c.ContentCommon.CID
	default:
		result = c.ContentCommon.Value(identifier)
	}
	return result
}

func (c *Image) SetValue(identifier string, value interface{}) error {
	switch identifier {

	case "attached_location":
		c.AttachedLocation = value.(int)

	case "path":
		c.Path = value.(fieldtype.TextField)

	case "title":
		c.Title = value.(fieldtype.TextField)

	case "type":
		c.Type = value.(fieldtype.TextField)

	default:
		err := c.ContentCommon.SetValue(identifier, value)
		if err != nil {
			return err
		}
	}
	//todo: check if identifier exist
	return nil
}

//Store content.
//Note: it will set id to CID after success
func (c *Image) Store(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.ToMap(), transaction...)
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.ToMap(), Cond("id", c.CID), transaction...)
		return err
	}
	return nil
}

//Delete content only
func (c *Image) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Image{}
	}

	newList := func() interface{} {
		return &[]Image{}
	}

	Register("image",
		ContentTypeRegister{
			New:     new,
			NewList: newList})
}
