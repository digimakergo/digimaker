//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

import (
	"strings"
)

//todo: better design relation between RichTextField and FieldValue.The new is
type RichTextField struct {
	FieldValue
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
type RichTextHandler struct{}

func (t RichTextHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t RichTextHandler) NewValueFromInput(input interface{}) interface{} {
	r := RichTextField{}
	r.Scan(input.(string))
	return r
}

func (t RichTextHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("richtext", RichTextHandler{})
}
