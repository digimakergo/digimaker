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
	rmdb.GetByID("article", 1, &article)

	text := article.Field("Title").(field.TextField)

	fmt.Printf("title: " + text.Value() + "\n")
	fmt.Println(article)
	fmt.Println(article.Section)
	t.Logf(text.Value())

}
