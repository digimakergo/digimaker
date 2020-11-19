//Author xc, Created on 2019-05-11 22:04
//{COPYRIGHTS}
package handlers

import (
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/handler"
	"fmt"
)

func init() {
	oHandller := handler.OperationHandler{Identifier: "test_handler",
		Execute: func(triggedEvent string, content contenttype.ContentTyper, params ...interface{}) error {
			fmt.Println("test handler invoked. trigger: " + triggedEvent)
			return nil
		}}
	handler.RegisterOperationHandler(oHandller)
}
