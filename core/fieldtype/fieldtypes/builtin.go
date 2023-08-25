package fieldtypes

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query/querier"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

/***** TextHandler ******/

type TextParameters struct {
	MinLength     int    `mapstructure:"min_length"`
	MaxLength     int    `mapstructure:"max_length"`
	RegExp        string `mapstructure:"regexp"`
	RegExpMessage string `mapstructure:"regexp_message"`
}

type TextHandler struct {
	fieldtype.FieldDef
	Params TextParameters
}

func (handler TextHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	str := fmt.Sprint(input)
	params := handler.Params
	strLength := len([]rune(str))

	//min length
	if params.MinLength > 0 && strLength < params.MinLength {
		return nil, fieldtype.NewValidationError(fmt.Sprintf("Input needs at least %v characters", params.MinLength))
	}

	//max length
	if params.MaxLength > 0 && strLength > params.MaxLength {
		return nil, fieldtype.NewValidationError(fmt.Sprintf("Input can not have more than %v characters", params.MinLength))
	}

	//regular expression match
	if params.RegExp != "" {
		matched, err := regexp.MatchString(params.RegExp, str)
		if err != nil {
			return nil, fmt.Errorf("Matching error: %v", err.Error())
		}
		if !matched {
			return nil, fieldtype.NewValidationError(params.RegExpMessage)
		}
	}
	return str, nil
}

func (handler TextHandler) ValidateDefinition() error {
	params := handler.Params
	if params.MinLength >= params.MaxLength {
		return errors.New("Min length should be less than max length")
	}
	return nil
}

func (handler TextHandler) DBField() string {
	rule := handler.Params
	maxLength := rule.MaxLength
	return fmt.Sprintf("VARCHAR (%v) NOT NULL DEFAULT ''", maxLength)
}

/***** RichTextHandler ******/
type RichTextHandler struct {
	fieldtype.FieldDef
}

func (handler RichTextHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	//replace html like <img src="var/f/fg/fge1ff.png" data-content="image;sdf319432424b432341" /> with real image path and size.
	return input, nil
}

func (handler RichTextHandler) DBField() string {
	return "TEXT"
}

func (r RichTextHandler) Output(ctx context.Context, querier querier.Querier, value interface{}) interface{} {
	//convert image to updated image path
	re := regexp.MustCompile(`<img[^>]+data-dm-content="[^"]+"[^>]+>`)

	strValue := value.(string)

	imagePrefix := viper.GetString("general.var_baseurl")
	replaceFunc := func(currentStr string) string {
		re2 := regexp.MustCompile(`([^ =]+)="([0-9a-zA-Z]|;)+"`)
		attributes := re2.FindAllString(currentStr, -1)
		attributeMap := map[string]string{}
		for _, attStr := range attributes {
			arr := strings.Split(attStr, "=")
			name := arr[0]
			value := strings.ReplaceAll(arr[1], `"`, "")
			attributeMap[name] = value
		}
		contentInfo := strings.Split(attributeMap["data-dm-content"], ";")
		if len(contentInfo) <= 1 {
			log.Warning("data-dm-content has wrong format, should be <contenttype>;<cuid>, no replace done. - "+currentStr, "output", ctx)
			return currentStr
		}

		condition := db.Cond("c._cuid", contentInfo[1])
		content, _ := querier.Fetch(ctx, contentInfo[0], condition)
		widthStr := ""
		if width, ok := attributeMap["width"]; ok {
			widthStr = `width="` + width + `"`
		}

		heightStr := ""
		if height, ok := attributeMap["height"]; ok {
			heightStr = `height="` + height + `"`
		}

		dataAttribute := `data-dm-content="` + attributeMap["data-dm-content"] + `"`

		if content == nil {
			//to do: check reason(might be missing access) and give log, and output different image
			return fmt.Sprintf(`<img src="not-available.png" %v %v %v />`, widthStr, heightStr, dataAttribute) //todo: make it configurable
		}

		path := imagePrefix + "/" + content.Value("image").(string)

		result := fmt.Sprintf(`<img src="%v" %v %v %v />`, path, widthStr, heightStr, dataAttribute)
		return result
	}
	result := re.ReplaceAllStringFunc(strValue, replaceFunc)
	return result
}

/** Password handler **/
type PasswordHandler struct {
	fieldtype.FieldDef
}

func (handler PasswordHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	str := fmt.Sprint(input)
	return util.HashPassword(str)
}

func (handler PasswordHandler) DBField() string {
	return "BINARY(60)"
}

/** Int handler **/
type IntHandler struct {
	fieldtype.FieldDef
}

func (handler IntHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	var i int
	switch input.(type) {
	case int:
		i = input.(int)
	case float64:
		i = int(input.(float64))
	case string:
		s := input.(string)
		if s == "" {
			return -1, nil
		}
		j, err := strconv.Atoi(s)
		if err != nil {
			return -1, fieldtype.NewValidationError(err.Error())
		}
		i = j
	default:
		return -1, fieldtype.NewValidationError("Not supported type as int")
	}
	return i, nil
}

func (handler IntHandler) DBField() string {
	return "INT NOT NULL DEFAULT -1" //todo: make default configurable
}

/** Checkbox handler ***/
type CheckboxHandler struct {
	fieldtype.FieldDef
}

// Only allow 1/0 or "1"/"0"
func (handler CheckboxHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	str := fmt.Sprint(input)
	valueInt, err := strconv.Atoi(str)
	if err != nil || (valueInt != 0 && valueInt != 1) {
		return nil, fieldtype.NewValidationError("Only allow 1 or 0")
	}
	return valueInt, nil
}

