package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
)

//CanLogin check if the username/email and password matches
func CanLogin(ctx context.Context, usernameEmail string, password string) (error, contenttype.ContentTyper) {
	//todo: use username instead of 'login'
	cond := db.Cond("c.login=", usernameEmail)
	if strings.Contains(usernameEmail, "@") {
		cond = db.Cond("c.email=", usernameEmail)
	}
	user, err := query.Fetch(ctx, "user", cond)
	if err != nil {
		return err, nil
	}
	if user == nil {
		return errors.New("User not found"), nil
	}

	disabled := user.Value("disabled")
	if disabled != nil {
		if disabled.(int) == 1 {
			return errors.New("User is disabled"), nil
		}
	}

	passwordField := user.Value("password")
	result := util.MatchPassword(password, passwordField.(string))
	if result {
		return nil, user
	} else {
		return errors.New("Password is wrong"), nil
	}
}

//enable or disable user(enable = false means disable)
//todo: create a user interface
func Enable(ctx context.Context, user contenttype.ContentTyper, enable bool, userId int) error {
	disabledField := user.Value("disabled")
	if disabledField == nil {
		return errors.New("No disabled feature")
	}
	disabled := disabledField.(int)
	if disabled == 1 && enable || disabled == 0 && !enable {
		disableInt := 1
		if enable {
			disableInt = 0
		}
		_, err := Update(ctx, user, map[string]interface{}{"disabled": disableInt}, userId)
		return err
	}
	return nil
}

func FetchUserRoles(ctx context.Context, userID int) ([]contenttype.ContentTyper, error) {
	userRoleList := []permission.UserRole{}
	_, err := db.BindEntity(ctx, &userRoleList, "dm_user_role", db.Cond("user_id", userID))
	if err != nil {
		return nil, fmt.Errorf("Can not get user role by user id %v: %w", userID, err)
	}

	roleIds := []int{}
	for _, item := range userRoleList {
		roleIds = append(roleIds, item.RoleID)
	}
	list, _, err := query.List(ctx, "role", db.Cond("c.id", roleIds))
	return list, err
}
