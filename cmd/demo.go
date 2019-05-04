package main

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/fieldtype"
	"dm/handler"
	"dm/query"
	"dm/util"
	"dm/util/debug"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

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
	//start request timing
	ctx := r.Context()
	debug.Debug(ctx, "Request started", "")
	debug.StartTiming(r.Context(), "request", "kernel")
	parser, err := template.ParseFiles("../web/template/view.html")
	queryDuration := 0
	templateDuration := 0
	if err != nil {
		debug.Error(r.Context(), err.Error(), "template")
	} else {
		tpl := template.Must(parser, err)

		//logic timing
		debug.StartTiming(r.Context(), "logic", "logic")
		rmdb := db.DBHanlder()
		article := entity.Article{}
		id, _ := strconv.Atoi(vars["id"])

		err := rmdb.GetByID("article", "dm_article", id, &article)

		if err != nil {
			fmt.Println(err)
		}

		//List of folder
		folders, _ := handler.Querier().List("folder", query.Cond("parent_id", 0))
		debug.Debug(ctx, "Got list of folder", "")

		//Get current Folder
		currentFolder, _ := handler.Querier().Fetch("folder", query.Cond("location.id", id))

		var variables map[string]interface{}
		c := currentFolder.(*entity.Folder)
		if c.ID != 0 {
			//Folder. Get list of article

			debug.Debug(ctx, "The folder is not empty. Trying to get folders and articles under.", "")
			variables = map[string]interface{}{"current": currentFolder,
				"current_def": contenttype.GetContentDefinition("folder"),
				"folders":     folders}

			folderType := currentFolder.Value("folder_type").(fieldtype.TextField)
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
			debug.Debug(ctx, "The folder is empty. Trying to get folders under.", "")
			currentArticle, _ := handler.Querier().Fetch("article", query.Cond("location.id", id))

			variables = map[string]interface{}{"current": currentArticle,
				"list":        nil,
				"current_def": contenttype.GetContentDefinition("article"),
				"folders":     folders}
		}

		//end Logic timing
		debug.EndTiming(r.Context(), "logic", "logic")
		queryDuration = debug.GetDebugger(r.Context()).Timers["logic"].Duration

		//template timing
		debug.StartTiming(r.Context(), "template", "all")

		folderList, _ := handler.Querier().List("folder", query.Cond("parent_id", id))
		variables["folder_list"] = folderList
		err = tpl.Execute(w, variables)
		if err != nil {
			debug.Error(r.Context(), err.Error(), "template")
		}

		//template timing end
		debug.EndTiming(r.Context(), "template", "all")
		templateDuration = debug.GetDebugger(r.Context()).Timers["template"].Duration
	}

	//system timing end
	debug.EndTiming(r.Context(), "request", "kernel")

	errorLog := ""
	for _, item := range debug.GetDebugger(r.Context()).List {
		errorLog += "<div class=info-" + item.Type + "><span class=category>[" + item.Category + "]</span><span>" + item.Type + "</span><span>" + item.Message + "</span></div>"
	}

	w.Write([]byte("<script>var dmtime={ 'total': " + strconv.Itoa(debug.GetDebugger(r.Context()).Timers["request"].Duration) +
		", 'query':" +
		strconv.Itoa(queryDuration) + ", 'template':" +
		strconv.Itoa(templateDuration) + "};" +
		"var errorLog='" + errorLog + "';" +
		"</script>" +
		"<link href='/static/css/debug.css' rel='stylesheet'>" +
		"<script src='https://ajax.googleapis.com/ajax/libs/jquery/3.4.0/jquery.min.js'></script>" +
		"<script src='/static/javascript/dmdebug.js'></script>"))
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

func Test(w http.ResponseWriter, r *http.Request) {
	debug.Debug(r.Context(), "This is wrong..", "")
	debug.Debug(r.Context(), "This is wrong2", "")
	debugger := debug.GetDebugger(r.Context())
	w.Write([]byte(debugger.List[0].Message))
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
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		vars := mux.Vars(r)
		Display(w, r, vars)
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		Display(w, r, map[string]string{"id": "1"})
	})
	// http.HandleFunc("/content/view/", func(w http.ResponseWriter, r *http.Request) {
	// 	Display(w, r)
	// })

	r.HandleFunc("/content/new/{type}/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		New(w, r)
	})

	r.HandleFunc("/content/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		Delete(w, r)
	})

	r.HandleFunc("/content/publish", func(w http.ResponseWriter, r *http.Request) {
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		Publish(w, r)
	})

	r.HandleFunc("/console/list", func(w http.ResponseWriter, r *http.Request) {
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		ModelList(w, r)
	})

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		Test(w, r)
	})

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "User-agent: * \nDisallow /")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../web"))))

	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}
