//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}

package rest

import (
	"context"
	"dm/core/handler"
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

func Children(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id, err := strconv.Atoi(params["id"])
	querier := handler.Querier()
	rootContent, err := querier.FetchByID(id)
	if err != nil {
		//todo: handle
	}
	context := debug.Init(context.Background())
	list, err := querier.Children(rootContent, 1, context)
	if err != nil {

	}
	data, _ := json.Marshal(list)
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
	tree, err := querier.SubTree(rootContent, 5, "folder", 1, context)
	if err != nil {
		//todo: handle error
		fmt.Println(err.Error())
	}

	data, _ := json.Marshal(tree)
	w.Write([]byte(data))
}
