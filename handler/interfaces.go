//Author xc, Created on 2019-05-11 12:33
//{COPYRIGHTS}
package handler

import (
	"database/sql"
	"dm/contenttype"
)

//This is a callback based on type.
//It's used for customzed content type. eg. image to set parent_id
type ContentTypeHandler interface {
	//When creating on server side
	Create(content contenttype.ContentTyper, tx *sql.Tx, parentID ...int) error

	//When created
	Created()

	//When after edit, the server handles the update
	Update()

	//when updated
	Updated()

	//when deleting
	Delete()
}

//Callback struct
type OperationHandler struct {
	Identifier string //Identifier for handler matching. see operation_handler.json/yaml
	Event      string //event type
	Execute    func(content contenttype.ContentTyper) error
}
