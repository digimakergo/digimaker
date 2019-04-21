package handler

import (
	"dm/contenttype"
	"dm/query"
	"fmt"
	"testing"
)

func TestValidtion(t *testing.T) {

	// handler := ContentHandler{}
	// params := map[string]interface{}{"title": "ff", "body": "Hello"}
	// passed, result := handler.Validate("article", params)
	// fmt.Println(result)
	// assert.Equal(t, passed, true)
	//
	// params = map[string]interface{}{"title": "", "body": "Hello"}
	// _, result = handler.Validate("article", params)
	// assert.Equal(t, result.Fields[0].Identifier, "title")
	//handler.Draft("article", 1)

	contenttype.LoadDefinition()
	list, _ := Querier().List("article", query.Cond("1", 1))
	fmt.Println(list)
}
