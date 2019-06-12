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
	baseFolder := "admin"

	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	err := Generate(baseFolder)
	if err != nil {
		fmt.Println("Fail to generate: " + err.Error())
	}
}

func Generate(baseFolder string) error {

	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles("dm/codegen/contenttypes/contenttype.tpl"))

	folder := baseFolder + "/entity"
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
