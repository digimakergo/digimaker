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

const CtxUserKey = key(1)

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
			success, err, claims := VerifyToken(r)
			if err != nil {
				log.Error("Error for authentication: "+err.Error(), "", ctx)
				HandleError(errors.New("Error for authenticatoin"), w)
				return
			}
			if !success {
				w.Write([]byte("Authorization failed."))
				return
			}

			userIDStr = strconv.Itoa(claims.UserID)
			ctx = context.WithValue(ctx, CtxUserKey, claims.UserID)
		}

		//start debug
		requestID := util.GenerateGUID()
		ctx = log.WithLogger(ctx, logrus.Fields{"ip": util.GetIP(r), "user": userIDStr, "request_id": requestID})
		log.StartTiming(ctx, "request")

		w.Header().Add("DM-Request-Id", requestID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		//close debug
		log.TickTiming(ctx, "request")
		log.LogTiming(ctx)
	})
}
