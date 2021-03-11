//Author xc, Created on 2019-05-11 22:04
//{COPYRIGHTS}
package handlers

import (
	"context"
	"fmt"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/handler"
)

func init() {
	oHandller := handler.OperationHandler{Identifier: "test_handler",
		Execute: func(ctx context.Context, triggedEvent string, content contenttype.ContentTyper, params ...interface{}) error {
			fmt.Println("test handler invoked. trigger: " + triggedEvent)
			return nil
		}}
	handler.RegisterOperationHandler(oHandller)
}
