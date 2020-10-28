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
	String   string
	existing string //existing value when changing. It's the 'orginal value' before LoadFromInput change the whole value.
}

//Validate validates input.
func (s String) Validate(rule VaidationRule) (bool, string) {
	return true, ""
}

//LoadFromInput load data from input before validation
func (s *String) LoadFromInput(input interface{}, params FieldParameters) error {
	if input == nil {
		s.String = ""
	} else {
		s.existing = s.String
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

func (s *String) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	s.Scan(str)
	return nil
}

// Int is a basic internal common type Int, which allows empty.
type Int struct {
	sql.NullInt64
	existing sql.NullInt64
}

//LoadFromInput loads data from input before validation
func (i *Int) LoadFromInput(input interface{}, params FieldParameters) error {
	if input == "" {
		input = nil
	}
	existing := i.NullInt64
	err := i.Scan(input)
	if err != nil {
		inputStr := fmt.Sprintln(input)
		log.Error("Input is not a nullable int:"+inputStr, "")
		return errors.New("Not a valid int: " + inputStr)
	}
	i.existing = existing
	return nil
}

//Get FieldValue
func (i Int) FieldValue() interface{} {
	if i.NullInt64.Valid {
		return int(i.Int64)
	} else {
		return nil
	}
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
	if !i.NullInt64.Valid {
		return json.Marshal("")
	}
	return json.Marshal(i.Int64)
}

func (i *Int) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err == nil {
		return i.LoadFromInput(str, FieldParameters{})
	} else {
		return i.Scan(data)
	}
}

// JSON is a json format string. Null is allowed
type JSON struct {
	sql.NullString
	existing sql.NullString
}

//LoadFromInput loads data from input before validation
func (j *JSON) LoadFromInput(input interface{}, params FieldParameters) error {
	//todo: validate if it's a json structure([],{})
	existing := j.NullString
	if str, ok := input.(string); ok {
		j.Scan(str)
		return nil
	} else {
		bytes, _ := json.Marshal(input)
		err := j.Scan(string(bytes))
		if err != nil {
			j.existing = existing
		}
		return err
	}
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
	if j.String == "" {
		return json.Marshal("")
	}
	return []byte(j.String), nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	j.Scan(data)
	return nil
}

func (j JSON) Type() string {
	return "json"
}
