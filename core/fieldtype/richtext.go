//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

import (
	"strings"
)

type RichTextField struct {
	FieldtypeValue
}

func (t *RichTextField) Scan(src interface{}) error {
	err := t.SetData(src, "richtext")
	return err
}

func (r *RichTextField) convertToOutput() {
	s := r.Raw
	s = strings.ReplaceAll(s, "fa", "FAG")
	r.Output = s
}

//implement FieldtypeHandler
type RichTextFieldHandler struct {
	*FieldtypeHandler
}

func (t RichTextFieldHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t RichTextFieldHandler) ToStorage(input interface{}) interface{} {
	r := RichTextField{}
	r.Raw = input.(string)
	return r
}

func (t RichTextFieldHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("richtext", RichTextFieldHandler{})
}
