//Package dm/codegen/main generate content entity model based on contenttype.json.
package main

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	_ "github.com/digimakergo/digimaker/codegen/contenttypes/temp"

	_ "github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/util"
)

//Generate content types
func main() {

	fmt.Println("Generating content entities for " + util.AbsHomePath())
	err := Generate("entity")
	if err != nil {
		fmt.Println("Fail to generate: " + err.Error())
	}
}

func Generate(subFolder string) error {
	tpl := template.Must(template.New("contenttype.tpl").
		Funcs(funcMap()).
		ParseFiles(util.DMPath() + "/codegen/contenttypes/contenttype.tpl"))

	tplEntity := template.Must(template.New("contenttype_entity.tpl").
		Funcs(funcMap()).
		ParseFiles(util.DMPath() + "/codegen/contenttypes/contenttype_entity.tpl"))

	fieldtypeMap := fieldtype.GetAllDefinition()

	contentTypeDef := definition.GetDefinitionList()["default"]
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtype.GetAllDefinition()
		vars["name"] = name
		imports := []string{}
		for _, ftype := range settings.FieldMap {
			typeStr := ftype.FieldType
			if _, ok := fieldtypeMap[typeStr]; !ftype.IsOutput && !ok {
				return errors.New("field type " + typeStr + " doesn't exist.")
			}

			importName := fieldtype.GetDef(ftype.FieldType).Import
			if importName != "" && !util.Contains(imports, importName) {
				imports = append(imports, importName)
			}
		}
		vars["imports"] = imports
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
		var err error
		if settings.HasLocation {
			err = tpl.Execute(file, vars)
		} else {
			err = tplEntity.Execute(file, vars)
		}
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
