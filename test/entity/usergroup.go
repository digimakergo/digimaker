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



type Usergroup struct{
     contenttype.ContentCommon `boil:",bind"`

     
    
         
         
         
            Summary  fieldtype.RichText `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
         
        
    
         
         
         
            Title  fieldtype.Text `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func (c *Usergroup ) TableName() string{
	 return "dm_usergroup"
}

func (c *Usergroup ) ContentType() string{
	 return "usergroup"
}

func (c *Usergroup ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *Usergroup) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}

func (c *Usergroup) ToMap() map[string]interface{}{
    result := map[string]interface{}{}
    for _, identifier := range c.IdentifierList(){
      result[identifier] = c.Value(identifier)
    }
    return result
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *Usergroup) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    

    
        
        
            result["summary"]=c.Summary
        
        
    
        
        
            result["title"]=c.Title
        
        
    
	for key, value := range c.ContentCommon.ToDBValues() {
		result[key] = value
	}
	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *Usergroup) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "summary","title",}...)
}

func (c *Usergroup) Definition(language ...string) contenttype.ContentType {
	def, _ := contenttype.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *Usergroup) Value(identifier string) interface{} {
    
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
func (c *Usergroup) SetValue(identifier string, value interface{}) error {
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
func (c *Usergroup) Store(transaction ...*sql.Tx) error {
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

func (c *Usergroup)StoreWithLocation(){

}

//Delete content only
func (c *Usergroup) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Usergroup{}
	}

	newList := func() interface{} {
		return &[]Usergroup{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]Usergroup)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("usergroup",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
