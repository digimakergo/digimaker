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
	typeStr := strings.TrimSpace(params["contentype"])

	containers := strings.Split(typeStr, "/")
	definition, _ := contenttype.GetDefinition(containers[0])

	data, _ := json.Marshal(definition)
	resultMap := map[string]interface{}{}
	json.Unmarshal(data, &resultMap)
	delete(resultMap, "table_name")
	delete(resultMap, "has_location")
	delete(resultMap, "has_version")
	result, _ := json.Marshal(resultMap)

	w.Header().Set("content-type", "application/json")

	w.Write(result)
}
