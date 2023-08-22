//Author xc, Created on 2019-06-01 20:00
//{COPYRIGHTS}

package permission

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
)

type AccessType int
type MatchData map[string]interface{}

const (
	AccessFull      = 1
	AccessNo        = 0
	AccessWithLimit = 2
)

var allowedFieldtypes []string = []string{"select", "radio", "checkbox"}

var cachedPolicyLocation map[int]contenttype.Location = map[int]contenttype.Location{}

//Get user's limits.
//empty result means no access - not no limit, while a empty limit(empty map) in the slice means no limit(can do anything)
// return access list, access type, error
// if accessType is AccessWithLimit, there must be valid values in the access list
func GetUserAccess(ctx context.Context, userID int, operation string) (AccessType, []AccessLimit, error) {
	policyList, err := GetUserPolicies(ctx, userID)
	log.Debug("Got policy list: "+fmt.Sprint(policyList), "permission", ctx)
	if err != nil {
		return AccessNo, nil, fmt.Errorf("Error when fetching policy list for user %v: %w", userID, err)
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

//If the user has acccess given data(targetData here)
//If realData is empty, just check if the user has given operation(can be full access or partly access to that operation)
func HasAccessTo(ctx context.Context, userID int, operation string, targetData ...MatchData) bool {
	result, _ := AccessMatched(ctx, userID, operation, targetData...)
	return result
}

//if it returns true, also it return the matched limit, false doesn't return limit. also full access doesn't return
func AccessMatched(ctx context.Context, userID int, operation string, targetData ...MatchData) (bool, AccessLimit) {
	//get permission limits
	accessType, limits, err := GetUserAccess(ctx, userID, operation)

	if err != nil {
		log.Error(err.Error(), "permission")
		return false, nil
	}

	if accessType == AccessFull {
		return true, nil
	}

	//When match data is not provided, if there is partly access it will be success.
	if len(targetData) == 0 {
		if accessType == AccessWithLimit {
			return true, nil
		} else {
			return false, nil
		}
	}

	if accessType == AccessNo {
		return false, nil
	}

	log.Debug("Access limits: "+fmt.Sprintln(limits), "permission", ctx)

	//match limits
	for i, limit := range limits {
		log.Debug("Matching limit "+strconv.Itoa(i+1)+"/"+strconv.Itoa(len(limits)), "permission", ctx)
		policyResult, matchLog := util.MatchCondition(limit, targetData[0])
		for _, item := range matchLog {
			log.Debug(item, "permission", ctx)
		}

		if policyResult {
			log.Debug("Policy matched.", "permission", ctx)
			return true, limit
		}
	}
	return false, nil
}

func getDefaultReadMatchData() MatchData {
	return MatchData{"contenttype": nil, "id": nil, "under": nil, "author": nil}
}

//If the use can read the content
//support keys: contenttype, id(locaton id), under, author(include "self")
func CanRead(ctx context.Context, userID int, content contenttype.ContentTyper) bool {
	def := content.Definition()
	data := getDefaultReadMatchData()
	data["contenttype"] = content.ContentType()
	if def.HasLocation {
		location := content.GetLocation()
		if def.HasLocation {
			data["id"] = location.ID
		}
		data["under"] = location.Path()
	}

	data["author"] = getAuthorMatchData(content, userID)
	result := HasAccessTo(ctx, userID, "content/read", data)
	return result
}

//support keys: contenttype, id(locaton id), under, author(include "self")
func CanDelete(ctx context.Context, content contenttype.ContentTyper, userId int) bool {
	return HasAccessTo(ctx, userId, "content/delete", getDeleteMatchData(content, userId))
}

//support keys: contenttype, id(parent locaton id), under, parent author(include "self")
func CanCreate(ctx context.Context, parent contenttype.ContentTyper, contenttype string, fields []string, userId int) bool {
	data := getCreateMatchData(parent, contenttype, fields, userId)
	return HasAccessTo(ctx, userId, "content/create", data)
}

//support keys: contenttype, id(locaton id), under, author(include "self")
func CanUpdate(ctx context.Context, content contenttype.ContentTyper, fields []string, userId int) bool {
	data := getUpdateMatchData(content, fields, userId)
	return HasAccessTo(ctx, userId, "content/update", data)
}

func GetUpdateFields(ctx context.Context, content contenttype.ContentTyper, userId int) ([]string, error) {
	data := getUpdateMatchData(content, []string{}, userId)
	matched, limit := AccessMatched(ctx, userId, "content/update", data)
	result := []string{}
	if matched {
		if limit == nil {
			result = content.Definition().FieldIdentifierList
		} else {
			if _, ok := limit["fields"]; ok {
				fields := limit["fields"].(map[string]interface{})
				subset := fields["subset"]
				for _, v := range subset.([]interface{}) {
					result = append(result, v.(string))
				}
			} else {
				result = content.Definition().FieldIdentifierList
			}
		}
	} else {
		return nil, errors.New("No access to update")
	}
	return result, nil
}

func getCreateMatchData(parent contenttype.ContentTyper, contenttype string, fields []string, userId int) MatchData {
	data := MatchData{}
	if parent != nil {
		def := parent.Definition()

		data["parent_contenttype"] = parent.ContentType()
		for key, v := range def.FieldMap {
			if util.Contains(allowedFieldtypes, v.FieldType) {
				data["parent/"+key] = parent.Value(key)
			}
		}
		data["parent_author"] = getAuthorMatchData(parent, userId)
		if def.HasLocation {
			location := parent.GetLocation()
			data["parent_id"] = location.ID
			data["under"] = location.Path()
		}
	}

	data["contenttype"] = contenttype
	data["fields"] = getFieldMatch(fields)
	return data
}

func getUpdateMatchData(content contenttype.ContentTyper, fields []string, userId int) MatchData {
	def := content.Definition()
	data := MatchData{}
	data["contenttype"] = content.ContentType()
	if def.HasLocation {
		location := content.GetLocation()
		data["id"] = location.ID
		data["under"] = location.Path()
	}
	data["author"] = getAuthorMatchData(content, userId)

	if content.ContentType() == "user" && content.GetID() == userId {
		data["user"] = "self"
	}
	data["fields"] = getFieldMatch(fields)
	return data
}

//if fields is empty, use nil - meaning alway match "fields" limit
func getFieldMatch(fields []string) interface{} {
	var result interface{}
	if len(fields) > 0 {
		result = fields
	} else {
		result = nil
	}
	return result
}

func getDeleteMatchData(content contenttype.ContentTyper, userId int) MatchData {
	def := content.Definition()
	data := MatchData{}
	data["contenttype"] = content.ContentType()
	if def.HasLocation {
		location := content.GetLocation()
		data["under"] = location.Path()
	}
	data["author"] = getAuthorMatchData(content, userId)
	return data
}

func getAuthorMatchData(content contenttype.ContentTyper, userID int) string {
	author := content.GetMetadata().Author
	if author != 0 && (userID == author) {
		return "self"
	} else {
		return strconv.Itoa(author)
	}
}

// add condition from permission.
// so if matched with limit, add that limit to condition
// if matches with a empty limit(if there is), return empty(meaning no limit)
// if doesn't match, return a False condition(no result in query)
func GetListCondition(ctx context.Context, userID int, contentType string, parent contenttype.ContentTyper) db.Condition {
	accessType, limits, err := GetUserAccess(ctx, userID, "content/read")
	if err != nil {
		log.Error(err, "permission", ctx)
		return db.FalseCond()
	}

	if accessType == AccessNo {
		return db.FalseCond()
	}

	if accessType == AccessFull {
		return db.EmptyCond()
	}

	result := db.EmptyCond()
	if accessType == AccessWithLimit {
		//add conditions based on limits
		matchData := getDefaultReadMatchData()
		matchData["contenttype"] = contentType

		for _, limit := range limits {
			itemMatched, matchLog := util.MatchCondition(limit, matchData)
			if !itemMatched {
				continue
			}

			log.Debug("Matching permission for query", "permission-query", ctx)
			for _, item := range matchLog {
				log.Debug(item, "permission-query", ctx)
			}

			log.Debug("Matched for query", "permission-query", ctx)

			currentCond := db.EmptyCond()
			//under
			if under, ok := limit["under"]; ok {

				switch under.(type) {
				case int, float64:
					if f, ok := under.(float64); ok {
						under = int(f)
					}
					underLocation := GetPolicyLocation(under.(int))
					//empty parent will use 'under' policy
					if parent == nil {
						currentCond = db.Cond("l.hierarchy like", underLocation.Hierarchy+"/%")
					} else {
						parentLocation := parent.GetLocation()
						if contenttype.IsUnderLocation(underLocation, *parentLocation) {
							currentCond = currentCond.Cond("l.hierarchy like", underLocation.Hierarchy+"/%")
						}
					}
				case []interface{}:
					for _, item := range under.([]interface{}) {
						//empty parent will use 'under' policy
						if parent == nil {
							underLocation := GetPolicyLocation(int(item.(float64))) //todo: improve to use []int instead of []interface{}
							currentCond = currentCond.Or(db.Cond("l.hierarchy like", underLocation.Hierarchy+"/%"))
						} else {
							parentLocation := parent.GetLocation()
							underLocation := GetPolicyLocation(int(item.(float64)))
							if contenttype.IsUnderLocation(underLocation, *parentLocation) {
								currentCond = currentCond.Or(db.Cond("l.hierarchy like", underLocation.Hierarchy+"/%"))
							}
						}
					}
				}
			}
			//author
			if author, ok := limit["author"]; ok {
				if author.(string) == "self" {
					currentCond = currentCond.Cond("c._author", userID)
				}
			}

			result = result.Or(currentCond)
		}
	}

	return result
}

func GetPolicyLocation(id int) contenttype.Location {
	if _, ok := cachedPolicyLocation[id]; !ok {
		location, _ := contenttype.GetLocationByID(id)
		if location.ID == 0 {
			log.Error("Location with id"+strconv.Itoa(id)+"Not found", "permission")
		}
		cachedPolicyLocation[id] = *location
	}
	return cachedPolicyLocation[id]
}
