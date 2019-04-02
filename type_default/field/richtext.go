//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package field

import "dm/model"

type RichTextField struct {
	*model.Field
	data string
}

func (t *RichTextField) Value() string {
	return t.data
}

func (t *RichTextField) Scan(src interface{}) error {
	t.data = "good2"
	return nil
}
