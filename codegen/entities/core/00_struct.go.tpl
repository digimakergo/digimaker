{{- $alias := .Aliases.Table .Table.Name -}}

// {{$alias.UpSingular}} is an object representing the database table.
// Implement dm.model.ContentTyper interface
type {{$alias.UpSingular}} struct {
	{{- range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{- if eq $.StructTagCasing "camel"}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | camelCase}}" yaml:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{- else -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{end -}}
}

func ( c *{{$alias.UpSingular}} ) Fields() map[string]model.Fielder{
	 return nil
}

func ( c *{{$alias.UpSingular}} ) Values() map[string]interface{}{
    result := make(map[string]interface{})
    {{range $column := .Table.Columns -}}
    {{- $colAlias := $alias.Column $column.Name -}}
      result["{{$column.Name}}"]= c.{{$colAlias}}
    {{end -}}
    return result
}

func ( c *{{$alias.UpSingular}} ) TableName() string{
	 return "{{.Table.Name}}"
}

func ( c *{{$alias.UpSingular}} ) Field( name string ) interface{}{
	  var result interface{}
		switch name {
			    {{range $column := .Table.Columns -}}
			    {{- $colAlias := $alias.Column $column.Name -}}
			    case "{{$column.Name}}","{{$colAlias}}" :
			      result = c.{{$colAlias}}
				  {{end -}}
			    default:
		}
		return result
}

func (c {{$alias.UpSingular}}) Store() error {
    handler := db.DBHanlder()
    if c.ID == 0 {
        id, err := handler.Insert(c.TableName(), c.Values())
        c.ID = id
        if err != nil {
            return err
        }
    } else {
        err := handler.Update(c.TableName(), c.Values(), Cond("id", c.ID))
        return err
    }
    return nil
}
