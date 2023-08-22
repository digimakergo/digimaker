package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

//Get current user
func CurrentUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	site := params["site"]

	ctx := r.Context()
	userID := util.CurrentUserID(ctx)

	user, _ := query.GetUser(userID)
	if user != nil {
		hasAccess := permission.HasAccessTo(ctx, userID, "site/access", map[string]interface{}{"site": site})
		if hasAccess {
			roles, _ := query.FetchUserRoles(ctx, userID)
			roleResult := []string{}
			for _, role := range roles {
				roleResult = append(roleResult, role.Value("identifier").(string))
			}

			result, _ := contenttype.ContentToMap(user)
			result["roles"] = roleResult

			WriteResponse(result, w)
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
	user, _ := query.Fetch(r.Context(), "user", db.Cond("c.email", email))
	if user == nil {
		HandleError(errors.New("User not found."), w)
		return
	}

	disabled := user.Value("disabled")
	if disabled != nil {
		if disabled.(int) == 1 {
			log.Info("User is disabled")
			HandleError(errors.New("User not found."), w)
			return
		}
	}

	//create hash
	activation := Activiation{}
	activation.Hash = util.GenerateUID()
	activation.Ref = strconv.Itoa(user.GetID())
	activation.Created = int(time.Now().Unix())
	activation.Type = "resetpassword"
	data, _ := json.Marshal(activation)
	var dbMap map[string]interface{}
	json.Unmarshal(data, &dbMap)
	_, err := db.Insert(r.Context(), "dm_activation", dbMap)
	if err == nil {
		//send password
		settings := viper.GetStringMapString("resetpassword_email")

		cond := util.Split(settings["condition"], ",")
		template, templateErr := query.Fetch(r.Context(), settings["type"], db.Cond(cond[0], cond[1]))
		if templateErr != nil {
			log.Error(templateErr, "", r.Context())
			return
		}

		if template == nil {
			log.Error(errors.New("No resset password template found"), "", r.Context())
			return
		}

		subject := template.Value(settings["subject"]).(string)
		body := template.Value(settings["body"]).(string)
		body = strings.ReplaceAll(body, "{hash}", activation.Hash)

		err = util.SendMail(
			[]string{email},
			subject,
			body)
		if err != nil {
			log.Error(fmt.Errorf("Mail not sent. %v", err), "", r.Context())
			return
		}
		//todo: rollback
		WriteResponse(true, w)
	}

}

//todo: move this to logic
func ResetPasswordDone(w http.ResponseWriter, r *http.Request) {
	//expire time
	// hours := 24
	params := mux.Vars(r)
	hash := params["hash"]
	activation := Activiation{}
	db.BindEntity(r.Context(), &activation, "dm_activation", db.Cond("hash", hash).Cond("type", "resetpassword"))
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
			user, _ := query.FetchByCID(r.Context(), "user", ref)
			if user == nil {
				HandleError(errors.New("No user found"), w)
				return
			}

			if password == "" {
				HandleError(errors.New("Password is empty"), w)
				return
			}

			_, err := handler.Update(r.Context(), user, handler.InputMap{"password": password}, 1)
			if err != nil {
				HandleError(err, w)
				return
			}

			db.Delete(r.Context(), "dm_activation", db.Cond("id", activation.ID))

			WriteResponse(true, w)
		}

	}

}

func EnableUser(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	ids := strings.Split(r.FormValue("id"), ",")
	enableType := params["type"]

	users := []contenttype.ContentTyper{}
	var err error
	for _, id := range ids {
		idInt, conErr := strconv.Atoi(id)
		if conErr != nil {
			HandleError(conErr, w)
			return
		}
		user, _ := query.FetchByCID(r.Context(), "user", idInt)
		if user == nil {
			err = errors.New("User not found: " + id)
		} else {
			users = append(users, user)
		}
	}

	if err != nil {
		HandleError(err, w)
		return
	}

	for _, user := range users {
		err = handler.Enable(r.Context(), user, enableType == "1", userID)
		if err != nil {
			HandleError(err, w)
			return
		}
	}

	WriteResponse(true, w)
}

func FetchUserRoles(w http.ResponseWriter, r *http.Request) {
	currentUserID := CheckUserID(r.Context(), w)
	if currentUserID == 0 {
		return
	}

	if !permission.HasAccessTo(r.Context(), currentUserID, ACCESS_MANAGE_ACCESS, permission.MatchData{}) {
		HandleError(errors.New("No access"), w)
		return
	}

	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["user"])
	list, err := handler.FetchUserRoles(r.Context(), userID)
	if err != nil {
		HandleError(err, w)
		return
	}
	WriteResponse(list, w)
}

func init() {
	RegisterRoute("/user/current/{site}", CurrentUser)
	RegisterRoute("/user/resetpassword/{email}", ResetPassword)
	RegisterRoute("/user/enable/{type}", EnableUser)
	RegisterRoute("/user/roles/{user}", FetchUserRoles)
	RegisterRoute("/user/resetpassword-confirm/{hash}", ResetPasswordDone, "POST")
}
