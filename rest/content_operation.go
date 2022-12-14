//Author xc, Created on 2019-08-11 17:07
//{COPYRIGHTS}

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"

	"github.com/gorilla/mux"
)

func Create(w http.ResponseWriter, r *http.Request) {

	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	parent := params["parent"]
	parentInt := 0
	if parent != "" {
		var err error
		parentInt, err = strconv.Atoi(parent)
		if err != nil {
			HandleError(errors.New("parent id should be integer."), w)
			return
		}
	}
	contentType := params["contenttype"]

	inputs := map[string]interface{}{}
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&inputs)

	if err != nil {
		HandleError(errors.New("Invalid input for json."), w)
		return
	}
	//todo: add value based on definition

	content, err := handler.Create(r.Context(), userID, contentType, inputs, parentInt)

	if err != nil {
		HandleError(err, w)
		return
	}

	WriteResponse(content, w)
}

//Save draft. url: /<id>/<ctype> where id can be 0 or parent location id
//request json: data: string(mostly json string)
//return new created time
func SaveDraft(w http.ResponseWriter, r *http.Request) {
	userId := CheckUserID(r.Context(), w)
	if userId == 0 {
		return
	}

	params := mux.Vars(r)
	ctype := params["type"]

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		HandleError(errors.New("Wrong id"), w)
		return
	}

	inputs := map[string]string{}
	decorder := json.NewDecoder(r.Body)
	err = decorder.Decode(&inputs)

	data, ok := inputs["data"]
	if !ok {
		HandleError(errors.New("need data"), w)
		return
	}

	newVersion, err := handler.SaveDraft(r.Context(), userId, data, ctype, id)
	if err != nil {
		HandleError(err, w)
		return
	}

	WriteResponse(newVersion.Created, w)
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

	contenttype := params["contenttype"]

	if contenttype == "" {
		_, err = handler.UpdateByID(r.Context(), idInt, inputs, userID)
	} else {
		_, err = handler.UpdateByContentID(r.Context(), contenttype, idInt, inputs, userID)
	}

	if err != nil {
		HandleError(err, w)
		return
	}

	WriteResponse(true, w)
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

		if !permission.CanUpdate(r.Context(), content, []string{}, userID) {
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
	WriteResponse(true, w)
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

	err := handler.Move(ctx, contentIds, targetId, userID)
	if err != nil {
		HandleError(err, w, 410)
		return
	}
	WriteResponse(true, w)
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
			//todo: use Delete by ids to support one transaction with roll back
			err := handler.DeleteByID(r.Context(), idInt, userID, true)
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
			err := handler.DeleteByCID(r.Context(), cidInt, contenttype, userID)
			if err != nil {
				HandleError(err, w)
				return
			}
		}
	} else {
		HandleError(errors.New("Invalid parameters"), w, StatusWrongParams)
		return
	}

	WriteResponse(true, w)
}

//copy old data
func Copy(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	parent := params["parent"]

	parentInt, err := strconv.Atoi(parent)
	if err != nil {
		HandleError(errors.New("parent id should be integer."), w)
		return
	}

	// To number
	idInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(errors.New("id should be integer."), w)
		return
	}
	ctype := params["contenttype"]

	if ctype == "" {
		contenter, err := query.FetchByID(r.Context(), idInt)
		if err != nil {
			HandleError(fmt.Errorf("Failed to get content via content id: %w", err), w)
			return
		}
		if contenter.GetCID() == 0 {
			HandleError(fmt.Errorf("Got empty : %w", err), w)
			return
		}
		_, err = handler.Copy(r.Context(), contenter, ctype, userID, parentInt)
		if err != nil {
			HandleError(err, w)
			return
		}
	} else {
		contenter, err := query.FetchByCID(r.Context(), ctype, idInt)
		if err != nil {
			HandleError(fmt.Errorf("Failed to get content via content id: %w", err), w)
			return
		}
		if contenter.GetCID() == 0 {
			HandleError(fmt.Errorf("Got empty content: %w", err), w)
			return
		}
		_, err = handler.Copy(r.Context(), contenter, ctype, userID, parentInt)
		if err != nil {
			HandleError(err, w)
			return
		}
	}

	WriteResponse(true, w)
}

func init() {
	RegisterRoute("/content/create/{contenttype}/{parent:[0-9]+}", Create, "POST")
	RegisterRoute("/content/create/{contenttype}", Create, "POST")
	RegisterRoute("/content/move/{contents}/{target}", Move)
	RegisterRoute("/content/update/{id:[0-9]+}", Update, "POST")
	RegisterRoute("/content/update/{contenttype}/{id:[0-9]+}", Update, "POST")
	RegisterRoute("/content/delete", Delete)
	RegisterRoute("/content/setpriority", SetPriority)
	RegisterRoute("/content/savedraft/{id:[0-9]+}/{type}", SaveDraft, "POST")
	RegisterRoute("/content/Copy/{id:[0-9]+}/{parent:[0-9]+}", Copy, "POST")
	RegisterRoute("/content/Copy/{contenttype}/{id:[0-9]+}/{parent:[0-9]+}", Copy, "POST")
}
