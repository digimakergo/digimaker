//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

//Package fieldtype implements build-in field types(value and fieldtype handler)
package fieldtype

import (
	"database/sql/driver"
	"errors"
)

//TextField is a field for normal text line. It implements Datatyper
type TextField struct {
	Data     string `json:"data"`
	viewData interface{}
}

//When update db.
func (t TextField) Value() (driver.Value, error) {
	return t.Data, nil
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

	t.Data = source
	return nil
}

//convert data to view data.
func (t TextField) ViewValue() string {
	return t.Data
}
