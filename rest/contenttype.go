//Author xc, Created on 2019-08-13 17:25
//{COPYRIGHTS}
package rest

import (
	"dm/core/contenttype"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func GetDefinition(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	identifier := strings.TrimSpace(params["contentype"])
	w.Header().Set("Access-Control-Allow-Origin", "*")

	definition := contenttype.GetContentDefinition(identifier)
	w.Header().Set("content-type", "application/json")
	data, _ := json.Marshal(definition)
	w.Write(data)
}
