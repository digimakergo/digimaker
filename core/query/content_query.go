package query

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/pkg/errors"
)

//TreeNode is a query result when querying SubTree
type TreeNode struct {
	*contenttype.Location
	Fields   interface{}              `json:"fields"` //todo: maybe more generaic attributes instead of hard coded 'Fields' here, or use custom MarshalJSON(remove *Locatoin then)?
	Content  contenttype.ContentTyper `json:"-"`
	Children []TreeNode               `json:"children"`
}

//Iterate loops all nodes in the tree
func (tn *TreeNode) Iterate(operation func(node *TreeNode)) {
	operation(tn)
	for i, child := range tn.Children {
		child.Iterate(operation)
		tn.Children[i] = child
	}
}

// FetchByID fetches content by location id.
//If no location found. it will return nil and error message.
func FetchByID(ctx context.Context, locationID int) (contenttype.ContentTyper, error) {
	//get type first by location.
	location := contenttype.Location{}
	_, err := db.BindEntity(ctx, &location, "dm_location", db.Cond("id", locationID))
	if err != nil {
		return nil, errors.Wrap(err, "[contentquery.fetchbyid]Can not fetch location by locationID "+strconv.Itoa(locationID))
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.fetchbyid]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.ContentType
	result, err := FetchByCID(ctx, contentType, contentID)
	return result, err
}

//FetchByUID fetches content by unique id
func FetchByUID(ctx context.Context, uid string) (contenttype.ContentTyper, error) {
	//get type first by location.
	location := contenttype.Location{}
	_, err := db.BindEntity(ctx, &location, "dm_location", db.Cond("uid", uid))
	if err != nil {
		return nil, errors.Wrap(err, "[contentquery.fetchbyuid]Can not fetch location by uid "+uid)
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.fetchbyid]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.ContentType
	result, err := FetchByCID(ctx, contentType, contentID)
	return result, err
}

// FetchByCID fetches a content by content id.
func FetchByCID(ctx context.Context, contentType string, contentID int) (contenttype.ContentTyper, error) {
	return Fetch(ctx, contentType, db.Cond("c.id", contentID))
}

// FetchByCUID fetches a content by content's uid(cuid)
func FetchByCUID(ctx context.Context, contentType string, cuid string) (contenttype.ContentTyper, error) {
	return Fetch(ctx, contentType, db.Cond("c.cuid", cuid))
}

// Fetch fetches first content based on condition.
func Fetch(ctx context.Context, contentType string, condition db.Condition) (contenttype.ContentTyper, error) {
	return DefaultQuerier.Fetch(ctx, contentType, condition)
}

func LocationList(condition db.Condition) ([]contenttype.Location, int, error) {
	locations := []contenttype.Location{}
	count, err := db.BindEntity(context.Background(), &locations, "dm_location", condition)
	if err != nil {
		return locations, 0, err
	}
	return locations, count, nil
}

// Children fetches children content directly under the given parentContent
func Children(ctx context.Context, userID int, parentContent contenttype.ContentTyper, contenttype string, cond db.Condition) ([]contenttype.ContentTyper, int, error) {
	contentTypeList := parentContent.Definition().AllowedTypes
	if !util.Contains(contentTypeList, contenttype) {
		return nil, -1, errors.New("content type " + contenttype + "doesn't exist or not allowed.")
	}
	result, countResult, err := SubList(ctx, userID, parentContent, contenttype, 1, cond)
	return result, countResult, err
}

// SubTree fetches content and return a tree result under rootContent, permission considered.
// See TreeNode for the tree structure
func SubTree(ctx context.Context, userID int, rootContent contenttype.ContentTyper, depth int, contentTypes string, sortby []string) (TreeNode, error) {
	contentTypeList := strings.Split(contentTypes, ",")
	var list []contenttype.ContentTyper
	for _, contentType := range contentTypeList {
		currentList, _, err := SubList(ctx, userID, rootContent, contentType, depth, db.EmptyCond().Sortby(sortby...))
		if err != nil {
			return TreeNode{}, err
		}
		for _, item := range currentList {
			list = append(list, item)
		}
	}

	treenode := TreeNode{Content: rootContent, Location: rootContent.GetLocation()}
	buildTree(&treenode, list)
	return treenode, nil
}

//build tree from list. Internal use.
//If there are items not in the tree(parent id is NOT equal to anyone in the list), they will not be attached to the tree.
func buildTree(treenode *TreeNode, list []contenttype.ContentTyper) {
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
		buildTree(&treenode.Children[i], list)
	}
}

