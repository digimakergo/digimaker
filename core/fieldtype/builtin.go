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

type Email String

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

func init() {
	RegisterFieldType(func() FieldTyper { return &Text{} })
	RegisterFieldType(func() FieldTyper { return &Number{} })
	RegisterFieldType(func() FieldTyper { return &RichText{} })
}
