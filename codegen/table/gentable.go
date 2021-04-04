package main

import (
	"fmt"
	"os"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/util"
)

func main() {
	contenttypes := ""
	if len(os.Args) >= 2 && os.Args[1] != "" {
		contenttypes = os.Args[1]
	}
	if contenttypes == "" {
		fmt.Println("Please use like './gentable article', separated by ,")
		return
	}

	fmt.Println("Generating table for " + contenttypes)
	for _, contenttype := range util.Split(contenttypes, ",") {
		err := GenerateTable(contenttype)
		if err != nil {
			fmt.Println("Fail to generate: " + err.Error())
		}
	}
}

func GenerateTable(ctype string) error {
	def, err := definition.GetDefinition(ctype)
	if err != nil {
		return err
	}
	result := "CREATE TABLE " + def.TableName + " ( id INT AUTO_INCREMENT PRIMARY KEY, "
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
