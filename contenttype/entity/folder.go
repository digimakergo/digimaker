//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
	. "dm/query"
)



type Folder struct{
     ContentCommon `boil:",bind"`
    
     
     Summary fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
    
     
     Title fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
    
     Location `boil:"location,bind"`
}

func ( Folder ) TableName() string{
	 return "dm_folder"
}

func (c Folder) contentValues() map[string]interface{} {
	result := make(map[string]interface{})
    
        result["summary"]=c.Summary
    
        result["title"]=c.Title
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c Folder) Values() map[string]interface{} {
    result := c.contentValues()

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

//Store content.
//Note: it will set id to CID after success
func (c *Folder) Store() error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.contentValues())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.contentValues(), Cond("id", c.CID))
		return err
	}
	return nil
}


func init() {
	new := func() contenttype.ContentTyper {
		return &Folder{}
	}

	newList := func() interface{} {
		return &[]Folder{}
	}

	convert := func(obj interface{}) []contenttype.ContentTyper {
		list := obj.(*[]Folder)
		var result []contenttype.ContentTyper
		for _, item := range *list {
			result = append(result, item)
		}
		return result
	}

	Register("folder",
		ContentTypeRegister{
			New:            new,
			NewList:        newList,
			ListToContentTyper: convert})
}
