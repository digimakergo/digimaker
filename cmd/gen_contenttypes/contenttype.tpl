//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
	. "dm/query"
)

{{$struct_name :=.name|UpperName}}

type {{$struct_name}} struct{
     ContentCommon `boil:",bind"`
    {{range $identifier, $fieldtype := .settings.Fields}}
     {{$type_settings := index $.def_fieldtype $fieldtype.FieldType}}
     {{$identifier|UpperName}} fieldtype.{{$type_settings.Value}} `boil:"{{$identifier}}" json:"{{$identifier}}" toml:"{{$identifier}}" yaml:"{{$identifier}}"`
    {{end}}
     Location `boil:"location,bind"`
}

func ( *{{$struct_name}} ) TableName() string{
	 return "{{.settings.TableName}}"
}

func (c *{{$struct_name}}) contentValues() map[string]interface{} {
	result := make(map[string]interface{})
    {{range $identifier, $fieldtype := .settings.Fields}}
        result["{{$identifier}}"]=c.{{$identifier|UpperName}}
    {{end}}
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

//todo: cache this(maybe cache map in a private property?)
//todo: maybe return all field identifers as []string?
func (c *{{$struct_name}}) Values() map[string]interface{} {
    result := c.contentValues()

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

func (c *{{$struct_name}}) Value(identifier string) interface{} {
	var result interface{}
	switch identifier {
    {{range $identifier, $fieldtype := .settings.Fields}}
    case "{{$identifier}}":
        result = c.{{$identifier|UpperName}}
    {{end}}
	case "cid":
		result = c.ContentCommon.CID
    default:
    	result = c.ContentCommon.Value( identifier )
    }
	return result
}


func (c *{{$struct_name}}) SetValue(identifier string, value interface{}) error {
	switch identifier {
        {{range $identifier, $fieldtype := .settings.Fields}}
             {{$type_settings := index $.def_fieldtype $fieldtype.FieldType}}
            case "{{$identifier}}":
            c.{{$identifier|UpperName}} = value.(fieldtype.{{$type_settings.Value}})
        {{end}}
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
func (c *{{$struct_name}}) Store() error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.contentValues())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.contentValues(), Cond("id", c.CID))
		return err
	}
	return nil
}


func init() {
	new := func() contenttype.ContentTyper {
		return &{{$struct_name}}{}
	}

	newList := func() interface{} {
		return &[]{{$struct_name}}{}
	}

	Register("{{.name}}",
		ContentTypeRegister{
			New:            new,
			NewList:        newList})
}
