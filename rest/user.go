package rest

import (
	"dm/eth/entity"
	"encoding/json"
	"net/http"
)

func CurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	value := ctx.Value("user")
	user := value.(*entity.User)
	data, _ := json.Marshal(user)
	w.Write(data)
}

func init() {
	RegisterRoute("/user/current", CurrentUser)
}
