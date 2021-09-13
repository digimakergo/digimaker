package rest

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Error bool        `json:"error"`
	Data  interface{} `json:"data"`
}

func WriteResponse(data interface{}, w http.ResponseWriter, isError ...bool) {
	res := response{}
	if len(isError) > 0 && isError[0] {
		res.Error = true
	}
	res.Data = data
	outputData, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error when outputing: " + err.Error()))
	}
	w.Header().Set("content-type", "application/json")
	w.Write(outputData)
}
