package handler

import (
	"dm/model"
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	model.LoadDefinition()

	handler := ContentHandler{}
	params := map[string]interface{}{"title": "ff", "body": "Hello"}
	result, _ := handler.Validate("article", params)
	fmt.Println(result)
	//handler.Draft("article", 1)

}
