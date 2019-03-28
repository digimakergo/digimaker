package models

//All the content type(eg. article, folder) will implement this interface.
type ContentTyper interface {
	//Return all fields
	Fields() map[string]Field

	//Visit  field dynamically
	Field(name string) Field

	//Visit all attribute dynamically including Fields + internal attribute eg. id, parent_id.
	Attr(name string) interface{}
}

//All of the fields will implements this interface
type Fielder interface {
	//Get value of
	Value()
	Create()
	Validate(Contenter)
	SetStoreData()
}
