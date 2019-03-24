package models

type Fielder interface {
	Create()
	Validate(Contenter)
	Store()
}

//Field is a general type for field. It needs to implement Fielder.
//
//A typical new field type(eg. isbn) needs implement Fielder, Datatyper(but not necessary both in a struct).
//
type Field struct {
	Identifier string
	InputData  string
	FieldType  Datatype
}

func (f Field) GetStoredData() (string, error) {
	return "", nil
}
