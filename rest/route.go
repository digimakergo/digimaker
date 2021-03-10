package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/log/httplog"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/gorilla/mux"
)

type key int

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
		ctx := r.Context()

		//set user_id to context
		userID := 0
		if r.Header.Get("Authorization") != "" {
			err, claims := VerifyAccessToken(r)

			if err != nil {
				if err == TokenErrorExpired {
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
