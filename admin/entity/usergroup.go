//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "dm/core/db"
    "dm/core/contenttype"
	"dm/core/fieldtype"
    
    "dm/core/util"
    
	. "dm/core/db"
)



type Usergroup struct{
     contenttype.ContentCommon `boil:",bind"`
    
         
         
         
            Summary  fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
         
        
    
         
         
         
            Title  fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func ( *Usergroup ) TableName() string{
	 return "dm_usergroup"
}

func ( *Usergroup ) ContentType() string{
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


//todo: cache this? (then you need a reload?)
func (c *Usergroup) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        
            result["summary"]=c.Summary
        
        
    
        
        
            result["title"]=c.Title
        
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Usergroup) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "summary","title",}...)
}

func (c *Usergroup) Definition() contenttype.ContentType {
	def, _ := contenttype.GetDefinition( c.ContentType() )
    return def
}

func (c *Usergroup) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    
    case "summary":
        
            result = c.Summary
        
    
    
    
    case "title":
        
            result = c.Title
        
    
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *Usergroup) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            
            case "summary":
            c.Summary = value.(fieldtype.RichTextField)
            
            
        
            
            
            
            case "title":
            c.Title = value.(fieldtype.TextField)
            
            
        
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
