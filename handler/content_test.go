package handler

import (
	"dm/model"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidtion(t *testing.T) {
	model.LoadDefinition()

	handler := ContentHandler{}
	params := map[string]interface{}{"title": "ff", "body": "Hello"}
	result := handler.Validate("article", params)
	fmt.Println(result)
	assert.Equal(t, result, ValidationResult{})

	params = map[string]interface{}{"title": "", "body": "Hello"}
	result = handler.Validate("article", params)
	assert.Equal(t, result.Fields[0].Identifier, "title")
	//handler.Draft("article", 1)

}
