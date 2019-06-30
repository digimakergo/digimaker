package sitekit

import (
	"dm/dm/handler"
	"dm/dm/util"
	"dm/niceurl"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/flosch/pongo2.v2"
)

//Route a site content with configuration.
//Config includes:
// route:
// template_folder
// root:
// default:
func RouteContent(r *mux.Router, config map[interface{}]interface{}) error {
	routesConfig := config["routes"].([]interface{})

	if _, ok := config["template_folder"]; !ok {
		return errors.New("Need template_folder setting.")
	}
	tempalteFolder := util.InterfaceToStringArray(config["template_folder"].([]interface{}))

	if _, ok := config["root"]; !ok {
		return errors.New("Need root setting.")
	}
	root := config["root"].(int)

	if _, ok := config["default"]; !ok {
		return errors.New("Need default setting.")
	}
	defaultContent := config["default"].(int)
	//go through all routes.
	for _, routeConfig := range routesConfig {
		var hosts []string
		route := routeConfig.(map[interface{}]interface{})
		if hostStr, ok := route["host"]; ok {
			hosts = util.Split(hostStr.(string))
		} else {
			return errors.New("Need host setting.")
		}

		var paths []string
		if value, ok := route["path"]; ok {
			if value == "" {
				paths = append(paths, "") //root
			} else {
				paths = util.Split(value.(string))
			}
		} else {
			paths = append(paths, "") //root
		}

		//hosts and paths should AND match, but host in hosts(also path in paths) is OR.
		for _, host := range hosts {
			for _, path := range paths {
				var s *mux.Router
				//use subrouter which is better for performance
				if path != "" {
					s = r.Host(host).PathPrefix("/" + path + "/").Subrouter()
				} else {
					s = r.Host(host).Subrouter()
				}
				routeContent(s, tempalteFolder[0], path, root, defaultContent)
			}
		}
	}

	return nil
}

func routeContent(r *mux.Router, templateFolder string, prefix string, root int, defaultContent int) {
	//todo: add route debug.
	r.HandleFunc("/content/view/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		viewContent(w, r, id, templateFolder, prefix)
	})

	r.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		viewContent(w, r, id, templateFolder, prefix)
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
