//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

package handler

import (
	"context"
	"dm/permission"
	"dm/util"
	"dm/util/debug"
	"fmt"
	"strconv"
)

//Check if module & action & data is under given policy list
//currentData example: {scope: "site1", under: [1,2,3,4,5], contenttype: "folder" }
func HasAccessTo(ugPolicyList []permission.UsergroupPolicy, module string, action string, currentData map[string]interface{}, context context.Context) bool {
	result := false
	for index, ugPolicy := range ugPolicyList {
		debug.Debug(context, "Matching policy "+strconv.Itoa(index+1)+"/"+strconv.Itoa(len(ugPolicyList))+" "+ugPolicy.Policy+" with module:"+module+", action:"+action+", data: "+fmt.Sprintln(currentData), "permission")
		policyResult := true
		policy := ugPolicy.GetPolicy()
		if ugPolicy.Under != "" { //if under is not "" and data doesn't match go to next policy

		}
		if ugPolicy.Scope != "" { //if scope is not "" and data doesn't match go to next policy

		}
		for i, permission := range policy.Permissions {
			debug.Debug(context, "Matching permission "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(policy.Permissions)), "permission")
			if permission.Action != action || permission.Module != module {
				policyResult = false
				debug.Debug(context, "Module or action doesn't match. expected module:"+permission.Module+", action:"+permission.Action, "permission")
			} else {
				var matchLog []string
				policyResult, matchLog = util.MatchCondition(permission.LimitedTo, currentData)
				for _, item := range matchLog {
					debug.Debug(context, item, "permission")
				}
			}
			if policyResult {
				result = true
				debug.Debug(context, "Policy matched.", "permission")
				break
			}
		}
		if result {
			break
		}
	}
	return result
}
