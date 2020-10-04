//Author xc, Created on 2019-05-11 12:33
//{COPYRIGHTS}
package handler

import (
	"github.com/xc/digimaker/core/contenttype"
)

//todo: think it might be good to use "struct with callback method" instead of interface.
//This is a callback based on type.
//It's used for customzed content type. eg. image to set parent_id
type ContentTypeHandlerValidate interface {
	ValidateCreate(inputs InputMap, parentID int) (bool, ValidationResult)
	ValidateUpdate(inputs InputMap, content contenttype.ContentTyper) (bool, ValidationResult)
}

type ContentTypeHandlerCreate interface {
	//When creating on server side, or import.
	Create(content contenttype.ContentTyper, inputs InputMap, parentID int) error
}

type ContentTypeHandlerUpdate interface {
	// //When after edit, the server handles the update
	Update(content contenttype.ContentTyper, inputs InputMap) error
}

type ContentTypeHandlerDelete interface {
	//
	// //when deleting
	Delete(content contenttype.ContentTyper) error
}

//Callback struct
type OperationHandler struct {
	Identifier string //Identifier for handler matching. see operation_handler.json/yaml
	//todo: use translation as parameter instead of optional interface{}
	Execute func(triggedEvent string, content contenttype.ContentTyper, params ...interface{}) error
}
