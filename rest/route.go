package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/digimakergo/digimaker/core/auth"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/log/httplog"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/gorilla/mux"
)

type key int

type routerItem struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
	Methods []string //GET/POST/PUT/empty - empty means all
}

var routeList []routerItem = []routerItem{}

func RegisterRoute(path string, funcHandler func(http.ResponseWriter, *http.Request), methods ...string) {
	item := routerItem{path, funcHandler, methods}
	routeList = append(routeList, item)
}

//Handle route with context. eg. user
//Loop registered route and listen handle function
func HandleRoute(router *mux.Router) {
	for _, item := range routeList {
		route := router.HandleFunc(item.Path, item.Handler)
		if len(item.Methods) > 0 {
			route.Methods(item.Methods...)
		} else {
			route.Methods("GET")
		}
	}
}

//Initialize request, including set context.
//todo: support callback
func InitRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//set user_id to context
		userID := 0
		authStr := r.Header.Get("Authorization")
		if authStr != "" && strings.HasPrefix(authStr, "Bearer ") {
			err, claims := VerifyAccessToken(r)

			if err != nil {
				if err == auth.TokenErrorExpired {
					HandleError(err, w, StatusExpired)
				} else {
					log.Error("Token verification error: "+err.Error(), "")
					HandleError(errors.New("Invalid token"), w, StatusUnauthed)
				}
				return
			}

			userID = claims.UserID
			ctx = context.WithValue(ctx, util.CtxKeyUserID, userID)
		}

		//start http log
		r = httplog.InitLog(r, ctx, userID)

		next.ServeHTTP(w, r)

		//write http log
		httplog.LogRequest(w, r)
	})
}
