//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

// Package permission implements all the permission data malipulating and generic access matching.
package permission

import (
	"context"
	"fmt"
	"strings"

	"github.com/digimakergo/digimaker/core/config"
	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype/fieldtypes"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"

	"errors"
)

/*************
Policy & Permission
*************/

type AccessLimit map[string]interface{}

type Policy struct {
	Operation []string    `json:"operation"`
	LimitedTo AccessLimit `json:"limited_to"` //todo: use a type Limitations/Limits?
}

//Policy collection. For merge Policy list, use []Policy
type PolicyList []Policy

var policyDefinition map[string]PolicyList
var rolePolicyMap map[string][]string
var roleVariables []string

func LoadPolicies() error {
	policyRoles := struct {
		Policies      map[string]PolicyList `json:"policies"`
		Roles         map[string][]string   `json:"roles"`
		RoleVariables []string              `json:"role_variables"`
	}{}

	err := util.UnmarshalData(config.ConfigPath()+"/policies.json", &policyRoles)
	if err != nil {
		return err
	}
	policyDefinition = policyRoles.Policies
	rolePolicyMap = policyRoles.Roles
	roleVariables = policyRoles.RoleVariables

	for _, policies := range policyRoles.Roles {
		for _, policyIdentifer := range policies {
			if _, ok := policyDefinition[policyIdentifer]; !ok {
				return errors.New("policelist " + policyIdentifer + " doen't exist.")
			}
		}
	}
	return nil
}

func GetRoles() []string {
	roles := []string{}
	for role, _ := range rolePolicyMap {
		roles = append(roles, role)
	}
	return roles
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

type Role struct {
	ID         int            `boil:"id" json:"id" toml:"id" yaml:"id"`
	Identifier string         `boil:"identifier" json:"identifier" toml:"identifier" yaml:"identifier"`
	Parameters fieldtypes.Map `boil:"parameters" json:"parameters" toml:"parameters" yaml:"parameters"`
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
//todo: Support field condition, eg: {"contenttype": "article","field_category": "news"} - policy that a user can read article whose category is news.
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
		return nil, fmt.Errorf("Can not get user role by user id %v: %w", userID, err)
	}

	//get permissions
	result := []Policy{}
	roleIDs := []int{}
	for _, userRole := range userRoleList {
		roleIDs = append(roleIDs, userRole.RoleID)
	}
	roles := contenttype.NewList("role")
	_, err = db.BindContent(ctx, roles, "role", db.Cond("c.id", roleIDs))
	if err != nil {
		return nil, err
	}

	roleList := contenttype.ToList("role", roles)
	for _, role := range roleList {
		policyList := GetRolePolicies(ctx, role.Value("identifier").(string))
		params := map[string]interface{}{}
		for _, field := range roleVariables {
			params[field] = role.Value(field)
		}
		for _, policy := range policyList {
			//assign role field values to policy variables
			for key, value := range policy.LimitedTo {
				if v, ok := value.(string); ok && strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
					vars := util.GetStrVar(v)
					varName := vars[0]
					if varValue, ok := params[varName]; ok {
						policy.LimitedTo[key] = varValue
						log.Debug("Set variable "+varName+" to "+key+": "+fmt.Sprint(varValue), "permission", ctx)
					}
				}
			}
			result = append(result, policy)
		}
	}
	return result, nil
}

// GetLimitsFromPolicy gets all limits from a policies
func GetLimitsFromPolicy(policyList []Policy, operation string) []AccessLimit {
	var result []AccessLimit
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
func GetRolePolicies(ctx context.Context, role string) []Policy {

	policyIdentifiers := []string{}
	policies := []Policy{}
	//loop policies under the role
	for _, policyIdentifier := range rolePolicyMap[role] {
		//todo: different roles may have different context(eg. target) under which the policies shouldn't merge
		if util.Contains(policyIdentifiers, policyIdentifier) {
			log.Warning("Policelist "+policyIdentifier+" is duplicated on roles. Ignored", "permission", ctx)
			continue
		}
		policyIdentifiers = append(policyIdentifiers, policyIdentifier)

		currentPolicies := policyDefinition[policyIdentifier]
		policies = append(policies, currentPolicies...)
	}

	log.Debug("Got policy identifiers: "+fmt.Sprintln(policyIdentifiers), "permission", ctx)
	return policies
}

// AssignToUser assigns a role to a user
func AssignToUser(ctx context.Context, roleID int, userID int) error {

	//todo: check if user exist.
	count, err := db.Count("dm_user_role", db.Cond("user_id", userID).Cond("role_id", roleID))
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Already assigned")
	}

	//todo: put db.Insert/update/delete into entity of UserRole(better generate automatically)
	_, err = db.Insert(ctx, "dm_user_role", map[string]interface{}{"user_id": userID, "role_id": roleID})
	if err != nil {
		log.Error("Assign to user: "+err.Error(), "")
		return errors.New("Error when inserting access data")
	}
	return nil
}

//RemoveAssignment removes a user from role assignment
func RemoveAssignment(ctx context.Context, userID int, role string) error {
	err := db.Delete(ctx, "dm_user_role", db.Cond("user_id", userID).Cond("role_id", role))
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
