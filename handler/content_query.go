package handler

import (
	"context"
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/permission"
	"dm/query"
	"dm/util"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ContentQuery struct{}

//TreeNode is a query result when querying SubTree
type TreeNode struct {
	Current  contenttype.ContentTyper
	Children []contenttype.ContentTyper
}

//Fetch content by location id.
//If no location found. it will return nil and error message.
func (cq ContentQuery) FetchByID(locationID int) (contenttype.ContentTyper, error) {
	//get type first by location.
	dbhandler := db.DBHanlder()
	location := contenttype.Location{}
	err := dbhandler.GetEntity("dm_location", query.Cond("id", locationID), &location)
	if err != nil {
		return nil, errors.Wrap(err, "[contentquery.fetchbyid]Can not fetch location by locationID "+strconv.Itoa(locationID))
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.fetchbyid]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.ContentType
	result, err := cq.FetchByContentID(contentType, contentID)
	return result, err
}

func (cq ContentQuery) FetchByUID(uid string) (contenttype.ContentTyper, error) {
	//get type first by location.
	dbhandler := db.DBHanlder()
	location := contenttype.Location{}
	err := dbhandler.GetEntity("dm_location", query.Cond("uid", uid), &location)
	if err != nil {
		return nil, errors.Wrap(err, "[contentquery.fetchbyuid]Can not fetch location by uid "+uid)
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.fetchbyid]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.ContentType
	result, err := cq.FetchByContentID(contentType, contentID)
	return result, err
}

//Fetch a content by content id.
func (cq ContentQuery) FetchByContentID(contentType string, contentID int) (contenttype.ContentTyper, error) {
	return cq.Fetch(contentType, query.Cond("content.id", contentID))
}

//Fetch a content by content's uid(cuid)
func (cq ContentQuery) FetchByCUID(contentType string, cuid string) (contenttype.ContentTyper, error) {
	return cq.Fetch(contentType, query.Cond("content.cuid", cuid))
}

//Fetch one content
func (cq ContentQuery) Fetch(contentType string, condition query.Condition) (contenttype.ContentTyper, error) {
	//todo: use limit in this case so it doesn't fetch more into memory.
	content := entity.NewInstance(contentType)
	err := cq.Fill(contentType, condition, content)
	if err != nil {
		return nil, err
	}
	if content.GetCID() == 0 {
		return nil, nil
	}
	return content, err
}

//Fetch a list of content, return eg. *[]Article
func (cq ContentQuery) List(contentType string, condition query.Condition) (interface{}, error) {
	contentList := entity.NewList(contentType)
	err := cq.Fill(contentType, condition, contentList)
	if err != nil {
		return nil, err
	}
	return contentList, err
}

//Get sub tree under rootContent, permission considered.
func (cq ContentQuery) SubTree(rootContent contenttype.ContentTyper, depth int, contentTypes string, userID int, context context.Context) (TreeNode, error) {
	contentTypeList := strings.Split(contentTypes, ",")
	for _, contentType := range contentTypeList {
		list, err := cq.SubList(rootContent, contentType, depth, userID, context)
		if err != nil {
			return TreeNode{}, err
		}
	}
	//todo: loop all the item and compose a tree
	return TreeNode{}, nil
}

//Get subtree with permission considered.
func (cq ContentQuery) SubList(rootContent contenttype.ContentTyper, contentType string, depth int, userID int, context context.Context) (interface{}, error) {
	limits, err := permission.GetUserLimits(userID, "content", "read", context)
	if err != nil {
		return nil, errors.Wrap(err, "Can not fetch permission.")
	}

	rootLocation := rootContent.GetLocation()
	rootHierarchy := rootLocation.Hierarchy
	rootDepth := rootLocation.Depth
	condition := query.Cond("hierarchy LIKE", rootHierarchy+"/%").Cond("depth <=", rootDepth+depth)

	//add condition based on limits
	var permissionCondition = query.Cond("1", "1")
	for _, limit := range limits {
		var currentCondition query.Condition
		if ctype, ok := limit["contenttype"]; ok {
			ctypeList := ctype.([]string)
			//if the limit doesn't include the type, get next limit.
			if !util.Contains(ctypeList, contentType) {
				continue
			}
		}

		if section, ok := limit["section"]; ok {
			currentCondition = query.Cond("section", section.(string))
		}

		//comment below out to have a better/different way of subtree limit, in that case currentCondition will be and.
		// if sTree, ok := limit["subtree"]; ok {
		// 	item := sTree.(string) //todo: support array
		// 	itemInt, _ := strconv.Atoi(item)
		// 	subtree = append(subtree, itemInt)
		// }
		if currentCondition.Children != nil {
			permissionCondition = permissionCondition.Or(currentCondition)
		}
	}
	condition = condition.And(permissionCondition)

	//fetch
	list, err := cq.List(contentType, condition)
	return list, err
}

//Fill all data into content which is a pointer
func (cq ContentQuery) Fill(contentType string, condition query.Condition, content interface{}) error {
	dbhandler := db.DBHanlder()
	tableName := contenttype.GetContentDefinition(contentType).TableName
	err := dbhandler.GetByFields(contentType, tableName, condition, content)
	if err != nil {
		message := "[List]Content Query error"
		util.Error(message, err.Error())
		return errors.Wrap(err, message)
	}
	return nil
}

//todo: use method instead of global variable
var querier ContentQuery = ContentQuery{}

func Querier() ContentQuery {
	return querier
}
