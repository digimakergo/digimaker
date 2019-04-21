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
     Location `boil:"dm_location,bind"`
     ContentCommon `boil:",bind"`
    {{range $identifier, $fieldtype := .settings.Fields}}
     {{$type_settings := index $.def_fieldtype $fieldtype.FieldType}}
     {{$identifier|UpperName}} fieldtype.{{$type_settings.Value}} `boil:"{{$identifier}}" json:"{{$identifier}}" toml:"{{$identifier}}" yaml:"{{$identifier}}"`
    {{end}}
}


func ( {{$struct_name}} ) TableName() string{
	 return "{{.settings.TableName}}"
}


func (c {{$struct_name}}) Values() map[string]interface{} {
	result := make(map[string]interface{})

    {{range $identifier, $fieldtype := .settings.Fields}}
        result["{{$identifier}}"]=c.{{$identifier|UpperName}}
    {{end}}

    for key, value := range c.ContentCommon.Values() {
        result[key] = value
    }

	for key, value := range c.Location.Values() {
		result[key] = value
	}
	return result
}

func (c {{$struct_name}}) Store() error {
	handler := db.DBHanlder()
	if c.CID == 0 {
		id, err := handler.Insert(c.TableName(), c.Values())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.Values(), Cond("id", c.CID))
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

	convert := func(obj interface{}) []contenttype.ContentTyper {
		list := obj.(*[]{{$struct_name}})
		var result []contenttype.ContentTyper
		for _, item := range *list {
			result = append(result, item)
		}
		return result
	}

	Register("{{.name}}",
		ContentTypeRegister{
			New:            new,
			NewList:        newList,
			ToContentTyper: convert})
}
