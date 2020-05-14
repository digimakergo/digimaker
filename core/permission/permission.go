//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

// Package permission implements all the permission data malipulating and generic access matching.
package permission

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/util"

	"github.com/pkg/errors"
)

/*************
Policy & Permission
*************/

type Permission struct {
	Operation []string               `json:"operation"`
	LimitedTo map[string]interface{} `json:"limited_to"`
}

type Policy struct {
	AssignType  []string     `json:"limited_to"`
	Permissions []Permission `json:"permissions"`
}

var policyDefinition map[string]Policy

func LoadPolicies() error {
	policies := map[string]Policy{}
	err := util.UnmarshalData(util.ConfigPath()+"/policies.json", &policies)
	if err != nil {
		return err
	}
	policyDefinition = policies
	return nil
}

func GetPolicy(identifier string) Policy {
	return policyDefinition[identifier]
}

/*************
User role
*************/
type UserRole struct {
	ID     int `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID int `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	RoleID int `boil:"role_id" json:"role_id" toml:"role_id" yaml:"role_id"`
}

func GetUserPolicies(userID int) ([]RolePolicy, error) {
	//get usergroups
	dbHandler := db.DBHanlder()

	list := []UserRole{}
	err := dbHandler.GetEntity("dm_user_role", db.Cond("user_id", userID), []string{}, &list)
	if err != nil {
		return nil, errors.Wrap(err, "Can not get user role by user id: "+strconv.Itoa(userID))
	}
	//get permissions
	policyList := []RolePolicy{}
	for _, userRole := range list {
		currentPolicyList, err := GetPermissions(userRole.RoleID)
		if err != nil {
			return nil, errors.Wrap(err, "Can not get permission on usergroup: "+strconv.Itoa(userRole.RoleID))
		}
		for _, policy := range currentPolicyList {
			policyList = append(policyList, policy)
		}
	}
	return policyList, nil
}

func GetLimitsFromPolicy(policyList []RolePolicy, operation string) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, ugPolicy := range policyList {
		policy := ugPolicy.GetPolicy()
		for _, permission := range policy.Permissions {
			for _, item := range permission.Operation {
				if item == operation {
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
	}
	//todo: merge all limits first.
	return result
}

/*************
Role policy
*************/
type RolePolicy struct {
	ID     int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	RoleID string `boil:"role_id" json:"role_id" toml:"role_id" yaml:"role_id"`
	Policy string `boil:"policy" json:"policy" toml:"policy" yaml:"policy"`
	Under  string `boil:"under" json:"under" toml:"under" yaml:"under"`
	Scope  string `boil:"scope" json:"scope" toml:"scope" yaml:"scope"`
	policy Policy `boil:"-"` //cache for Policy instance
}

//Get policy detail of current usergroup policy.
func (rp RolePolicy) GetPolicy() Policy {
	if len(rp.policy.Permissions) == 0 {
		rp.policy = GetPolicy(rp.Policy)
	}
	return rp.policy
}

//Get UsergroupPolicy slice based on usergroupID including inhertated permissions.
func GetPermissions(usergroupID int) ([]RolePolicy, error) {
	dbHandler := db.DBHanlder()
	location := contenttype.Location{}
	//todo: maybe better to use content id
	err := dbHandler.GetEntity("dm_location", db.Cond("id", usergroupID), []string{}, &location) //note: use this instead of handler.Querier() to avoid cycle dependency because handler package rely on permission
	if err != nil {
		fmt.Println(err) //todo: make it generic
	}

	hierarchy := location.Hierarchy
	ids := strings.Split(hierarchy, "/")

	roleIDs := []int{}
	for _, item := range ids {
		itemInt, _ := strconv.Atoi(item)
		roleIDs = append(roleIDs, itemInt)
	}
	policyList := []RolePolicy{}
	err = dbHandler.GetEntity("dm_role_policy", db.Cond("role_id", roleIDs), []string{}, &policyList)
	if err != nil {
		return nil, errors.Wrap(err, "Can not fetch dm_usergroup_policy. usergroup_id :"+strings.Join(ids, ","))
	}
	return policyList, nil
}
