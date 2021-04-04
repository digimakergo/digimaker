package fieldtype

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/mitchellh/mapstructure"
)

/***** TextHandler ******/

type TextParameters struct {
	MinLength     int    `mapstructure:"min_length"`
	MaxLength     int    `mapstructure:"max_length"`
	RegExp        string `mapstructure:"regexp"`
	RegExpMessage string `mapstructure:"regexp_message"`
}

type TextHandler struct {
	definition.FieldDef
	params TextParameters
}

func (handler *TextHandler) getParams() (TextParameters, error) {
	//cache it
	emptyParams := TextParameters{}
	if handler.params == emptyParams {
		if handler.Parameters != nil {
			rule := TextParameters{}
			err := mapstructure.Decode(handler.Parameters, &rule)
			if err != nil {
				returnError := errors.New("Validation rule error:" + err.Error())
				return emptyParams, returnError
			}
			handler.params = rule
		}
	}
	return handler.params, nil
}

func (handler TextHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	str := fmt.Sprint(input)
	params, err := handler.getParams()
	if err != nil {
		return params, err
	}
	strLength := len([]rune(str))

	//min length
	if params.MinLength > 0 && strLength < params.MinLength {
		return nil, NewValidationError(fmt.Sprintf("Input needs at least %v characters", params.MinLength))
	}

	//max length
	if params.MaxLength > 0 && strLength > params.MaxLength {
		return nil, NewValidationError(fmt.Sprintf("Input can not have more than %v characters", params.MinLength))
	}

	//regular expression match
	if params.RegExp != "" {
		matched, err := regexp.MatchString(params.RegExp, str)
		if err != nil {
			return nil, fmt.Errorf("Matching error: %v", err.Error())
		}
		if !matched {
			return nil, NewValidationError(params.RegExpMessage)
		}
	}
	return str, nil
}

func (handler TextHandler) DBField() string {
	rule, _ := handler.getParams()
	maxLength := rule.MaxLength
	return fmt.Sprintf("VARCHAR (%v) DEFAULT ''", maxLength)
}

/***** RichTextHandler ******/
type RichTextHandler struct {
	definition.FieldDef
}

func (handler RichTextHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	//replace html like <img src="var/f/fg/fge1ff.png" data-content="image;sdf319432424b432341" /> with real image path and size.
	return input, nil
}

func (handler RichTextHandler) DBField() string {
	return "text"
}

func (r RichTextHandler) ConvertOuput() interface{} {
	//replace html like <img src="var/f/fg/fge1ff.png" data-content="image;sdf319432424b432341" /> with real image path and size.
	return ""
}

/** Checkbox handler ***/
type CheckboxHandler struct {
	definition.FieldDef
}

//Only allow 1/0 or "1"/"0"
func (handler CheckboxHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	str := fmt.Sprint(input)
	valueInt, err := strconv.Atoi(str)
	if err != nil || (valueInt != 0 && valueInt != 1) {
		return nil, NewValidationError("Only allow 1 or 0")
	}
	return valueInt, nil
}

func (handler CheckboxHandler) DBField() string {
	return "TINYINT 1"
}

/** Radio handler ***/
type RadioHandler struct {
	definition.FieldDef
}

//max 30 length
func (handler RadioHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	str := fmt.Sprint(input)
	length := len(str)
	if length > 30 {
		return nil, NewValidationError("Radio value can not be more than 30 characters")
	}
	return str, nil
}

func (handler RadioHandler) DBField() string {
	return "VARCHAR (30) NOT NULL DEFAULT ''"
}

/** Datetime handler ***/

type DatetimeParameters struct {
	Dateonly bool `mapstructure:"dateonly"`
}

type DatetimeHandler struct {
	definition.FieldDef
}

//support 3 format: 2020-01-10, 2020-01-10 12:12:00(server side timezone), 2021-01-10T11:45:26+02:00, 1722145855(unix)
func (handler DatetimeHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	str := strings.TrimSpace(fmt.Sprint(input))

	//unix format
	if ok, _ := regexp.Match(`\d+`, []byte(str)); ok {
		unixInt, _ := strconv.Atoi(str)
		value := time.Unix(int64(unixInt), 0)
		return value, nil
	}

	if !strings.Contains(str, ":") {
		value, err := time.Parse("2006-01-02", str)
		if err != nil {
			return nil, NewValidationError("Wrong format, only allow 2006-01-02")
		}
		return value, nil
	} else {
		if strings.Contains(str, " ") {
			value, err := time.Parse("2006-01-02 15:04:05", str)
			if err != nil {
				return nil, NewValidationError("Wrong format, only allow like 2006-01-02 15:04:05")
			}
			return value, nil
		} else {
			value, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return nil, NewValidationError("Wrong format, only allow RFC3339 format")
			}
			return value, nil
		}

	}
}

func (handler DatetimeHandler) DBField() string {
	return "DATETIME"
}

func init() {
	Register(
		Definition{
			Name:     "text",
			DataType: "string",
			NewHandler: func(def definition.FieldDef) Handler {
				return TextHandler{FieldDef: def}
			}})

	Register(
		Definition{Name: "richtext",
			DataType: "string",
			NewHandler: func(def definition.FieldDef) Handler {
				return RichTextHandler{FieldDef: def}
			}})

	Register(
		Definition{Name: "password",
			DataType: "string",
			NewHandler: func(def definition.FieldDef) Handler {
				return TextHandler{FieldDef: def}
			}})

	Register(
		Definition{Name: "checkbox",
			DataType: "int",
			NewHandler: func(def definition.FieldDef) Handler {
				return CheckboxHandler{FieldDef: def}
			}})

	Register(
		Definition{Name: "radio",
			DataType: "string",
			NewHandler: func(def definition.FieldDef) Handler {
				return RadioHandler{FieldDef: def}
			}})

	Register(
		Definition{Name: "datetime",
			DataType: "time.Time",
			Package:  "time",
			NewHandler: func(def definition.FieldDef) Handler {
				return DatetimeHandler{FieldDef: def}
			}})

	// Register("number", "int", TextHandler{})                       //number is postive int
	// Register("float", "float", TextHandler{})                      //number is postive float
	// Register("full_float", "fieldtype.NilablFloat", TextHandler{}) //all float
	// Register("select", "string", TextHandler{})
	// Register("json", "json", TextHandler{})
	// Register(
	// 	Definition{
	// 		Name:     "image",
	// 		DataType: "string",
	// 	})
	// Register("file", "string", TextHandler{})
	// Register("relation", "string", TextHandler{})
}
