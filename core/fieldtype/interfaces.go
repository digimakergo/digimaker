package fieldtype

// FieldTyper defines a field type data
// There are 2 types of datas in a field: input data, output data. They can be the same, but sometime can be different.
// eg. for a text field, they are all the same
//     for a richtext field, input can contains some like <a href="uid:dfdsf213123llkjjj">Test</a>, output data is <a href="test/test-22"></a>
//There are 3 ways to init a field:
//1) by NewFromInput
//2) Scan - from db or import
//3) UnMarshall - from json
//
type FieldTyper interface {
	// //Init from db
	// //Scan(src interface{}) error
	// sql.Scanner
	//
	// //Get value when insert/update to DB.
	// //same as database.sql.driver.Valuer
	// //Value() (Value, error)
	// driver.Valuer

	//Init from input(http input or api input)
	//The return must be a basic type. Invoke after validation succeed.
	LoadFromInput(input interface{}) error

	//Get value when doing 'internal exchange'.same as Value() in driver.Valuer, but return a basic type(eg. string, int, datetime, or nil)
	FieldValue() interface{}

	//Validate the field, return false, error message when fails
	Validate(rule VaidationRule) (bool, string)

	//If the field is empty. eg. in a selection, 0/-1 can mean empty(not selected)
	IsEmpty() bool

	//Type of the FieldTyper, eg. text.
	//Because go doesn't support constructor, it makes it hard to set type when:
	//1) creating like Article{}(ok we can force to use a method then we have to loop every new field)
	//2) when doing unmarshall from json. ok we can do it in UnMarshall method to init
	//3) when init from DB. ok we can do it in Scan method..
	// Above all, creating a Type() to always return a fixed type is an easier option, than fining all the possible creating way(and set type in the method of that way).
	Type() string
}

type FieldTypeEvent interface {
	BeforeSaving() error
}
