//Author xc, Created on 2019-08-11 17:07
//{COPYRIGHTS}

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/util"

	"github.com/gorilla/mux"
)

func New(w http.ResponseWriter, r *http.Request) {

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
	//todo: more permission check
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}

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

	agent := r.UserAgent()
	folderName := util.HomePath() + "/log/draft/" + strconv.Itoa(userId)
	if _, existError := os.Stat(folderName); os.IsNotExist(existError) {
		os.Mkdir(folderName, 0775)
	}
	logPath := folderName + "/" + time.Now().Format("20060102_150405.log")
	ip := r.Header.Get("X-Forwarded-For")

	inputs := map[string]string{}
	decorder := json.NewDecoder(r.Body)
	err = decorder.Decode(&inputs)
	if err != nil {
		HandleError(errors.New("wrong format"), w)
		ioutil.WriteFile(logPath, []byte(ip+","+agent+"\n"+err.Error()), 0775)
		return
	}

	logConent := []byte(ip + "," + agent + "\n" + fmt.Sprint(inputs))
	ioutil.WriteFile(logPath, logConent, 0775)

	data, ok := inputs["data"]
	if !ok {
		HandleError(errors.New("need data"), w)
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
	if newVersion.ID == 0 || newVersion.ID > 0 && newVersion.Created != createdTime {
		message := "Not saved. new id:" + strconv.Itoa(newVersion.ID) + " .new saved time:" +
			strconv.Itoa(newVersion.Created) + ", should saved time:" + strconv.Itoa(createdTime)
		HandleError(errors.New(message), w)
		return
	}

	w.Write([]byte(strconv.Itoa(newVersion.Created)))
}

func Update(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

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

	idStr := r.FormValue("id")
	idSlice := strings.Split(idStr, ",")
	for _, id := range idSlice {
		_, err := strconv.Atoi(id)
		if err != nil {
			HandleError(errors.New("Illegal id"), w, StatusWrongParams)
			return
		}
	}

	for _, id := range idSlice {
		idInt, _ := strconv.Atoi(id)
		handler := handler.ContentHandler{}
		handler.Context = r.Context()
		//todo: use Delete by ids to support one transaction with roll back
		err := handler.DeleteByID(idInt, userID, true)
		if err != nil {
			HandleError(err, w)
			return
		}
	}

	w.Write([]byte("1"))
}

func init() {

	RegisterRoute("/content/new/{parent}/{contenttype}", New)
	RegisterRoute("/content/update/{id}", Update)
	RegisterRoute("/content/delete", Delete)
	RegisterRoute("/content/savedraft/{id}/{type}", SaveDraft)
}
