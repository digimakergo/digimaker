//Package dm/codegen/main generate content entity model based on contenttype.json.
package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/util"
)

//Generate content types
func main() {
	homePath := ""
	if len(os.Args) >= 2 && os.Args[1] != "" {
		homePath = os.Args[1]
		util.InitHomePath(homePath)
	}

	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	fmt.Println("Generating content entities for " + homePath)
	err := Generate(homePath, "entity")
	if err != nil {
		fmt.Println("Fail to generate: " + err.Error())
	}
}

func Generate(homePath string, subFolder string) error {

	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles(os.Getenv("GOPATH") + "/src/dm/core/codegen/contenttypes/contenttype.tpl"))

	contentTypeDef := contenttype.GetDefinitionList()["default"]
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtype.GetAllDefinition()
		vars["name"] = name
		vars["fields"] = settings.FieldMap

		vars["settings"] = settings

		path := util.HomePath() + "/" + subFolder + "/" + name + ".go"
		//todo: genereate to a template folder first and then copy&override target,
		//and if there is error remove that folder
		fmt.Println("Generating " + name)
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
