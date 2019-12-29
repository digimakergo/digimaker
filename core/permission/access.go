//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

package permission

import (
	"context"
	"dm/core/contenttype"
	"dm/core/util"
	"dm/core/util/debug"
	"fmt"
	"log"
	"strconv"
)

//If the user has acccess given matchedData(realData here)
func HasAccessTo(userID int, operation string, realData map[string]interface{}, context context.Context) bool {
	//get permission limits
	limits, err := GetUserLimits(userID, operation, context)
	debug.Debug(context, "Limits: "+fmt.Sprintln(limits), "permission")

	if err != nil {
		log.Println(err.Error())
		return false
	}

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
	return result
}

//If the use can read the content
func CanRead(userID int, content contenttype.ContentTyper, context context.Context) bool {
	location := content.GetLocation()
	data := map[string]interface{}{
		"id":          location.ID,
		"contenttype": content.ContentType(),
		"under":       location.Path(),
	}
	if userID == content.GetLocation().Author {
		data["author"] = "self"
	}
	result := HasAccessTo(userID, "content/read", data, context)
	return result
}
