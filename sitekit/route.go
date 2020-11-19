package sitekit

import (
	"github.com/digimakergo/digimaker/core/util"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func SiteRouterHandle(r *mux.Router, identifierStr string, pattern string, handle func(w http.ResponseWriter, r *http.Request)) {
	list := util.Split(identifierStr)
	for _, identifier := range list {
		SiteRouter(r, identifier, func(s *mux.Router, site string) {
			s.HandleFunc(pattern, func(wr http.ResponseWriter, re *http.Request) {
				handle(wr, re)
			})
		})
	}
}

//Route a site content with configuration.
//reoutesConfig includes:
//- host
//- path
//- combined
func SiteRouter(r *mux.Router, identifier string, handler func(s *mux.Router, site string)) error {
	//go through all routes.
	routesConfig := GetSiteSettings(identifier).Routes
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
				handler(s, identifier)
			}
		}
	}

	return nil
}
