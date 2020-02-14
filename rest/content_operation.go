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
	ctx := debug.Init(r.Context())
	r = r.WithContext(ctx)

	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}

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
	content, validationResult, err := handler.Create(contentType, inputs, userId, parentInt)

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

	//todo: more permission check
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}

	version := contenttype.Version{}
	dbHandler := db.DBHanlder()
	dbHandler.GetEntity(version.TableName(),
		db.Cond("author", userId).Cond("content_id", id).Cond("version", 0).Cond("content_type", ctype),
		[]string{}, &version)
	if version.ID == 0 {
		version.ContentType = ctype
		version.ContentID = id
		version.Author = userId
		version.Version = 0
	}
	version.Data = data
	createdTime := int(time.Now().Unix())
	version.Created = createdTime
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

	//refetch to make sure it's saved
	newVersion := contenttype.Version{}
	dbHandler.GetEntity(version.TableName(),
		db.Cond("author", userId).Cond("content_id", id).Cond("version", 0).Cond("content_type", ctype),
		[]string{}, &newVersion)
	if newVersion.ID == 0 || newVersion.ID > 0 && newVersion.Created == createdTime {
		HandleError(errors.New("Not saved."), w)
		return
	}

	w.Write([]byte(strconv.Itoa(newVersion.Created)))
}

func Update(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

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
	result, validationResult, err := handler.UpdateByID(idInt, inputs, userID)

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

func Delete(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	idInt, _ := strconv.Atoi(id)
	handler := handler.ContentHandler{}
	ctx := debug.Init(r.Context())
	r = r.WithContext(ctx)
	handler.Context = ctx
	err := handler.DeleteByID(idInt, userID, true)
	if err != nil {
		HandleError(err, w)
		return
	}
	w.Write([]byte("1"))
}

func init() {

	RegisterRoute("/content/new/{parent}/{contenttype}", New)
	RegisterRoute("/content/update/{id}", Update)
	RegisterRoute("/content/delete/{id}", Delete)
	RegisterRoute("/content/savedraft/{id}/{type}", SaveDraft)
}
