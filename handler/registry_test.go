//Author xc, Created on 2019-05-11 13:28
//{COPYRIGHTS}
package handler

import (
	"dm/contenttype"
	"fmt"
	"testing"
)

type TestOperationHandler struct {
}

func TestRegistry(t *testing.T) {
	condition := map[string]interface{}{"id": 12, "type": "image"}
	testHandler := OperationHandler{Identifier: "test_handler",
		Event: "change", Execute: func(content contenttype.ContentTyper) error {
			return nil
		}}
	RegisterOperationHandler(testHandler)
	handlers, log := GetOperationHandlerByCondition("change", condition)
	for _, handler := range handlers {
		fmt.Println(handler.Identifier)
	}
	for _, message := range log {
		fmt.Println(message)
	}
}
