//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}

package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"

	"github.com/gorilla/mux"
)

func GetContent(w http.ResponseWriter, r *http.Request) {
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
		data, _ := contenttype.ContentToJson(content) //todo: use export for same serilization?
		w.Write([]byte(data))
	}

}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

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

	if !permission.CanRead(userID, content, r.Context()) {
		HandleError(errors.New("No permisison to the content."), w)
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

func buildCondition(userid int, def contenttype.ContentType, query url.Values) (db.Condition, error) {
	author := query.Get("author")
	condition := db.EmptyCond()
	if author != "" {
		if author == "self" {
			condition = condition.Cond("author", userid)
		} else {
			authorInt, err := strconv.Atoi(author)
			if err != nil {
				return db.EmptyCond(), errors.New("wrong author format")
			}
			condition = condition.Cond("author", authorInt)
		}
	}

	//id
	idStr := query.Get("id")
	if idStr != "" {
		ids, err := util.ArrayStrToInt(strings.Split(idStr, ","))
		if err != nil {
			return db.EmptyCond(), errors.New("Wrong id format")
		}
		condition = condition.And("location.id", ids)
	}

	//cid
	cidStr := query.Get("cid")
	if cidStr != "" {
		cids, err := util.ArrayStrToInt(strings.Split(cidStr, ","))
		if err != nil {
			return db.EmptyCond(), errors.New("Wrong cid format")
		}
		condition = condition.And("c.id", cids)
	}

	for field := range def.FieldMap {
		value := query.Get("field." + field)
		if value != "" {
			condition = condition.And("c."+field, value) //todo: support operator //todo: support more value type(eg. array)
		}
	}

	return condition, nil
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
	ctype := params["contenttype"]
	def, err := contenttype.GetDefinition(ctype)
	if err != nil {
		HandleError(errors.New("Cann't get content type"), w)
		return
	}

	querier := handler.Querier()
	rootContent, err := querier.FetchByID(id)
	if err != nil {
		//todo: handle
	}
	cxt := r.Context()
	userid := CheckUserID(cxt, w)
	if userid == 0 {
		return
	}

	//filter
	condition, err := buildCondition(userid, def, r.URL.Query())
	if err != nil {
		HandleError(err, w, 410)
		return
	}

	limitArr := []int{}
	if offsetStr != "" && limitStr != "" {
		limitArr = []int{offset, limit}
	}

	list, count, err := querier.Children(rootContent, ctype, userid, condition, limitArr, sortbyArr, true, cxt)
	// list, count, err := querier.SubList(rootContent, ctype, 100, userid, condition, limitArr, sortbyArr, true, cxt)
	if err != nil {
		HandleError(err, w)
		return
	}

	result := struct {
		List  interface{} `json:"list"`
		Count int         `json:"count"`
	}{Count: count}

	configFields := util.GetConfigArr("rest_list_fields", ctype)
	if configFields != nil {
		//output needed fields
		outputList := []map[string]interface{}{}
		for _, content := range list {
			//get a map based content
			outputContent, err := contenttype.ContentToMap(content)
			if err != nil {
				log.Error("Marshall content error: "+err.Error(), "", cxt)
				HandleError(errors.New("Error when converting data."), w)
				return
			}

			for _, field := range content.IdentifierList() {
				if !util.Contains(configFields, field) {
					delete(outputContent, field)
				}
			}
			outputList = append(outputList, outputContent)
		}
		result.List = outputList
	} else {
		result.List = list
	}

	data, _ := json.Marshal(result)
	w.Write([]byte(data))
}

//List
//todo: merge with Children/allback
func List(w http.ResponseWriter, r *http.Request) {
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

	ctype := params["contenttype"]
	def, err := contenttype.GetDefinition(ctype)
	if err != nil {
		HandleError(errors.New("Cann't get content type"), w)
		return
	}

	querier := handler.Querier()
	rootContent, err := querier.FetchByID(3)
	if err != nil {
		//todo: handle
	}
	cxt := r.Context()
	userid := CheckUserID(cxt, w)
	if userid == 0 {
		return
	}

	//filter
	condition, err := buildCondition(userid, def, r.URL.Query())
	if err != nil {
		HandleError(err, w, 410)
		return
	}

	limitArr := []int{}
	if offsetStr != "" && limitStr != "" {
		limitArr = []int{offset, limit}
	}

	list, count, err := querier.SubList(rootContent, ctype, 0, userid, condition, limitArr, sortbyArr, true, cxt)
	if err != nil {
		HandleError(err, w)
		return
	}

	result := struct {
		List  interface{} `json:"list"`
		Count int         `json:"count"`
	}{Count: count}

	result.List = list

	data, _ := json.Marshal(result)
	w.Write([]byte(data))
}

//Get tree menu under a node
func TreeMenu(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	types := r.FormValue("type")
	if types == "" {
		HandleError(errors.New("Wrong parameters"), w, StatusWrongParams)
		return
	}
	typeList := util.Split(types)
	contenttypeList := contenttype.GetDefinitionList()["default"]
	for _, ctype := range typeList {
		if _, ok := contenttypeList[ctype]; !ok {
			HandleError(errors.New("Invalid type"), w, StatusWrongParams)
			return
		}
	}

	querier := handler.Querier()
	rootContent, err := querier.FetchByID(id)

	if err != nil {
		HandleError(err, w)
		return
	}

	if rootContent == nil {
		HandleError(errors.New("content Not found"), w, 403)
		return
	}

	if !permission.CanRead(userID, rootContent, r.Context()) {
		HandleError(errors.New("No permission"), w)
		return
	}

	tree, err := querier.SubTree(rootContent, 5, strings.Join(typeList, ","), userID, []string{"id"}, r.Context())
	if err != nil {
		HandleError(err, w)
		return
	}

	//todo: make this configurable
	tree.Iterate(func(node *handler.TreeNode) {
		if node.ContentType == "folder" {
			node.Fields = map[string]interface{}{"subtype": node.Content.Value("folder_type")}
		}
	})

	data, _ := json.Marshal(tree)
	w.Write([]byte(data))
}

func init() {
	RegisterRoute("/content/get/{id:[0-9]+}", GetContent)
	RegisterRoute("/content/version/{id:[0-9]+}/{version:[0-9]+}", GetVersion)

	RegisterRoute("/content/treemenu/{id:[0-9]+}", TreeMenu)
	RegisterRoute("/content/list/{id:[0-9]+}/{contenttype}", Children)
	RegisterRoute("/content/list/{contenttype}", List)
}
