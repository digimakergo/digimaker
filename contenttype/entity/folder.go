//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
    "dm/util"
	. "dm/query"
)



type Folder struct{
     ContentCommon `boil:",bind"`
    
     
     
        Summary fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
     
    
     
     
        Title fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
     
    
     Location `boil:"location,bind"`
}

func ( *Folder ) TableName() string{
	 return "dm_folder"
}

func ( *Folder ) ContentType() string{
	 return "folder"
}


//todo: cache this? (then you need a reload?)
func (c *Folder) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        result["summary"]=c.Summary
        
    
        
        result["title"]=c.Title
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Folder) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "summary","title",}...)
}

func (c *Folder) Value(identifier string) interface{} {
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


func (c *Folder) SetValue(identifier string, value interface{}) error {
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
func (c *Folder) Store() error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.ToMap())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.ToMap(), Cond("id", c.CID))
		return err
	}
	return nil
}


func init() {
	new := func() contenttype.ContentTyper {
		return &Folder{}
	}

	newList := func() interface{} {
		return &[]Folder{}
	}

	Register("folder",
		ContentTypeRegister{
			New:            new,
			NewList:        newList})
}
