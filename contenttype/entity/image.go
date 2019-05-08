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



type Image struct{
     ContentCommon `boil:",bind"`
    
     
     
        Imagetype fieldtype.TextField `boil:"imagetype" json:"imagetype" toml:"imagetype" yaml:"imagetype"`
     
    
     
     
        Path fieldtype.TextField `boil:"path" json:"path" toml:"path" yaml:"path"`
     
    
     
     
        Title fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
     
    
     contenttype.Location `boil:"location,bind"`
}

func ( *Image ) TableName() string{
	 return "dm_image"
}

func ( *Image ) ContentType() string{
	 return "image"
}

func (c *Image) GetLocation() *contenttype.Location{
    return &c.Location
}


//todo: cache this? (then you need a reload?)
func (c *Image) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        result["imagetype"]=c.Imagetype
        
    
        
        result["path"]=c.Path
        
    
        
        result["title"]=c.Title
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Image) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "imagetype","path","title",}...)
}

func (c *Image) DisplayIdentifierList() []string {
	return []string{ "imagetype","title","path",}
}

func (c *Image) Value(identifier string) interface{} {
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    var result interface{}
	switch identifier {
    
    case "imagetype":
        
            result = c.Imagetype
        
    
    case "path":
        
            result = c.Path
        
    
    case "title":
        
            result = c.Title
        
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *Image) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            case "imagetype":
            c.Imagetype = value.(fieldtype.TextField)
            
        
            
            
            case "path":
            c.Path = value.(fieldtype.TextField)
            
        
            
            
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
func (c *Image) Store(transaction ...*sql.Tx) error {
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
func (c *Image) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Image{}
	}

	newList := func() interface{} {
		return &[]Image{}
	}

	Register("image",
		ContentTypeRegister{
			New:            new,
			NewList:        newList})
}
