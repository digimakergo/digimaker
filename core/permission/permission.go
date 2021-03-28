//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

// Package permission implements all the permission data malipulating and generic access matching.
package permission

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digimakergo/digimaker/core/db"
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

type key int

func getCtxPolicyKey(userID int) key {
	return key(userID)
}

//InitPolicyContext cache the policies into provided context
func InitPolicyContext(ctx context.Context, userID int) (context.Context, error) {
	policies, err := GetUserPolicies(ctx, userID)
	if err != nil {
		return ctx, err
	}

	cacheKey := getCtxPolicyKey(userID)
	result := context.WithValue(ctx, cacheKey, policies)
	return result, nil
}

//GetUserPolicies returns policies of a user, if it's already cached in the context, return it.
//todo: Will be a powerful to support variables in policies. eg:under:"{role.under}", contenttype: "role.contenttypes"
func GetUserPolicies(ctx context.Context, userID int) ([]Policy, error) {
	//first get cache from context.
	cachedPolicies := ctx.Value(getCtxPolicyKey(userID))
	if cachedPolicies != nil {
		return cachedPolicies.([]Policy), nil
	}
	//get roles of user
	userRoleList := []UserRole{}
	_, err := db.BindEntity(context.Background(), &userRoleList, "dm_user_role", db.Cond("user_id", userID))
	if err != nil {
		return nil, errors.Wrap(err, "Can not get user role by user id: "+strconv.Itoa(userID))
	}
	//get permissions
	policyList := []Policy{}
	roleIDs := []int{}
	for _, userRole := range userRoleList {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	currentPolicyList := GetRolePolicies(ctx, roleIDs)
	for _, policy := range currentPolicyList {
		policyList = append(policyList, policy)
	}

	return policyList, nil
}

// GetLimitsFromPolicy gets all limits from a policies
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

// GetRolePolicies returns policies of role ids
func GetRolePolicies(ctx context.Context, roleIDs []int) []Policy {
	roles := db.DatamapList{}
	_, err := db.BindEntity(ctx, &roles, "dm_role", db.Cond("id", roleIDs))
	if err != nil {
		log.Error("Can not get role on role ids: "+fmt.Sprint(roleIDs), "", ctx)
		return nil
	}

	policyIdentifiers := []string{}
	policies := []Policy{}
	for _, role := range roles {
		roleIdentifier := role["identifier"].(string)

		//loop policies under the role
		for _, policyIdentifier := range rolePolicyMap[roleIdentifier] {
			//todo: different roles may have different context(eg. target) under which the policies shouldn't merge
			if util.Contains(policyIdentifiers, policyIdentifier) {
				log.Warning("Policelist "+policyIdentifier+" is duplicated on roles. Ignored", "permission", ctx)
				continue
			}
			policyIdentifiers = append(policyIdentifiers, policyIdentifier)

			currentPolicies := policyDefinition[policyIdentifier]
			policies = append(policies, currentPolicies...)
		}
	}

	log.Debug("Got policy identifiers: "+fmt.Sprintln(policyIdentifiers), "permission", ctx)
	return policies
}

// AssignToUser assigns a role to a user
func AssignToUser(ctx context.Context, roleID int, userID int) error {
	//todo: check if role exisit. maybe need role entity?
	//todo: check if user exist.
	count, err := db.Count("dm_user_role", db.Cond("user_id", userID).Cond("role_id", roleID))
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Already assigned.")
	}
	//todo: put db.Insert/update/delete into entity of UserRole(better generate automatically)
	_, err = db.Insert(ctx, "dm_user_role", map[string]interface{}{"user_id": userID, "role_id": roleID})
	if err != nil {
		return err
	}
	return nil
}

//RemoveAssignment removes a user from role assignment
func RemoveAssignment(ctx context.Context, userID int, roleID int) error {
	err := db.Delete(ctx, "dm_user_role", db.Cond("user_id", userID).Cond("role_id", roleID))
	if err != nil {
		return err
	}
	return nil
}

func init() {
	err := LoadPolicies()
	if err != nil {
		log.Fatal("Loading policies error: " + err.Error())
	}
}
