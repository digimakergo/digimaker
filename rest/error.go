package rest

import "net/http"

func HandleError(err error, w http.ResponseWriter) {
	w.WriteHeader(500)    //todo: make it as parameter?
	errStr := err.Error() //todo: use error code here
	w.Write([]byte(errStr))
}
