package handler

import (
	"dm/contenttype"
	"dm/fieldtype"
	"dm/query"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateRelation(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	handler := RelationHandler{}
	currentArticle, _ := Querier().Fetch("article", query.Cond("location.id", 6))

	article, _ := Querier().Fetch("article", query.Cond("location.id", 42))

	priority, _ := strconv.Atoi(time.Now().Format("0102150405"))
	err := handler.AddContent(currentArticle, article, "related_articles", priority, time.Now().String())
	fmt.Println(err)
	assert.Equal(t, nil, err)
}
