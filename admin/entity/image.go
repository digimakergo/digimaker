//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "dm/core/db"
    "dm/core/contenttype"
	"dm/core/fieldtype"
    
	. "dm/core/db"
)



type Image struct{
     contenttype.ContentCommon `boil:",bind"`
    
         
         
         
            Image  fieldtype.TextField `boil:"image" json:"image" toml:"image" yaml:"image"`
         
        
    
         
         
         
            Imagetype  string `boil:"imagetype" json:"imagetype" toml:"imagetype" yaml:"imagetype"`
         
        
    
         
         
         
            ParentId  int `boil:"parent_id" json:"parent_id" toml:"parent_id" yaml:"parent_id"`
         
        
    
         
         
         
            Title  fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
    
}

func ( *Image ) TableName() string{
	 return "dm_image"
}

func ( *Image ) ContentType() string{
	 return "image"
}

func (c *Image ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *Image) GetLocation() *contenttype.Location{
    
    return nil
    
}


//todo: cache this? (then you need a reload?)
func (c *Image) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        
            result["image"]=c.Image
        
        
    
        
        
            result["imagetype"]=c.Imagetype
        
        
    
        
        
            result["parent_id"]=c.ParentId
        
        
    
        
        
            result["title"]=c.Title
        
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Image) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "image","imagetype","parent_id","title",}...)
}

func (c *Image) Definition() contenttype.ContentTypeSetting {
	return contenttype.GetContentDefinition( c.ContentType() )
}

func (c *Image) Value(identifier string) interface{} {
    
    var result interface{}
	switch identifier {
    
    
    case "image":
        
            result = c.Image
        
    
    
    
    case "imagetype":
        
            result = c.Imagetype
        
    
    
    
    case "parent_id":
        
            result = c.ParentId
        
    
    
    
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
        
            
            
            
            case "image":
            c.Image = value.(fieldtype.TextField)
            
            
        
            
            
            
            case "imagetype":
            c.Imagetype = value.(string)
            
            
        
            
            
            
            case "parent_id":
            c.ParentId = value.(int)
            
            
        
            
            
            
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

func (c *Image)StoreWithLocation(){

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

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]Image)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("image",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
