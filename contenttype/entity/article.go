//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
	"dm/contenttype"
	"dm/db"
	"dm/fieldtype"
	. "dm/query"
)

type Article struct {
	ContentCommon `boil:",bind"`

	Body fieldtype.RichTextField `boil:"body" json:"body" toml:"body" yaml:"body"`

	Summary fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`

	Title fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`

	Location `boil:"location,bind"`
}

func (*Article) TableName() string {
	return "dm_article"
}

func (*Article) ContentType() string {
	return "article"
}

//todo: cache this? (then you need a reload?)
func (c *Article) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	result["body"] = c.Body

	result["summary"] = c.Summary

	result["title"] = c.Title

	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Article) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(), []string{"body", "coverimage", "related_articles", "summary", "title"}...)
}

func (c *Article) Value(identifier string) interface{} {
	var result interface{}
	switch identifier {

	case "body":

		result = c.Body

	case "coverimage":

		result = c.Relations.Value["coverimage"]

	case "related_articles":

		result = c.Relations.Value["related_articles"]

	case "summary":

		result = c.Summary

	case "title":

		result = c.Title

	case "cid":
		result = c.ContentCommon.CID
	default:
		result = c.ContentCommon.Value(identifier)
	}
	return result
}

func (c *Article) SetValue(identifier string, value interface{}) error {
	switch identifier {

	case "body":
		c.Body = value.(fieldtype.RichTextField)

	case "summary":
		c.Summary = value.(fieldtype.RichTextField)

	case "title":
		c.Title = value.(fieldtype.TextField)

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
func (c *Article) Store() error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.ToMap())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.ToMap(), Cond("id", c.CID))
		return err
	}
	return nil
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Article{}
	}

	newList := func() interface{} {
		return &[]Article{}
	}

	Register("article",
		ContentTypeRegister{
			New:     new,
			NewList: newList})
}
