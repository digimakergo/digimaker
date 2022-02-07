//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}

package rest

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

//Result item in the list
type ResultItem map[string]interface{}

//Result list
type ResultList []ResultItem

func GetContent(w http.ResponseWriter, r *http.Request) {
	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	params := mux.Vars(r)
	var content contenttype.ContentTyper
	contentType := params["contenttype"]
	var err error
	if contentType != "" {
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			if len(params["id"]) != 20 {
				HandleError(errors.New("Not valid cuid"), w)
				return
			}
			content, err = query.FetchByCUID(r.Context(), contentType, params["id"])
		} else {
			content, err = query.FetchByCID(r.Context(), contentType, id)
		}
	} else {
		if params["id"] != "" {
			id, err := strconv.Atoi(params["id"])
			if err != nil {
				HandleError(errors.New("Invalid id"), w)
				return
			}
			content, err = query.FetchByID(r.Context(), id)
		} else {
			path := r.URL.Query().Get("path")
			if path != "" {
				content, err = query.FetchByPath(r.Context(), path)
			} else {
				err = errors.New("Parameter not supported")
			}
		}
	}
	if err != nil {
		HandleError(err, w)
		return
	}

	if content != nil && !permission.CanRead(r.Context(), userID, content) {
		HandleError(errors.New("Doesn't have permission."), w, 403)
		return
	}

	outputContent, err := query.Output(r.Context(), content)
	if err != nil {
		HandleError(err, w)
		return
	}
	WriteResponse(outputContent, w)
}

func GetVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID := CheckUserID(r.Context(), w)
	if userID == 0 {
		return
	}

	id, _ := strconv.Atoi(params["id"])

	versionNo, _ := strconv.Atoi(params["version"])

	content, err := query.FetchByID(r.Context(), id)
	if err != nil {
		HandleError(errors.New("Can not find content. error: "+err.Error()), w)
		return
	}

	if content == nil {
		HandleError(errors.New("Content doesn't exist"), w)
		return
	}

	if !permission.CanRead(r.Context(), userID, content) {
		HandleError(errors.New("No permisison to the content."), w)
		return
	}

	maxVersion := content.Value("version").(int)

	if versionNo > maxVersion {
		HandleError(errors.New("version doesn't exist."), w)
		return
	}

	version := contenttype.Version{}
	db.BindEntity(r.Context(),
		&version,
		version.TableName(),
		db.Cond("content_id", content.GetCID()).Cond("content_type", content.ContentType()).Cond("version", versionNo))
	if version.ID == 0 {
		HandleError(errors.New("version doesn't exist."), w)
		return
	}

	WriteResponse(version, w)
}

