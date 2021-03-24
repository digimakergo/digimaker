package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/digimakergo/digimaker/core/util"
)

const StatusUnauthed = 403
const StatusWrongParams = 400
const StatusExpired = 440
const StatusNotFound = 404

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type responseError struct {
	Error errorBody `json:"error"`
}

func HandleError(err error, w http.ResponseWriter, httpCode ...int) {
	//todo: output debug info if needed.
	if len(httpCode) == 0 {
		w.WriteHeader(StatusWrongParams)
	} else {
		w.WriteHeader(httpCode[0])
	}

	resError := responseError{}
	resError.Error = errorBody{Code: "10001", Message: err.Error()} //todo: use error code here
	errStr, _ := json.Marshal(resError)
	w.Write([]byte(errStr))
}

//Check if there is user id, if not output error and return 0
func CheckUserID(context context.Context, w http.ResponseWriter) int {
	userId := util.CurrentUserID(context)
	if userId == 0 {
		HandleError(errors.New("Need to login"), w, 401)
		return 0
	} else {
		return userId
	}
}
