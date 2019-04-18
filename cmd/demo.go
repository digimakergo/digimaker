package main

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/fieldtype"
	"dm/handler"
	"dm/query"
	"dm/util"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func BootStrap() {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		path := os.Args[1]
		util.SetConfigPath(path + "/configs")
	}
	contenttype.LoadDefinition()
	fieldtype.LoadDefinition()
}

//This is a initial try which use template to do basic feature.

func Display(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("../web/template/view.html"))
	rmdb := db.DBHanlder()
	article := entity.Article{}
	idStr := r.FormValue("id")
	id := 2
	if idStr != "" {
		id, _ = strconv.Atoi(idStr)
	}

	err := rmdb.GetByID("article", id, &article)

	if err != nil {
		fmt.Println(err)
	}

	var articles []entity.Article
	handler.Query.List("article", query.Cond("1", "1"), &articles)

	var folders []entity.Folder
	handler.Query.List("folder", query.Cond("is_invisible", 0), &folders)

	tpl.Execute(w, map[string]interface{}{"article": article, "articles": articles, "folders": folders})
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
