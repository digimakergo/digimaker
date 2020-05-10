package fieldtype

import (
	"fmt"
	"strconv"
)

/**
* Initial types END
 */
//Text
type Text struct {
	String
}

func (t Text) Type() string {
	return "text"
}

type RichText struct {
	String
}

func (rt RichText) Type() string {
	return "richtext"
}

type Number struct {
	Int
}

func (n Number) Type() string {
	return "number"
}

func (n Number) Validate(input interface{}, rule VaidationRule) (bool, string) {
	s := fmt.Sprint(input)
	if s != "" {
		_, err := strconv.Atoi(s)
		if err != nil {
			return false, s + " is not a number."
		}
	}
	return true, ""
}

//Radio struct represent radio type
type Radio struct {
	Int
}

func (r Radio) Type() string {
	return "radio"
}

//Radio struct represent radio type
type Checkbox struct {
	Int
}

func (c Checkbox) Type() string {
	return "checkbox"
}

//Password struct represent password type
type Password struct {
	String
}

func (r Password) Type() string {
	return "password"
}

//Password struct represent password type
type Image struct {
	String
}

func (i Image) Type() string {
	return "image"
}

type RelationList struct {
	JSON
}

func (r RelationList) Type() string {
	return "relationlist"
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "text", Value: "Text"},
		func() FieldTyper { return &Text{} })
	RegisterFieldType(
		FieldtypeDef{Type: "json", Value: "JSON"},
		func() FieldTyper { return &JSON{} })
	RegisterFieldType(
		FieldtypeDef{Type: "radio", Value: "Radio"},
		func() FieldTyper { return &Radio{} })
	RegisterFieldType(
		FieldtypeDef{Type: "checkbox", Value: "Checkbox"},
		func() FieldTyper { return &Checkbox{} })
	RegisterFieldType(
		FieldtypeDef{Type: "richtext", Value: "RichText"},
		func() FieldTyper { return &RichText{} })
	RegisterFieldType(
		FieldtypeDef{Type: "number", Value: "Number"},
		func() FieldTyper { return &Number{} })
	RegisterFieldType(
		FieldtypeDef{Type: "password", Value: "Password"},
		func() FieldTyper { return &Password{} })
	RegisterFieldType(
		FieldtypeDef{Type: "image", Value: "Image"},
		func() FieldTyper { return &Image{} })
	RegisterFieldType(
		FieldtypeDef{Type: "relationlist", Value: "RelationList", IsRelation: true},
		func() FieldTyper { return &RelationList{} })
}
