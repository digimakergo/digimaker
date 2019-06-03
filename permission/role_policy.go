//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

import (
	"dm/db"
	"dm/handler"
	"dm/query"
	"fmt"
	"strconv"
	"strings"
)

type RolePolicy struct {
	ID      int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	RoleID  string `boil:"role_id" json:"role_id" toml:"role_id" yaml:"role_id"`
	Policy  string `boil:"policy" json:"policy" toml:"policy" yaml:"policy"`
	Subtree string `boil:"subtree" json:"subtree" toml:"subtree" yaml:"subtree"`
	Scope   string `boil:"scope" json:"scope" toml:"scope" yaml:"scope"`
}

func GetPermissions(locationID int) {
	content, err := handler.Querier().FetchByID(locationID)

	if err != nil {
		fmt.Println(err)
	}
	hierarchy := content.GetLocation().Hierarchy
	ids := strings.Split(hierarchy, "/")
	ids = append(ids, strconv.Itoa(locationID))
	dbHandler := db.DBHanlder()
	//todo: use join here
	list := []RoleAssignment{}
	err = dbHandler.GetEnity("dm_role_assignment", query.Cond("assign_to", ids), &list)
	if err != nil {
		fmt.Println(err)
	}

	roleIDs := []int{}
	for _, item := range list {
		roleIDs = append(roleIDs, item.RoleID)
	}
	policyList := []RolePolicy{}
	err = dbHandler.GetEnity("dm_role_policy", query.Cond("id", roleIDs), &policyList)
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range policyList {
		policyIdentifier := item.Policy
		policy := GetPolicy(policyIdentifier)
		fmt.Println(policy.Permissions)
	}
	//todo: cache all policy and permissions
	fmt.Println(policyList)
}
