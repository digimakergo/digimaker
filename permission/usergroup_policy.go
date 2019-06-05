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

	"github.com/pkg/errors"
)

type UsergroupPolicy struct {
	ID          int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	UsergroupID string `boil:"usergroup_id" json:"usergroup_id" toml:"usergroup_id" yaml:"usergroup_id"`
	Policy      string `boil:"policy" json:"policy" toml:"policy" yaml:"policy"`
	Subtree     string `boil:"subtree" json:"subtree" toml:"subtree" yaml:"subtree"`
	Scope       string `boil:"scope" json:"scope" toml:"scope" yaml:"scope"`
	policy      Policy `boil:"-"` //cache for Policy instance
}

//Get policy detail of current usergroup policy.
func (ugPolicy UsergroupPolicy) GetPolicy() Policy {
	if len(ugPolicy.policy.Permissions) == 0 {
		ugPolicy.policy = GetPolicy(ugPolicy.Policy)
	}
	return ugPolicy.policy
}

//Get UsergroupPolicy slice based on usergroupID including inhertated permissions.
func GetPermissions(usergroupID int) ([]UsergroupPolicy, error) {
	content, err := handler.Querier().FetchByID(usergroupID) //todo: maybe better to

	if err != nil {
		fmt.Println(err) //todo: make it generic
	}
	hierarchy := content.GetLocation().Hierarchy
	ids := strings.Split(hierarchy, "/")
	dbHandler := db.DBHanlder()

	usergroupIDs := []int{}
	for _, item := range ids {
		itemInt, _ := strconv.Atoi(item)
		usergroupIDs = append(usergroupIDs, itemInt)
	}
	policyList := []UsergroupPolicy{}
	err = dbHandler.GetEnity("dm_usergroup_policy", query.Cond("usergroup_id", usergroupIDs), &policyList)
	if err != nil {
		return nil, errors.Wrap(err, "Can not fetch dm_usergroup_policy. usergroup_id :"+strings.Join(ids, ","))
	}
	return policyList, nil
}
