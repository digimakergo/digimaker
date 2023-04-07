//This file is generated automatically, DO NOT EDIT IT.
//Generated time:

package entity

import (
    "context"
    "database/sql"
    "github.com/digimakergo/digimaker/core/db"
    "github.com/digimakergo/digimaker/core/contenttype"
    "github.com/digimakergo/digimaker/core/definition"    
    {{if .settings.HasLocation}}
    "github.com/digimakergo/digimaker/core/util"
    {{end}}
	. "github.com/digimakergo/digimaker/core/db"
    {{range $i, $import := .imports }}
    "{{$import}}"
    {{end}}
)

{{$struct_name :=.name|UpperName}}

type {{$struct_name}} struct{
     contenttype.Metadata `boil:",bind" json:"metadata"`

     ID int `boil:"id" json:"id" toml:"id" yaml:"id"`

     {{range $identifier, $fieldtype := .data_fields}}
            {{if not ($identifier|InternalIdentifier)}}
                {{$identifier|UpperName}}  {{$fieldtype}} `boil:"{{$identifier}}" json:"{{$identifier}}" toml:"{{$identifier}}" yaml:"{{$identifier}}"`
            {{end}}
     {{end}}
    {{range $identifier, $fieldtype := .fields}}
         {{$type_settings := index $.def_fieldtype $fieldtype.FieldType}}    
         {{if not $fieldtype.IsOutput}}
            {{$identifier|UpperName}}  {{$type_settings.DataType}} `boil:"{{if not (eq $fieldtype.FieldType "relationlist" ) }}{{$identifier}}{{else}}-{{end}}" json:"{{if eq $fieldtype.FieldType "password"}}-{{else}}{{$identifier}}{{end}}" toml:"{{$identifier}}" yaml:"{{$identifier}}"`
         {{end}}        
    {{end}}
}

func (c {{$struct_name}} ) GetID() int{
        return c.ID
}

func (c *{{$struct_name}} ) GetMetadata() *contenttype.Metadata{
        return &c.Metadata
}

func (c *{{$struct_name}}) GetLocation() *contenttype.Location{
    return nil
}

//Get map of the all fields(including data_fields)
//todo: cache this? (then you need a reload?)
func (c *{{$struct_name}}) ToDBValues() map[string]interface{} {
	result := make(map[string]interface{})
    {{range $identifier, $fieldtype := .data_fields}}
         {{if not ($identifier|InternalIdentifier)}}
         result["{{$identifier}}"]=c.{{$identifier|UpperName}}
         {{end}}
    {{end}}

    {{range $identifier, $fieldtype := .fields}}
        {{if not ( eq $fieldtype.FieldType "relationlist" ) }}
        {{if not $fieldtype.IsOutput}}
            result["{{$identifier}}"]=c.{{$identifier|UpperName}}
        {{end}}
        {{end}}
    {{end}}

    for key, value := range c.Metadata.ToDBValues() {
		result[key] = value
	}

	return result
}

//Get identifier list of fields(NOT including data_fields )
func (c *{{$struct_name}}) IdentifierList() []string {
	return []string{ {{range $identifier, $fieldtype := .fields}}{{if not $fieldtype.IsOutput}}"{{$identifier}}",{{end}}{{end}}}
}

func (c *{{$struct_name}}) Definition( language ...string) definition.ContentType {
	def, _ := definition.GetDefinition( c.ContentType(), language... )
    return def
}

//Get field value
func (c *{{$struct_name}}) Value(identifier string) interface{} {
    {{if .settings.HasLocation}}
    if util.Contains( c.Location.IdentifierList(), identifier ) {
        return c.Location.Field( identifier )
    }
    {{end}}
    var result interface{}
	switch identifier {
    {{range $identifier, $fieldtype := .data_fields}}
      {{if not ($identifier|InternalIdentifier)}}
      case "{{$identifier}}":
         result = c.{{$identifier|UpperName}}
      {{end}}
    {{end}}
    {{range $identifier, $fieldtype := .fields}}
    {{if not $fieldtype.IsOutput}}
    case "{{$identifier}}":        
            result = (c.{{$identifier|UpperName}})        
    {{end}}
    {{end}}

    default:
    }
	return result
}

//Set value to a field
func (c *{{$struct_name}}) SetValue(identifier string, value interface{}) error {
	switch identifier {
        {{range $identifier, $fieldtype := .data_fields}}
         {{if not ($identifier|InternalIdentifier)}}
          case "{{$identifier}}":
             c.{{$identifier|UpperName}} = value.({{$fieldtype}})
        {{end}}
        {{end}}
        {{range $identifier, $fieldtype := .fields}}
            {{$type_settings := index $.def_fieldtype $fieldtype.FieldType}}
            {{if not $fieldtype.IsOutput}}
            case "{{$identifier}}":
            c.{{$identifier|UpperName}} = value.({{$type_settings.DataType}})
            {{end}}         
        {{end}}
	}
	//todo: check if identifier exist
	return nil
}

//Store content.
//Note: it will set id to ID after success
func (c *{{$struct_name}}) Store(ctx context.Context, transaction ...*sql.Tx) error {
	if c.ID == 0 {
		id, err := db.Insert(ctx, "{{.settings.TableName}}", c.ToDBValues(), transaction...)
		c.ID = id
		if err != nil {
			return err
		}
	} else {
		err := db.Update(ctx, "{{.settings.TableName}}", c.ToDBValues(), Cond("id", c.ID), transaction...)
		return err
	}
	return nil
}


func (c *{{$struct_name}})StoreWithLocation(){

}

//Delete content only
func (c *{{$struct_name}}) Delete(ctx context.Context, transaction ...*sql.Tx) error {
	contentError := db.Delete(ctx, "{{.settings.TableName}}", Cond("id", c.ID), transaction...)
	return contentError
}

func init() {
	new := func() contenttype.ContentTyper {
    entity := &{{$struct_name}}{}
    entity.Metadata.Contenttype = "{{.name}}"
    return entity
	}

	newList := func() interface{} {
		return &[]{{$struct_name}}{}
	}

    toList := func(obj interface{}) []contenttype.ContentTyper {
        contentList := *obj.(*[]{{$struct_name}})
        list := make([]contenttype.ContentTyper, len(contentList))
        for i, _ := range contentList {
            list[i] = &contentList[i]
        }
        return list
    }

	contenttype.Register( "{{.name}}",
		contenttype.ContentTypeRegister{
			New:            new,
			NewList:        newList,
            ToList:         toList})
}
