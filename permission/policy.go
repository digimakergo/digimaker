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

type Permission struct {
	Module     string
	Action     []string
	Limitation interface{}
}

type Policy struct {
	Identifier  string
	Name        string
	LimitedTo   []string
	Permissions []Permission
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

	//todo: cache all policy and permissions
	fmt.Println(policyList)
}
