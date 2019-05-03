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
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//go:generate go run gen_contenttypes/gen.go
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
	queryStart := time.Now()
	rmdb := db.DBHanlder()
	article := entity.Article{}
	id, _ := strconv.Atoi(vars["id"])

	err := rmdb.GetByID("article", "dm_article", id, &article)

	if err != nil {
		fmt.Println(err)
	}

	//List of folder
	folders, _ := handler.Querier().List("folder", query.Cond("parent_id", 0))

	//Get current Folder
	currentFolder, _ := handler.Querier().Fetch("folder", query.Cond("location.id", id))

	var variables map[string]interface{}
	c := currentFolder.(*entity.Folder)
	if c.ID != 0 {
		//Folder. Get list of article

		variables = map[string]interface{}{"current": currentFolder,
			"current_def": contenttype.GetContentDefinition("folder"),
			"folders":     folders}

		folderType := currentFolder.Value("folder_type").(fieldtype.TextField)
		fmt.Println(folderType)
		if folderType.Data == "image" {
			images := &[]entity.Image{}
			handler := db.DBHanlder()
			handler.GetEnity("dm_image", query.Cond("attached_location", currentFolder.GetLocation().ID), images)
			variables["list"] = images
			fmt.Println(images)
		} else {
			articles, _ := handler.Querier().List("article", query.Cond("parent_id", id))
			variables["list"] = articles
		}

	} else {
		currentArticle, _ := handler.Querier().Fetch("article", query.Cond("location.id", id))

		variables = map[string]interface{}{"current": currentArticle,
			"list":        nil,
			"current_def": contenttype.GetContentDefinition("article"),
			"folders":     folders}
	}

	queryEnd := time.Now()
	queryDuration := queryEnd.Sub(queryStart)

	templateStart := time.Now()
	folderList, _ := handler.Querier().List("folder", query.Cond("parent_id", id))
	variables["folder_list"] = folderList
	tpl.Execute(w, variables)
	templateEnd := time.Now()
	period := templateEnd.Sub(templateStart).Nanoseconds()
	w.Write([]byte("<script>var dmtime={ 'total': " + strconv.Itoa(int(templateEnd.Sub(queryStart).Nanoseconds()/1000000)) +
		", 'query':" +
		strconv.Itoa(int(queryDuration/1000000)) + ", 'template':" +
		strconv.Itoa(int(period/1000000)) + "};</script>"))
}

func New(w http.ResponseWriter, r *http.Request) {
	// handler := handler.ContentHandler{}

	vars := mux.Vars(r)

	variables := map[string]interface{}{}
	variables["id"] = vars["id"]
	variables["type"] = vars["type"]
	variables["posted"] = false
	if r.Method == "POST" {
		variables["posted"] = true
		parentID, _ := strconv.Atoi(vars["id"])
		params := map[string]interface{}{}
		r.ParseForm()
		for key, value := range r.PostForm {
			if key != "id" && key != "type" {
				params[key] = value[0]
			}
		}
		contentType := r.PostFormValue("type")
		handler := handler.ContentHandler{}
		success, result, error := handler.Create(parentID, contentType, params)
		fmt.Println(success, result, error)
		if !success {
			variables["success"] = false
			if error != nil {
				variables["error"] = error.Error()
			}
			variables["validation"] = result
		} else {
			variables["success"] = true
		}

	}
	tpl := template.Must(template.ParseFiles("../web/template/new_" + vars["type"] + ".html"))
	//variables := map[string]interface{}{}
	tpl.Execute(w, variables)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	handler := handler.ContentHandler{}
	id, _ := strconv.Atoi(vars["id"])
	err := handler.DeleteByID(id, false)
	if err != nil {
		w.Write([]byte(("error:" + err.Error())))
	} else {
		w.Write([]byte("success!"))
	}
}

func Publish(w http.ResponseWriter, r *http.Request) {

}

func ModelList(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("../web/template/console/list.html"))
	variables := map[string]interface{}{}
	variables["definition"] = contenttype.GetDefinition()
	tpl.Execute(w, variables)
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

	r.HandleFunc("/content/new/{type}/{id}", func(w http.ResponseWriter, r *http.Request) {
		New(w, r)
	})

	r.HandleFunc("/content/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		Delete(w, r)
	})

	r.HandleFunc("/content/publish", func(w http.ResponseWriter, r *http.Request) {
		Publish(w, r)
	})

	r.HandleFunc("/console/list", func(w http.ResponseWriter, r *http.Request) {
		ModelList(w, r)
	})

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User-agent: * \nDisallow /")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../web"))))

	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}
