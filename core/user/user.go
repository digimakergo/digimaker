package user

import (
	"dm/core/contenttype"
	"dm/core/db"
	"dm/core/fieldtype"
	"dm/core/handler"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//CanLogin check if the username/email and password matches
func CanLogin(usernameEmail string, password string) (error, contenttype.ContentTyper) {
	querier := handler.Querier()

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
	passwordField := user.Value("password").(fieldtype.TextField)
	result := MatchPassword(password, passwordField.Raw)
	if result {
		return nil, user
	} else {
		return errors.New("Password is wrong"), nil
	}
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14) //todo: make cost configable
	return string(hash), err
}

func MatchPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}
