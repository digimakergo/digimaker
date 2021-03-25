//Author xc, Created on 2019-04-28 18:11
//{COPYRIGHTS}

package query

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/digimakergo/digimaker/core/db"
	_ "github.com/digimakergo/digimaker/test"
	"github.com/digimakergo/digimaker/test/entity"

	"github.com/stretchr/testify/assert"
)

func ExampleFetchByCID() {
	content, err := FetchByCID(context.Background(), "folder", 1)
	if err == nil && content != nil {
		fmt.Println(content.(*entity.Folder).ContentID)
	}
	//output: 1
}

func ExampleFetchByID() {
	content, err := FetchByID(context.Background(), 1)
	if err == nil && content != nil {
		fmt.Println(content.(*entity.Folder).ContentID)
	}
	//output: 1
}

func TestSubList(t *testing.T) {
	rootContent, _ := FetchByID(context.Background(), 1)
	fmt.Println("=========")
	list, _, _ := SubList(context.Background(), rootContent, "article", 2, 7, db.TrueCond(), false)

	//
	for _, item := range list {
		article := item.(*entity.Article)
		fmt.Println(strconv.Itoa(article.ID) + ":" + article.GetName())
	}
}

func TestSubTree(t *testing.T) {
	rootContent, _ := FetchByID(context.Background(), 1)
	treenode, _ := SubTree(context.Background(), rootContent, 3, "folder,article", 7, []string{})

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
	folders, _, _ := List(context.Background(), "folder", db.Cond("1", "1"))
	fmt.Println(folders)
}

func TestQueryImage(t *testing.T) {
	images := &[]entity.Image{}
	_, err := db.BindContent(context.Background(), images, "dm_image", db.Cond("1", 1))
	assert.Nil(t, err)
	assert.NotNil(t, images)
}
