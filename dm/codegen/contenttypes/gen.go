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
	packageName := ""
	if len(os.Args) >= 2 && os.Args[1] != "" {
		packageName = os.Args[1]
		util.SetPackageName(packageName)
	}

	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	fmt.Println("Generating content entities for " + packageName)
	err := Generate(packageName, "entity")
	if err != nil {
		fmt.Println("Fail to generate: " + err.Error())
	}
}

func Generate(packageName string, subFolder string) error {

	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles(os.Getenv("GOPATH") + "/src/dm/dm/codegen/contenttypes/contenttype.tpl"))

	contentTypeDef := contenttype.GetDefinition()
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtype.GetDefinition()
		vars["name"] = name
		vars["settings"] = settings

		path := util.HomePath() + "/" + subFolder + "/" + name + ".go"
		//todo: genereate to a template folder first and then copy&override target,
		//and if there is error remove that folder
		file, _ := os.Create(path)
		err := tpl.Execute(file, vars)
		if err != nil {
			return err
		}
		fmt.Println("Generated " + name + ".go")
	}
	return nil
}

func funcMap() template.FuncMap {
	funcMap := template.FuncMap{
		"UpperName": util.UpperName,
	}
	return funcMap
}
