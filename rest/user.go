package rest

import (
	"dm/core/db"
	"dm/core/fieldtype"
	"dm/core/handler"
	"dm/core/util"
	"dm/eth/entity"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := ctx.Value("user")
	if value != nil {
		user := value.(*entity.User)
		data, _ := json.Marshal(user)
		w.Write(data)
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
	dbHanldler.GetEntity("dm_activation", db.Cond("hash", hash).Cond("type", "resetpassword"), []string{}, &activation)
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
			pHandler := fieldtype.GetHandler("password")
			valid, message := pHandler.Validate(password)
			if !valid {
				HandleError(errors.New(message), w)
				return
			}
			value := pHandler.NewValue(password)
			user.SetValue("password", value)
			user.Store()

			dbHanldler.Delete("dm_activation", db.Cond("id", activation.ID))
			//todo: use transaction

			w.Write([]byte("1"))
		}

	}

}

func init() {
	RegisterRoute("/user/current", CurrentUser)
	RegisterRoute("/user/resetpassword/{email}", ResetPassword)
	RegisterRoute("/user/resetpassword-confirm/{hash}", ResetPasswordDone)
}
