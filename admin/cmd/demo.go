package main

import (
	"dm/admin/entity"
	"dm/core"
	"dm/core/contenttype"
	"dm/core/db"
	"dm/core/fieldtype"
	"dm/core/handler"
	_ "dm/core/handler/handlers"
	"dm/core/util/debug"
	"dm/rest"
	"dm/sitekit/niceurl"
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
		success := core.Bootstrap(path)
		if !success {
			fmt.Println("Failed to start. Exiting.")
			os.Exit(1)
		}
	} else {
		fmt.Println("Need a path parameter. Exiting.")
		os.Exit(1)
	}
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
		folders, _ := handler.Querier().List("folder", db.Cond("parent_id", 0))
		debug.Debug(ctx, "Got list of folder", "system")

		//Get current Folder
		current, _ := handler.Querier().FetchByID(id)

		variables := map[string]interface{}{}
		variables["current"] = current
		variables["current_def"] = current.Definition()
		variables["folders"] = folders

		switch current.ContentType() {
		case "folder":
			switch current.Value("folder_type").(fieldtype.TextField).Data {
			//image folder
			case "image":
				debug.Debug(ctx, "Trying to get images", "system")
				images := &[]entity.Image{}
				fmt.Println(current.GetLocation().ID)
				handler := db.DBHanlder()
				handler.GetEntity("dm_image", db.Cond("parent_id", current.GetLocation().ID), images)
				variables["list"] = images
			//user folder
			case "user":
				users, err := handler.Querier().List("user", db.Cond("parent_id", id))
				fmt.Println(err)
				variables["list"] = users
			}
		}

		if _, ok := variables["list"]; !ok {
			allowedTypes := []string{"article"}
			if current.ContentType() != "folder" {
				allowedTypes = current.Definition().AllowedTypes
			}
			list := []contenttype.ContentTyper{}
			for _, allowedType := range allowedTypes {
				currentList, _ := handler.Querier().Children(current, allowedType, 7, r.Context())
				list = append(list, currentList...)
			}
			variables["list"] = list
		}

		rootID := current.GetLocation().Path()[0]
		rootContent, err := handler.Querier().FetchByID(rootID)
		tree, err := handler.Querier().SubTree(rootContent, 4, "folder,usergroup", 7, r.Context())
		fmt.Println(tree)
		variables["tree"] = tree

		//end Logic timing
		debug.EndTiming(r.Context(), "logic", "logic")

		//template timing
		debug.StartTiming(r.Context(), "template", "all")

		folderList, _ := handler.Querier().List("folder", db.Cond("parent_id", id))
		variables["folder_list"] = folderList

		variables["format_time"] = func(unix int) string {
			return time.Unix(int64(unix), 0).Format("02.01.2006 15:04")
		}

		variables["generate_url"] = niceurl.GenerateUrl

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
		content, result, error := handler.Create(contentType, params, parentID)
		fmt.Println(content, result, error)
		if content == nil {
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

func Export(w http.ResponseWriter, r *http.Request, vars map[string]string) {
	debug.StartTiming(r.Context(), "query", "kernel")
	id, _ := strconv.Atoi(vars["id"])
	mh := handler.ExportHandler{}
	content, _ := handler.Querier().FetchByID(id)
	parent, _ := handler.Querier().FetchByID(content.Value("parent_id").(int))
	data, _ := mh.Export(content, parent)
	w.Write([]byte("<html><body style=\"font-family: monospace\">" + data + "</body></html>"))
	debug.EndTiming(r.Context(), "query", "kernel")
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

	r.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	r.HandleFunc("/content/export/{id}", func(w http.ResponseWriter, r *http.Request) {
		DMHandle(w, r, func(w http.ResponseWriter, r *http.Request, vars map[string]string) {
			Export(w, r, vars)
		})
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

	//rest api
	restRouter := r.PathPrefix("/api").Subrouter()
	rest.Route(restRouter)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../web"))))

	r.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("w write."))
		fmt.Fprintf(w, "hello world")
	})

	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}
