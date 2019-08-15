package fieldtype

import "strings"

//All field types need to implement FieldtypeHandler interface.
type FieldtypeHandler interface {
	//Create a value instance from input(http input or api input)
	NewValueFromInput(input interface{}) interface{}

	//Validate the field, return false, error message when fails
	Validate(input interface{}) (bool, string)

	//If the field is empty based on input
	IsEmpty(input interface{}) bool
}

//This struct use generic field handling which include a FieldtypeHandler(composition patter)
type FieldHandler struct {
	Fieldtype string
	handler   FieldtypeHandler
}

func (f FieldHandler) Validate(input interface{}) (bool, string) {
	return f.handler.Validate(input)
}

func (f FieldHandler) NewValue(input interface{}) interface{} {
	return f.handler.NewValueFromInput(input)
}

func (f FieldHandler) IsEmpty(input interface{}) bool {
	result := false
	if strings.TrimSpace(input.(string)) == "" {
		result = true
	}
	fieldIsEmpty := f.handler.IsEmpty(input)
	result = result && fieldIsEmpty
	return result
}

//Relation field handler can convert relations into RelationField
type RelationFieldHandler interface {
	ToStorage(contents interface{}) interface{}
	UpdateOne(toContent interface{}, identifier string, from interface{})
}
