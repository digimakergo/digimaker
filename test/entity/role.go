//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "github.com/digimakergo/digimaker/core/db"
    "github.com/digimakergo/digimaker/core/contenttype"
	  "github.com/digimakergo/digimaker/core/fieldtype"
    
    "github.com/digimakergo/digimaker/core/util"
    
	. "github.com/digimakergo/digimaker/core/db"
    
)



type Role struct{
     contenttype.ContentCommon `boil:",bind"`

     
    
         
         
         
            Summary  fieldtype.RichText `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
         
        
    
         
         
         
            Title  fieldtype.Text `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func (c *Role ) TableName() string{
	 return "dm_role"
}

func (c *Role ) ContentType() string{
	 return "role"
}

func (c *Role ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *Role) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}

func (c *Role) ToMap() map[string]interface{}{
    result := map[string]interface{}{}
    for _, identifier := range c.IdentifierList(){
      result[identifier] = c.Value(identifier)
    }
    return result
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *Role) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    

    
        
        
            result["summary"]=c.Summary
        
        
    
        
        
            result["title"]=c.Title
        
        
    
	for key, value := range c.ContentCommon.ToDBValues() {
		result[key] = value
	}
	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *Role) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "summary","title",}...)
}

func (c *Role) Definition(language ...string) contenttype.ContentType {
	def, _ := contenttype.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *Role) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    
    
    case "summary":
        
            result = &(c.Summary)
        
    
    
    
    case "title":
        
            result = &(c.Title)
        
    
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}

//Set value to a field
func (c *Role) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
        
            
            
            
            case "summary":
            c.Summary = value.(fieldtype.RichText)
            
            
        
            
            
            
            case "title":
            c.Title = value.(fieldtype.Text)
            
            
        
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
func (c *Role) Store(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.ToDBValues(), transaction...)
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.ToDBValues(), Cond("id", c.CID), transaction...)
		return err
	}
	return nil
}

func (c *Role)StoreWithLocation(){

}

//Delete content only
func (c *Role) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Role{}
	}

	newList := func() interface{} {
		return &[]Role{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]Role)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("role",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
