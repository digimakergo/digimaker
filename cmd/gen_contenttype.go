package main

import (
	"dm/contenttype"
	"dm/fieldtype"
	"fmt"
	"html/template"
	"os"
	"strings"
)

//Generate content types
func main() {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()
	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles("../codegen/contenttype.tpl"))

	folder := "/Users/xc/go/caf-prototype/src/dm/contenttype/entity"
	contentTypeDef := contenttype.GetDefinition()
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtype.GetDefinition()
		vars["name"] = name + "f"
		vars["settings"] = settings

		path := folder + "/" + name + ".go"
		file, _ := os.Create(path)
		err := tpl.Execute(file, vars)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func funcMap() template.FuncMap {
	funcMap := template.FuncMap{
		"Upper": strings.Title,
	}
	return funcMap
}
