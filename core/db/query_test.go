package db

import (
	"context"
	"fmt"
	"testing"
	// "github.com/digimakergo/digimaker/core/db/mysql"
)

// func TestMain(m *testing.M) {
// 	//contenttype.LoadDefinition()
// 	// m.Run()
// }

func ExampleBindContent() {
	content := []struct {
		Title string `boil:"title"`
	}{}
	//In real case for list, it should be list = contenttype.NewList("folder")
	//In real case for content it should be content = contenttype.New("folder")
	BindContent(context.Background(), &content, "folder", TrueCond())
}

func ExampleBindEntity() {
	//Fetch all orders into a struct list
	orderlist := []struct {
		ID   string `boil:"id"`
		Name string `boil:"name"`
	}{}
	BindEntity(context.Background(), &orderlist, "order", Cond("author", 5).Limit(0, 10))

	//Fetch all orders to a map list
	datalist := DatamapList{}
	count, err := BindEntity(context.Background(), &datalist, "order", TrueCond())

	//Fetch one order to a map
	order := Datamap{}
	_, err = BindEntity(context.Background(), &order, "order", Cond("id", 2))

}

func ExampleCount() {
	//Count all orders
	count, err := Count("order", TrueCond())
}

func TestContent(t *testing.T) {
	fmt.Println("test contest")
	one := SingleQuery{}
	one.Condition = Cond("cl.id", 1)
	one.Table = "dm_article"
	// one.Alias = "c1"

	query := Query{}
	query.Queries = []SingleQuery{one}
	query.SortArr = []string{"modified desc"}
	entity := "test"
	BindContentWithQuery(context.Background(), entity, "article", query, ContentOption{})

	//BindContent(context.Background(), entity, "article", db.Cond("author", 1)
	// BindContent(context.Background(), entity, "article", db.Cond("author", 1).Sortby("modified desc").Limit(0, 10))
	//
	// BindContent(context.Background(), entity, "article, image", db.Cond("author", 1).Cond("article.image ==", "image.id").Sortby("modified desc").Limit(0, 10))
	//
	// BindContent(context.Background(), entity, db.Query("article", "article1" db.Cond("test", 1)).
	// 				Join("article", "article2", db.Cond("test", 1), db.Cond("article1.image ==", "image.id")).
	// 				LeftJoin("")
	// 				Sortby("").
	// 				Limit(""))
}
