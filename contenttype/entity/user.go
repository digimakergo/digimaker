//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
    "dm/util"
	. "dm/query"
)



type User struct{
     ContentCommon `boil:",bind"`
    
     
     
        Firstname fieldtype.TextField `boil:"firstname" json:"firstname" toml:"firstname" yaml:"firstname"`
     
    
     
     
        Lastname fieldtype.TextField `boil:"lastname" json:"lastname" toml:"lastname" yaml:"lastname"`
     
    
     
     
        Login fieldtype.TextField `boil:"login" json:"login" toml:"login" yaml:"login"`
     
    
     
     
        Password fieldtype.TextField `boil:"password" json:"password" toml:"password" yaml:"password"`
     
    
     contenttype.Location `boil:"location,bind"`
}

func ( *User ) TableName() string{
	 return "dm_user"
}

func ( *User ) ContentType() string{
	 return "user"
}

func (c *User) GetLocation() *contenttype.Location{
    return &c.Location
}


//todo: cache this? (then you need a reload?)
func (c *User) ToMap() map[string]interface{} {
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

func (c *User) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "firstname","lastname","login","password",}...)
}

func (c *User) DisplayIdentifierList() []string {
	return []string{ "firstname","lastname",}
}

func (c *User) Value(identifier string) interface{} {
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    var result interface{}
	switch identifier {
    
    case "firstname":
        
            result = c.Firstname
        
    
    case "lastname":
        
            result = c.Lastname
        
    
    case "login":
        
            result = c.Login
        
    
    case "password":
        
            result = c.Password
        
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *User) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            case "firstname":
            c.Firstname = value.(fieldtype.TextField)
            
        
            
            
            case "lastname":
            c.Lastname = value.(fieldtype.TextField)
            
        
            
            
            case "login":
            c.Login = value.(fieldtype.TextField)
            
        
            
            
            case "password":
            c.Password = value.(fieldtype.TextField)
            
        
	default:
		err := c.ContentCommon.SetValue(identifier, value)
        if err != nil{
            return err
        }
	}
	//todo: check if identifier exist
	return nil
}

//Store content.
//Note: it will set id to CID after success
func (c *User) Store(transaction ...*sql.Tx) error {
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
func (c *User) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &User{}
	}

	newList := func() interface{} {
		return &[]User{}
	}

	Register("user",
		ContentTypeRegister{
			New:            new,
			NewList:        newList})
}
