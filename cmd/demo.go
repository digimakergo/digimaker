package main

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/fieldtype"
	"dm/handler"
	_ "dm/handler/handlers"
	"dm/query"
	"dm/util"
	"dm/util/debug"
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
	//start request timing
	ctx := r.Context()
	parser, err := template.ParseFiles("../web/template/view.html")

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
		debug.Debug(ctx, "Got list of folder", "system")

		//Get current Folder
		currentFolder, _ := handler.Querier().Fetch("folder", query.Cond("location.id", id))

		var variables map[string]interface{}
		c := currentFolder.(*entity.Folder)
		if c.ID != 0 {
			//Folder. Get list of article

			debug.Debug(ctx, "It is a folder. Trying to get folders and articles under.", "system")
			variables = map[string]interface{}{"current": currentFolder,
				"current_def": contenttype.GetContentDefinition("folder"),
				"folders":     folders}

			folderType := currentFolder.Value("folder_type").(fieldtype.TextField)
			if folderType.Data == "image" {
				debug.Debug(ctx, "Trying to get images", "system")
				images := &[]entity.Image{}
				fmt.Println(currentFolder.GetLocation().ID)
				handler := db.DBHanlder()
				handler.GetEnity("dm_image", query.Cond("parent_id", currentFolder.GetLocation().ID), images)
				variables["list"] = images
				fmt.Println(images)
			} else {
				articles, _ := handler.Querier().List("article", query.Cond("parent_id", id))
				variables["list"] = articles
			}

		} else {
			debug.Debug(ctx, "Not a folder. Trying to get folders under.", "system")
			currentArticle, _ := handler.Querier().Fetch("article", query.Cond("location.id", id))

			variables = map[string]interface{}{"current": currentArticle,
				"list":        nil,
				"current_def": contenttype.GetContentDefinition("article"),
				"folders":     folders}
		}

		//end Logic timing
		debug.EndTiming(r.Context(), "logic", "logic")

		//template timing
		debug.StartTiming(r.Context(), "template", "all")

		folderList, _ := handler.Querier().List("folder", query.Cond("parent_id", id))
		variables["folder_list"] = folderList

		variables["format_time"] = func(unix int) string {
			return time.Unix(int64(unix), 0).Format("02.01.2006 15:04")
		}

		err = tpl.Execute(w, variables)
		if err != nil {
			debug.Error(r.Context(), err.Error(), "template")
		}

		//template timing end
		debug.EndTiming(r.Context(), "template", "all")
	}

}

func New(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	// handler := handler.ContentHandler{}

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
		handler := handler.ContentHandler{Context: r.Context()}
		success, result, error := handler.Create(contentType, params, parentID)
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
	debug.StartTiming(r.Context(), "template", "kernel")
	contentType := vars["type"]
	def := contenttype.GetContentDefinition(contentType)
	variables["definition"] = def
	variables["contenttype"] = contentType
	tpl := template.Must(template.ParseFiles("../web/template/new.html"))
	//variables := map[string]interface{}{}
	tpl.Execute(w, variables)
	debug.EndTiming(r.Context(), "template", "kernel")
}

func Edit(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	// handler := handler.ContentHandler{}

	variables := map[string]interface{}{}
	id, _ := strconv.Atoi(vars["id"])
	variables["id"] = id
	content, _ := handler.Querier().FetchByID(id)
	variables["content"] = content
	variables["posted"] = false
	if r.Method == "POST" {
		variables["posted"] = true
		id, _ := strconv.Atoi(vars["id"])
		params := map[string]interface{}{}
		r.ParseForm()
		for key, value := range r.PostForm {
			if key != "id" && key != "type" {
				params[key] = value[0]
			}
		}
		// contentType := r.PostFormValue("type")
		cHandler := handler.ContentHandler{Context: r.Context()}
		success, result, error := cHandler.UpdateByID(id, params)
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
		content, _ = handler.Querier().FetchByID(id)
		variables["content"] = content
	}
	debug.StartTiming(r.Context(), "template", "kernel")
	contentType := content.ContentType()
	def := contenttype.GetContentDefinition(contentType)
	variables["definition"] = def
	variables["contenttype"] = contentType
	tpl := template.Must(template.ParseFiles("../web/template/edit.html"))
	//variables := map[string]interface{}{}
	tpl.Execute(w, variables)
	debug.EndTiming(r.Context(), "template", "kernel")
}

