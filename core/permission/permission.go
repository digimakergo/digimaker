//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

// Package permission implements all the permission data malipulating and generic access matching.
package permission

import (
	"strconv"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"

	"github.com/pkg/errors"
)

/*************
Policy & Permission
*************/

type Policy struct {
	Operation []string               `json:"operation"`
	LimitedTo map[string]interface{} `json:"limited_to"` //todo: use a type Limitations/Limits?
}

//Policy collection. For merge Policy list, use []Policy
type PolicyList []Policy

var policyDefinition map[string]PolicyList

func LoadPolicies() error {
	policies := map[string]PolicyList{}
	err := util.UnmarshalData(util.ConfigPath()+"/policies.json", &policies)
	if err != nil {
		return err
	}
	policyDefinition = policies
	return nil
}

func GetPolicy(identifier string) PolicyList {
	return policyDefinition[identifier]
}

func GetPolicyDefinition() map[string]PolicyList {
	return policyDefinition
}

/*************
User role
*************/
type UserRole struct {
	ID     int `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID int `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	RoleID int `boil:"role_id" json:"role_id" toml:"role_id" yaml:"role_id"`
}

func GetUserPolicies(userID int) ([]Policy, error) {
	//get usergroups
	dbHandler := db.DBHanlder()

	userRoleList := []UserRole{}
	err := dbHandler.GetEntity("dm_user_role", db.Cond("user_id", userID), nil, nil, &userRoleList)
	if err != nil {
		return nil, errors.Wrap(err, "Can not get user role by user id: "+strconv.Itoa(userID))
	}
	//get permissions
	policyList := []Policy{}
	for _, userRole := range userRoleList {
		currentPolicyList := GetRolePolicies(userRole.RoleID)
		for _, policy := range currentPolicyList {
			policyList = append(policyList, policy)
		}
	}
	return policyList, nil
}

func GetLimitsFromPolicy(policyList []Policy, operation string) []map[string]interface{} {
	var result []map[string]interface{}
	for _, policy := range policyList {
		for _, item := range policy.Operation {
			if item == operation {
				limit := policy.LimitedTo
				result = append(result, limit) //todo: nil limit is handled?
			}
		}
	}
	//todo: merge all limits first.
	return result
}

func GetRolePolicies(roleID int) PolicyList {
	role := contenttype.NewInstance("role")
	dbHandler := db.DBHanlder()
	dbHandler.GetByFields("role", "dm_role", db.Cond("c.id", roleID), nil, nil, role, false)
	if role == nil {
		log.Warning("Role doesn't exist on ID"+strconv.Itoa(roleID), "")
		return PolicyList{}
	}

	policyField := role.Value("policies").(*fieldtype.Text)
	policyStr := policyField.String.String

	policies := GetPolicy(policyStr)
	return policies
}

func AssignToUser(roleID int, userID int) error {
	//todo: check if role exisit. maybe need role entity?
	//todo: check if user exist.
	useRole := UserRole{}
	dbHandler := db.DBHanlder()
	dbHandler.GetEntity("dm_user_role", db.Cond("user_id", userID).Cond("role_id", roleID), nil, nil, &useRole)
	if useRole.ID > 0 {
		return errors.New("Already assigned.")
	}
	_, err := dbHandler.Insert("dm_user_role", map[string]interface{}{"user_id": userID, "role_id": roleID})
	if err != nil {
		return err
	}
	return nil
}
