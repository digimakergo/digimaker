//Author xc, Created on 2019-08-13 17:25
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/xc/digimaker/core/contenttype"

	"github.com/gorilla/mux"
)

func GetDefinition(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	typeStr := strings.TrimSpace(params["contentype"])
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "default"
	}

	containers := strings.Split(typeStr, "/")
	definition, _ := contenttype.GetDefinition(containers[0], language)

	data, _ := json.Marshal(definition)
	resultMap := map[string]interface{}{}
	json.Unmarshal(data, &resultMap)
	fields := resultMap["fields"].([]interface{})
	dataFieldsObj := resultMap["data_fields"]
	if dataFieldsObj != nil {
		dataFields := dataFieldsObj.([]interface{})
		fields = append(fields, dataFields...)
	}
	resultMap["fields"] = fields
	delete(resultMap, "table_name")
	delete(resultMap, "has_version")
	delete(resultMap, "data_fields")
	result, _ := json.Marshal(resultMap)

	w.Header().Set("content-type", "application/json")

	w.Write(result)
}

func init() {
	RegisterRoute("/contenttype/get/{contentype}", GetDefinition)
}
