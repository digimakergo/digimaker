package models

type Fielder interface {
	Create()
	Validate(Contenter)
	Store()
}

type Field struct {
	Identifier string
	InputData  string
	FieldType  Datatype
}

func (f Field) GetStoredData() (string, error) {
	return "", nil
}
