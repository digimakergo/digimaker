package fieldtype

import (
	"database/sql"
	"database/sql/driver"
)

// FieldType defines a field type data
// There are 2 types of datas in a field: input data, output data. They can be the same, but sometime can be different.
// eg. for a text field, they are all the same
//     for a richtext field, input can contains some like <a href="uid:dfdsf213123llkjjj">Test</a>, output data is <a href="test/test-22"></a>
//There are 2 ways to init a field: by NewFromInput, Scan - from db or import
//
type FieldTyper interface {
	//Init from db
	//Scan(src interface{}) error
	sql.Scanner

	//Get value when insert/update to DB.
	//same as database.sql.driver.Valuer
	//Value() (Value, error)
	driver.Valuer

	//Init from input(http input or api input)
	//The return must be a basic type. Invoke after validation succeed.
	LoadFromInput(input interface{})

	//Get value when doing 'internal exchange'.same as Value() in driver.Valuer, but return a basic type(eg. string, int, datetime, or nil)
	FieldValue() interface{}

	//Validate the field, return false, error message when fails
	Validate(input interface{}, rule VaidationRule) (bool, string)

	//If the field is empty. eg. in a selection, 0/-1 can mean empty(not selected)
	IsEmpty() bool

	Type() string
}
