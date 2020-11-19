package handler

import (
	"errors"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/util"
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
	passwordField := user.Value("password").(*fieldtype.Password)
	result := util.MatchPassword(password, passwordField.FieldValue().(string))
	if result {
		return nil, user
	} else {
		return errors.New("Password is wrong"), nil
	}
}
