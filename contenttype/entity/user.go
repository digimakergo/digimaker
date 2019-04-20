//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
	"dm/fieldtype"
	. "dm/query"
)



type User struct{
     Location `boil:"dm_location,bind"`
     ContentCommon `boil:",bind"`
    
     
     Firstname fieldtype.TextField `boil:"firstname" json:"firstname" toml:"firstname" yaml:"firstname"`
    
     
     Lastname fieldtype.TextField `boil:"lastname" json:"lastname" toml:"lastname" yaml:"lastname"`
    
     
     Login fieldtype.TextField `boil:"login" json:"login" toml:"login" yaml:"login"`
    
     
     Password fieldtype.TextField `boil:"password" json:"password" toml:"password" yaml:"password"`
    
}


func ( User ) TableName() string{
	 return "dm_user"
}


func (c User) Values() map[string]interface{} {
	result := make(map[string]interface{})

    
        result["firstname"]=c.Firstname
    
        result["lastname"]=c.Lastname
    
        result["login"]=c.Login
    
        result["password"]=c.Password
    

    for key, value := range c.ContentCommon.Values() {
        result[key] = value
    }

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

func (c User) Store() error {
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
