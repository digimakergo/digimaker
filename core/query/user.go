package query

import (
	"context"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/permission"
)

func FetchUserRoles(ctx context.Context, userID int) ([]contenttype.ContentTyper, error) {
	userRoleList := []permission.UserRole{}
	_, err := db.BindEntity(context.Background(), &userRoleList, "dm_user_role", db.Cond("user_id", userID))
	if err != nil {
		return nil, err
	}
	roleIds := []int{}
	for _, item := range userRoleList {
		roleIds = append(roleIds, item.RoleID)
	}
	list, _, err := List(ctx, "role", db.Cond("c.id", roleIds))
	return list, err
}
