//Author xc, Created on 2020-10-1 12:50
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/log"
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

func AssignedUsers(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	if !permission.HasAccessTo(userID, "access/assigned-users", permission.MatchData{}, r.Context()) {
		HandleError(errors.New("No access"), w)
		return
	}

	parentID, err := strconv.Atoi(r.FormValue("parent"))
	if err != nil {
		HandleError(err, w)
		return
	}

	querier := handler.Querier()
	role, _ := querier.FetchByID(parentID)
	if role == nil {
		HandleError(errors.New("Roles doesn't exist"), w)
		return
	}

	roleID := role.GetCID()

	userRoles := []permission.UserRole{}
	dbHandler := db.DBHanlder()

	//todo: use one query with join
	//todo: support order, pagnation params
	dbHandler.GetEntity("dm_user_role", db.Cond("role_id", roleID), nil, nil, &userRoles)
	userIDs := []int{}
	for _, userRole := range userRoles {
		userIDs = append(userIDs, userRole.UserID)
	}

	list, count, _ := querier.List("user", db.Cond("c.id", userIDs), nil, nil, true)

	data, _ := json.Marshal(map[string]interface{}{"list": list, "count": count})
	w.Write(data)
}

func AssignUser(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	if !permission.HasAccessTo(userID, "access/assign-user", permission.MatchData{}, r.Context()) {
		HandleError(errors.New("No access"), w)
		return
	}

	params := mux.Vars(r)
	roleID, _ := strconv.Atoi(params["role"])
	assignedUserID, _ := strconv.Atoi(params["user"])
	err := permission.AssignToUser(roleID, assignedUserID)

	if err != nil {
		log.Error("Error when assigning: "+err.Error(), "")
		HandleError(errors.New("Error when assigning"), w, 400)
		return
	}
	w.Write([]byte("1"))
}

func init() {
	RegisterRoute("/access/update-fields/current-user", CurrentUserEditField)
	RegisterRoute("/access/assigned-users", AssignedUsers)
	RegisterRoute("/access/assign/{role}/{user}", AssignUser)
}
