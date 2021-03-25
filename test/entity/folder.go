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



type Folder struct{
     contenttype.ContentCommon `boil:",bind"`

     
    
         
         
         
            FolderType  fieldtype.Text `boil:"folder_type" json:"folder_type" toml:"folder_type" yaml:"folder_type"`
         
        
    
         
         
         
            Summary  fieldtype.RichText `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
         
        
    
         
         
         
            Title  fieldtype.Text `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func (c *Folder ) TableName() string{
	 return "dm_folder"
}

func (c *Folder ) ContentType() string{
	 return "folder"
}

func (c *Folder ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *Folder) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}

func (c *Folder) ToMap() map[string]interface{}{
    result := map[string]interface{}{}
    for _, identifier := range c.IdentifierList(){
      result[identifier] = c.Value(identifier)
    }
    return result
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *Folder) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    

    
        
        
            result["folder_type"]=c.FolderType
        
        
    
        
        
            result["summary"]=c.Summary
        
        
    
        
        
            result["title"]=c.Title
        
        
    
	for key, value := range c.ContentCommon.ToDBValues() {
		result[key] = value
	}
	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *Folder) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "folder_type","summary","title",}...)
}

func (c *Folder) Definition(language ...string) definition.ContentType {
	def, _ := definition.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *Folder) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    
    
    case "folder_type":
        
            result = &(c.FolderType)
        
    
    
    
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
func (c *Folder) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
        
            
            
            
            case "folder_type":
            c.FolderType = value.(fieldtype.Text)
            
            
        
            
            
            
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
func (c *Folder) Store(ctx context.Context, transaction ...*sql.Tx) error {
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

func (c *Folder)StoreWithLocation(){

}

//Delete content only
func (c *Folder) Delete(ctx context.Context, transaction ...*sql.Tx) error {
	contentError := db.Delete(ctx, c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Folder{}
	}

	newList := func() interface{} {
		return &[]Folder{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]Folder)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("folder",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
