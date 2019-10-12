//Author xc, Created on 2019-08-11 17:07
//{COPYRIGHTS}

package rest

import (
	"dm/core/handler"
	"dm/core/util/debug"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := debug.Init(r.Context())
	r = r.WithContext(ctx)

	params := mux.Vars(r)
	parent := params["parent"]
	parentInt, err := strconv.Atoi(parent)
	if err != nil {
		HandleError(errors.New("parent id should be integer."), w)
		return
	}
	contentType := params["contenttype"]

	inputs := map[string]interface{}{}
	decorder := json.NewDecoder(r.Body)
	err = decorder.Decode(&inputs)

	fmt.Println(err)
	fmt.Println(inputs)
	//todo: add value based on definition

	handler := handler.ContentHandler{Context: r.Context()}
	content, validationResult, err := handler.Create(contentType, inputs, parentInt)

	w.Header().Set("content-type", "application/json")
	if !validationResult.Passed() {
		data, _ := json.Marshal(validationResult)
		w.WriteHeader(400)
		w.Write(data)
		return
	}

	if err != nil {
		HandleError(err, w)
		return
	}

	data, _ := json.Marshal(content)
	w.Write(data)
}

func Update() {

}

func Delete() {

}
