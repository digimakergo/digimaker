//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "github.com/xc/digimaker/core/db"
    "github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/fieldtype"
    
    "github.com/xc/digimaker/core/util"
    
	. "github.com/xc/digimaker/core/db"
)



type Article struct{
     contenttype.ContentCommon `boil:",bind"`

     
    
         
         
         
            Body  fieldtype.RichText `boil:"body" json:"body" toml:"body" yaml:"body"`
         
        
    
         
         
         
            Coverimage  fieldtype.Text `boil:"coverimage" json:"coverimage" toml:"coverimage" yaml:"coverimage"`
         
        
    
         
         
         
            Editors  fieldtype.Text `boil:"editors" json:"editors" toml:"editors" yaml:"editors"`
         
        
    
         
         
    
         
         
         
        
    
         
         
         
            Summary  fieldtype.RichText `boil:"summary" json:"summary" toml:"summary" yaml:"summary"`
         
        
    
         
         
         
            Title  fieldtype.Text `boil:"title" json:"title" toml:"title" yaml:"title"`
         
        
    
         
         
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func (c *Article ) TableName() string{
	 return "dm_article"
}

func (c *Article ) ContentType() string{
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

func (c *Article) ToMap() map[string]interface{}{
    result := map[string]interface{}{}
    for _, identifier := range c.IdentifierList(){
      result[identifier] = c.Value(identifier)
    }
    return result
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *Article) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    

    
        
        
            result["body"]=c.Body
        
        
    
        
        
            result["coverimage"]=c.Coverimage
        
        
    
        
        
            result["editors"]=c.Editors
        
        
    
        
    
        
        
        
    
        
        
            result["summary"]=c.Summary
        
        
    
        
        
            result["title"]=c.Title
        
        
    
        
    
	for key, value := range c.ContentCommon.ToDBValues() {
		result[key] = value
	}
	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *Article) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "body","coverimage","editors","related_articles","summary","title","useful_resources",}...)
}

func (c *Article) Definition(language ...string) contenttype.ContentType {
	def, _ := contenttype.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *Article) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    
    
    case "body":
        
            result = c.Body
        
    
    
    
    case "coverimage":
        
            result = c.Coverimage
        
    
    
    
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

//Set value to a field
func (c *Article) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
        
            
            
            
            case "body":
            c.Body = *(value.(*fieldtype.RichText))
            
            
        
            
            
            
            case "coverimage":
            c.Coverimage = *(value.(*fieldtype.Text))
            
            
        
            
            
            
            case "editors":
            c.Editors = *(value.(*fieldtype.Text))
            
            
        
            
            
        
            
            
            
            
        
            
            
            
            case "summary":
            c.Summary = *(value.(*fieldtype.RichText))
            
            
        
            
            
            
            case "title":
            c.Title = *(value.(*fieldtype.Text))
            
            
        
            
            
        
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
		id, err := handler.Insert(c.TableName(), c.ToDBValues(), transaction...)
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.ToDBValues(), Cond("id", c.CID), transaction...)
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
