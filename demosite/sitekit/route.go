package sitekit

import (
	"dm/dm/handler"
	"dm/dm/util"
	"dm/dm/website"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/flosch/pongo2.v2"
)

func SiteRouterHandle(r *mux.Router, identifier string, pattern string, handle func(w http.ResponseWriter, r *http.Request)) {
	SiteRouter(r, identifier, func(s *mux.Router) {
		s.HandleFunc(pattern, func(wr http.ResponseWriter, re *http.Request) {
			handle(wr, re)
		})
	})
}

//Route a site content with configuration.
//reoutesConfig includes:
//- host
//- path
//- combined
func SiteRouter(r *mux.Router, identifier string, handler func(s *mux.Router)) error {
	//go through all routes.
	routesConfig := website.GetSiteSettings(identifier).Routes
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
					s = r.Host(host).PathPrefix("/{path:" + path + "}" + "/").Subrouter()
				} else {
					s = r.Host(host).Subrouter()
				}
				handler(s)
			}
		}
	}

	return nil
}

func OutputContent(w http.ResponseWriter, r *http.Request, id int, siteIdentifier string, prefix string) {
	// Execute the template per HTTP request
	pongo2.DefaultSet.Debug = true
	siteSettings := website.GetSiteSettings(siteIdentifier)
	pongo2.DefaultSet.SetBaseDirectory("../templates/" + siteSettings.TemplateBase)
	tplExample := pongo2.Must(pongo2.FromCache("../default/content/view.html"))
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	//todo: handle error, template compiling much better.
	if err != nil {
		fmt.Println(err)
	}
	if !util.ContainsInt(content.GetLocation().Path(), siteSettings.RootContent.GetLocation().ID) {
		w.Write([]byte("Not valid content"))
		return
	}

	err = tplExample.ExecuteWriter(pongo2.Context{"content": content,
		"root":     siteSettings.RootContent,
		"viewmode": "full",
		"site":     "demosite",
		"prefix":   prefix}, w)

}
