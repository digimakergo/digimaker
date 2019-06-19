package main

import (
	"dm/dm"
	"dm/dm/handler"
	"dm/niceurl"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "dm/demosite/entity"
	_ "dm/demosite/pongo2"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func BootStrap() {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		path := os.Args[1]
		success := dm.Bootstrap(path)
		if !success {
			fmt.Println("Failed to start. Exiting.")
			os.Exit(1)
		}
	} else {
		fmt.Println("Need a path parameter. Exiting.")
		os.Exit(1)
	}
}

func main() {
	BootStrap()

	r := mux.NewRouter()

	r.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		viewContent(w, r, id)
	})

	r.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		viewContent(w, r, id)
	})

	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}

func viewContent(w http.ResponseWriter, r *http.Request, id int) {
	// Execute the template per HTTP request
	pongo2.DefaultSet.Debug = true
	tplExample := pongo2.Must(pongo2.FromCache("../templates/default/viewcontent.html"))
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	root, err := querier.FetchByID(55)
	fmt.Println(err)
	err = tplExample.ExecuteWriter(pongo2.Context{"content": content, "root": root, "viewmode": "full", "site": "demosite"}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
