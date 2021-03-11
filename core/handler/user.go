package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
)

//CanLogin check if the username/email and password matches
func CanLogin(ctx context.Context, usernameEmail string, password string) (error, contenttype.ContentTyper) {
	//todo: use username instead of 'login'
	cond := db.Cond("login=", usernameEmail)
	if strings.Contains(usernameEmail, "@") {
		cond = db.Cond("email=", usernameEmail)
	}
	user, err := query.Fetch(ctx, "user", cond)
	if err != nil {
		// return false, err
		//todo: log it
	}
	if user == nil {
		//todo: user error code.
		return errors.New("User not found"), nil
	}

	disabled := user.Value("disabled")
	if disabled != nil {
		if disabled.(*fieldtype.Checkbox).FieldValue().(int) == 1 {
			return errors.New("User is disabled"), nil
		}
	}

	passwordField := user.Value("password").(*fieldtype.Password)
	result := util.MatchPassword(password, passwordField.FieldValue().(string))
	if result {
		return nil, user
	} else {
		return errors.New("Password is wrong"), nil
	}
}

//enable or disable user(enable = false means disable)
func Enable(ctx context.Context, user contenttype.ContentTyper, enable bool, userId int) error {
	disabledField := user.Value("disabled")
	if disabledField == nil {
		return errors.New("No disabled feature")
	}
	disabled := disabledField.(*fieldtype.Checkbox).FieldValue().(int)
	if disabled == 1 && enable || disabled == 0 && !enable {
		handler := ContentHandler{}
		disableInt := 1
		if enable {
			disableInt = 0
		}
		result, _, err := handler.Update(ctx, user, map[string]interface{}{"disabled": disableInt}, userId)
		if result {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}
