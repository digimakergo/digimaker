//Author xc, Created on 2020-10-1 12:50
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/gorilla/mux"
)

//Get current user's updatefields on him/herself
func CurrentUserEditField(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}
	content, _ := query.GetUser(userID)

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

	if !permission.HasAccessTo(r.Context(), userID, "access/assigned-users", permission.MatchData{}) {
		HandleError(errors.New("No access"), w)
		return
	}

	parentID, err := strconv.Atoi(r.FormValue("parent"))
	if err != nil {
		HandleError(err, w)
		return
	}

	role, _ := query.FetchByID(r.Context(), parentID)
	if role == nil {
		HandleError(errors.New("Roles doesn't exist"), w)
		return
	}

	roleID := role.GetCID()

	userRoles := []permission.UserRole{}

	//todo: use one query with join
	//todo: support order, pagnation params
	db.BindEntity(r.Context(), &userRoles, "dm_user_role", db.Cond("role_id", roleID))
	userIDs := []int{}
	for _, userRole := range userRoles {
		userIDs = append(userIDs, userRole.UserID)
	}

	list, count, _ := query.List(r.Context(), "user", db.Cond("c.id", userIDs))

	resultList := []interface{}{}
	for _, item := range list {
		cmap, _ := contenttype.ContentToMap(item)
		cmap["role_id"] = roleID
		resultList = append(resultList, cmap)
	}

	data, _ := json.Marshal(map[string]interface{}{"list": resultList, "count": count})
	w.Write(data)
}

func AssignUser(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	roleID, _ := strconv.Atoi(params["role"])
	assignedUserID, _ := strconv.Atoi(params["user"])

	if !permission.HasAccessTo(r.Context(), userID, "access/assign-user", permission.MatchData{}) {
		HandleError(errors.New("No access"), w)
		return
	}

	role, _ := query.FetchByID(r.Context(), roleID)
	if role == nil {
		HandleError(errors.New("Role not found"), w, 400)
		return
	}

	err := permission.AssignToUser(r.Context(), role.GetCID(), assignedUserID)

	if err != nil {
		log.Error("Error when assigning: "+err.Error(), "")
		HandleError(errors.New("Error when assigning"), w, 400)
		return
	}
	w.Write([]byte("1"))
}

//unassign user from role
func UnassignUser(w http.ResponseWriter, r *http.Request) {
	loginUserID := CheckUserID(r.Context(), w)
	if loginUserID == 0 {
		return
	}

	//todo: move all this to handler
	if !permission.HasAccessTo(r.Context(), loginUserID, "access/unassign-user", permission.MatchData{}) {
		HandleError(errors.New("No access"), w)
		return
	}

	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["user"])
	roleID, _ := strconv.Atoi(params["role"])

	err := permission.RemoveAssignment(r.Context(), userID, roleID)
	if err != nil {
		HandleError(err, w)
		return
	}
	w.Write([]byte("1"))
}

func init() {
	RegisterRoute("/access/update-fields/current-user", CurrentUserEditField)
	RegisterRoute("/access/assigned-users", AssignedUsers)
	RegisterRoute("/access/assign/{role}/{user}", AssignUser)
	RegisterRoute("/access/unassign/{user}/{role}", UnassignUser)
}
