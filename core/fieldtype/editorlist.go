//Author xc, Created on 2019-03-26 20:44
//{COPYRIGHTS}

package fieldtype

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type EditorList struct {
	FieldValue
}

func (t *EditorList) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Incompatible type")
	}

	t.Raw = source
	return nil
}

//when update db
func (t EditorList) Value() (driver.Value, error) {
	return t.Data, nil
}

//implement FieldtypeHandler
type EditorListHandler struct{}

func (e EditorListHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (e EditorListHandler) NewValueFromInput(input interface{}) interface{} {
	r := EditorList{}
	r.Raw = input.(string)
	return r
}

func (e EditorListHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("editorlist", EditorListHandler{})
}
