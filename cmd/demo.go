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
	"github.com/gorilla/mux"
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

func Display(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	tpl := template.Must(template.ParseFiles("../web/template/view.html"))
	rmdb := db.DBHanlder()
	article := entity.Article{}
	id, _ := strconv.Atoi(vars["id"])

	err := rmdb.GetByID("article", id, &article)

	if err != nil {
		fmt.Println(err)
	}

	//List of folder
	var folders []entity.Folder
	handler.Query.List("folder", query.Cond("parent_id", 0), &folders)

	//Get current Folder
	var currentFolder entity.Folder
	handler.Query.List("folder", query.Cond("dm_location.id", id), &currentFolder)

	var variables map[string]interface{}
	if currentFolder.CID != 0 {
		//Get list of article
		var articles []entity.Article
		handler.Query.List("article", query.Cond("parent_id", id), &articles)

		variables = map[string]interface{}{"current": currentFolder,
			"list":    articles,
			"folders": folders}
	} else {
		var currentArticle entity.Article
		handler.Query.List("article", query.Cond("dm_location.id", id), &currentArticle)
		variables = map[string]interface{}{"current": currentArticle,
			"list":    nil,
			"folders": folders}
	}

	var folderList []entity.Folder
	handler.Query.List("folder", query.Cond("parent_id", id), &folderList)
	variables["folder_list"] = folderList
	tpl.Execute(w, variables)
}

func Draft(w http.ResponseWriter, r *http.Request) {
	// handler := handler.ContentHandler{}
}

func Publish(w http.ResponseWriter, r *http.Request) {

}

func main() {

	BootStrap()
	r := mux.NewRouter()
	r.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Display(w, r, vars)
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Display(w, r, map[string]string{"id": "1"})
	})
	// http.HandleFunc("/content/view/", func(w http.ResponseWriter, r *http.Request) {
	// 	Display(w, r)
	// })

	r.HandleFunc("/content/draft", func(w http.ResponseWriter, r *http.Request) {
		Draft(w, r)
	})

	r.HandleFunc("/content/publish", func(w http.ResponseWriter, r *http.Request) {
		Publish(w, r)
	})

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User-agent: * \nDisallow /")
	})
	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}
