package fieldtype

import (
	"database/sql/driver"
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

//Radio struct represent radio type
type Datetime struct {
	String
}

func (dt Datetime) Type() string {
	return "datetime"
}

func (dt Datetime) Value() (driver.Value, error) {
	if dt.String.String == "" {
		return nil, nil
	} else {
		return dt.String.String, nil
	}
}

type RelationList struct {
	JSON
}

func (r RelationList) Type() string {
	return "relationlist"
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "text", Value: "fieldtype.Text"},
		func() FieldTyper { return &Text{} })
	RegisterFieldType(
		FieldtypeDef{Type: "json", Value: "fieldtype.JSON"},
		func() FieldTyper { return &JSON{} })
	RegisterFieldType(
		FieldtypeDef{Type: "radio", Value: "fieldtype.Radio"},
		func() FieldTyper { return &Radio{} })
	RegisterFieldType(
		FieldtypeDef{Type: "checkbox", Value: "fieldtype.Checkbox"},
		func() FieldTyper { return &Checkbox{} })
	RegisterFieldType(
		FieldtypeDef{Type: "richtext", Value: "fieldtype.RichText"},
		func() FieldTyper { return &RichText{} })
	RegisterFieldType(
		FieldtypeDef{Type: "number", Value: "fieldtype.Number"},
		func() FieldTyper { return &Number{} })
	RegisterFieldType(
		FieldtypeDef{Type: "datetime", Value: "fieldtype.Datetime"},
		func() FieldTyper { return &Datetime{} })
}
