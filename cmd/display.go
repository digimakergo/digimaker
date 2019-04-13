package main

import (
	"context"
	"database/sql"
	"dm/handler"
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

func Draft(w http.ResponseWriter, r *http.Request) {
	handler := handler.ContentHandler{}
	handler.Create(title, 1)
}

func Publish(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Display(w)
	})

	http.HandleFunc("/content/draft", func(w http.ResponseWriter, r *http.Request) {
		Draft(w, r)
	})

	http.HandleFunc("/content/publish", func(w http.ResponseWriter, r *http.Request) {
		Publish(w, r)
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User-agent: * \nDisallow /")
	})
	http.ListenAndServe(":8089", nil)
}

func GetDB() *sql.DB {
	db, err := sql.Open("mysql", "test:test123@tcp(185.35.187.91)/dev_emf")

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
