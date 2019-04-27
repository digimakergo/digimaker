//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "dm/db"
    "dm/contenttype"
	"dm/fieldtype"
    "dm/util"
	. "dm/query"
)

{{$struct_name :=.name|UpperName}}

type {{$struct_name}} struct{
     ContentCommon `boil:",bind"`
    {{range $identifier, $fieldtype := .settings.Fields}}
     {{$type_settings := index $.def_fieldtype $fieldtype.FieldType}}
     {{if not $type_settings.IsRelation }}
        {{$identifier|UpperName}} fieldtype.{{$type_settings.Value}} `boil:"{{$identifier}}" json:"{{$identifier}}" toml:"{{$identifier}}" yaml:"{{$identifier}}"`
     {{end}}
    {{end}}
     Location `boil:"location,bind"`
}

func ( *{{$struct_name}} ) TableName() string{
	 return "{{.settings.TableName}}"
}

func ( *{{$struct_name}} ) ContentType() string{
	 return "{{.name}}"
}


//todo: cache this? (then you need a reload?)
func (c *{{$struct_name}}) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
    {{range $identifier, $fieldtype := .settings.Fields}}
        {{if not (index $.def_fieldtype $fieldtype.FieldType).IsRelation}}
        result["{{$identifier}}"]=c.{{$identifier|UpperName}}
        {{end}}
    {{end}}
	for key, value := range c.ContentCommon.Values() {
		result[key] = value
	}
	return result
}

func (c *{{$struct_name}}) IdentifierList() []string {
	return append(c.ContentCommon.IdentifierList(),[]string{ {{range $identifier, $fieldtype := .settings.Fields}}"{{$identifier}}",{{end}}}...)
}

func (c *{{$struct_name}}) Value(identifier string) interface{} {
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    var result interface{}
	switch identifier {
    {{range $identifier, $fieldtype := .settings.Fields}}
    case "{{$identifier}}":
        {{if not (index $.def_fieldtype $fieldtype.FieldType).IsRelation}}
            result = c.{{$identifier|UpperName}}
        {{else}}
            result = c.Relations.Value["{{$identifier}}"]
        {{end}}
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
            {{if not $type_settings.IsRelation}}
            case "{{$identifier}}":
            c.{{$identifier|UpperName}} = value.(fieldtype.{{$type_settings.Value}})
            {{end}}
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
		id, err := handler.Insert(c.TableName(), c.ToMap())
		c.CID = id
		if err != nil {
			return err
		}
	} else {
		err := handler.Update(c.TableName(), c.ToMap(), Cond("id", c.CID))
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
