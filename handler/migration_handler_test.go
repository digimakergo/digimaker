//Author xc, Created on 2019-05-12 15:25
//{COPYRIGHTS}
package handler

import (
	"dm/contenttype"
	"dm/fieldtype"
	"fmt"
	"testing"
)

func TestExport(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	content, _ := Querier().FetchByID(76)

	mh := MigrationHandler{}
	parent, _ := Querier().FetchByID(content.Value("parent_id").(int))
	str, _ := mh.Export(content, parent)
	fmt.Println("Exporting:")
	fmt.Println(str)
}

func TestImport(t *testing.T) {
	testData := `{"article":{"body":{"Data":"Hello world..."},"parent_uid":"bj8a5rq53akpc5md8b45","cuid":"bj8a5rq53akpc5md8b45","editors":{"Data":""},"location":{"author":0,"content_type":"article","content_uid":"bj8a5rq23akpc5md8b40","identifier_path":"","is_hidden":false,"is_invisible":false,"language":"","main_id":76,"name":"Good morning2233","p":"","priority":0,"section":"","uid":"bj8a5rq23akpc5md8b4g"},"modified":1557642283,"published":1557177071,"relations":{"Value":{}},"summary":{"Data":"updated3"},"title":{"Data":"Good morning2233"},"version":25}}`
	mh := MigrationHandler{}
	err := mh.ImportALine([]byte(testData))
	fmt.Println(err)
}
