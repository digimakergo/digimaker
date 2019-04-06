package db

import (
	"dm/model"
	"dm/model/entity"
	"dm/type_default/field"
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	model.LoadDefinition()

	rmdb := new(RMDB)
	var article entity.Article
	rmdb.GetByID("article", 2, &article)

	text := article.Field("Title").(field.TextField)
	//contentType := article.Field("ContentType").(string)
	fmt.Printf("title: " + text.Value() + "\n")
	fmt.Print(article)
	t.Logf(text.Value())

}
