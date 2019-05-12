//Author xc, Created on 2019-05-11 12:33
//{COPYRIGHTS}
package handler

import (
	"database/sql"
	"dm/contenttype"
)

//todo: think it might be good to use "struct with callback method" instead of interface.
//This is a callback based on type.
//It's used for customzed content type. eg. image to set parent_id
type ContentTypeHandler interface {
	Validate(inputs map[string]interface{}, result *ValidationResult)

	//When creating on server side, or import.
	// This is low level and should be used for eg. set parent_id when new record is inserted.
	New(content contenttype.ContentTyper, tx *sql.Tx, parentID ...int) error

	// //When created
	// Created()

	// //When after edit, the server handles the update
	// Update()
	//
	// //when updated
	// Updated()
	//
	// //when deleting
	// Delete()
}

//Callback struct
type OperationHandler struct {
	Identifier string //Identifier for handler matching. see operation_handler.json/yaml
	Execute    func(triggedEvent string, content contenttype.ContentTyper, params ...interface{}) error
}
