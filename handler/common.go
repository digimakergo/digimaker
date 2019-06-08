//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

package handler

import (
	"context"
	"dm/contenttype"
	"dm/permission"
	"dm/util"
	"dm/util/debug"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

//If the user has acccess given matchedData(realData here)
func HasAccessTo(userID int, module string, action string, realData map[string]interface{}, context context.Context) (bool, error) {
	//get permission limits
	policyList, err := permission.GetUserPolicies(userID)
	debug.Debug(context, "Policy list: "+fmt.Sprintln(policyList), "permission")
	if err != nil {
		return false, errors.Wrap(err, "Error when fetching policy list for user:"+strconv.Itoa(userID))
	}
	limits := permission.GetLimitsFromPolicy(policyList, module, action)
	debug.Debug(context, "Limits: "+fmt.Sprintln(limits), "permission")

	//match limits
	result := false
	for i, limit := range limits {
		debug.Debug(context, "Matching limit "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(limits)), "permission")
		policyResult, matchLog := util.MatchCondition(limit, realData)
		for _, item := range matchLog {
			debug.Debug(context, item, "permission")
		}

		if policyResult {
			result = true
			debug.Debug(context, "Policy matched.", "permission")
			break
		}
	}
	return result, nil
}

//If the use can read the content
func CanRead(userID int, content contenttype.ContentTyper, context context.Context) (bool, error) {
	location := content.GetLocation()
	data := map[string]interface{}{
		"id":          location.ID,
		"contenttype": content.ContentType(),
		"under":       location.Path(),
	}
	result, err := HasAccessTo(userID, "content", "read", data, context)
	if err != nil {
		return false, err
	}
	return result, nil
}
