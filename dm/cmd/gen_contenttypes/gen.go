package main

import (
	"dm/dm/contenttype"
	"dm/dm/fieldtype"
	"dm/dm/util"
	"fmt"
	"os"
	"text/template"
)

//Generate content types
func main() {
	baseFolder := "/Users/xc/go/caf-prototype/src/dm"

	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	err := Generate(baseFolder)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Generate(baseFolder string) error {
	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles(baseFolder + "/cmd/gen_contenttypes/contenttype.tpl"))

	folder := baseFolder + "/contenttype/entity"
	contentTypeDef := contenttype.GetDefinition()
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtype.GetDefinition()
		vars["name"] = name
		vars["settings"] = settings

		path := folder + "/" + name + ".go"
		//todo: genereate to a template folder first and then copy&override target,
		//and if there is error remove that folder
		file, _ := os.Create(path)
		err := tpl.Execute(file, vars)
		if err != nil {
			return err
		}
	}
	return nil
}

func funcMap() template.FuncMap {
	funcMap := template.FuncMap{
		"UpperName": util.UpperName,
	}
	return funcMap
}
