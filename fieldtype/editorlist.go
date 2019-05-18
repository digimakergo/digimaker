//Author xc, Created on 2019-03-26 20:44
//{COPYRIGHTS}

package fieldtype

import (
	"database/sql/driver"
	"errors"
)

type EditorList struct {
	Data string `json:"data"`
}

func (t *EditorList) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Incompatible type")
	}

	t.Data = source
	return nil
}

//when update db
func (t EditorList) Value() (driver.Value, error) {
	return t.Data, nil
}
