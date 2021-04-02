package fieldtype

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/mitchellh/mapstructure"
)

type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}

func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

/***** TextHandler ******/

type TextValidationRule struct {
	MinLength     int    `mapstructure:"min_length"`
	MaxLength     int    `mapstructure:"max_length"`
	RegExp        string `mapstructure:"regexp"`
	RegExpMessage string `mapstructure:"regexp_message"`
}

type TextHandler struct {
	Def definition.FieldDef
}

func (handler TextHandler) ConvertInput(input interface{}) (interface{}, error) {
	str := fmt.Sprint(input)
	if handler.Def.Validation != nil {
		rule := TextValidationRule{}
		err := mapstructure.Decode(handler.Def.Validation, &rule)
		if err != nil {
			returnError := errors.New("Validation rule error:" + err.Error())
			return nil, returnError
		}

		strLength := len([]rune(str))
		//min length
		if rule.MinLength > 0 && strLength < rule.MinLength {
			return nil, NewValidationError(fmt.Sprintf("Input needs at least %v characters", rule.MinLength))
		}
		//max length
		if rule.MaxLength > 0 && strLength > rule.MaxLength {
			return nil, NewValidationError(fmt.Sprintf("Input can not have more than %v characters", rule.MinLength))
		}

		//regular expression match
		if rule.RegExp != "" {
			matched, err := regexp.MatchString(rule.RegExp, str)
			if err != nil {
				return nil, fmt.Errorf("Matching error: %v", err.Error())
			}
			if !matched {
				return nil, NewValidationError(rule.RegExpMessage)
			}
		}
	}

	return str, nil
}

func (TextHandler) IsEmpty(input interface{}) bool {
	str := fmt.Sprint(input)
	trimed := strings.TrimSpace(str)
	return trimed == ""
}

func (TextHandler) BeforeSave(existing interface{}) bool {
	return true
}

/***** RichText ******/
type RichtextHandler struct {
	String string
	ouput  string
}

func (r RichtextHandler) ConvertOuput() interface{} {
	//replace html like <img src="var/f/fg/fge1ff.png" data-content="image;sdf319432424b432341" /> with real image path and size.
	return ""
}

func init() {
	Register(Definition{Name: "text",
		BaseType: "string",
		NewHandler: func(def definition.FieldDef) Handler {
			return TextHandler{Def: def}
		}})

	// Register("richtext", "string", TextHandler{})
	// Register("number", "int", TextHandler{})                       //number is postive int
	// Register("float", "float", TextHandler{})                      //number is postive float
	// Register("full_float", "fieldtype.NilablFloat", TextHandler{}) //all float
	// Register("checkbox", "int", TextHandler{})
	// Register("radio", "int", TextHandler{})
	// Register("select", "string", TextHandler{})
	// Register("datetime", "int", TextHandler{})
	// Register("json", "json", TextHandler{})
	// Register("image", "string", TextHandler{})
	// Register("file", "string", TextHandler{})
	// Register("relation", "string", TextHandler{})
}
