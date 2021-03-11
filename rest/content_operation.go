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

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"

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
	dbHandler.GetEntity(r.Context(), version.TableName(),
		db.Cond("author", userId).Cond("content_id", id).Cond("version", 0).Cond("content_type", ctype),
		[]string{}, nil, &version, false)
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
	err = version.Store(r.Context(), tx)
	tx.Commit()
	if err != nil {
		HandleError(err, w)
		return
	}

	//refetch to make sure it's saved
	newVersion := contenttype.Version{}
	dbHandler.GetEntity(r.Context(), version.TableName(),
		db.Cond("author", userId).Cond("content_id", id).Cond("version", 0).Cond("content_type", ctype),
		[]string{}, nil, &newVersion, false)
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

	updateHandler := handler.ContentHandler{Context: r.Context()}

	contenttype := params["contenttype"]
	var result bool
	var validationResult handler.ValidationResult

	if contenttype == "" {
		result, validationResult, err = updateHandler.UpdateByID(idInt, inputs, userID)
	} else {
		result, validationResult, err = updateHandler.UpdateByContentID(contenttype, idInt, inputs, userID)
	}

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

func SetPriority(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := r.FormValue("params")
	paramArr := strings.Split(params, ";")
	tx, _ := db.CreateTx()
	for _, paramStr := range paramArr {
		arr := strings.Split(paramStr, ",")
		id, err := strconv.Atoi(arr[0])
		if err != nil {
			HandleError(errors.New("Invalid params format on id"), w)
			return
		}
		priority, err := strconv.Atoi(arr[1])
		if err != nil {
			HandleError(errors.New("Invalid params format on priority"), w)
			return
		}
		content, err := query.FetchByID(r.Context(), id)
		if err != nil {
			log.Error(err.Error(), "")
			HandleError(errors.New("Can't find content"), w)
			return
		}

		if !permission.CanUpdate(r.Context(), content, userID) {
			HandleError(errors.New("No permision for "+strconv.Itoa(id)), w)
			tx.Rollback()
			return
		}

		//todo: create transction once.
		//update priority
		location := content.GetLocation()
		location.Priority = priority
		err = location.Store(r.Context(), tx)
		if err != nil {
			tx.Rollback()
			log.Error(err.Error(), "")
			HandleError(errors.New("Something wrong"), w)
			return
		}
	}
	tx.Commit()
	w.Write([]byte("1"))
}

func Move(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := CheckUserID(ctx, w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	contentStr := params["contents"]
	contentIds := []int{}

	contentIdStrs := strings.Split(contentStr, ",")
	for _, str := range contentIdStrs {
		contentId, err := strconv.Atoi(str)
		if err != nil {
			HandleError(errors.New("Invalid content format"), w, 410)
			return
		}
		contentIds = append(contentIds, contentId)
	}
	targetId, _ := strconv.Atoi(params["target"])

	handler := handler.ContentHandler{}
	err := handler.Move(ctx, contentIds, targetId, userID)
	if err != nil {
		HandleError(err, w, 410)
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
	cidStr := r.FormValue("cid")
	contenttype := r.FormValue("type")
	if idStr != "" {
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
	} else if cidStr != "" && contenttype != "" {
		cids := strings.Split(cidStr, ",")
		for _, cid := range cids {
			_, err := strconv.Atoi(cid)
			if err != nil {
				HandleError(errors.New("Invalid cid"), w, StatusWrongParams)
				return
			}
		}
		for _, cid := range cids {
			cidInt, _ := strconv.Atoi(cid)
			handler := handler.ContentHandler{}
			handler.Context = r.Context()
			err := handler.DeleteByCID(cidInt, contenttype, userID)
			if err != nil {
				HandleError(err, w)
				return
			}
		}
	} else {
		HandleError(errors.New("Invalid parameters"), w, StatusWrongParams)
		return
	}

	w.Write([]byte("1"))
}

func init() {

	RegisterRoute("/content/new/{parent}/{contenttype}", New)
	RegisterRoute("/content/move/{contents}/{target}", Move)
	RegisterRoute("/content/update/{id}", Update)
	RegisterRoute("/content/update/{contenttype}/{id}", Update)
	RegisterRoute("/content/delete", Delete)
	RegisterRoute("/content/setpriority", SetPriority)
	RegisterRoute("/content/savedraft/{id}/{type}", SaveDraft)
}
