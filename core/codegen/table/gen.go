package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/util"
)

func main() {
	homePath := ""
	tableName := ""
	if len(os.Args) >= 3 && os.Args[1] != "" {
		homePath = os.Args[1]
		util.InitHomePath(homePath)
		tableName = os.Args[2]
	}

	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	fmt.Println("Generating table " + tableName + " in " + homePath)
	err := GenerateTable(tableName)
	if err != nil {
		fmt.Println("Fail to generate: " + err.Error())
	}
}

func GenerateTable(ctype string) error {
	def, err := contenttype.GetDefinition(ctype)
	if err != nil {
		return err
	}
	fmap := def.FieldMap
	result := "CREATE TABLE " + def.TableName + " ( id INT AUTO_INCREMENT PRIMARY KEY, "
	for identifier, value := range fmap {
		dbType := ""
		switch value.FieldType {
		case "checkbox":
			dbType = "INT NOT NULL DEFAULT 0"
		case "text":
			dbType = "varchar(255) NOT NULL DEFAULT ''"
		case "richtext":
			dbType = "TEXT NOT NULL DEFAULT ''"
		case "container":
			continue
		case "eth_indicator":
			continue
		default:
			return errors.New("Not supported fieldtype." + value.FieldType)
		}
		result += identifier + " " + dbType + ", "
	}
	result += "published int NOT NULL DEFAULT 0, modified int NOT NULL DEFAULT 0, cuid varchar(30) NOT NULL DEFAULT '')"
	fmt.Println(result)
	return nil
}
