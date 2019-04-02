//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package field

import (
	"dm/model"
	"errors"
)

//TextField is a field for normal text line. It implements Datatyper
type TextField struct {
	*model.Field
	data string
}

func (t *TextField) Value() string {
	return t.data
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
