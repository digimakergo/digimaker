//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
    
	. "dm/query"
)



type File struct{
     ContentCommon `boil:",bind"`
    
     
     
        Filetype  fieldtype.TextField `boil:"filetype" json:"filetype" toml:"filetype" yaml:"filetype"`
     
    
     
     
        Path  fieldtype.TextField `boil:"path" json:"path" toml:"path" yaml:"path"`
     
    
     
     
        Title  fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
     
    
    
}

func ( *File ) TableName() string{
	 return "dm_file"
}

func ( *File ) ContentType() string{
	 return "file"
}

func (c *File ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *File) GetLocation() *contenttype.Location{
    
    return nil
    
}


//todo: cache this? (then you need a reload?)
func (c *File) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        result["filetype"]=c.Filetype
        
    
        
        result["path"]=c.Path
        
    
        
        result["title"]=c.Title
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *File) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "filetype","path","title",}...)
}

func (c *File) Definition() contenttype.ContentTypeSetting {
	return contenttype.GetContentDefinition( c.ContentType() )
}

func (c *File) Value(identifier string) interface{} {
    
    var result interface{}
	switch identifier {
    
    case "filetype":
        
            result = c.Filetype
        
    
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


func (c *File) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            case "filetype":
            c.Filetype = value.(fieldtype.TextField)
            
        
            
            
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
func (c *File) Store(transaction ...*sql.Tx) error {
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

func (c *File)StoreWithLocation(){

}

//Delete content only
func (c *File) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &File{}
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

	Register("file",
		ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
