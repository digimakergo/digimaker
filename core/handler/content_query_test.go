//Author xc, Created on 2019-04-28 18:11
//{COPYRIGHTS}

package handler

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/util"
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

func TestUpdate(t *testing.T) {
	rmdb := db.DBHanlder()

	var article entity.Article
	rmdb.GetByFields("article", "dm_article", db.Cond("content_id", 1), nil, nil, &article, false)
	//Update remote id of the article
	fmt.Println(article)
	uid := util.GenerateUID()
	println(uid)
	article.UID = uid
	err := article.Store()
	fmt.Println(err)

	/*
		id, error := rmdb.Insert("dm_article", map[string]interface{}{"modified": 231213})
		if error != nil {
			fmt.Println(id, error.Error())
		}
	*/

	err = rmdb.Update(article.TableName(), map[string]interface{}{"body": "test" + time.Now().String()}, db.Cond("id", 1))
	assert.Nil(t, err)
	var article2 entity.Article
	rmdb.GetByFields("article", "dm_article", db.Cond("content_id", 1), nil, nil, &article2, false)

	//assert.Equal(t, article2.RemoteID, uid)

	// //insert
	// article3 := new(entity.Article)
	// article3.Modified = 5555555
	// err = article3.Store()

	articles, _, err := Querier().List("article", db.Cond("1", "1"), nil, nil, false)
	fmt.Println(articles)

	fmt.Println("New article")
	// article4, err := Querier().Fetch("article", db.Cond("location.id", 43))
	// fmt.Println(article4.(*entity.Article).Editors)

}
