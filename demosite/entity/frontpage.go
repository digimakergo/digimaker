//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "dm/dm/db"
    "dm/dm/contenttype"
	"dm/dm/fieldtype"
    
    "dm/dm/util"
    
	. "dm/dm/db"
)



type Frontpage struct{
     contenttype.ContentCommon `boil:",bind"`
    
     
     
        Mainarea  fieldtype.TextField `boil:"mainarea" json:"mainarea" toml:"mainarea" yaml:"mainarea"`
     
    
     
     
    
     
     
        Sidearea  fieldtype.TextField `boil:"sidearea" json:"sidearea" toml:"sidearea" yaml:"sidearea"`
     
    
     
     
    
     
     
    
     
     
        Title  fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
     
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func ( *Frontpage ) TableName() string{
	 return "dm_frontpage"
}

func ( *Frontpage ) ContentType() string{
	 return "frontpage"
}

func (c *Frontpage ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *Frontpage) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}


//todo: cache this? (then you need a reload?)
func (c *Frontpage) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        result["mainarea"]=c.Mainarea
        
    
        
    
        
        result["sidearea"]=c.Sidearea
        
    
        
    
        
    
        
        result["title"]=c.Title
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Frontpage) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "mainarea","mainarea_blocks","sidearea","sidearea_blocks","slideshow","title",}...)
}

func (c *Frontpage) Definition() contenttype.ContentTypeSetting {
	return contenttype.GetContentDefinition( c.ContentType() )
}

func (c *Frontpage) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    case "mainarea":
        
            result = c.Mainarea
        
    
    case "mainarea_blocks":
        
            result = c.Relations.Map["mainarea_blocks"]
        
    
    case "sidearea":
        
            result = c.Sidearea
        
    
    case "sidearea_blocks":
        
            result = c.Relations.Map["sidearea_blocks"]
        
    
    case "slideshow":
        
            result = c.Relations.Map["slideshow"]
        
    
    case "title":
        
            result = c.Title
        
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *Frontpage) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            case "mainarea":
            c.Mainarea = value.(fieldtype.TextField)
            
        
            
            
        
            
            
            case "sidearea":
            c.Sidearea = value.(fieldtype.TextField)
            
        
            
            
        
            
            
        
            
            
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
func (c *Frontpage) Store(transaction ...*sql.Tx) error {
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

func (c *Frontpage)StoreWithLocation(){

}

//Delete content only
func (c *Frontpage) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Frontpage{}
	}

	newList := func() interface{} {
		return &[]Frontpage{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]Frontpage)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("frontpage",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
