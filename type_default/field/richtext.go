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
