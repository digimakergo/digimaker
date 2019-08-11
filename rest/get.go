//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}

package rest

import (
	"dm/dm/handler"
	"encoding/json"
	"net/http"
)

func GetContent(id int, w http.ResponseWriter) {
	querier := handler.Querier()
	content, err := querier.FetchByID(id)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("content-type", "application/json")
		data, _ := json.Marshal(content) //todo: use export for same serilization?
		w.Write(data)
	}

}

func Children() {

}

func SubTree() {

}
