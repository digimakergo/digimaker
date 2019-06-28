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

	"github.com/gorilla/mux"
	"gopkg.in/flosch/pongo2.v2"
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

	sitelist := []string{"demosite", "test"} //todo: load dynamically either from api or config file.
	defaultSite := "demosite"

	r := mux.NewRouter()
	//route subsites
	for _, site := range sitelist {
		s := r.PathPrefix("/" + site + "/").Subrouter()
		RouteContent(s, site, site)
	}
	//route default site
	RouteContent(r, defaultSite, "")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}

func RouteContent(r *mux.Router, site string, prefix string) {
	r.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		viewContent(w, r, id, site, prefix) //todo: use default site.
	})

	r.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		viewContent(w, r, id, site, prefix)
	})
}

func viewContent(w http.ResponseWriter, r *http.Request, id int, templateFolder string, prefix string) {
	// Execute the template per HTTP request
	pongo2.DefaultSet.Debug = true
	pongo2.DefaultSet.SetBaseDirectory("../templates/" + templateFolder)
	tplExample := pongo2.Must(pongo2.FromCache("../default/viewcontent.html"))
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	root, err := querier.FetchByID(55)
	fmt.Println(err)
	fmt.Println(content)
	fmt.Println(root)
	err = tplExample.ExecuteWriter(pongo2.Context{"content": content, "root": root, "viewmode": "full", "site": "demosite", "prefix": prefix}, w)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
