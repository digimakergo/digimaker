//Author xc, Created on 2019-05-11 13:28
//{COPYRIGHTS}
package handler

import (
	"fmt"
	"testing"
)

type TestOperationHandler struct {
}

func TestRegistry(t *testing.T) {
	condition := map[string]interface{}{"id": 12, "type": "image"}
	testHandler := OperationHandler{Identifier: "test_handler",
		Execute: func(event string, params ...interface{}) error {
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
