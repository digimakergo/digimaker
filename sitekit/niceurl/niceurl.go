//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

//Package niceurl provides nice url feature for dm framework
package niceurl

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

func GenerateUrl(content contenttype.ContentTyper, root contenttype.ContentTyper, prefix string) string { //todo: add prefix
	location := content.GetLocation()
	result := ""
	if location != nil {
		rootDepth := root.GetLocation().Depth
		path := strings.Join(strings.Split(location.IdentifierPath, "/")[rootDepth:], "/")
		if util.Contains(getPathContenttypes(), content.ContentType()) {
			result = path
		} else {
			result = path + "-" + strconv.Itoa(location.ID)
		}

	} else {
		//todo: give a warning.
	}
	if prefix != "" {
		result = "/" + prefix + result
	} else {
		result = "/" + result
	}
	return result
}

func getPathContenttypes() []string {
	return viper.GetStringSlice("site_settings.niceurl_contenttype")
}

//Matches pattern *-1231 and set mux Vars["id"] as 1231 if matches
func ViewContentMatcher(r *http.Request, rm *mux.RouteMatch) bool {
	uri := r.URL.Path
	result := false
	if strings.HasPrefix(uri, "/") {
		index := strings.LastIndex(uri, "-") //todo: read pattern from config file related to GenerateUrl
		if index != -1 {
			str := uri[index+1:]
			_, err := strconv.Atoi(str)
			if err == nil {
				rm.Vars = map[string]string{"id": str}
				result = true
			} else {
				result = MatchPath(r, rm)
			}
		} else {
			result = MatchPath(r, rm)
		}
	}
	return result
}

func MatchPath(r *http.Request, rm *mux.RouteMatch) bool {
	uri := r.RequestURI
	if ok, _ := regexp.Match("[a-zA-Z0-9_\\-\\/+]$", []byte(uri)); ok {
		for _, cType := range getPathContenttypes() {
			root, _ := query.FetchLocationByID(r.Context(), 3) //todo:
			path := root.IdentifierPath + uri
			location, _ := query.FetchLocation(r.Context(), db.Cond("content_type", cType).Cond("identifier_path", path))
			if location.ID > 0 {
				rm.Vars = map[string]string{"id": strconv.Itoa(location.ID)}
				return true
			}
		}
	}
	return false
}
