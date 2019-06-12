//Author xc, Created on 2019-03-26 20:44
//{COPYRIGHTS}

package fieldtype

import (
	"errors"
	"fmt"
)

type StringValue struct {
	Data string
}

func (t *StringValue) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Incompatible type")
	}

	fmt.Println(source)
	t.Data = source
	return nil
}
