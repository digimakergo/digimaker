//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package field

import (
	"database/sql/driver"
	"dm/model"
)

type RichTextField struct {
	*model.Field
	data string
}

//when update db
func (t RichTextField) Value() (driver.Value, error) {
	return t.data, nil
}

func (t *RichTextField) Scan(src interface{}) error {
	t.data = "good2"
	return nil
}
