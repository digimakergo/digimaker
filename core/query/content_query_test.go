//Author xc, Created on 2019-04-28 18:11
//{COPYRIGHTS}

package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/digimakergo/digimaker/core/db"
	_ "github.com/digimakergo/digimaker/test"
	"github.com/digimakergo/digimaker/test/entity"

	"github.com/stretchr/testify/assert"
)

func ExampleFetch() {
	//Fetch folder whose name is "Content", only first will be returned
	content, err := Fetch(context.Background(), "folder", db.Cond("l.name", "Content"))
	if err == nil && content != nil {
		fmt.Println(content.(*entity.Folder).ID)
	}
	//output: 1
}

func ExampleFetchByCID() {
	//Fetch folder which is content id(cid) 1
	content, err := FetchByCID(context.Background(), "folder", 1)
	if err == nil && content != nil {
		fmt.Println(content.(*entity.Folder).ContentID)
	}
	//output: 1
}

func ExampleFetchByLID() {
	//Fetch folder which has location id 1
	content, err := FetchByLID(context.Background(), 1)
	if err == nil && content != nil {
		fmt.Println(content.(*entity.Folder).ContentID)
	}
	//output: 1
}

func ExampleSubList() {
	rootContent, _ := FetchByID(context.Background(), 1)
	//Fetch articles(level in 3) under root
	list, _, _ := SubList(context.Background(), 1, rootContent, "article", 3, db.EmptyCond())

	fmt.Println(len(list) > 0)
	//output: true
}

func ExampleSubList_withSortLimit() {
	rootContent, _ := FetchByID(context.Background(), 1)
	//Fetch articles(level in 3) under root

	list, _, _ := SubList(context.Background(), 1, rootContent, "article", 3, db.EmptyCond().Sortby("c.modified desc").Limit(0, 3))

	fmt.Println(len(list))
	//output: 3
}

func ExampleList() {
	//fetch list of folders where id larger than 2
	list, _, _ := List(context.Background(), "folder", db.Cond("l.id>", 2))
	fmt.Println(list[0].GetLocation().ID)
	//output: 3
}

func ExampleList_sortLimit() {
	//fetch list ordered by modified desc, limit 0,4
	list, _, _ := List(context.Background(), "folder", db.EmptyCond().Sortby("modified desc").Limit(0, 4))

	fmt.Println(len(list))
	//output: 4
}

func ExampleListWithUser() {
	//fetch list of folders where id larger than 2
	list, _, _ := ListWithUser(context.Background(), 1, "folder", db.Cond("l.id>", 2))
	fmt.Println(list[0].GetLocation().ID)
	//output: 3
}

func ExampleSubTree() {
	rootContent, _ := FetchByID(context.Background(), 1)
	//fetch trees of root for 3 levels
	treenode, _ := SubTree(context.Background(), 1, rootContent, 3, "folder,article", []string{})

	level := 0
	treenode.Iterate(func(node *TreeNode) {
		level++
	})
	fmt.Println(level > 3)
	//output: true
}

func ExampleChildren() {
	rootContent, _ := FetchByID(context.Background(), 1)

	//fetch all folders directly under root. Can filter if valid condition parameter is provided.
	children, _, _ := Children(context.Background(), 1, rootContent, "folder", db.EmptyCond())

	fmt.Println(len(children) > 0)
	//output: true
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
