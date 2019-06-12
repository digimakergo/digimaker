//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

import (
	"dm/dm/util"
)

type Permission struct {
	Module    string                 `json:"module"`
	Action    string                 `json:"action"`
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
