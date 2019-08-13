package main

import (
	"dm/core"
	"dm/core/handler"
	"fmt"
	"net/http"
	"os"

	"github.com/flosch/pongo2"
)

//Generate database schema based on contenttype.json
func GenerateSchema() {

}

func VerifySchema() {

}

//Generate entities based on contenttype.json
func GenerateEnities() {

}

var tplExample = pongo2.Must(pongo2.FromFile("example.html"))

func examplePage(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	querier := handler.Querier()
	content, err := querier.FetchByID(1)
	fmt.Println(err)
	err = tplExample.ExecuteWriter(pongo2.Context{"content": content}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
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

	http.HandleFunc("/", examplePage)
	http.ListenAndServe(":8080", nil)
}
