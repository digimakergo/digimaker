package models

import "errors"

type Fielder interface {
	Create()
	Validate(Contenter)
	SetStoreData()
}

//Field is a general type for field. It needs to implement Fielder.
//
//A typical new field type(eg. isbn) needs implement Fielder, Datatyper(but not necessary both in a struct).
//
type Field struct {
	FieldType  string //do not use DataType because this is better to do instance with json.
	InputData  string
	StoredData string
}

func (f Field) GetStoredData() (string, error) {
	return "", nil
}

//SetStoreData converts InputData to StoredData with validation.
func (f *Field) SetStoreData(c *Content, identifer string) error {
	err := f.Validate(c, identifer)
	if err != nil {
		return errors.New("Validation error. " + err.Error())
	}

	f.StoredData = f.InputData //todo: use specific field to convert. By default it store what's given.

	return nil
}

func (f *Field) Validate(c *Content, identifer string) error {
	return nil
}
