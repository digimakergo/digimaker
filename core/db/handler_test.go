//Author xc, Created on 2019-04-27 20:33
//{COPYRIGHTS}
package db

import (
	"context"
	"fmt"
)

func ExampleInsert() {
	values := map[string]interface{}{
		"name":        "Test",
		"description": "Test1"}

	//Insert without transaction.
	result, _ := Insert(context.Background(), "demo_order", values)

	fmt.Println(result > 0)
	//output: true
}

func ExampleInsert_withTransaction() {
	values := map[string]interface{}{
		"name":        "Test",
		"description": "Test2"}

	database, _ := DB()
	tx, _ := database.Begin()

	//Insert with transation
	result, err := Insert(context.Background(), "demo_order", values)
	tx.Commit()

	if err != nil {
		tx.Rollback()
	}
	fmt.Println(result > 0)
	//output: true
}

func ExampleUpdate() {
	//Update only name
	values := map[string]interface{}{
		"name": "Order updated"}
	Update(context.Background(), "demo_order", values, Cond("id", 1))

	//fetch updated
	order := struct {
		Name string `boil:"name"`
	}{}
	BindEntity(context.Background(), &order, "demo_order", Cond("id", 1))
	fmt.Println(order.Name)
	//output: Order updated
}

func ExampleDelete() {
	//Update only name
	values := map[string]interface{}{
		"name": "Order"}
	id, _ := Insert(context.Background(), "demo_order", values)

	Delete(context.Background(), "demo_order", Cond("id", id))

	count, _ := Count("demo_order", Cond("id", id))
	fmt.Println(count)
	//output: 0
}
