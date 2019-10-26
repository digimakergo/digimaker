//Author xc, Created on 2019-10-14 20:15
//{COPYRIGHTS}

//Package fieldtype implements build-in field types(value and fieldtype handler).
package fieldtype

import (
	"strings"
)

//TextField is a field for normal text line. It implements Datatyper
type CheckboxField struct {
	FieldValue
}

// //when update value to db
// func (t *CheckboxField) Value() (driver.Value, error) {
// 	str, err := strconv.Atoi(t.Raw)
// 	fmt.Println("Good")
// 	fmt.Println(str)
// 	return str, err
// }

func (t *CheckboxField) Scan(src interface{}) error {
	err := t.SetData(src, "checkbox")
	return err
}

//convert data to view data.
func (t CheckboxField) ViewValue() string {
	return t.Raw
}

//implement FieldtypeHandler
type CheckboxHandler struct{}

func (t CheckboxHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t CheckboxHandler) NewValueFromInput(input interface{}) interface{} {
	r := CheckboxField{}
	// i, _ := strconv.Atoi(input.(string))
	r.Scan(input.(string))
	return r
}

func (t CheckboxHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("checkbox", CheckboxHandler{})
}
