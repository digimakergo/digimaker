package fieldtype

// A field type
// There are 3 types of datas in a field: input data, raw data, output data. They can be the same, but sometime can be different.
// eg. for a text field, they are all the same
//     for a richtext field, input can contains some like <a href="uid:dfdsf213123llkjjj">Test</a>, raw data can be the same as input, and output data is <a href="test/test-22"></a>
//There are 2 ways to init a field: by NewFromInput, Scan - from db or import
//
type Fieldtype interface {
	//Validate for input. Do not store
	Validate(interface{}) (bool, string)

	//Validation if input is empty. Do not store
	IsEmtpy(input interface{}) bool

	//Create a field from input. must be validated.
	NewFromInput(interface{})

	//When binding data from db
	Scan(src interface{}) error

	GetDefinition()

	//Get Raw data - the data from db
	GetRaw() interface{}

	//Get the output data. It can be string, int, boolean, or object(eg. json)
	GetOutput() interface{}
}
