//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}

package rest

import (
	"context"
	"dm/core/contenttype"
	"dm/core/db"
	"dm/core/handler"
	"dm/core/permission"
	"dm/core/util"
	"dm/core/util/debug"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetContent(w http.ResponseWriter, r *http.Request) {
	ctx := debug.Init(r.Context())
	r = r.WithContext(ctx)

	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	querier := handler.Querier()
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		HandleError(errors.New("Invalid id"), w)
		return
	}
	content, err := querier.FetchByID(id)
	if err != nil {
		HandleError(err, w)
		return
	} else {
		if !permission.CanRead(userID, content, r.Context()) {
			HandleError(errors.New("Doesn't have permission."), w, 403)
			return
		}
		w.Header().Set("content-type", "application/json")
		data, _ := json.Marshal(content) //todo: use export for same serilization?
		w.Write(data)
	}
}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	querier := handler.Querier()
	id, _ := strconv.Atoi(params["id"])

	versionNo, _ := strconv.Atoi(params["version"])

	querier = handler.Querier()
	content, err := querier.FetchByID(id)
	if err != nil {
		HandleError(errors.New("Can not find content. error: "+err.Error()), w)
		return
	}

	if content == nil {
		HandleError(errors.New("Content doesn't exist"), w)
		return
	}

	maxVersion := content.Value("version").(int)

	if versionNo > maxVersion {
		HandleError(errors.New("version doesn't exist."), w)
		return
	}

	dbHandler := db.DBHanlder()
	version := contenttype.Version{}
	dbHandler.GetEntity(version.TableName(),
		db.Cond("content_id", content.GetCID()).Cond("content_type", content.ContentType()).Cond("version", versionNo),
		[]string{},
		&version)
	if version.ID == 0 {
		HandleError(errors.New("version doesn't exist."), w)
		return
	}

	data, _ := json.Marshal(version)
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(data))
}

//Get children of a content(eg. folder)
func Children(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	getParams := r.URL.Query()

	//offset and limit
	offsetStr := getParams.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if offsetStr != "" && err != nil {
		HandleError(errors.New("Invalid offset"), w)
		return
	}

	limitStr := getParams.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if limitStr != "" && err != nil {
		HandleError(errors.New("Invalid limit"), w)
		return
	}

	//sort by
	sortbyStr := getParams.Get("sortby")
	sortbyArr := util.Split(sortbyStr, ";")

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		HandleError(errors.New("Invalid id"), w)
		return
	}
	contenttype := params["contenttype"]
	querier := handler.Querier()
	rootContent, err := querier.FetchByID(id)
	if err != nil {
		//todo: handle
	}
	context := debug.Init(context.Background())
	cxt := r.Context()
	if cxt.Value("user_id") == nil {
		HandleError(errors.New("No user found"), w, 410)
		return
	}
	userid := cxt.Value("user_id").(int)

	//filter
	author := getParams.Get("author")
	condition := db.Cond("1", "1")
	if author != "" {
		if author == "self" {
			condition = condition.Cond("location.author", userid)
		} else {
			authorInt, err := strconv.Atoi(author)
			if err != nil {
				HandleError(errors.New("wrong author format"), w, 410)
				return
			}
			condition = condition.Cond("location.author", authorInt)
		}
	}
	//todo: add more filters including field filter.

	limitArr := []int{}
	if offsetStr != "" && limitStr != "" {
		limitArr = []int{offset, limit}
	}

	list, count, err := querier.Children(rootContent, contenttype, userid, condition, limitArr, sortbyArr, true, context)
	if err != nil {
		HandleError(err, w)
		return
	}

	result := struct {
		List  interface{} `json:"list"`
		Count int         `json:"count"`
	}{list, count}

	data, _ := json.Marshal(result)
	w.Write([]byte(data))
}

func SubTree() {

}

//Get tree menu under a node
func TreeMenu(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	querier := handler.Querier()
	rootContent, err := querier.FetchByID(id)
	if err != nil {
		//todo: handle
	}

	context := debug.Init(context.Background())
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}
	tree, err := querier.SubTree(rootContent, 5, "folder,role,usergroup", userID, []string{"id"}, context)
	if err != nil {
		//todo: handle error
		fmt.Println(err.Error())
	}

	data, _ := json.Marshal(tree)
	w.Write([]byte(data))
}

func init() {
	RegisterRoute("/content/get/{id:[0-9]+}", GetContent)
	RegisterRoute("/content/get/{id:[0-9]+}/{version:[0-9]+}", GetVersion)

	RegisterRoute("/content/treemenu/{id:[0-9]+}", TreeMenu)
	RegisterRoute("/content/list/{id:[0-9]+}/{contenttype}", Children)
}
