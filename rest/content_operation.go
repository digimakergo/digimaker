//Author xc, Created on 2019-08-11 17:07
//{COPYRIGHTS}

package rest

import (
	"dm/core/contenttype"
	"dm/core/db"
	"dm/core/handler"
	"dm/core/util/debug"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func New(w http.ResponseWriter, r *http.Request) {
	//todo: permission
	ctx := debug.Init(r.Context())
	r = r.WithContext(ctx)

	userId := r.Context().Value("user_id")
	if userId == nil {
		HandleError(errors.New("Need to login"), w) //todo: add status code
		return
	}
	userIdInt := userId.(int)

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
	content, validationResult, err := handler.Create(contentType, inputs, userIdInt, parentInt)

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

func SaveDraft(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctype := params["type"]
	_, err := contenttype.GetDefinition(ctype)
	if err != nil {
		HandleError(errors.New("type doesn't exist."), w)
		return
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		HandleError(errors.New("Wrong id"), w)
		return
	}

	inputs := map[string]string{}
	decorder := json.NewDecoder(r.Body)
	err = decorder.Decode(&inputs)
	if err != nil {
		HandleError(errors.New("wrong format"), w)
		return
	}
	data, ok := inputs["data"]
	if !ok {
		HandleError(errors.New("need data"), w)
		return
	}

	//todo: permission check
	userIdI := r.Context().Value("user_id")
	if userIdI == nil {
		HandleError(errors.New("need login"), w)
		return
	}
	userId := userIdI.(int)

	version := contenttype.Version{}
	dbHandler := db.DBHanlder()
	dbHandler.GetEntity(version.TableName(),
		db.Cond("author", userId).Cond("content_id", id).Cond("content_type", ctype),
		[]string{}, &version)
	if version.ID == 0 {
		version.ContentType = ctype
		version.ContentID = id
		version.Author = userId
		version.Version = 0
	}
	version.Data = data
	version.Created = int(time.Now().Unix())
	tx, _ := db.CreateTx()
	if err != nil {
		HandleError(err, w)
		return
	}
	err = version.Store(tx)
	tx.Commit()
	if err != nil {
		HandleError(err, w)
		return
	}

	w.Write([]byte(strconv.Itoa(version.Created)))
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
	RegisterRoute("/content/savedraft/{id}/{type}", SaveDraft)
}
