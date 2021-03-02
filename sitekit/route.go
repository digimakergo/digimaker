package sitekit

import (
	"errors"

	"github.com/digimakergo/digimaker/core/util"

	"github.com/gorilla/mux"
)

//Route a site content with configuration.
//reoutesConfig includes:
//- host
//- path
//- combined
func HandleOnSite(r *mux.Router, identifier string, handler func(s *mux.Router, site string)) error {
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
				var subRouter *mux.Router
				//use subrouter which is better for performance
				if path != "" {
					subRouter = r.Host(host).PathPrefix("/{path:" + path + "}").Subrouter()
				} else {
					subRouter = r.Host(host).Subrouter()
				}
				handler(subRouter, identifier)
			}
		}
	}

	return nil
}
