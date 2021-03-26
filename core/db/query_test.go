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
	BindEntity(context.Background(), &orderlist, "demo_order", Cond("author", 5).Limit(0, 10))

	//Fetch all orders to a map list
	datalist := DatamapList{}
	BindEntity(context.Background(), &datalist, "demo_order", TrueCond())
	fmt.Print(len(datalist) > 0)

	//Fetch one order to a map list
	order := DatamapList{}
	BindEntity(context.Background(), &order, "demo_order", Cond("id>", 1))
	fmt.Print(len(order) > 0)
	//output: truetrue

}

func ExampleCount() {
	//Count all orders
	count, _ := Count("order", TrueCond())
	fmt.Println(count > 0)
	//output: true
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
