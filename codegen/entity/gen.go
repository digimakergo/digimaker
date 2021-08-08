//Package dm/codegen/main generate content entity model based on contenttype.json.
package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"text/template"

	_ "github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	_ "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/digimakergo/digimaker/core/util"
)

//go:embed contenttype.tpl
//go:embed contenttype_entity.tpl
var fs embed.FS

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
		Funcs(funcMap()).ParseFS(fs, "contenttype.tpl"))

	tplEntity := template.Must(template.New("contenttype_entity.tpl").
		Funcs(funcMap()).ParseFS(fs, "contenttype_entity.tpl"))

	fieldtypeMap := fieldtype.GetAllFieldtype()

	contentTypeDef := definition.GetDefinitionList()["default"]
	for name, settings := range contentTypeDef {
		vars := map[string]interface{}{}
		vars["def_fieldtype"] = fieldtypeMap
		vars["name"] = name
		imports := []string{}
		for _, ftype := range settings.FieldMap {
			typeStr := ftype.FieldType
			if _, ok := fieldtypeMap[typeStr]; !ftype.IsOutput && !ok {
				return errors.New("field type " + typeStr + " doesn't exist.")
			}

			packagePath := fieldtype.GetFieldtype(ftype.FieldType).Package
			if packagePath != "" && !util.Contains(imports, packagePath) {
				imports = append(imports, packagePath)
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
