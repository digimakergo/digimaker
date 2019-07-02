package main

import (
	"dm/dm"
	"dm/dm/util"
	"fmt"
	"net/http"
	"os"

	_ "dm/demosite/entity"
	"dm/sitekit"
	_ "dm/sitekit/pongo2"

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

	//Init sites
	config := util.GetConfigSectionAll("sites", "site").(map[string]interface{})
	err := sitekit.InitSites(r, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Handle custom module.
	//This should be before hanlding content, otherwise it can be routed by content first.
	sitekit.SiteRouterHandle(r, "test1", "/user/list", func(w http.ResponseWriter, re *http.Request) {
		sitekit.Output(w, re, "test1", "user/list", map[string]interface{}{}, map[string]interface{}{"test": "hello"})
	})

	//Handle content
	err = sitekit.HandleContent(r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.Handle("/", r)
	fmt.Println("success!")
	http.ListenAndServe(":8092", nil)

}
