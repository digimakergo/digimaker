package fieldtype

import "github.com/xc/digimaker/core/util"

//Password struct represent password type
type Password struct {
	String
}

func (p Password) Type() string {
	return "password"
}

//LoadFromInput load data from input before validation
func (p *Password) LoadFromInput(input interface{}) error {
	if input != nil {
		str := input.(string)
		if str != "" {
			p.String.existing = p.String.String
			hashedStr, err := util.HashPassword(str)
			if err != nil {
				return err
			}
			p.String.String = hashedStr
		}
	}
	return nil
}

//Validate the field, return false, error message when fails
func (p *Password) Validate(rule VaidationRule) (bool, string) {
	if p.existing == "" && p.String.String == "" {
		return false, "Password is empty"
	}

	return true, ""
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "password", Value: "fieldtype.Password"},
		func() FieldTyper { return &Password{} })
}
