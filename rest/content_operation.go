//Author xc, Created on 2019-08-11 17:07
//{COPYRIGHTS}

package rest

import (
	"dm/core/handler"
	"dm/core/util/debug"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func New(w http.ResponseWriter, r *http.Request) {
	//todo: permission
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

	if err != nil {
		HandleError(errors.New("Invalid input for json."), w)
		return
	}
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

func Update(w http.ResponseWriter, r *http.Request) {
	//todo: permission
	ctx := debug.Init(r.Context())
	r = r.WithContext(ctx)

	params := mux.Vars(r)
	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(errors.New("id should be integer."), w)
		return
	}

	inputs := map[string]interface{}{}
	decorder := json.NewDecoder(r.Body)
	err = decorder.Decode(&inputs)

	if err != nil {
		HandleError(errors.New("Invalid input for json."), w)
		return
	}
	//todo: add value based on definition

	handler := handler.ContentHandler{Context: r.Context()}
	result, validationResult, err := handler.UpdateByID(idInt, inputs)

	if !result {
		w.Header().Set("content-type", "application/json")
		if !validationResult.Passed() {
			data, _ := json.Marshal(validationResult)
			w.WriteHeader(400)
			w.Write(data)
			return
		}
	}

	if err != nil {
		HandleError(err, w)
		return
	}

	w.Write([]byte("1"))
}

func Delete() {

}

func init() {

	RegisterRoute("/content/new/{parent}/{contenttype}", New)
	RegisterRoute("/content/update/{id}", Update)

}
