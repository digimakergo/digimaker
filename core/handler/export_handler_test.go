//Author xc, Created on 2019-05-12 15:25
//{COPYRIGHTS}
package handler

import (
	"fmt"
	"testing"

	"github.com/digimakergo/digimaker/core/contenttype"
)

func TestExport(t *testing.T) {
	contenttype.LoadDefinition()

	content, _ := Querier().FetchByID(3)

	mh := ExportHandler{}
	parent, _ := Querier().FetchByID(content.Value("parent_id").(int))
	str, err := mh.Export(content, parent)
	fmt.Println("Exporting:")
	fmt.Println(str, err)
}

func TestImport(t *testing.T) {
	// testData := `{"data":{"body":{"Data":"Hello world..."},"cuid":"bj8a5rq23akpc5md8b46","editors":{"Data":""},"location":{"author":0,"content_type":"article","identifier_path":"","is_hidden":false,"is_invisible":false,"language":"","main_id":76,"name":"Good morning2233","p":"","priority":0,"section":"","uid":"bj8a5rq23akpc5md8b41"},"modified":1557642283,"published":1557177071,"relations":{"Value":{}},"summary":{"Data":"updated3"},"title":{"Data":"Good morning2233"},"version":25},"content_type":"article","cuid":"bj8a5rq23akpc5md8b46","parent_uid":"bj8a5rq53akpc5md8b45"}`
	// mh := MigrationHandler{}
	// err := mh.ImportALine([]byte(testData))
	// fmt.Println(err)
}
