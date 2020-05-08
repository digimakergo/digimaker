//Author xc, Created on 2020-05-07 11:00 (merged from 2019-03-25 20:00 )
//{COPYRIGHTS}
package fieldtype

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

/**
* Basisc internal types
 */
//Basic String(note: null string is not allowed, all empty string should use "" - also in db)
type String struct {
	String string
}

//Validate validates input.
func (s String) Validate(input interface{}, rule VaidationRule) (bool, string) {
	return true, ""
}

//LoadFromInput load data from input after validation
func (s *String) LoadFromInput(input interface{}) {
	if input == nil {
		s.String = ""
	} else {
		s.String = input.(string)
	}
}

//Value returns string to db
func (s *String) Value() (driver.Value, error) {
	return s.String, nil
}

//Scan scan data from db.
func (s *String) Scan(src interface{}) error {
	switch src.(type) {
	case string:
		s.String = src.(string)
	case []byte:
		s.String = string(src.([]byte))
	default:
		return errors.New("Can not scan type in ." + fmt.Sprint(src))
	}
	return nil
}

//FieldValue return string value
func (s String) FieldValue() interface{} {
	return s.String
}

//IsEmpty checks if it is a emtpy input
func (s String) IsEmpty() bool {
	return s.String == ""
}

//MarshalJSON marshals to []byte
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String)
}

// Int is a basic internal common type Int
type Int struct {
	sql.NullInt64
}

//LoadFromInput loads data from input after validation
func (i Int) LoadFromInput(input interface{}) {
	i.Scan(input)
}

//Get FieldValue
func (i Int) FieldValue() interface{} {
	return int(i.Int64)
}

//IsEmpty checks if the input is empty. so ""/nil is empty,
func (i Int) IsEmpty() bool {
	return false
}

//Validate validates input.
func (i Int) Validate(input interface{}, rule VaidationRule) (bool, string) {
	return true, ""
}

//MarshalJSON marshals to []byte
func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Int64)
}

func (i *Int) UnMarshalJSON(data []byte) error {
	i.Scan(data)
	return nil
}
