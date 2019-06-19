package main

import (
	"dm/dm"
	"dm/dm/handler"
	"fmt"
	"net/http"
	"os"

	_ "dm/demosite/entity"

	"github.com/flosch/pongo2"
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
	http.HandleFunc("/", viewContent)

	http.ListenAndServe(":8089", nil)
}

func viewContent(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	pongo2.DefaultSet.Debug = true
	tplExample := pongo2.Must(pongo2.DefaultSet.FromCache("../web/default/viewcontent.html"))
	querier := handler.Querier()
	content, err := querier.FetchByID(1)
	fmt.Println(err)
	err = tplExample.ExecuteWriter(pongo2.Context{"content": content, "viewmode": "full", "site": "demosite"}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
