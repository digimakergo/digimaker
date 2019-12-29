package handler

import (
	"dm/core/contenttype"
	"dm/core/db"
	"dm/core/fieldtype"
	"dm/core/util"
	"errors"
	"strings"
)

//CanLogin check if the username/email and password matches
func CanLogin(usernameEmail string, password string) (error, contenttype.ContentTyper) {
	querier := Querier()

	//todo: use username instead of 'login'
	cond := db.Cond("login=", usernameEmail)
	if strings.Contains(usernameEmail, "@") {
		cond = db.Cond("email=", usernameEmail)
	}
	user, err := querier.Fetch("user", cond)
	if err != nil {
		// return false, err
		//todo: log it
	}
	if user == nil {
		//todo: user error code.
		return errors.New("User not found"), nil
	}
	passwordField := user.Value("password").(fieldtype.PasswordField)
	result := util.MatchPassword(password, passwordField.Raw)
	if result {
		return nil, user
	} else {
		return errors.New("Password is wrong"), nil
	}
}
