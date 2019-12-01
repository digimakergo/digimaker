package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

var routeMap map[string]func(http.ResponseWriter, *http.Request) = map[string]func(http.ResponseWriter, *http.Request){}

func RegisterRoute(path string, funcHandler func(http.ResponseWriter, *http.Request)) {
	routeMap[path] = funcHandler
}

//Handle route with context. eg. user
//Loop registered route and listen handle function
func HandleRoute(router *mux.Router) {
	for path, handleFunc := range routeMap {
		router.HandleFunc(path, handleFunc)
	}

}

//Initialize request, including set context.
//todo: support callback
func InitRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
