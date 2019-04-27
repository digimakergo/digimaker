package handler

import (
	"dm/contenttype"
	"dm/fieldtype"
	"testing"
)

func TestCreateRelation(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	// handler := RelationHandler{}
	// currentArticle, _ := Querier().Fetch("article", query.Cond("location.id", 6))
	//
	// article, _ := Querier().Fetch("article", query.Cond("location.id", 42))
	//
	// priority, _ := strconv.Atoi(time.Now().Format("0102150405"))
	// handler.Add(currentArticle, article, "related_articles", priority, time.Now().String())
}
