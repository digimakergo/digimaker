//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
	. "dm/query"
)



type User struct{
     ContentCommon `boil:",bind"`
    
     
     Firstname fieldtype.TextField `boil:"firstname" json:"firstname" toml:"firstname" yaml:"firstname"`
    
     
     Lastname fieldtype.TextField `boil:"lastname" json:"lastname" toml:"lastname" yaml:"lastname"`
    
     
     Login fieldtype.TextField `boil:"login" json:"login" toml:"login" yaml:"login"`
    
     
     Password fieldtype.TextField `boil:"password" json:"password" toml:"password" yaml:"password"`
    
     Location `boil:"location,bind"`
}

func ( User ) TableName() string{
	 return "dm_user"
}

func (c User) contentValues() map[string]interface{} {
	result := make(map[string]interface{})
    
        result["firstname"]=c.Firstname
    
        result["lastname"]=c.Lastname
    
        result["login"]=c.Login
    
        result["password"]=c.Password
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c User) Values() map[string]interface{} {
    result := c.contentValues()

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

//Store content.
//Note: it will set id to CID after success
func (c *User) Store() error {
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
		return &User{}
	}

	newList := func() interface{} {
		return &[]User{}
	}

	convert := func(obj interface{}) []contenttype.ContentTyper {
		list := obj.(*[]User)
		var result []contenttype.ContentTyper
		for _, item := range *list {
			result = append(result, item)
		}
		return result
	}

	Register("user",
		ContentTypeRegister{
			New:            new,
			NewList:        newList,
			ListToContentTyper: convert})
}
