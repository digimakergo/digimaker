//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
	"dm/fieldtype"
	. "dm/query"
)



type Folder struct{
     Location `boil:"dm_location,bind"`
     ContentCommon `boil:",bind"`
    
     
     Summary fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
    
     
     Title fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
    
}


func ( Folder ) TableName() string{
	 return "dm_folder"
}


func (c Folder) Values() map[string]interface{} {
	result := make(map[string]interface{})

    
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

func (c Folder) Store() error {
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
