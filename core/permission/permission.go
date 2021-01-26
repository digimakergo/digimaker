//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

// Package permission implements all the permission data malipulating and generic access matching.
package permission

import (
	"fmt"
	"strconv"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"

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
var rolePolicyMap map[string][]string

func LoadPolicies() error {
	policyRoles := struct {
		Policies map[string]PolicyList `json:"policies"`
		Roles    map[string][]string   `json:"roles"`
	}{}

	err := util.UnmarshalData(util.ConfigPath()+"/policies.json", &policyRoles)
	if err != nil {
		return err
	}
	policyDefinition = policyRoles.Policies
	rolePolicyMap = policyRoles.Roles

	for _, policies := range policyRoles.Roles {
		for _, policyIdentifer := range policies {
			if _, ok := policyDefinition[policyIdentifer]; !ok {
				return errors.New("policelist " + policyIdentifer + " doen't exist.")
			}
		}
	}
	return nil
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

//todo: cache in context?
func GetUserPolicies(userID int) ([]Policy, error) {
	dbHandler := db.DBHanlder()

	//get roles of user
	userRoleList := []UserRole{}
	err := dbHandler.GetEntity("dm_user_role", db.Cond("user_id", userID), nil, nil, &userRoleList)
	if err != nil {
		return nil, errors.Wrap(err, "Can not get user role by user id: "+strconv.Itoa(userID))
	}
	//get permissions
	policyList := []Policy{}
	roleIDs := []int{}
	for _, userRole := range userRoleList {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	currentPolicyList := GetRolePolicies(roleIDs)
	for _, policy := range currentPolicyList {
		policyList = append(policyList, policy)
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

func GetRolePolicies(roleIDs []int) []Policy {
	roles := contenttype.NewList("role")
	dbHandler := db.DBHanlder()
	dbHandler.GetByFields("role", "dm_role", db.Cond("c.id", roleIDs), nil, nil, roles, false)
	if roles == nil {
		log.Warning("Role doesn't exist on ID(s)"+fmt.Sprint(roleIDs), "")
		return PolicyList{}
	}

	roleList := contenttype.ToList("role", roles)
	policyIdentifiers := []string{}
	policies := []Policy{}
	for _, role := range roleList {
		roleIdentifierField := role.Value("identifier").(*fieldtype.Text)
		roleIdentifier := roleIdentifierField.String.String

		//loop policies under the role
		for _, policyIdentifier := range rolePolicyMap[roleIdentifier] {
			//todo: different roles may have different context(eg. target) under which the policies shouldn't merge
			if util.Contains(policyIdentifiers, policyIdentifier) {
				log.Debug("Policelist "+policyIdentifier+" is duplicated on roles. Ignored", "")
				continue
			}
			policyIdentifiers = append(policyIdentifiers, policyIdentifier)

			currentPolicies := policyDefinition[policyIdentifier]
			policies = append(policies, currentPolicies...)
		}
	}

	log.Debug("Got policylist: "+fmt.Sprintln(policyIdentifiers), "permission")
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

//Remove a user from role assignment
func RemoveAssignment(userID int, roleID int) error {
	dbHandler := db.DBHanlder()
	err := dbHandler.Delete("dm_user_role", db.Cond("user_id", userID).Cond("role_id", roleID))
	if err != nil {
		return err
	}
	return nil
}