func Delete(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	handler := handler.ContentHandler{Context: r.Context()}
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
	debug.StartTiming(r.Context(), "template", "kernel")
	tpl := template.Must(template.ParseFiles("../web/template/console/list.html"))
	variables := map[string]interface{}{}
	variables["definition"] = contenttype.GetDefinition()
	tpl.Execute(w, variables)
	debug.EndTiming(r.Context(), "template", "kernel")
}

func DMHandle(w http.ResponseWriter, r *http.Request, functionHandler func(http.ResponseWriter, *http.Request, map[string]string)) {
	ctx := debug.Init(r.Context())
	r = r.WithContext(ctx)

	debug.StartTiming(r.Context(), "request", "kernel")
	debug.Debug(ctx, "Request started", "request")

	vars := mux.Vars(r)
	functionHandler(w, r, vars)

	debug.Debug(ctx, "Request ended", "request")
	debug.EndTiming(ctx, "request", "kernel")

	errorLog := ""
	errorCount := 0
	for _, item := range debug.GetDebugger(r.Context()).List {
		if item.Type == "error" {
			errorCount++
		}
		errorLog += "<div class=info-" + item.Type + "><span class=category>[" + item.Category + "]</span><span>" + item.Message + "</span></div>"
	}

	queryDuration, err := debug.GetDuration(ctx, "logic")
	queryStr := "'query': "
	if err == nil {
		queryStr += strconv.Itoa(queryDuration)
	} else {
		queryStr += "null"
	}
	templateDuration, err := debug.GetDuration(ctx, "template")
	templateStr := "'template': "
	if err == nil {
		templateStr += strconv.Itoa(templateDuration)
	} else {
		templateStr += "null"
	}
	total, _ := debug.GetDuration(ctx, "request")

	w.Write([]byte("<script>var dmtime={ 'total': " + strconv.Itoa(total) +
		", " + queryStr +
		"," + templateStr +
		", errors:" + strconv.Itoa(errorCount) +
		"};" +
		"var errorLog=\"" + errorLog + "\";" +
		"</script>" +
		"<link href='/static/css/debug.css' rel='stylesheet'>" +
		"<script src='https://ajax.googleapis.com/ajax/libs/jquery/3.4.0/jquery.min.js'></script>" +
		"<script src='/static/javascript/dmdebug.js'></script>"))
}

func main() {

	BootStrap()
	r := mux.NewRouter()
	r.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
		DMHandle(w, r, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
			Display(w, r, vars)
		})
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		DMHandle(w, r, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
			Display(w, r, map[string]string{"id": "1"})
		})
	})

	r.HandleFunc("/content/new/{type}/{id}", func(w http.ResponseWriter, r *http.Request) {
		DMHandle(w, r, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
			New(w, r, vars)
		})
	})

	r.HandleFunc("/content/edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		DMHandle(w, r, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
			Edit(w, r, vars)
		})
	})

	r.HandleFunc("/content/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		DMHandle(w, r, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
			Delete(w, r, vars)
		})
	})

	r.HandleFunc("/content/publish", func(w http.ResponseWriter, r *http.Request) {
		ctx := debug.Init(r.Context())
		r = r.WithContext(ctx)
		Publish(w, r)
	})

	r.HandleFunc("/console/list", func(w http.ResponseWriter, r *http.Request) {
		DMHandle(w, r, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
			ModelList(w, r)
		})
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
