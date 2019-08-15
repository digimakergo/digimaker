//Author xc, Created on 2019-08-15 17:37
//{COPYRIGHTS}

package fieldtype

import "strings"

type ImageField struct {
	FieldValue
}

func (t *ImageField) Scan(src interface{}) error {
	err := t.SetData(src, "image")
	return err
}

//implement FieldtypeHandler
type ImageHandler struct{}

func (t ImageHandler) Validate(input interface{}) (bool, string) {
	return true, ""
}

func (t ImageHandler) NewValueFromInput(input interface{}) interface{} {
	r := ImageField{}
	r.Scan(input.(string))
	return r
}

func (t ImageHandler) IsEmpty(input interface{}) bool {
	if strings.TrimSpace(input.(string)) == "" {
		return true
	}
	return false
}

func init() {
	RegisterHandler("image", RichTextHandler{})
}
