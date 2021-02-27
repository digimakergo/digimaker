package httplog

import (
	"context"
	"net/http"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
)

func InitLog(r *http.Request, ctx context.Context, userID int) *http.Request {
	requestID := util.GenerateGUID() //todo: use debug id
	debugKey := "DebugID"

	debugID := ""
	if r.Header.Get(debugKey) != "" {
		debugID = r.Header.Get(debugKey)
	} else {
		cookie, err := r.Cookie(debugKey)
		if err != nil {
			debugID = cookie.String()
		}
	}

	info := log.ContextInfo{
		RequestID: requestID,
		IP:        util.GetIP(r),
		UserID:    userID,
		URI:       r.RequestURI,
		DebugID:   debugID,
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
