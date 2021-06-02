package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	_ "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/digimakergo/digimaker/core/util"
)

func main() {
	contenttypes := []string{}
	if len(os.Args) >= 2 && os.Args[1] != "" {
		contenttypes = util.Split(os.Args[1], ",")
	} else {
		contenttypeList := definition.GetDefinitionList()["default"]
		for identifier, _ := range contenttypeList {
			contenttypes = append(contenttypes, identifier)
		}
		sort.Strings(contenttypes)
	}

	fmt.Println("Generating table for " + strings.Join(contenttypes, ","))
	fmt.Println("----------")
	for _, contenttype := range contenttypes {
		err := GenerateTable(contenttype)
		if err != nil {
			fmt.Println("Fail to generate: " + err.Error())
		}
		fmt.Println("")
	}
}

func GenerateTable(ctype string) error {
	def, err := definition.GetDefinition(ctype)
	if err != nil {
		return err
	}
	result := "CREATE TABLE " + def.TableName + " (id INT AUTO_INCREMENT PRIMARY KEY, "
	identifierList := def.FieldIdentifierList
	for _, identifier := range identifierList {
		fieldDef := def.FieldMap[identifier]
		handler := fieldtype.GethHandler(fieldDef)
		if handler != nil {
			dbType := handler.DBField()
			if dbType != "" {
				result += identifier + " " + dbType + ", "
			}
		}
	}
	result += "author INT, published INT NOT NULL DEFAULT 0, modified INT NOT NULL DEFAULT 0, cuid VARCHAR(30) NOT NULL DEFAULT '')"
	fmt.Println(result)
	return nil
}
