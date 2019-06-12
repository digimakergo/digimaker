//Author xc, Created on 2019-04-28 18:11
//{COPYRIGHTS}

package handler

import (
	"context"
	"dm/dm/contenttype"
	"dm/dm/contenttype/entity"
	"dm/dm/fieldtype"
	"dm/dm/permission"
	"dm/dm/util/debug"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchByContent(t *testing.T) {
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()
	permission.LoadPolicies()

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
		fmt.Println(strconv.Itoa(article.ID) + ":" + article.GetName())
	}
	fmt.Println(debug.GetDebugger(context).List)
}

func TestSubTree(t *testing.T) {
	querier := Querier()
	rootContent, _ := querier.FetchByID(1)
	context := debug.Init(context.Background())
	fmt.Println("TREEEEEEEEEEEEE")
	treenode, _ := querier.SubTree(rootContent, 3, "folder,article", 7, context)

	fmt.Println(treenode.Content.GetName())
	children := treenode.Children
	for _, child := range children {
		fmt.Println(child.Content.GetName())
		for _, child2 := range child.Children {
			fmt.Println("- " + child2.Content.GetName())
			for _, child3 := range child2.Children {
				fmt.Println("-- " + child3.Content.GetName())
			}
		}
	}

}
