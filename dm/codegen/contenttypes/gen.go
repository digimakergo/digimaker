//Package dm/codegen/main generate content entity model based on contenttype.json.
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
	projectName := ""
	if len(os.Args) >= 2 && os.Args[1] != "" {
		projectName = os.Args[1]
		util.SetConfigPath(projectName + "/configs")
	}

	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	err := Generate(projectName + "/entity")
	if err != nil {
		fmt.Println("Fail to generate: " + err.Error())
	}
}

func Generate(outputFolder string) error {

	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles("dm/codegen/contenttypes/contenttype.tpl"))

	contentTypeDef := contenttype.GetDefinition()
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtype.GetDefinition()
		vars["name"] = name
		vars["settings"] = settings

		path := outputFolder + "/" + name + ".go"
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
