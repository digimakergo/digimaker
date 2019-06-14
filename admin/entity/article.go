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



type Article struct{
     ContentCommon `boil:",bind"`
    
     
     
        Body  fieldtype.RichTextField `boil:"body" json:"body" toml:"body" yaml:"body"`
     
    
     
     
    
     
     
        Editors  fieldtype.EditorList `boil:"editors" json:"editors" toml:"editors" yaml:"editors"`
     
    
     
     
    
     
     
        Summary  fieldtype.RichTextField `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
     
    
     
     
        Title  fieldtype.TextField `boil:"title" json:"title" toml:"title" yaml:"title"`
     
    
     
     
    
    
     contenttype.Location `boil:"location,bind"  json:"location"`
    
}

func ( *Article ) TableName() string{
	 return "dm_article"
}

func ( *Article ) ContentType() string{
	 return "article"
}

func (c *Article ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *Article) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}


//todo: cache this? (then you need a reload?)
func (c *Article) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        result["body"]=c.Body
        
    
        
    
        
        result["editors"]=c.Editors
        
    
        
    
        
        result["summary"]=c.Summary
        
    
        
        result["title"]=c.Title
        
    
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Article) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "body","coverimage","editors","related_articles","summary","title","useful_resources",}...)
}

func (c *Article) Definition() contenttype.ContentTypeSetting {
	return contenttype.GetContentDefinition( c.ContentType() )
}

func (c *Article) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    case "body":
        
            result = c.Body
        
    
    case "coverimage":
        
            result = c.Relations.Map["coverimage"]
        
    
    case "editors":
        
            result = c.Editors
        
    
    case "related_articles":
        
            result = c.Relations.Map["related_articles"]
        
    
    case "summary":
        
            result = c.Summary
        
    
    case "title":
        
            result = c.Title
        
    
    case "useful_resources":
        
            result = c.Relations.Map["useful_resources"]
        
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *Article) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            case "body":
            c.Body = value.(fieldtype.RichTextField)
            
        
            
            
        
            
            
            case "editors":
            c.Editors = value.(fieldtype.EditorList)
            
        
            
            
        
            
            
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
func (c *Article) Store(transaction ...*sql.Tx) error {
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

func (c *Article)StoreWithLocation(){

}

//Delete content only
func (c *Article) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Article{}
	}

	newList := func() interface{} {
		return &[]Article{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]Article)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("article",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
