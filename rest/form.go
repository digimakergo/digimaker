//Author xc, Created on 2019-10-11 17:07
//{COPYRIGHTS}

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/handler"

	"github.com/gorilla/mux"
)

func Validate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contentType := params["contenttype"]
	//todo: add permission check for which form/container can be.
	contentType = strings.ReplaceAll(contentType, ",", "/")
	fmt.Println(len(strings.Split(contentType, "/")))
	if len(strings.Split(contentType, "/")) > 2 {
		w.WriteHeader(400)
		w.Write([]byte("Only support one level container."))
		return
	}

	inputs := map[string]interface{}{}
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&inputs)
	if err != nil {
		HandleError(err, w)
		return
	}

	fieldMap, err := definition.GetFields(contentType)
	if err != nil {
		HandleError(err, w)
		return
	}

	arr := strings.Split(contentType, "/")
	result, validationResult := handler.Validate(arr[0], fieldMap, inputs, true)
	if result {
		w.Write([]byte("1"))
	} else {
		w.WriteHeader(400)
		data, _ := json.Marshal(validationResult)
		w.Write(data)
	}
}

func init() {
	RegisterRoute("/form/validate/{contenttype}", Validate)
}
