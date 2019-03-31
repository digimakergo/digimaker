package main

import (
	"context"
	"database/sql"
	"dm/model/entity"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

//This is a initial try which use template to do basic feature.

func Display(w http.ResponseWriter) {
	tpl := template.Must(template.ParseFiles("../web/template/view.html"))
	db := GetDB()
	locations, err := entity.Locations(Where("parent_id != -1")).All(context.Background(), db)
	if err != nil {
		panic(err)
	}

	tpl.Execute(w, locations)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Display(w)
	})

	http.ListenAndServe(":8089", nil)
}

func GetDB() *sql.DB {
	db, err := sql.Open("mysql", "test:test@tcp(185.35.187.91)/dev_emf")

	if err != nil {
		fmt.Printf(err.Error())
		return nil
	}

	if db.Ping() != nil {
		fmt.Printf(err.Error())
		return nil
	}
	return db
}
