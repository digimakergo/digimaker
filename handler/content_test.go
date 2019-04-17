package handler

import (
	"dm/def"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidtion(t *testing.T) {
	def.LoadDefinition()

	handler := ContentHandler{}
	params := map[string]interface{}{"title": "ff", "body": "Hello"}
	passed, result := handler.Validate("article", params)
	fmt.Println(result)
	assert.Equal(t, passed, true)

	params = map[string]interface{}{"title": "", "body": "Hello"}
	_, result = handler.Validate("article", params)
	assert.Equal(t, result.Fields[0].Identifier, "title")
	//handler.Draft("article", 1)

}