func (handler CheckboxHandler) DBField() string {
	return "TINYINT(1)"
}

/** Radio handler ***/
type RadioHandler struct {
	fieldtype.FieldDef
}

// max 30 length
func (handler RadioHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	str := fmt.Sprint(input)
	length := len(str)
	if length > 30 {
		return nil, fieldtype.NewValidationError("Radio value can not be more than 30 characters")
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
	fieldtype.FieldDef
}

// support 3 format: 2020-01-10, 2020-01-10 12:12:00(server side timezone), 2021-01-10T11:45:26+02:00, 1722145855(unix)
func (handler DatetimeHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	str := strings.TrimSpace(fmt.Sprint(input))

	if str == "" {
		return time.Time{}, nil
	}

	//unix format
	if ok, _ := regexp.Match(`^\d+$`, []byte(str)); ok {
		unixInt, _ := strconv.Atoi(str)
		value := time.Unix(int64(unixInt), 0)
		return value, nil
	}

	if !strings.Contains(str, ":") {
		value, err := time.Parse("2006-01-02", str)
		if err != nil {
			return nil, fieldtype.NewValidationError("Wrong format, only allow 2006-01-02")
		}
		return value, nil
	} else {
		if strings.Contains(str, " ") {
			value, err := time.Parse("2006-01-02 15:04", str)
			if err != nil {
				return nil, fieldtype.NewValidationError("Wrong format, only allow like 2006-01-02 15:04")
			}
			return value, nil
		} else {
			value, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return nil, fieldtype.NewValidationError("Wrong format, only allow RFC3339 format")
			}
			return value, nil
		}

	}
}

func (handler DatetimeHandler) DBField() string {
	return "DATETIME"
}

// convert parameters(definition.FieldParameters - map[string]interface{}) to struct using mapstructure
func ConvertParameters(params fieldtype.FieldParameters, paramStruct interface{}) error {
	if params != nil {
		err := mapstructure.Decode(params, &paramStruct)
		if err != nil {
			returnError := errors.New("Parameters error:" + err.Error())
			return returnError
		}
	}
	return nil
}

/** Select handler ***/

const selectMultMax = 255 //max characters of multi select

type SelectParameters struct {
	Multi bool                `mapstructure:"multi"`
	List  []map[string]string `mapstructure:"list"`
}

type SelectHandler struct {
	fieldtype.FieldDef
	Params SelectParameters
}

// Only allow 1/0 or "1"/"0"
func (handler SelectHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	value := fmt.Sprint(input)
	params := handler.Params

	if value != "" {
		if params.Multi {
			if len(value) > selectMultMax {
				return nil, fieldtype.NewValidationError(fmt.Sprintf("Value can not be longer than, %v", selectMultMax))
			}
		}

		valueArr := util.Split(value, ";")
		for _, v := range valueArr {
			if !util.ListContains(params.List, "value", v) {
				return nil, fieldtype.NewValidationError("Value is not defined in list: " + v)
			}
		}
	}

	return value, nil
}

func (handler SelectHandler) Output(ctx context.Context, querier querier.Querier, value interface{}) interface{} {
	params := handler.Params
	result := []map[string]string{}
	values := util.Split(value.(string), ";")

	if params.Multi {
		for _, item := range params.List {
			if util.Contains(values, item["value"]) {
				result = append(result, item)
			}
		}
		return result
	} else {
		for _, item := range params.List {
			if util.Contains(values, item["value"]) {
				return item
			}
		}
		return nil
	}
}

func (handler SelectHandler) DBField() string {
	if handler.Params.Multi {
		return fmt.Sprintf("VARCHAR(%v) NOT NULL DEFAULT ''", selectMultMax)
	}
	return "varchar(50) NOT NULL DEFAULT ''"
}

func init() {
	fieldtype.Register(
		fieldtype.Definition{
			Name:     "text",
			DataType: "string",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				params := TextParameters{}
				err := ConvertParameters(def.Parameters, &params)
				if err != nil {
					log.Error("Definition error on text, parameters ignored: "+err.Error(), "")
				}
				if params.MaxLength == 0 {
					params.MaxLength = 255 //default text max length
				}

				return TextHandler{FieldDef: def, Params: params}
			}})

	fieldtype.Register(
		fieldtype.Definition{Name: "richtext",
			DataType: "string",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return RichTextHandler{FieldDef: def}
			}})

	fieldtype.Register(
		fieldtype.Definition{Name: "password",
			DataType: "string",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return PasswordHandler{FieldDef: def}
			}})

	fieldtype.Register(
		fieldtype.Definition{Name: "checkbox",
			DataType: "int",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return CheckboxHandler{FieldDef: def}
			}})

	fieldtype.Register(
		fieldtype.Definition{Name: "radio",
			DataType: "string",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return RadioHandler{FieldDef: def}
			}})

	fieldtype.Register(
		fieldtype.Definition{Name: "datetime",
			DataType: "time.Time",
			Package:  "time",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return DatetimeHandler{FieldDef: def}
			}})

	fieldtype.Register(
		fieldtype.Definition{Name: "select",
			DataType: "string",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				params := SelectParameters{}
				err := ConvertParameters(def.Parameters, &params)
				if err != nil {
					log.Error("Definition error on select, parameters ignored: "+err.Error(), "")
				}
				return SelectHandler{FieldDef: def, Params: params}
			}})

	fieldtype.Register(
		fieldtype.Definition{Name: "int",
			DataType: "int", //in practise it's for positive int
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return IntHandler{FieldDef: def}
			}})

	// Register("float", "float", TextHandler{})                      //number is postive float
}
