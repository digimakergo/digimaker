package rest

import "net/http"

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
