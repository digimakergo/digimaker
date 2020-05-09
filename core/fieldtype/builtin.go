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

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "text"},
		func() FieldTyper { return &Text{} })
	RegisterFieldType(
		FieldtypeDef{Type: "radio"},
		func() FieldTyper { return &Radio{} })
	RegisterFieldType(
		FieldtypeDef{Type: "checkbox"},
		func() FieldTyper { return &Checkbox{} })
	RegisterFieldType(
		FieldtypeDef{Type: "richtext"},
		func() FieldTyper { return &RichText{} })
	RegisterFieldType(
		FieldtypeDef{Type: "number"},
		func() FieldTyper { return &Number{} })
	RegisterFieldType(
		FieldtypeDef{Type: "password"},
		func() FieldTyper { return &Password{} })
	RegisterFieldType(
		FieldtypeDef{Type: "image"},
		func() FieldTyper { return &Image{} })

}
