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
	err := rmdb.GetByID("article", 2, &article)

	text := article.Field("Title").(field.TextField)
	fmt.Printf("title: " + text.Value() + "\n")
	fmt.Println(err == nil)

	t.Logf(text.Value())
	var article2 entity.Article
	err = rmdb.GetByFields("article", map[string]interface{}{"content_id": 1}, &article2)
	fmt.Println(article2)
	fmt.Println(err == nil)
}
