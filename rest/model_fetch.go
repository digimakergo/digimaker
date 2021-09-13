//Author xc, Created on 2019-08-13 17:25
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"

	"github.com/gorilla/mux"
)

//todo: check permission
func GetDefinition(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	typeStr := strings.TrimSpace(params["contentype"])
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "default"
	}

	containers := strings.Split(typeStr, "/")
	definition, _ := definition.GetDefinition(containers[0], language)

	resultMap := filterDefinition(definition)

	WriteResponse(resultMap, w)
}

func filterDefinition(definition definition.ContentType) map[string]interface{} {
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
	return resultMap
}

func GetAllDefinitions(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	if language == "" {
		language = "default"
	}

	definitionList := definition.GetDefinitionList()
	list := definitionList[language]
	result := map[string]interface{}{}
	for contenttype, definition := range list {
		resultMap := filterDefinition(definition)
		result[contenttype] = resultMap
	}

	WriteResponse(result, w)
}

func init() {
	RegisterRoute("/contenttype/get", GetAllDefinitions)
	RegisterRoute("/contenttype/get/{contentype}", GetDefinition)
}
