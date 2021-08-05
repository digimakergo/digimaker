package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/util"
)

const StatusUnauthed = 403
const StatusWrongParams = 400
const StatusExpired = 440
const StatusNotFound = 404
const StatusServer = 500

type errorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail"`
}

type responseError struct {
	Error errorBody `json:"error"`
}

func HandleError(err error, w http.ResponseWriter, httpCode ...int) {
	//todo: output debug info if needed.
	var statusCode int = StatusServer

	body := errorBody{Code: "10001", Message: err.Error(), Detail: ""} //todo: use error code here
	var validation handler.ValidationResult
	if errors.As(err, &validation) {
		body.Detail = err
		body.Code = "20001"
		statusCode = StatusWrongParams
	}
	resError := responseError{}
	resError.Error = body
	errStr, _ := json.Marshal(resError)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
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
