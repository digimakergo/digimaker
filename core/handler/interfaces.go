//Author xc, Created on 2019-05-11 12:33
//{COPYRIGHTS}
package handler

import (
	"context"

	"github.com/digimakergo/digimaker/core/contenttype"
)

//todo: think it might be good to use "struct with callback method" instead of interface.
//This is a callback based on type.
//It's used for customzed content type. eg. image to set parent_id
type ContentTypeHandlerValidate interface {
	ValidateCreate(ctx context.Context, inputs InputMap, parentID int) (bool, ValidationResult)
	ValidateUpdate(ctx context.Context, inputs InputMap, content contenttype.ContentTyper) (bool, ValidationResult)
}

type ContentTypeHandlerCreate interface {
	//When creating on server side, or import.
	Create(ctx context.Context, content contenttype.ContentTyper, inputs InputMap, parentID int) error
}

type ContentTypeHandlerUpdate interface {
	// //When after edit, the server handles the update
	Update(ctx context.Context, content contenttype.ContentTyper, inputs InputMap) error
}

type ContentTypeHandlerDelete interface {
	//
	// //when deleting
	Delete(ctx context.Context, content contenttype.ContentTyper) error
}

//Callback struct
type OperationHandler struct {
	Identifier string //Identifier for handler matching. see operation_handler.json/yaml
	//todo: use translation as parameter instead of optional interface{}
	Execute func(ctx context.Context, triggedEvent string, content contenttype.ContentTyper, params ...interface{}) error
}
