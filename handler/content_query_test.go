//Author xc, Created on 2019-04-28 18:11
//{COPYRIGHTS}

package handler

import (
	"context"
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/fieldtype"
	"dm/util/debug"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchByContent(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()

	querier := Querier()
	content, err := querier.FetchByContentID("article", 2)
	assert.Equal(t, nil, err)
	fmt.Println(content.ToMap()["content_id"])
	// assert.Equal(t, 2, content.ToMap()["content_id"].(int))

	content, err = querier.FetchByID(1)
	assert.Equal(t, 1, content.Value("id"))
}

func TestSubList(t *testing.T) {
	querier := Querier()
	rootContent, _ := querier.FetchByID(1)
	context := debug.Init(context.Background())
	fmt.Println("=========")
	list, _ := querier.SubList(rootContent, "article", 2, 7, context)

	//
	for _, item := range list {
		article := item.(*entity.Article)
		fmt.Println(strconv.Itoa(article.ID) + ":" + article.Name)
	}
	fmt.Println(debug.GetDebugger(context).List)
}
