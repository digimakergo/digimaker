//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "dm/dm/db"
    "dm/dm/contenttype"
	"dm/dm/fieldtype"
    
    "dm/dm/util"
    
	. "dm/dm/query"
)



type Folder struct{
     ContentCommon `boil:",bind"`
    
     
     
        FolderType  fieldtype.TextField `boil:"folder_type" json:"folder_type" toml:"folder_type" yaml:"folder_type"`
     
    
     
     
        Summary  fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
     
    
     
     
        Title  fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
     
    
    
     contenttype.Location `boil:"location,bind"  json:"location"`
    
}

func ( *Folder ) TableName() string{
	 return "dm_folder"
}

func ( *Folder ) ContentType() string{
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


//todo: cache this? (then you need a reload?)
func (c *Folder) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        result["folder_type"]=c.FolderType
        
    
        
        result["summary"]=c.Summary
        
    
        
        result["title"]=c.Title
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Folder) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "folder_type","summary","title",}...)
}

func (c *Folder) Definition() contenttype.ContentTypeSetting {
	return contenttype.GetContentDefinition( c.ContentType() )
}

func (c *Folder) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    case "folder_type":
        
            result = c.FolderType
        
    
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
        
            
            
            case "folder_type":
            c.FolderType = value.(fieldtype.TextField)
            
        
            
            
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
func (c *Folder) Store(transaction ...*sql.Tx) error {
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

func (c *Folder)StoreWithLocation(){

}

//Delete content only
func (c *Folder) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
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
