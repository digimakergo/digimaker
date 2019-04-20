//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
	"dm/fieldtype"
	. "dm/query"
)



type Article struct{
     Location `boil:"dm_location,bind"`
     ContentCommon `boil:",bind"`
    
     
     Body fieldtype.RichTextField `boil:"body" json:"body" toml:"body" yaml:"body"`
    
     
     Summary fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
    
     
     Title fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
    
}


func ( Article ) TableName() string{
	 return "dm_article"
}


func (c Article) Values() map[string]interface{} {
	result := make(map[string]interface{})

    
        result["body"]=c.Body
    
        result["summary"]=c.Summary
    
        result["title"]=c.Title
    

    for key, value := range c.ContentCommon.Values() {
        result[key] = value
    }

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

func (c Article) Store() error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.Values())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.Values(), Cond("id", c.CID))
		return err
	}
	return nil
}
