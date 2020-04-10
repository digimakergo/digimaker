package handler

import (
	"context"
	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"
	"github.com/xc/digimaker/core/log"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ContentQuery struct{}

//TreeNode is a query result when querying SubTree
type TreeNode struct {
	*contenttype.Location
	Content  contenttype.ContentTyper `json:"-"`
	Children []TreeNode               `json:"children"`
}

//Fetch content by location id.
//If no location found. it will return nil and error message.
func (cq ContentQuery) FetchByID(locationID int) (contenttype.ContentTyper, error) {
	//get type first by location.
	dbhandler := db.DBHanlder()
	location := contenttype.Location{}
	err := dbhandler.GetEntity("dm_location", db.Cond("id", locationID), []string{}, &location)
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
	err := dbhandler.GetEntity("dm_location", db.Cond("uid", uid), []string{}, &location)
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
	return cq.Fetch(contentType, db.Cond("content.id", contentID))
}

//Fetch a content by content's uid(cuid)
func (cq ContentQuery) FetchByCUID(contentType string, cuid string) (contenttype.ContentTyper, error) {
	return cq.Fetch(contentType, db.Cond("content.cuid", cuid))
}

//Fetch one content
func (cq ContentQuery) Fetch(contentType string, condition db.Condition) (contenttype.ContentTyper, error) {
	//todo: use limit in this case so it doesn't fetch more into memory.
	content := contenttype.NewInstance(contentType)
	count := -1
	err := cq.Fill(contentType, condition, []int{}, []string{}, content, &count)
	if err != nil {
		return nil, err
	}
	if content.GetCID() == 0 {
		return nil, nil
	}
	return content, err
}

//Fetch a list of content based on conditions. This is a database level 'list'. Return eg. *[]Article
func (cq ContentQuery) List(contentType string, condition db.Condition, limit []int, sortby []string, withCount bool) ([]contenttype.ContentTyper, int, error) {
	contentList := contenttype.NewList(contentType)
	count := -1
	if withCount {
		count = 0
	}
	err := cq.Fill(contentType, condition, limit, sortby, contentList, &count)
	if err != nil {
		return nil, count, err
	}
	result := contenttype.ToList(contentType, contentList)
	return result, count, err
}

//Fetch children
func (cq ContentQuery) Children(parentContent contenttype.ContentTyper, contenttype string, userID int, cond db.Condition, limit []int, sortby []string, withCount bool, context context.Context) ([]contenttype.ContentTyper, int, error) {
	contentTypeList := parentContent.Definition().AllowedTypes
	if !util.Contains(contentTypeList, contenttype) {
		return nil, -1, errors.New("content type " + contenttype + "doesn't exist or not allowed.")
	}
	result, countResult, err := cq.SubList(parentContent, contenttype, 1, userID, cond, limit, sortby, withCount, context)
	return result, countResult, err
}

//Get sub tree under rootContent, permission considered.
func (cq ContentQuery) SubTree(rootContent contenttype.ContentTyper, depth int, contentTypes string, userID int, sortby []string, context context.Context) (TreeNode, error) {
	contentTypeList := strings.Split(contentTypes, ",")
	var list []contenttype.ContentTyper
	for _, contentType := range contentTypeList {
		currentList, _, err := cq.SubList(rootContent, contentType, depth, userID, db.Cond("1", "1"), []int{}, sortby, false, context)
		if err != nil {
			return TreeNode{}, err
		}
		for _, item := range currentList {
			list = append(list, item)
		}
	}

	treenode := TreeNode{Content: rootContent, Location: rootContent.GetLocation()}
	cq.buildTree(&treenode, list)
	return treenode, nil
}

//build tree from list. Internal use.
//If there are items not in the tree(parent id is NOT equal to anyone in the list), they will not be attached to the tree.
func (cq ContentQuery) buildTree(treenode *TreeNode, list []contenttype.ContentTyper) {
	//Add current level contents
	parentLocation := treenode.Content.GetLocation()
	for _, item := range list {
		location := item.GetLocation()
		if location.Depth == parentLocation.Depth+1 && location.ParentID == parentLocation.ID {
			treenode.Children = append(treenode.Children, TreeNode{Content: item, Location: location})
		}
	}

	//Add sub level. If it's leaf node it will not run the loop.
	for i, _ := range treenode.Children {
		cq.buildTree(&treenode.Children[i], list)
	}
}

func permCondition(userID int, contenttype string, context context.Context) db.Condition {
	limits, err := permission.GetUserLimits(userID, "content/read", context)
	if err != nil {
		//todo: debug messsage it
	}

	//add conditions based on limits
	var result db.Condition
	for _, limit := range limits {
		if ctype, ok := limit["contenttype"]; ok {
			ctypeList := ctype.([]interface{})
			ctypeMatched := false
			for _, value := range ctypeList {
				if value.(string) == contenttype {
					ctypeMatched = true
					break
				}
			}
			//if the limit doesn't include the type, get next limit.
			if !ctypeMatched {
				continue
			}
		}

		var sectionCond db.Condition
		if section, ok := limit["section"]; ok {
			sectionCond = db.Cond("location.section", util.InterfaceToStringArray(section.([]interface{})))
		} else {
			sectionCond = db.Cond("1", "1")
		}

		//comment below out to have a better/different way of subtree limit, in that case currentCondition will be and.
		// if sTree, ok := limit["subtree"]; ok {
		// 	item := sTree.(string) //todo: support array
		// 	itemInt, _ := strconv.Atoi(item)
		// 	subtree = append(subtree, itemInt)
		// }
		if result.Children == nil {
			result = sectionCond
		} else {
			result = result.Or(sectionCond)
		}

		//todo: current self author will override the other policy. to be fixed.
		if author, ok := limit["author"]; ok {
			if author.(string) == "self" {
				authorCond := db.Cond("location.author", userID)
				if result.Children == nil {
					result = authorCond
				} else {
					result = result.And(authorCond)
				}
			}
		}

	}
	return result
}

//Get subtree with permission considered.
func (cq ContentQuery) SubList(rootContent contenttype.ContentTyper, contentType string, depth int, userID int, condition db.Condition, limit []int, sortby []string, withCount bool, context context.Context) ([]contenttype.ContentTyper, int, error) {

	rootLocation := rootContent.GetLocation()
	if depth == 1 {
		//Direct children
		condition = condition.Cond("location.parent_id", rootLocation.ID)
	} else {
		rootHierarchy := rootLocation.Hierarchy
		rootDepth := rootLocation.Depth
		condition = condition.Cond("location.hierarchy like", rootHierarchy+"/%").Cond("location.depth <=", rootDepth+depth)
	}

	//permission condition
	permCondition := permCondition(userID, contentType, context)
	if permCondition.Children != nil {
		condition = condition.And(permCondition)
	}

	//fetch
	list, count, err := cq.List(contentType, condition, limit, sortby, withCount)
	return list, count, err
}

//Fill all data into content which is a pointer
func (cq ContentQuery) Fill(contentType string, condition db.Condition, limit []int, sortby []string, content interface{}, count *int) error {
	dbhandler := db.DBHanlder()
	def, _ := contenttype.GetDefinition(contentType)
	tableName := def.TableName
	hasCount := *count != -1
	countResult, err := dbhandler.GetByFields(contentType, tableName, condition, limit, sortby, content, hasCount)
	if err != nil {
		message := "[List]Content Query error"
		log.Error(message+err.Error(), "")
		return errors.Wrap(err, message)
	}
	*count = countResult
	return nil
}

//Get version where version number is 0
func (cq ContentQuery) Draft(author int) {

}

//return a version content
func (cq ContentQuery) Version(contentType string, condition db.Condition) (contenttype.Version, contenttype.ContentTyper, error) {

	def, err := contenttype.GetDefinition(contentType)
	if err != nil {
		return contenttype.Version{}, nil, err
	}

	version := contenttype.Version{}
	dbHandler := db.DBHanlder()
	err = dbHandler.GetEntity("dm_version", condition.Cond("content_type", contentType), []string{}, &version)
	if err != nil {
		return contenttype.Version{}, nil, err
	}
	if version.ID == 0 {
		return contenttype.Version{}, nil, nil
	}

	data := []byte(version.Data)
	content := contenttype.NewInstance(contentType)
	author := version.Author
	obj := map[string]interface{}{}
	json.Unmarshal(data, &obj)

	for name := range def.FieldMap {
		value := obj[name]
		if value != nil {
			//valueStr, _ := json.Marshal(value)
			fHandler := fieldtype.GetHandler(def.FieldMap[name].FieldType)
			fieldValue := fHandler.NewValue(value)
			content.SetValue(name, fieldValue)
		}
	}

	content.SetValue("author", author)
	return version, content, nil
}

func (cq ContentQuery) GetUser(id int) (contenttype.ContentTyper, error) {
	querier := Querier()
	user, err := querier.FetchByContentID("user", id)
	return user, err
}

func (cq ContentQuery) GetContentAuthor(content contenttype.ContentTyper) (contenttype.ContentTyper, error) {
	authorID := content.Value("author").(int)
	user, err := cq.GetUser(authorID)
	return user, err
}

//todo: use method instead of global variable
var querier ContentQuery = ContentQuery{}

func Querier() ContentQuery {
	return querier
}