// add condition from permission.
// so if matched with limit, add that limit to condition
// if matches with a empty limit(if there is), return empty(meaning no limit)
// if doesn't match, return a False condition(no result in query)
func accessCondition(userID int, contenttype string, context context.Context) db.Condition {
	accessType, limits, err := permission.GetUserAccess(context, userID, "content/read")
	if err != nil {
		//todo: debug messsage it
	}

	if accessType == permission.AccessNo {
		return db.FalseCond()
	}

	if accessType == permission.AccessFull {
		return db.EmptyCond()
	}

	//add conditions based on limits
	result := db.EmptyCond()
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

		currentCond := db.EmptyCond()

		if section, ok := limit["section"]; ok {
			currentCond = currentCond.Cond("l.section", util.InterfaceToStringArray(section.([]interface{})))
		}

		//comment below out to have a better/different way of subtree limit, in that case currentCondition will be and.
		// if sTree, ok := limit["subtree"]; ok {
		// 	item := sTree.(string) //todo: support array
		// 	itemInt, _ := strconv.Atoi(item)
		// 	subtree = append(subtree, itemInt)
		// }

		//todo: current self author will override the other policy. to be fixed.
		if author, ok := limit["author"]; ok {
			if author.(string) == "self" {
				currentCond = currentCond.Cond("author", userID)
			}
		}

		result = result.Or(currentCond)
	}

	return result
}

// SubList fetches content list with permission considered(only return contents the user has access to).
func SubList(ctx context.Context, userID int, rootContent contenttype.ContentTyper, contentType string, depth int, condition db.Condition) ([]contenttype.ContentTyper, int, error) {
	rootLocation := rootContent.GetLocation()
	def, _ := definition.GetDefinition(contentType)
	if def.HasLocation {
		if depth == 1 {
			//Direct children
			condition = condition.Cond("l.parent_id", rootLocation.ID)
		} else {
			rootHierarchy := rootLocation.Hierarchy
			rootDepth := rootLocation.Depth
			condition = condition.Cond("l.hierarchy like", rootHierarchy+"/%")
			if depth > 0 {
				condition = condition.Cond("l.depth <=", rootDepth+depth)
			}
		}
	} else {
		condition = condition.Cond("location_id", rootLocation.ID)
	}

	//fetch
	list, count, err := ListWithUser(ctx, userID, contentType, condition)
	return list, count, err
}

// ListWithUser fetches a list of content which the user has read permission to.
func ListWithUser(ctx context.Context, userID int, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error) {
	return DefaultQuerier.ListWithUser(ctx, userID, contentType, condition)
}

func List(ctx context.Context, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error) {
	return DefaultQuerier.List(ctx, contentType, condition)
}

//return a version content
func Version(ctx context.Context, contentType string, condition db.Condition) (contenttype.Version, contenttype.ContentTyper, error) {
	_, err := definition.GetDefinition(contentType)
	if err != nil {
		return contenttype.Version{}, nil, err
	}

	version := contenttype.Version{}
	_, err = db.BindEntity(ctx, &version, "dm_version", condition.Cond("content_type", contentType))
	if err != nil {
		return contenttype.Version{}, nil, err
	}
	if version.ID == 0 {
		return contenttype.Version{}, nil, nil
	}

	data := []byte(version.Data)
	content := contenttype.NewInstance(contentType)
	author := version.Author

	err = json.Unmarshal(data, &content)
	if err != nil {
		return version, nil, errors.New("Not valid data for content")
	}

	content.SetValue("author", author)
	return version, content, nil
}

func GetUser(id int) (contenttype.ContentTyper, error) {
	user, err := FetchByCID(context.Background(), "user", id)
	return user, err
}

func GetContentAuthor(content contenttype.ContentTyper) (contenttype.ContentTyper, error) {
	authorID := content.Value("author").(int)
	user, err := GetUser(authorID)
	return user, err
}

//GetRelationOptions get content list based on relation parameters
func RelationOptions(ctx context.Context, ctype string, identifier string, limit []int, sortby []string, hasCount bool) ([]contenttype.ContentTyper, int, error) {
	contentDef, err := definition.GetDefinition(ctype)
	if err != nil {
		return nil, 0, errors.New("Can not get content defintion of " + ctype)
	}
	field, ok := contentDef.FieldMap[identifier]
	if !ok {
		return nil, 0, errors.New("Can not find filed " + identifier + " from content " + ctype)
	}
	if field.FieldType != "relation" {
		return nil, 0, errors.New("Not a relation type")
	}
	params, err := fieldtype.ConvertRelationParams(field.Parameters)
	if err != nil {
		return nil, 0, err
	}
	condition := db.EmptyCond()
	for cKey, cValue := range params.Condition {
		condition = condition.And(cKey, cValue)
	}
	return List(ctx, params.Type, condition)
}
