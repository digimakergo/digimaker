//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package field

import (
	"database/sql/driver"
	"dm/model"
	"errors"
)

//TextField is a field for normal text line. It implements Datatyper
type TextField struct {
	*model.Field
	data string
}

//When update db.
func (t TextField) Value() (driver.Value, error) {
	return t.data, nil
}

func (t *TextField) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Incompatible type for GzippedText")
	}

	t.data = source
	return nil
}
