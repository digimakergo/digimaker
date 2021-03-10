//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

package permission

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/pkg/errors"
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
func GetUserAccess(ctx context.Context, userID int, operation string) (AccessType, []map[string]interface{}, error) {
	policyList, err := GetUserPolicies(ctx, userID)
	log.Debug("Got policy list: "+fmt.Sprint(policyList), "permission", ctx)
	if err != nil {
		return AccessNo, nil, errors.Wrap(err, "Error when fetching policy list for user:"+strconv.Itoa(userID))
	}
	//todo: cache limits to user, and cache anoymous globally.
	limitList := GetLimitsFromPolicy(policyList, operation)
	log.Debug("Got access list of "+operation+": "+fmt.Sprint(limitList), "permission", ctx)

	//empty access list
	if limitList == nil {
		log.Debug("No access.", "permission", ctx)
		return AccessNo, limitList, nil
	}
	//check if there is an access with no limit.
	for i, limit := range limitList {
		if limit == nil {
			log.Debug("Full access on "+strconv.Itoa(i+1), "permission", ctx)
			return AccessFull, limitList, nil
		}
	}
	return AccessWithLimit, limitList, nil
}

//If the user has acccess given matchedData(realData here)
func HasAccessTo(ctx context.Context, userID int, operation string, realData MatchData) bool {
	//get permission limits
	accessType, limits, err := GetUserAccess(ctx, userID, operation)

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

	log.Debug("Access limits: "+fmt.Sprintln(limits), "permission", ctx)

	//match limits
	for i, limit := range limits {
		log.Debug("Matching limit "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(limits)), "permission", ctx)
		policyResult, matchLog := util.MatchCondition(limit, realData)
		for _, item := range matchLog {
			log.Debug(item, "permission", ctx)
		}

		if policyResult {
			log.Debug("Policy matched.", "permission", ctx)
			return true
		}
	}
	return false
}

//todo: support more
//If the use can read the content
func CanRead(ctx context.Context, userID int, content contenttype.ContentTyper) bool {
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
	result := HasAccessTo(ctx, userID, "content/read", data)
	return result
}

func CanDelete(ctx context.Context, content contenttype.ContentTyper, userId int) bool {
	return HasAccessTo(ctx, userId, "content/delete", GetMatchData(content, userId))
}

func CanCreate(ctx context.Context, parent contenttype.ContentTyper, contenttype string, userId int) bool {
	data := GetMatchData(parent, userId)
	data["contenttype"] = contenttype
	return HasAccessTo(ctx, userId, "content/create", data)
}

func CanUpdate(ctx context.Context, content contenttype.ContentTyper, userId int) bool {
	data := GetMatchData(content, userId)
	return HasAccessTo(ctx, userId, "content/update", data)
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

//Get update fields for this user. If content is a user content, it supports "cid":"self"
//return fields list, if all matches, return ["*"]
//Note: fields can not be set in different rules, meaning first matches will get the fields
func GetUpdateFields(ctx context.Context, content contenttype.ContentTyper, userID int) ([]string, error) {
	accessType, accessMap, err := GetUserAccess(ctx, userID, "content/update_fields")
	if err != nil {
		return nil, err
	}
	result := []string{}
	if accessType == AccessFull {
		result = append(result, "*")
	} else if accessType == AccessWithLimit {
		matchData := GetMatchData(content, userID)
		fmt.Println(matchData)
		if content.ContentType() == "user" && content.GetCID() == userID { //todo: make "user" configurable
			matchData["cid"] = "self"
		}
		matchData["fields"] = nil          //todo: maybe a better way to get fields instead of using nil to match pass
		for _, limits := range accessMap { //todo: is it sure it will be the first one(golang's random order on map)?
			matched, _ := util.MatchCondition(limits, matchData) //todo: set log
			if matched {
				fieldsI := limits["fields"]
				result = util.InterfaceToStringArray(fieldsI.([]interface{}))
				break
			}
		}
	}
	return result, nil
}
