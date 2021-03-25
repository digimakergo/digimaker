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
    
	. "github.com/digimakergo/digimaker/core/db"
    
)



type Image struct{
     contenttype.ContentEntity `boil:",bind"`

     
          Imagetype  string `boil:"imagetype" json:"imagetype" toml:"imagetype" yaml:"imagetype"`
     
          ParentId  int `boil:"parent_id" json:"parent_id" toml:"parent_id" yaml:"parent_id"`
     
    
         
         
         
            Path  fieldtype.Text `boil:"path" json:"path" toml:"path" yaml:"path"`
         
        
    
         
         
         
            Title  fieldtype.Text `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
}

func (c *Image ) TableName() string{
	 return "dm_image"
}

func (c *Image ) ContentType() string{
	 return "image"
}

func (c *Image ) GetName() string{
	 return ""
}

func (c *Image) GetLocation() *contenttype.Location{
    return nil
}

func (c *Image) ToMap() map[string]interface{}{
    result := map[string]interface{}{}
    for _, identifier := range c.IdentifierList(){
      result[identifier] = c.Value(identifier)
    }
    return result
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *Image) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    
         result["imagetype"]=c.Imagetype
    
         result["parent_id"]=c.ParentId
    

    
        
        
            result["path"]=c.Path
        
        
    
        
        
            result["title"]=c.Title
        
        
    
	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *Image) IdentifierList() []string {
	return []string{ "path","title",}
}

func (c *Image) Definition(language ...string) definition.ContentType {
	def, _ := definition.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *Image) Value(identifier string) interface{} {
    
    var result interface{}
	switch identifier {
    
      case "imagetype":
         result = c.Imagetype
    
      case "parent_id":
         result = c.ParentId
    
    
    
    case "path":
        
            result = &(c.Path)
        
    
    
    
    case "title":
        
            result = &(c.Title)
        
    
    

    default:
    }
	return result
}

//Set value to a field
func (c *Image) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
          case "imagetype":
             c.Imagetype = value.(string)
        
          case "parent_id":
             c.ParentId = value.(int)
        
        
            
            
            
            case "path":
            c.Path = value.(fieldtype.Text)
            
            
        
            
            
            
            case "title":
            c.Title = value.(fieldtype.Text)
            
            
        
	default:

	}
	//todo: check if identifier exist
	return nil
}

//Store content.
//Note: it will set id to ID after success
func (c *Image) Store(ctx context.Context, transaction ...*sql.Tx) error {
	if c.ID == 0 {
		id, err := db.Insert(ctx, c.TableName(), c.ToDBValues(), transaction...)
		c.ID = id
		if err != nil {
			return err
		}
	} else {
		err := db.Update(ctx, c.TableName(), c.ToDBValues(), Cond("id", c.ID), transaction...)
		return err
	}
	return nil
}


func (c *Image)StoreWithLocation(){

}

//Delete content only
func (c *Image) Delete(ctx context.Context, transaction ...*sql.Tx) error {
	contentError := db.Delete(ctx, c.TableName(), Cond("id", c.ID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
    entity := &Image{}
    entity.ContentEntity.ContentType = "Image"
    return entity
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
