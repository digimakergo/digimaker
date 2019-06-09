//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

import (
	"context"
	"dm/db"
	"dm/query"
	"dm/util/debug"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type UserUsergroup struct {
	ID          int `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID      int `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	UsergroupID int `boil:"usergroup_id" json:"usergroup_id" toml:"usergroup_id" yaml:"usergroup_id"`
}

func GetUserPolicies(userID int) ([]UsergroupPolicy, error) {
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

//Get user's limits
func GetUserLimits(userID int, module string, action string, context context.Context) ([]map[string]interface{}, error) {
	policyList, err := GetUserPolicies(userID)
	debug.Debug(context, "Got policy list: "+fmt.Sprintln(policyList), "permission")
	if err != nil {
		return nil, errors.Wrap(err, "Error when fetching policy list for user:"+strconv.Itoa(userID))
	}
	//todo: cache limits to user, and cache anoymous globally.
	result := GetLimitsFromPolicy(policyList, module, action)
	return result, nil
}

func GetLimitsFromPolicy(policyList []UsergroupPolicy, module string, action string) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, ugPolicy := range policyList {
		policy := ugPolicy.GetPolicy()
		for _, permission := range policy.Permissions {
			if permission.Module == module && permission.Action == action {
				limit := permission.LimitedTo
				if ugPolicy.Scope != "" {
					limit["scope"] = ugPolicy.Scope
				}
				if ugPolicy.Under != "" {
					limit["under"] = ugPolicy.Under
				}
				result = append(result, limit)
			}
		}
	}
	return result
}
