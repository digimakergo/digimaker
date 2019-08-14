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
