package rest

import (
	"errors"
	"net/http"

	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/viper"
)

func GenerateDebugToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := CheckUserID(ctx, w)
	if userID == 0 {
		return
	}

	if permission.HasAccessTo(ctx, userID, "util/debug") {
		duration := viper.GetDuration("general.debug_token_last")
		token := util.NewDebugToken(duration)
		WriteResponse(token, w)
	} else {
		HandleError(errors.New("No access"), w)
	}
}

func ClearDebugToken(w http.ResponseWriter, r *http.Request) {
	if util.GetDebugToken() == "" {
		WriteResponse("It's empty already", w)
	} else {
		util.ClearDebugToken()
		WriteResponse(true, w)
	}
}

func init() {
	RegisterRoute("/debug/generate-token", GenerateDebugToken)
	RegisterRoute("/debug/clear-token", ClearDebugToken)
}
