//Package dm/codegen/main generate content entity model based on contenttype.json.
package main

import (
	"errors"
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

	fmt.Println("Generating content entities for " + homePath)
	err := Generate(homePath, "entity")
	if err != nil {
		fmt.Println("Fail to generate: " + err.Error())
	}
}

func Generate(homePath string, subFolder string) error {

	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles(util.DMPath() + "/codegen/contenttypes/contenttype.tpl"))

	fieldtypeMap := fieldtype.GetAllDefinition()

	contentTypeDef := contenttype.GetDefinitionList()["default"]
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtype.GetAllDefinition()
		vars["name"] = name
		for _, ftype := range settings.FieldMap {
			typeStr := ftype.FieldType
			if _, ok := fieldtypeMap[typeStr]; !ftype.IsOutput && !ok {
				return errors.New("field type " + typeStr + " doesn't exist.")
			}
		}
		vars["fields"] = settings.FieldMap
		datafieldMap := map[string]string{}
		for _, item := range settings.DataFields {
			datafieldMap[item.Identifier] = item.FieldType
		}
		vars["data_fields"] = datafieldMap

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
