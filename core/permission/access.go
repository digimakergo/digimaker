//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

package permission

import (
	"context"
	"fmt"
	"strconv"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

//If the user has acccess given matchedData(realData here)
func HasAccessTo(userID int, operation string, realData map[string]interface{}, context context.Context) bool {
	//get permission limits
	limits, err := GetUserLimits(userID, operation, context)
	log.Debug("Limits: "+fmt.Sprintln(limits), "permission", context)

	if err != nil {
		log.Error(err.Error(), "permission")
		return false
	}

	//match limits
	result := false
	for i, limit := range limits {
		log.Debug("Matching limit "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(limits)), "permission", context)
		policyResult, matchLog := util.MatchCondition(limit, realData)
		for _, item := range matchLog {
			log.Debug(item, "permission")
		}

		if policyResult {
			result = true
			log.Debug("Policy matched.", "permission", context)
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
	if userID == content.Value("author").(int) {
		data["author"] = "self"
	}
	result := HasAccessTo(userID, "content/read", data, context)
	return result
}
