package models

// Datatype is the base struct for all explect datatypes
type Datatype struct {
	Definition map[string]string
	Identifier string
	Name       string
	Searchable bool
}
