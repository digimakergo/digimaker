//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}

package rest

import (
	"context"
	"dm/core/handler"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
		w.Header().Set("content-type", "application/json")
		data, _ := json.Marshal(content) //todo: use export for same serilization?
		w.Write(data)
	}
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	list, count, err := querier.Children(rootContent, contenttype, 1, []int{offset, limit}, sortbyArr, true, context)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	querier := handler.Querier()
	rootContent, err := querier.FetchByID(id)
	if err != nil {
		//todo: handle
	}

	context := debug.Init(context.Background())
	tree, err := querier.SubTree(rootContent, 5, "folder,usergroup", 1, []string{"id"}, context)
	if err != nil {
		//todo: handle error
		fmt.Println(err.Error())
	}

	data, _ := json.Marshal(tree)
	w.Write([]byte(data))
}

func init() {
	RegisterRoute("/contenttype/get/{contentype}", GetDefinition)
	RegisterRoute("/content/get/{id}", GetContent)

	RegisterRoute("/content/treemenu/{id}", TreeMenu)
	RegisterRoute("/content/list/{id}/{contenttype}", Children)
}
