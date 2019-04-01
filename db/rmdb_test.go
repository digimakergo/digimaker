package db

import (
	"dm/type_default/field"
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	rmdb := new(RMDB)
	article := rmdb.GetByID("article", 2)

	text := article.Field("Title").(field.TextField)
	contentType := article.Field("ContentType").(string)
	fmt.Printf("title: " + text.Value() + ", content_type: " + contentType + "\n")
	t.Logf(text.Value())

}
