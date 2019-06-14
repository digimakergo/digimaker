//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

import (
	"dm/dm/contenttype"
	"dm/dm/db"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type UsergroupPolicy struct {
	ID          int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	UsergroupID string `boil:"usergroup_id" json:"usergroup_id" toml:"usergroup_id" yaml:"usergroup_id"`
	Policy      string `boil:"policy" json:"policy" toml:"policy" yaml:"policy"`
	Under       string `boil:"under" json:"under" toml:"under" yaml:"under"`
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
	dbHandler := db.DBHanlder()
	location := contenttype.Location{}
	//todo: maybe better to use content id
	err := dbHandler.GetEntity("dm_location", db.Cond("id", usergroupID), &location) //note: use this instead of handler.Querier() to avoid cycle dependency because handler package rely on permission
	if err != nil {
		fmt.Println(err) //todo: make it generic
	}

	hierarchy := location.Hierarchy
	ids := strings.Split(hierarchy, "/")

	usergroupIDs := []int{}
	for _, item := range ids {
		itemInt, _ := strconv.Atoi(item)
		usergroupIDs = append(usergroupIDs, itemInt)
	}
	policyList := []UsergroupPolicy{}
	err = dbHandler.GetEntity("dm_usergroup_policy", db.Cond("usergroup_id", usergroupIDs), &policyList)
	if err != nil {
		return nil, errors.Wrap(err, "Can not fetch dm_usergroup_policy. usergroup_id :"+strings.Join(ids, ","))
	}
	return policyList, nil
}
