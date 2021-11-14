//Author xc, Created on 2020-10-1 12:50
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/gorilla/mux"
)

var ACCESS_MANAGE_ACCESS = "access/manage"

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

	WriteResponse(fields, w)
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

	data := map[string]interface{}{"list": resultList, "count": count}
	WriteResponse(data, w)
}

func UserRoles(w http.ResponseWriter, r *http.Request) {
	currentUserID := CheckUserID(r.Context(), w)
	if currentUserID == 0 {
		return
	}
	//todo: check permission
	params := mux.Vars(r)
	userID, _ := strconv.Atoi(params["user"])
	list, err := query.FetchUserRoles(r.Context(), userID)
	if err != nil {
		HandleError(err, w)
		return
	}
	WriteResponse(list, w)
}

func AssignCreateRole(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	assignedUserID, _ := strconv.Atoi(params["user"])

	assignParams := struct {
		Parameters fieldtypes.Map `json:"parameters"`
		Role       string         `json:"role"`
		Title      string         `json:"title"`
	}{}

	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&assignParams)
	if err != nil {
		HandleError(errors.New("Assignment parameters wrong format: "+err.Error()), w, 400)
		return
	}

	if assignParams.Role == "" || assignParams.Title == "" {
		HandleError(errors.New("Please input role & title"), w, 400)
		return
	}

	if !permission.HasAccessTo(r.Context(), userID, "access/assign-user", permission.MatchData{}) {
		HandleError(errors.New("No access"), w)
		return
	}

	input := handler.InputMap{}
	input["title"] = assignParams.Title
	input["identifier"] = assignParams.Role
	input["parameters"] = assignParams.Parameters
	role, createErr := handler.Create(r.Context(), userID, "role", input, 7) //todo: make parent id as optional
	if createErr != nil {
		log.Error(createErr.Error(), "permission")
		HandleError(err, w)
		return
	}
	err = permission.AssignToUser(r.Context(), role.GetCID(), assignedUserID)

	if err != nil {
		HandleError(errors.New("Error when assigning: "+err.Error()), w, 400)
		return
	}
	WriteResponse(true, w)
}

func AssignUserRoles(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	assignedUserID, _ := strconv.Atoi(params["user"])

	if !permission.HasAccessTo(r.Context(), userID, "access/assign-user", permission.MatchData{}) {
		HandleError(errors.New("No access"), w)
		return
	}

	roles := strings.Split(params["roles"], ",")
	roleIDs := []int{}
	for _, roleStr := range roles {
		roleID, _ := strconv.Atoi(roleStr)
		roleContent, _ := query.FetchByCID(r.Context(), "role", roleID)
		if roleContent == nil {
			HandleError(errors.New("Role not found "+roleStr), w, 400)
			return
		}
		roleIDs = append(roleIDs, roleID)
	}

	for _, roleID := range roleIDs {
		err := permission.AssignToUser(r.Context(), roleID, assignedUserID)
		if err != nil {
			HandleError(errors.New("Error when assigning: "+err.Error()), w, 400)
			return
		}
	}

	WriteResponse(true, w)
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
	role := params["role"]

	err := permission.RemoveAssignment(r.Context(), userID, role)
	if err != nil {
		HandleError(err, w)
		return
	}
	WriteResponse(true, w)
}

func AllRoles(w http.ResponseWriter, r *http.Request) {
	//todo: add check
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}
	roles := permission.GetRoles()
	WriteResponse(roles, w)
}

func init() {
	RegisterRoute("/access/update-fields/current-user", CurrentUserEditField)
	RegisterRoute("/access/assigned-users", AssignedUsers)
	RegisterRoute("/access/roles/{user}", UserRoles)
	RegisterRoute("/access/allroles", AllRoles)
	RegisterRoute("/access/assign/{user}/{roles}", AssignUserRoles)
	RegisterRoute("/access/assign-create/{user}", AssignCreateRole, "POST")
	RegisterRoute("/access/unassign/{user}/{role}", UnassignUser)
}