func buildCondition(userid int, cType string, def definition.ContentType, query url.Values) (db.Condition, error) {
	author := query.Get("author")
	condition := db.EmptyCond()
	if author != "" {
		if author == "self" {
			condition = condition.Cond("c.author", userid)
		} else {
			authorInt, err := strconv.Atoi(author)
			if err != nil {
				return db.EmptyCond(), errors.New("wrong author format")
			}
			condition = condition.Cond("c.author", authorInt)
		}
	}

	//id
	idStr := query.Get("id")
	if idStr != "" {
		ids, err := util.ArrayStrToInt(strings.Split(idStr, ","))
		if err != nil {
			return db.EmptyCond(), errors.New("Wrong id format")
		}
		condition = condition.And("l.id", ids)
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

	//contain. todo: merge with filter below
	for key, _ := range query {
		value := query.Get(key)
		if strings.HasPrefix(value, "contain:") {
			enabledFields := viper.GetStringSlice("rest.like_fields")
			if util.Contains(enabledFields, cType+"/"+key) {
				fieldValue := strings.TrimPrefix(value, "contain:")
				if def.HasLocation && util.Contains(definition.LocationColumns, key) {
					condition = condition.And("l."+key+" like", fieldValue)
				} else {
					condition = condition.And("c."+key+" like", fieldValue)
				}
			}
		}
	}

	//filter
	for field := range def.FieldMap {
		value := query.Get("field." + field)
		if value != "" {
			operation, val := extraOperation(value)
			condition = condition.And("c."+field+" "+operation, val)
		}
	}

	return condition, nil
}

//extract operation and return operation+value
func extraOperation(param string) (string, interface{}) {
	//todo: support operator //todo: support more value type(eg. array)
	if strings.HasPrefix(param, "contain:") {
		value := strings.TrimPrefix(param, "contain:")
		return "like", "%" + value + "%"
	}
	return "", param
}

func BuildSortby(r *http.Request) []string {
	getParams := r.URL.Query()
	sortbyStr := getParams.Get("sortby")
	sortbyArr := util.Split(sortbyStr, ";")
	return sortbyArr
}

func BuildLimit(r *http.Request) ([]int, error) {
	getParams := r.URL.Query()
	offsetStr := getParams.Get("offset")
	limitStr := getParams.Get("limit")

	//if there is no offset&limit, use default limit
	if offsetStr == "" && limitStr == "" {
		return []int{0, 10}, nil //todo: use configuration
	}

	offset, err := strconv.Atoi(offsetStr)
	if offsetStr != "" && err != nil {
		return nil, errors.New("Invalid offset")
	}

	limit, err := strconv.Atoi(limitStr)
	if limitStr != "" && err != nil {
		return nil, errors.New("Invalid limit")
	}

	//if limit is 0, set no limit
	if limit == 0 {
		return []int{offset, 1000000}, nil //max one fetch without limit is 1000000
	}

	return []int{offset, limit}, nil
}

//List
func List(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	getParams := r.URL.Query()

	limit, err := BuildLimit(r)
	if err != nil {
		HandleError(err, w)
		return
	}

	sortby := BuildSortby(r)

	ctype := params["contenttype"]
	def, err := definition.GetDefinition(ctype)
	if err != nil {
		HandleError(errors.New("Can't get content type"), w)
		return
	}

	ctx := r.Context()
	userID := CheckUserID(ctx, w)
	if userID == 0 {
		return
	}

	//filter
	condition, err := buildCondition(userID, ctype, def, r.URL.Query())
	if err != nil {
		HandleError(err, w, 410)
		return
	}

	//sort by
	condition = condition.Sortby(sortby...)
	if len(limit) == 2 {
		condition = condition.Limit(limit[0], limit[1])
	}

	rootStr := getParams.Get("parent")
	var list []contenttype.ContentTyper
	var count int
	if rootStr != "" {
		var rootContent contenttype.ContentTyper
		rootID, err := strconv.Atoi(rootStr)
		if err != nil {
			HandleError(errors.New("Invalid parent"), w)
			return
		}
		rootContent, err = query.FetchByID(r.Context(), rootID)
		if err != nil {
			log.Error(err.Error(), "", r.Context())
			HandleError(errors.New("Can't get parent"), w, 410)
			return
		}

		//level
		levelStr := getParams.Get("level")
		level := 0
		if levelStr != "" {
			level, err = strconv.Atoi(levelStr)
			if err != nil {
				HandleError(errors.New("Invalid level"), w)
				return
			}
		}

		list, count, err = query.SubList(ctx, userID, rootContent, ctype, level, condition.WithCount())
		if err != nil {
			HandleError(err, w)
			return
		}
	} else {
		list, count, err = query.ListWithUser(ctx, userID, ctype, condition.WithCount())
		if err != nil {
			HandleError(err, w)
			return
		}
	}
	result := ResultItem{}
	outputList, err := query.OutputList(ctx, list)
	if err != nil {
		HandleError(err, w)
		return
	}
	result["list"] = outputList
	result["count"] = count

	//columns
	withColumns := r.URL.Query().Get("columns") == "true"
	if withColumns {
		columns := getColumns(ctype)
		result["columns"] = columns
	}

	WriteResponse(result, w)
}

//Get content list from relation definition
func RelationOptionList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contenttype := params["contenttype"]
	fieldIdentifier := params["field"]
	if contenttype == "" || fieldIdentifier == "" {
		HandleError(errors.New("Need type and field"), w, 410)
		return
	}

	ctx := r.Context()
	userid := CheckUserID(ctx, w)
	if userid == 0 {
		return
	}

	list, count, err := query.RelationOptions(ctx, contenttype, fieldIdentifier, nil, nil, false)
	if err != nil {
		HandleError(err, w)
		return
	}

	result := struct {
		List  interface{} `json:"list"`
		Count int         `json:"count"`
	}{Count: count}

	result.List = list

	WriteResponse(result, w)
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
	contenttypeList := definition.GetDefinitionList()["default"]
	for _, ctype := range typeList {
		if _, ok := contenttypeList[ctype]; !ok {
			HandleError(errors.New("Invalid type"), w, StatusWrongParams)
			return
		}
	}

	rootContent, err := query.FetchByID(r.Context(), id)

	if err != nil {
		HandleError(err, w)
		return
	}

	if rootContent == nil {
		HandleError(errors.New("content Not found"), w, 403)
		return
	}

	if !permission.CanRead(r.Context(), userID, rootContent) {
		HandleError(errors.New("No permission"), w)
		return
	}

	tree, err := query.SubTree(r.Context(), userID, rootContent, 5, strings.Join(typeList, ","), []string{"priority desc", "id"})
	if err != nil {
		HandleError(err, w)
		return
	}

	//todo: make this configurable
	tree.Iterate(func(node *query.TreeNode) {
		if node.ContentType == "folder" {
			node.Fields = map[string]interface{}{"subtype": node.Content.Value("folder_type")}
		}
	})

	WriteResponse(tree, w)
}

func init() {
	RegisterRoute("/content/get/{id:[0-9]+}", GetContent)
	RegisterRoute("/content/get", GetContent)
	RegisterRoute("/content/get/{contenttype}/{id}", GetContent)
	RegisterRoute("/content/version/{id:[0-9]+}/{version:[0-9]+}", GetVersion)

	RegisterRoute("/content/treemenu/{id:[0-9]+}", TreeMenu)
	RegisterRoute("/content/list/{contenttype}", List)
	RegisterRoute("/relation/optionlist/{contenttype}/{field}", RelationOptionList)
}
