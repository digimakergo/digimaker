package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"

	"github.com/gorilla/mux"
)

//Get current user
func CurrentUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	site := params["site"]

	ctx := r.Context()
	userID := ctx.Value(CtxKeyUserID).(int)

	user, _ := handler.Querier().GetUser(userID)
	if user != nil {
		hasAccess := permission.HasAccessTo(userID, "site/access", map[string]interface{}{"site": site}, ctx)
		if hasAccess {
			data, err := json.Marshal(user)
			if err != nil {
				log.Error(err.Error(), "")
			}
			w.Write(data)
		} else {
			HandleError(errors.New("Doesn't have access"), w, 403)
		}
	} else {
		HandleError(errors.New("No user in context"), w)
	}
}

//todo: move this into entity folder
type Activiation struct {
	ID      int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Created int    `boil:"created" json:"created" toml:"created" yaml:"created"`
	Hash    string `boil:"hash" json:"hash" toml:"hash" yaml:"hash"`

	//type. eg. resetpassword
	Type string `boil:"type" json:"type" toml:"type" yaml:"type"`
	//reference. eg. userid
	Ref string `boil:"ref" json:"ref" toml:"ref" yaml:"ref"`
}

//todo: move this into user logic under core/user or handler/user.go folder
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	//todo: valid if it's a email
	querier := handler.Querier()
	user, _ := querier.Fetch("user", db.Cond("email", email))
	if user == nil {
		HandleError(errors.New("User not found."), w)
		return
	}

	//create hash
	activation := Activiation{}
	activation.Hash = util.GenerateUID()
	fmt.Println(activation.Hash)
	activation.Ref = strconv.Itoa(user.GetCID())
	activation.Created = int(time.Now().Unix())
	activation.Type = "resetpassword"
	dbHandler := db.DBHanlder()
	data, _ := json.Marshal(activation)
	var dbMap map[string]interface{}
	json.Unmarshal(data, &dbMap)
	_, err := dbHandler.Insert("dm_activation", dbMap)
	if err == nil {
		//send password
		url := "http://xxxx.com/api/user/resetpassword-confirm/" + activation.Hash
		util.SendMail(
			[]string{email},
			"reset password",
			"Click to reset password: <a href="+url+">url</a>")
		//todo: rollback
		w.Write([]byte("1"))
	}

}

//todo: move this to logic
func ResetPasswordDone(w http.ResponseWriter, r *http.Request) {
	//expire time
	// hours := 24
	params := mux.Vars(r)
	hash := params["hash"]
	dbHanldler := db.DBHanlder()
	activation := Activiation{}
	dbHanldler.GetEntity("dm_activation", db.Cond("hash", hash).Cond("type", "resetpassword"), []string{}, nil, &activation)
	if activation.ID == 0 {
		HandleError(errors.New("Wrong hash?"), w)
		return
	} else {
		now := time.Now()
		if time.Unix(int64(activation.Created), 0).Add(time.Hour * 24).Before(now) {
			HandleError(errors.New("Token expired."), w)
			return
		}

		inputs := map[string]interface{}{}
		decorder := json.NewDecoder(r.Body)
		err := decorder.Decode(&inputs)
		if err != nil {
			return
		}
		if password, ok := inputs["password"]; ok {
			ref, _ := strconv.Atoi(activation.Ref)
			querier := handler.Querier()
			user, _ := querier.FetchByContentID("user", ref)
			if user == nil {
				HandleError(errors.New("No user found"), w)
				return
			}

			if password == "" {
				HandleError(errors.New("Password is empty"), w)
				return
			}

			cHandler := handler.ContentHandler{}
			cHandler.Context = r.Context()
			success, validateResult, err := cHandler.Update(user, handler.InputMap{"password": password}, 1)

			if !success {
				if !validateResult.Passed() {
					HandleError(errors.New(validateResult.Fields["password"]), w)
					return
				}
				HandleError(err, w)
				return
			}

			dbHanldler.Delete("dm_activation", db.Cond("id", activation.ID))

			w.Write([]byte("1"))
		}

	}

}

func init() {
	RegisterRoute("/user/current/{site}", CurrentUser)
	RegisterRoute("/user/resetpassword/{email}", ResetPassword)
	RegisterRoute("/user/resetpassword-confirm/{hash}", ResetPasswordDone)
}
