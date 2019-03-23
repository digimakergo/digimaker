package base

type Datatyper interface {
	Create()
	Validate()
	Store()
}

// Datatype is the base struct for all explect datatypes
type Datatype struct {
	identifier string
	name       string
	searchable bool
}

func loadDefinition() {

}

func (d Datatype) Create() {

}

func (d Datatype) Store() error {
	return nil
}
