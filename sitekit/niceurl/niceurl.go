//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

//Package niceurl provides nice url feature for dm framework
package niceurl

import (
	"github.com/xc/digimaker/core/contenttype"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GenerateUrl(content contenttype.ContentTyper) string { //todo: add prefix
	location := content.GetLocation()
	result := ""
	if location != nil {
		path := location.IdentifierPath
		pattern := "digit" //todo: read from config file.
		switch pattern {
		case "digit":
			result = path + "-" + strconv.Itoa(location.ID)
		default:

		}
	} else {
		//todo: give a warning.
	}
	return result
}

//Matches pattern *-1231 and set mux Vars["id"] as 1231 if matches
func ViewContentMatcher(r *http.Request, rm *mux.RouteMatch) bool {
	uri := r.RequestURI
	index := strings.LastIndex(uri, "-") //todo: read pattern from config file related to GenerateUrl
	result := false
	if index != -1 {
		str := uri[index+1:]
		_, err := strconv.Atoi(str)
		if err == nil {
			rm.Vars = map[string]string{"id": str}
			result = true
		}
	}
	return result
}
