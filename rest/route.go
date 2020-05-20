package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

type key int

//CtxKeyUserID defines context key for user id
const CtxKeyUserID = key(1)

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
		userIDStr := ""
		if r.Header.Get("Authorization") != "" {
			err, claims := VerifyToken(r)

			if err != nil {
				if err == TokenErrorExpired {
					HandleError(err, w, StatusExpired)
				} else {
					log.Error("Token verification error: "+err.Error(), "")
					HandleError(errors.New("Invalid token"), w, StatusUnauthed)
				}
				return
			}

			userIDStr = strconv.Itoa(claims.UserID)
			ctx = context.WithValue(ctx, CtxKeyUserID, claims.UserID)
		}

		//start debug
		requestID := util.GenerateGUID()
		ctx = log.WithLogger(ctx, logrus.Fields{"ip": util.GetIP(r), "user": userIDStr, "request_id": requestID})
		log.StartTiming(ctx, "request")

		w.Header().Add("DM-Request-Id", requestID)
		w.Header().Set("Access-Control-Allow-Origin", "*") //todo: make host configurable

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		//close debug
		log.TickTiming(ctx, "request")
		log.LogTiming(ctx)
	})
}
