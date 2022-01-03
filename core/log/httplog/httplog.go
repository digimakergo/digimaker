package httplog

import (
	"context"
	"net/http"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
)

func InitLog(r *http.Request, ctx context.Context, userID int) *http.Request {

	//check if is debuggable
	canDebug := false
	debugHeader := util.GetConfig("general", "debug_header")
	debugToken := r.Header.Get(debugHeader)
	if debugToken != "" && debugToken == util.GetDebugToken() {
		canDebug = true
	}

	requestID := util.GenerateGUID()
	info := log.ContextInfo{
		RequestID: requestID,
		IP:        util.GetIP(r),
		UserID:    userID,
		URI:       r.RequestURI,
		Debug:     canDebug,
	}

	ctx = log.InitContext(ctx, &info)
	log.StartTiming(ctx, "total")

	result := r.WithContext(ctx)
	return result
}

func LogRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.EndTiming(ctx, "total")
	log.LogTiming(ctx)
	requestID := log.GetContextInfo(ctx).RequestID
	w.Header().Add("DM-Request-Id", requestID)
}
