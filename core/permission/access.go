//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

package permission

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

type AccessType int
type MatchData map[string]interface{}

const (
	AccessFull      = 1
	AccessNo        = 0
	AccessWithLimit = 2
)

//Get user's limits.
//empty result means no access - not no limit, while a empty limit(empty map) in the slice means no limit(can do anything)
// return access list, access type, error
// if accessType is AccessWithLimit, there must be valid values in the access list
func GetUserAccess(userID int, operation string, context context.Context) (AccessType, []map[string]interface{}, error) {
	policyList, err := GetUserPolicies(userID)
	log.Debug("Got policy list: "+fmt.Sprint(policyList), "permission")
	if err != nil {
		return AccessNo, nil, errors.Wrap(err, "Error when fetching policy list for user:"+strconv.Itoa(userID))
	}
	//todo: cache limits to user, and cache anoymous globally.
	accessList := GetLimitsFromPolicy(policyList, operation)
	log.Debug("Got access list of "+operation+": "+fmt.Sprint(accessList), "permission")

	//empty access list
	if accessList == nil {
		log.Debug("No access.", "permission")
		fmt.Println(userID)
		return AccessNo, accessList, nil
	}
	//check if there is an access with no limit.
	for i, access := range accessList {
		if access == nil {
			log.Debug("Full access on "+strconv.Itoa(i+1), "permission")
			return AccessFull, accessList, nil
		}
	}
	return AccessWithLimit, accessList, nil
}

//If the user has acccess given matchedData(realData here)
func HasAccessTo(userID int, operation string, realData MatchData, context context.Context) bool {
	//get permission limits
	accessType, limits, err := GetUserAccess(userID, operation, context)

	if err != nil {
		log.Error(err.Error(), "permission")
		return false
	}

	if accessType == AccessFull {
		return true
	}

	if accessType == AccessNo {
		return false
	}

	log.Debug("Access limits: "+fmt.Sprintln(limits), "permission", context)

	//match limits
	for i, limit := range limits {
		log.Debug("Matching limit "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(limits)), "permission", context)
		policyResult, matchLog := util.MatchCondition(limit, realData)
		for _, item := range matchLog {
			log.Debug(item, "permission", context)
		}

		if policyResult {
			log.Debug("Policy matched.", "permission", context)
			return true
		}
	}
	return false
}

//todo: support more
//If the use can read the content
func CanRead(userID int, content contenttype.ContentTyper, context context.Context) bool {
	location := content.GetLocation()
	data := map[string]interface{}{"contenttype": content.ContentType()}
	if location != nil {
		data["id"] = location.ID
		data["under"] = location.Path()
	}

	author := content.Value("author")
	if author != nil && (userID == author.(int)) {
		data["author"] = "self"
	}
	result := HasAccessTo(userID, "content/read", data, context)
	return result
}

func CanDelete(ctx context.Context, content contenttype.ContentTyper, userId int) bool {
	return HasAccessTo(userId, "content/delete", GetMatchData(content, userId), ctx)
}

func CanCreate(ctx context.Context, parent contenttype.ContentTyper, contenttype string, userId int) bool {
	data := GetMatchData(parent, userId)
	data["contenttype"] = contenttype
	return HasAccessTo(userId, "content/create", data, ctx)
}

func GetMatchData(content contenttype.ContentTyper, userId int) MatchData {
	def := content.Definition()
	data := MatchData{}
	data["contenttype"] = content.ContentType()
	if def.HasLocation {
		location := content.GetLocation()
		data["id"] = location.ID
		data["under"] = location.Path()
		author := content.Value("author")
		if author != nil && (userId == author.(int)) {
			data["author"] = "self"
		}
	}
	return data
}
