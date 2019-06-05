//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

import (
	"dm/db"
	"dm/query"
	"strconv"

	"github.com/pkg/errors"
)

type UserUsergroup struct {
	ID          int `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID      int `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	UsergroupID int `boil:"usergroup_id" json:"usergroup_id" toml:"usergroup_id" yaml:"usergroup_id"`
}

func GetUserPermission(userID int) ([]UsergroupPolicy, error) {
	//get usergroups
	dbHandler := db.DBHanlder()

	list := []UserUsergroup{}
	err := dbHandler.GetEntity("dm_user_usergroup", query.Cond("user_id", userID), &list)
	if err != nil {
		return nil, errors.Wrap(err, "Can not get user group by user id: "+strconv.Itoa(userID))
	}
	//get permissions
	policyList := []UsergroupPolicy{}
	for _, userUsergroup := range list {
		currentPolicyList, err := GetPermissions(userUsergroup.UsergroupID)
		if err != nil {
			return nil, errors.Wrap(err, "Can not get permission on usergroup: "+strconv.Itoa(userUsergroup.UsergroupID))
		}
		for _, policy := range currentPolicyList {
			policyList = append(policyList, policy)
		}
	}
	return policyList, nil
}
