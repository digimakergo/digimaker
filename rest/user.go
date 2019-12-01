package rest

import (
	"dm/eth/entity"
	"encoding/json"
	"errors"
	"net/http"
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

func init() {
	RegisterRoute("/user/current", CurrentUser)
}
