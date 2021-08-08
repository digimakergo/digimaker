package table

import (
	"fmt"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	_ "github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
)

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
