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
	str, _ := mh.Export(content)
	fmt.Println("hello")
	fmt.Println(str)
}

func TestImport(t *testing.T) {

}
