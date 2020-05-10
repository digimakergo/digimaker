//Author xc, Created on 2019-04-28 18:11
//{COPYRIGHTS}

package handler

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/test/entity"

	"github.com/stretchr/testify/assert"
)

func TestFetchByContent(t *testing.T) {

	querier := Querier()
	content, err := querier.FetchByContentID("folder", 1)
	assert.Equal(t, nil, err)
	fmt.Println(content.ToMap()["content_id"])
	// assert.Equal(t, 2, content.ToMap()["content_id"].(int))

	content, err = querier.FetchByID(1)
	assert.Equal(t, 1, content.Value("id"))
}

func TestSubList(t *testing.T) {
	querier := Querier()
	rootContent, _ := querier.FetchByID(1)
	fmt.Println("=========")
	list, _, _ := querier.SubList(rootContent, "article", 2, 7, db.EmptyCond(), []int{}, []string{}, false, ctx)

	//
	for _, item := range list {
		article := item.(*entity.Article)
		fmt.Println(strconv.Itoa(article.ID) + ":" + article.GetName())
	}
}

func TestSubTree(t *testing.T) {
	querier := Querier()
	rootContent, _ := querier.FetchByID(1)
	fmt.Println("TREEEEEEEEEEEEE")
	treenode, _ := querier.SubTree(rootContent, 3, "folder,article", 7, []string{}, ctx)

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

func TestQuery(t *testing.T) {

	dbHandler := db.DBHanlder()
	var article entity.Article
	dbHandler.GetByID("article", "dm_article", 2, &article)

	assert.NotNil(t, article)

	folders, _, _ := Querier().List("folder", db.Cond("1", "1"), []int{}, []string{}, false)
	fmt.Println("HELLO")
	fmt.Println(folders)
}
