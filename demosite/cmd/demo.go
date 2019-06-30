package main

import (
	"dm/dm"
	"fmt"
	"net/http"
	"os"

	_ "dm/demosite/entity"
	"dm/demosite/sitekit"

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
	//read from config file, route content.
	err := sitekit.RouteFromFile(r, "sites", "site")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("success!")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.Handle("/", r)
	http.ListenAndServe(":8089", nil)
}
