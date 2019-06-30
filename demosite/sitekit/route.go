package sitekit

import (
	"dm/dm/handler"
	"dm/dm/util"
	"dm/dm/website"
	"dm/niceurl"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/flosch/pongo2.v2"
)

func RouteFromFile(r *mux.Router, section string, configFile string) error {
	config := util.GetConfigSectionAll(section, configFile).(map[string]interface{})
	for siteIdentifier, siteConfig := range config {
		fmt.Println(siteIdentifier)
		err := RouteContent(r, siteIdentifier, siteConfig.(map[string]interface{}))
		if err != nil {
			return err
		}
	}
	return nil
}

//Route a site content with configuration.
//Config includes:
// route:
// template_folder
// root:
// default:
func RouteContent(r *mux.Router, identifier string, config map[string]interface{}) error {
	routesConfig := config["routes"].([]interface{})

	if _, ok := config["template_folder"]; !ok {
		return errors.New("Need template_folder setting.")
	}
	templateFolder := util.InterfaceToStringArray(config["template_folder"].([]interface{}))

	if _, ok := config["root"]; !ok {
		return errors.New("Need root setting.")
	}
	root := config["root"].(int)

	if _, ok := config["default"]; !ok {
		return errors.New("Need default setting.")
	}
	defaultContent := config["default"].(int)
	siteSettings := website.SiteSettings{TemplateBase: templateFolder[0],
		TemplateFolders: templateFolder}
	website.InitSiteSettings(identifier, siteSettings)
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
				routeContent(s, templateFolder[0], path, root, defaultContent)
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
		renderContent(w, r, id, templateFolder, prefix, root, defaultContent)
	})

	r.MatcherFunc(niceurl.ViewContentMatcher).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		renderContent(w, r, id, templateFolder, prefix, root, defaultContent)
	})
}

func renderContent(w http.ResponseWriter, r *http.Request, id int, templateFolder string, prefix string, rootContent int, defaultContent int) {
	// Execute the template per HTTP request
	pongo2.DefaultSet.Debug = true
	pongo2.DefaultSet.SetBaseDirectory("../templates/" + templateFolder)
	tplExample := pongo2.Must(pongo2.FromCache("../default/content/view.html"))
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	//todo: handle error, template compiling much better.
	if err != nil {
		fmt.Println(err)
	}
	if !util.ContainsInt(content.GetLocation().Path(), rootContent) {
		w.Write([]byte("Not valid content"))
		return
	}
	root, err := querier.FetchByID(rootContent)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = tplExample.ExecuteWriter(pongo2.Context{"content": content,
		"root":     root,
		"viewmode": "full",
		"site":     "demosite",
		"prefix":   prefix}, w)

}
