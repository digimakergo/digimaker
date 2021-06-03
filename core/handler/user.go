package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/permission"
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
func Enable(user contenttype.ContentTyper, enable bool, userId int, ctx context.Context) error {
	disabledField := user.Value("disabled")
	if disabledField == nil {
		return errors.New("No disabled feature")
	}
	disabled := disabledField.(*fieldtype.Checkbox).FieldValue().(int)
	if disabled == 1 && enable || disabled == 0 && !enable {
		handler := ContentHandler{}
		handler.Context = ctx
		disableInt := 1
		if enable {
			disableInt = 0
		}
		result, _, err := handler.Update(user, map[string]interface{}{"disabled": disableInt}, userId)
		if result {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

func FetchUserRoles(ctx context.Context, userID int) ([]contenttype.ContentTyper, error) {
	userRoleList := []permission.UserRole{}

	dbHandler := db.DBHanlder()
	err := dbHandler.GetEntity("dm_user_role", db.Cond("user_id", userID), nil, nil, &userRoleList)
	if err != nil {
		return nil, err
	}
	roleIds := []int{}
	for _, item := range userRoleList {
		roleIds = append(roleIds, item.RoleID)
	}
	querier := Querier()
	list, _, err := querier.List("role", db.Cond("c.id", roleIds), nil, nil, false)
	return list, err
}
