//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "context"
    "database/sql"
    "github.com/digimakergo/digimaker/core/db"
    "github.com/digimakergo/digimaker/core/definition"
    "github.com/digimakergo/digimaker/core/contenttype"
	  "github.com/digimakergo/digimaker/core/fieldtype"
    
    "github.com/digimakergo/digimaker/core/util"
    
	. "github.com/digimakergo/digimaker/core/db"
    
)



type User struct{
     contenttype.ContentCommon `boil:",bind"`

     
    
         
         
         
            Email  fieldtype.Text `boil:"email" json:"email" toml:"email" yaml:"email"`
         
        
    
         
         
         
            Firstname  fieldtype.Text `boil:"firstname" json:"firstname" toml:"firstname" yaml:"firstname"`
         
        
    
         
         
         
            Lastname  fieldtype.Text `boil:"lastname" json:"lastname" toml:"lastname" yaml:"lastname"`
         
        
    
         
         
         
            Login  fieldtype.Text `boil:"login" json:"login" toml:"login" yaml:"login"`
         
        
    
         
         
         
            Password  fieldtype.Password `boil:"password" json:"-" toml:"password" yaml:"password"`
         
        
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func (c *User ) TableName() string{
	 return "dm_user"
}

func (c *User ) ContentType() string{
	 return "user"
}

func (c *User ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *User) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}

func (c *User) ToMap() map[string]interface{}{
    result := map[string]interface{}{}
    for _, identifier := range c.IdentifierList(){
      result[identifier] = c.Value(identifier)
    }
    return result
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *User) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    

    
        
        
            result["email"]=c.Email
        
        
    
        
        
            result["firstname"]=c.Firstname
        
        
    
        
        
            result["lastname"]=c.Lastname
        
        
    
        
        
            result["login"]=c.Login
        
        
    
        
        
            result["password"]=c.Password
        
        
    
	for key, value := range c.ContentCommon.ToDBValues() {
		result[key] = value
	}
	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *User) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "email","firstname","lastname","login","password",}...)
}

func (c *User) Definition(language ...string) definition.ContentType {
	def, _ := definition.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *User) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    
    
    case "email":
        
            result = &(c.Email)
        
    
    
    
    case "firstname":
        
            result = &(c.Firstname)
        
    
    
    
    case "lastname":
        
            result = &(c.Lastname)
        
    
    
    
    case "login":
        
            result = &(c.Login)
        
    
    
    
    case "password":
        
            result = &(c.Password)
        
    
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}

//Set value to a field
func (c *User) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
        
            
            
            
            case "email":
            c.Email = value.(fieldtype.Text)
            
            
        
            
            
            
            case "firstname":
            c.Firstname = value.(fieldtype.Text)
            
            
        
            
            
            
            case "lastname":
            c.Lastname = value.(fieldtype.Text)
            
            
        
            
            
            
            case "login":
            c.Login = value.(fieldtype.Text)
            
            
        
            
            
            
            case "password":
            c.Password = value.(fieldtype.Password)
            
            
        
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
func (c *User) Store(ctx context.Context, transaction ...*sql.Tx) error {
	if c.CID == 0 {
		id, err := db.Insert(ctx, c.TableName(), c.ToDBValues(), transaction...)
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := db.Update(ctx, c.TableName(), c.ToDBValues(), Cond("id", c.CID), transaction...)
    if err != nil {
			return err
		}
	}

	err := c.StoreRelations(ctx, c.ContentType(), transaction...)
	if err != nil {
		return err
	}

	return nil
}

func (c *User)StoreWithLocation(){

}

//Delete content only
func (c *User) Delete(ctx context.Context, transaction ...*sql.Tx) error {
	contentError := db.Delete(ctx, c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &User{}
	}

	newList := func() interface{} {
		return &[]User{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]User)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("user",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
