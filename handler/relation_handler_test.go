package handler

import (
	"dm/contenttype"
	"dm/fieldtype"
	"dm/query"
	"testing"
)

func TestCreateRelation(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	handler := RelationHandler{}

	currentArticle, _ := Querier().Fetch("article", query.Cond("location.id", 6))

	article, _ := Querier().Fetch("article", query.Cond("location.id", 42))

	handler.AddTo(currentArticle, article, "related_articles", 0, "")
}
