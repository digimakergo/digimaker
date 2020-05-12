//Author xc, Created on 2020-05-07 11:00 (merged from 2019-03-25 20:00 )
//{COPYRIGHTS}
package fieldtype

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xc/digimaker/core/log"
)

/**
* Basisc internal types
 */
//Basic String(note: null string is not allowed, all empty string should use "" - also in db)
type String struct {
	String string
}

//Validate validates input.
func (s String) Validate(rule VaidationRule) (bool, string) {
	return true, ""
}

//LoadFromInput load data from input before validation
func (s *String) LoadFromInput(input interface{}) error {
	if input == nil {
		s.String = ""
	} else {
		s.String = input.(string)
	}
	return nil
}

//Value returns string to db
func (s String) Value() (driver.Value, error) {
	return s.String, nil
}

//Scan scan data from db.
func (s *String) Scan(src interface{}) error {
	if src == nil {
		s.String = ""
		return nil
	}
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

// Int is a basic internal common type Int, which allows empty.
type Int struct {
	sql.NullInt64
}

//LoadFromInput loads data from input before validation
func (i Int) LoadFromInput(input interface{}) error {
	err := i.Scan(input)
	if err != nil {
		inputStr := fmt.Sprintln(input)
		log.Error("Input is not a nullable int:"+inputStr, "")
		return errors.New("Not a valid int: " + inputStr)
	}
	return nil
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
func (i Int) Validate(rule VaidationRule) (bool, string) {
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

// JSON is a json format string. Null is allowed
type JSON struct {
	sql.NullString
}

//LoadFromInput loads data from input before validation
func (j JSON) LoadFromInput(input interface{}) error {
	return j.Scan(input)
}

//Get FieldValue
func (j JSON) FieldValue() interface{} {
	if !j.Valid {
		return nil
	} else {
		return j.String //todo: improve this.
	}
}

//IsEmpty checks if the input is empty. so ""/nil is empty,
func (j JSON) IsEmpty() bool {
	return !j.Valid
}

//Validate validates input.
func (j JSON) Validate(rule VaidationRule) (bool, string) {
	return true, ""
}

//MarshalJSON marshals to []byte
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.String)
}

func (j *JSON) UnMarshalJSON(data []byte) error {
	j.Scan(data)
	return nil
}

func (j JSON) Type() string {
	return "json"
}
