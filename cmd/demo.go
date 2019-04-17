package main

import (
	"dm/db"
	"dm/model"
	"dm/model/entity"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func BootStrap() {
	if os.Args[1] != "" {
		model.DMPath = os.Args[1]
	}
	model.LoadDefinition()
}

//This is a initial try which use template to do basic feature.

func Display(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("../web/template/view.html"))
	rmdb := db.DBHanlder()
	article := entity.Article{}
	idStr := r.FormValue("id")
	id := 1
	if idStr != "" {
		id, _ = strconv.Atoi(idStr)
	}

	err := rmdb.GetByID("article", id, &article)
	if err != nil {
		fmt.Println(err)
	}

	tpl.Execute(w, article)
}

func Draft(w http.ResponseWriter, r *http.Request) {
	// handler := handler.ContentHandler{}
}

func Publish(w http.ResponseWriter, r *http.Request) {

}

func main() {
	BootStrap()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Display(w, r)
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
