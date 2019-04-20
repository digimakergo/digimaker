//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
	"dm/fieldtype"
	. "dm/query"
)



type Userf struct{
     *Location
     *ContentCommon
    
     
     Firstname fieldtype.TextField `boil:"firstname" json:"firstname" toml:"firstname" yaml:"firstname"`
    
     
     Lastname fieldtype.TextField `boil:"lastname" json:"lastname" toml:"lastname" yaml:"lastname"`
    
     
     Login fieldtype.TextField `boil:"login" json:"login" toml:"login" yaml:"login"`
    
     
     Password fieldtype.TextField `boil:"password" json:"password" toml:"password" yaml:"password"`
    
}


func ( Userf ) TableName() string{
	 return "dm_user"
}


func (c Userf) Values() map[string]interface{} {
	result := make(map[string]interface{})

    for key, value := range c.ContentCommon.Values() {
        result[key] = value
    }

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

func (c Userf) Store() error {
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
