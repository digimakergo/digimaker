package main

import (
	"context"
	"database/sql"
	orm "dm/model/entity"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "test:test@tcp(185.35.187.91)/dev_emf")

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	if db.Ping() != nil {
		fmt.Printf(err.Error())
		return
	}

	//content := new(base.Content)

	count, err := entity.Locations().Count(context.Background(), db)

	fmt.Printf("Count: %d \n", count)
}
