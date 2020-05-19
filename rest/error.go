package rest

import (
	"context"
	"errors"
	"net/http"
)

func HandleError(err error, w http.ResponseWriter, httpCode ...int) {
	//todo: output debug info if needed.
	if len(httpCode) == 0 {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(httpCode[0])
	}
	errStr := err.Error() //todo: use error code here
	w.Write([]byte(errStr))
}

//Check if there is user id, if not output error and return 0
func CheckUserID(context context.Context, w http.ResponseWriter) int {
	userId := context.Value(CtxKeyUserID)
	if userId == nil {
		HandleError(errors.New("Need to login"), w, 401)
		return 0
	} else {
		return userId.(int)
	}
}
