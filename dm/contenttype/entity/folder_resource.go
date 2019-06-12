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



type FolderResource struct{
     ContentCommon `boil:",bind"`
    
     
     
        ResourceType  fieldtype.TextField `boil:"resource_type" json:"resource_type" toml:"resource_type" yaml:"resource_type"`
     
    
     
     
        Summary  fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
     
    
     
     
        Title  fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
     
    
    
     contenttype.Location `boil:"location,bind"  json:"location"`
    
}

func ( *FolderResource ) TableName() string{
	 return "dm_folder_resource"
}

func ( *FolderResource ) ContentType() string{
	 return "folder_resource"
}

func (c *FolderResource ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *FolderResource) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}


//todo: cache this? (then you need a reload?)
func (c *FolderResource) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        result["resource_type"]=c.ResourceType
        
    
        
        result["summary"]=c.Summary
        
    
        
        result["title"]=c.Title
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *FolderResource) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "resource_type","summary","title",}...)
}

func (c *FolderResource) Definition() contenttype.ContentTypeSetting {
	return contenttype.GetContentDefinition( c.ContentType() )
}

func (c *FolderResource) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    case "resource_type":
        
            result = c.ResourceType
        
    
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


func (c *FolderResource) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            case "resource_type":
            c.ResourceType = value.(fieldtype.TextField)
            
        
            
            
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
func (c *FolderResource) Store(transaction ...*sql.Tx) error {
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

func (c *FolderResource)StoreWithLocation(){

}

//Delete content only
func (c *FolderResource) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &FolderResource{}
	}

	newList := func() interface{} {
		return &[]FolderResource{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]FolderResource)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("folder_resource",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
