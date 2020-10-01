//Author xc, Created on 2020-10-1 12:50
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"net/http"

	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/permission"
)

//Get current user's updatefields on him/herself
func CurrentUserEditField(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}
	querier := handler.Querier()
	content, _ := querier.GetUser(userID)

	fields, err := permission.GetUpdateFields(r.Context(), content, userID)
	if err != nil {
		HandleError(err, w)
		return
	}

	data, _ := json.Marshal(fields)
	w.Write(data)
}

func init() {
	RegisterRoute("/access/update-fields/current-user", CurrentUserEditField)
}
