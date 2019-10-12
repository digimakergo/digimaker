//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "database/sql"
    "dm/core/db"
    "dm/core/contenttype"
	"dm/core/fieldtype"
    
    "dm/core/util"
    
	. "dm/core/db"
)



type Report struct{
     contenttype.ContentCommon `boil:",bind"`
    
         
         
         
            Address  fieldtype.TextField `boil:"address" json:"address" toml:"address" yaml:"address"`
         
        
    
         
         
         
        
    
         
         
         
        
    
         
         
         
            Description  fieldtype.RichTextField `boil:"description" json:"description" toml:"description" yaml:"description"`
         
        
    
         
         
         
            Email  fieldtype.TextField `boil:"email" json:"email" toml:"email" yaml:"email"`
         
        
    
         
         
         
            Email1  fieldtype.TextField `boil:"email1" json:"email1" toml:"email1" yaml:"email1"`
         
        
    
         
         
         
            Mobile  fieldtype.TextField `boil:"mobile" json:"mobile" toml:"mobile" yaml:"mobile"`
         
        
    
         
         
         
            Mobile1  fieldtype.TextField `boil:"mobile1" json:"mobile1" toml:"mobile1" yaml:"mobile1"`
         
        
    
         
         
         
            Name  fieldtype.TextField `boil:"name" json:"name" toml:"name" yaml:"name"`
         
        
    
         
         
         
        
    
         
         
         
        
    
    
     contenttype.Location `boil:"location,bind"`
    
}

func ( *Report ) TableName() string{
	 return "eh_report"
}

func ( *Report ) ContentType() string{
	 return "report"
}

func (c *Report ) GetName() string{
	 location := c.GetLocation()
     if location != nil{
         return location.Name
     }else{
         return ""
     }
}

func (c *Report) GetLocation() *contenttype.Location{
    
    return &c.Location
    
}


//todo: cache this? (then you need a reload?)
func (c *Report) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    
        
        
            result["address"]=c.Address
        
        
    
        
        
        
    
        
        
        
    
        
        
            result["description"]=c.Description
        
        
    
        
        
            result["email"]=c.Email
        
        
    
        
        
            result["email1"]=c.Email1
        
        
    
        
        
            result["mobile"]=c.Mobile
        
        
    
        
        
            result["mobile1"]=c.Mobile1
        
        
    
        
        
            result["name"]=c.Name
        
        
    
        
        
        
    
        
        
        
    
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *Report) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ "address","basic","contact","description","email","email1","mobile","mobile1","name","step1","step2",}...)
}

func (c *Report) Definition() contenttype.ContentType {
	def, _ := contenttype.GetDefinition( c.ContentType() )
    return def
}

func (c *Report) Value(identifier string) interface{} {
    
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    
    var result interface{}
	switch identifier {
    
    
    case "address":
        
            result = c.Address
        
    
    
    
    
    
    
    
    case "description":
        
            result = c.Description
        
    
    
    
    case "email":
        
            result = c.Email
        
    
    
    
    case "email1":
        
            result = c.Email1
        
    
    
    
    case "mobile":
        
            result = c.Mobile
        
    
    
    
    case "mobile1":
        
            result = c.Mobile1
        
    
    
    
    case "name":
        
            result = c.Name
        
    
    
    
    
    
    
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *Report) SetValue(identifier string, value interface{}) error {
	switch identifier {
        
            
            
            
            case "address":
            c.Address = value.(fieldtype.TextField)
            
            
        
            
            
            
            
        
            
            
            
            
        
            
            
            
            case "description":
            c.Description = value.(fieldtype.RichTextField)
            
            
        
            
            
            
            case "email":
            c.Email = value.(fieldtype.TextField)
            
            
        
            
            
            
            case "email1":
            c.Email1 = value.(fieldtype.TextField)
            
            
        
            
            
            
            case "mobile":
            c.Mobile = value.(fieldtype.TextField)
            
            
        
            
            
            
            case "mobile1":
            c.Mobile1 = value.(fieldtype.TextField)
            
            
        
            
            
            
            case "name":
            c.Name = value.(fieldtype.TextField)
            
            
        
            
            
            
            
        
            
            
            
            
        
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
func (c *Report) Store(transaction ...*sql.Tx) error {
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

func (c *Report)StoreWithLocation(){

}

//Delete content only
func (c *Report) Delete(transaction ...*sql.Tx) error {
	handler := db.DBHanlder()
	contentError := handler.Delete(c.TableName(), Cond("id", c.CID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
		return &Report{}
	}

	newList := func() interface{} {
		return &[]Report{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]Report)
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register("report",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
