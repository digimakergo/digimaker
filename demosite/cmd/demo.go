package main

import (
	"dm/dm"
	"dm/dm/util"
	"dm/niceurl"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
	//read from config file, route content.
	config := util.GetConfigSectionAll("sites", "site").(map[string]interface{})

	//init
	err := sitekit.Init(r, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//example of route
	sitekit.SiteRouterHandle(r, "test1", "/user/list", func(w http.ResponseWriter, re *http.Request) {
		sitekit.OutputTemplate(w, re, "test1", "user/list")
	})

	//loop sites and route
	for identifier, _ := range config {
		var handleContentView = func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])
			prefix := ""
			if path, ok := vars["path"]; ok {
				prefix = path
			}
			sitekit.OutputContent(w, r, id, identifier, prefix)
		}

		//site route and get sub route
		err := sitekit.SiteRouter(r, identifier, func(s *mux.Router) {
			s.HandleFunc("/content/view/{id}", handleContentView)
			s.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(handleContentView)
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	http.Handle("/", r)
	fmt.Println("success!")
	http.ListenAndServe(":8089", nil)

}
