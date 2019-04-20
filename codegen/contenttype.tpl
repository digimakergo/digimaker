//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
	"dm/fieldtype"
	. "dm/query"
)

{{$struct_name :=.name|Upper}}

type {{$struct_name}} struct{
     *Location
     *ContentCommon
    {{range $identifier, $fieldtype := .settings.Fields}}
     {{$type_settings := index $.def_fieldtype $fieldtype.FieldType}}
     {{$identifier|Upper}} fieldtype.{{$type_settings.Value}} `boil:"{{$identifier}}" json:"{{$identifier}}" toml:"{{$identifier}}" yaml:"{{$identifier}}"`
    {{end}}
}


func ( {{$struct_name}} ) TableName() string{
	 return "{{.settings.TableName}}"
}


func (c {{$struct_name}}) Values() map[string]interface{} {
	result := make(map[string]interface{})

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
