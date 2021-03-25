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



type File struct{
     contenttype.ContentEntity `boil:",bind"`

     
    
         
         
         
            Filetype  fieldtype.Text `boil:"filetype" json:"filetype" toml:"filetype" yaml:"filetype"`
         
        
    
         
         
         
            Path  fieldtype.Text `boil:"path" json:"path" toml:"path" yaml:"path"`
         
        
    
         
         
         
            Title  fieldtype.Text `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
}

func (c *File ) TableName() string{
	 return "dm_file"
}

func (c *File ) ContentType() string{
	 return "file"
}

func (c *File ) GetName() string{
	 return ""
}

func (c *File) GetLocation() *contenttype.Location{
    return nil
}

func (c *File) ToMap() map[string]interface{}{
    result := map[string]interface{}{}
    for _, identifier := range c.IdentifierList(){
      result[identifier] = c.Value(identifier)
    }
    return result
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *File) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    

    
        
        
            result["filetype"]=c.Filetype
        
        
    
        
        
            result["path"]=c.Path
        
        
    
        
        
            result["title"]=c.Title
        
        
    
	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *File) IdentifierList() []string {
	return []string{ "filetype","path","title",}
}

func (c *File) Definition(language ...string) definition.ContentType {
	def, _ := definition.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *File) Value(identifier string) interface{} {
    
    var result interface{}
	switch identifier {
    
    
    
    case "filetype":
        
            result = &(c.Filetype)
        
    
    
    
    case "path":
        
            result = &(c.Path)
        
    
    
    
    case "title":
        
            result = &(c.Title)
        
    
    

    default:
    }
	return result
}

//Set value to a field
func (c *File) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
        
            
            
            
            case "filetype":
            c.Filetype = value.(fieldtype.Text)
            
            
        
            
            
            
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
func (c *File) Store(ctx context.Context, transaction ...*sql.Tx) error {
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


func (c *File)StoreWithLocation(){

}

//Delete content only
func (c *File) Delete(ctx context.Context, transaction ...*sql.Tx) error {
	contentError := db.Delete(ctx, c.TableName(), Cond("id", c.ID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
    entity := &File{}
    entity.ContentEntity.ContentType = "File"
    return entity
	}

	newList := func() interface{} {
		return &[]File{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]File)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("file",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
