package query

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/permission"
)

//TreeNode is a query result when querying SubTree
type TreeNode struct {
	*contenttype.Location
	Name     string                   `json:"name"`
	Fields   interface{}              `json:"fields"` //todo: maybe more generaic attributes instead of hard coded 'Fields' here, or use custom MarshalJSON(remove *Locatoin then)?
	Content  contenttype.ContentTyper `json:"content"`
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

// FetchByLID fetches content by location id.
//If no location found. it will return nil and error message.
func FetchByLID(ctx context.Context, locationID int) (contenttype.ContentTyper, error) {
	//get type first by location.
	location, err := FetchLocationByID(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("[contentquery.FetchByLID]Can not fetch location by locationID %v: %w", locationID, err)
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.FetchByLID]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.Contenttype
	result, err := FetchByCID(ctx, contentType, contentID)

	return result, err
}

func FetchLocationByID(ctx context.Context, locationID int) (contenttype.Location, error) {
	return FetchLocation(ctx, db.Cond("id", locationID))
}

func FetchLocation(ctx context.Context, cond db.Condition) (contenttype.Location, error) {
	locations, _, err := LocationList(cond.Limit(0, 1))
	if err != nil {
		return contenttype.Location{}, err
	}
	if len(locations) > 0 {
		return locations[0], nil
	} else {
		return contenttype.Location{}, nil
	}
}

func FetchByPath(ctx context.Context, path string) (contenttype.ContentTyper, error) {
	//get type first by location.
	location, err := FetchLocation(ctx, db.Cond("identifier_path", path))
	if err != nil {
		return nil, fmt.Errorf("[contentquery.fetchbypath]Can not fetch location by location path %v: %w", path, err)
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.fetchbypath]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.Contenttype
	result, err := FetchByCID(ctx, contentType, contentID)
	return result, err
}

//FetchByUID fetches content by unique id
func FetchByUID(ctx context.Context, uid string) (contenttype.ContentTyper, error) {
	//get type first by location.
	location := contenttype.Location{}
	_, err := db.BindEntity(ctx, &location, "dm_location", db.Cond("uid", uid))
	if err != nil {
		return nil, fmt.Errorf("[contentquery.fetchbyuid]Can not fetch location by uid %v: %w", uid, err)
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.FetchByLID]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.Contenttype
	result, err := FetchByCID(ctx, contentType, contentID)
	return result, err
}

// FetchByCID is duplicate of FetchByCID
func FetchByID(ctx context.Context, contentType string, contentID int) (contenttype.ContentTyper, error) {
	return FetchByID(ctx, contentType, contentID)
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

	treenode := TreeNode{Content: rootContent, Name: rootContent.GetName(), Location: rootContent.GetLocation()}
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
			treenode.Children = append(treenode.Children, TreeNode{Content: item, Name: item.GetName(), Location: location})
		}
	}

	//Add sub level. If it's leaf node it will not run the loop.
	for i, _ := range treenode.Children {
		buildTree(&treenode.Children[i], list)
	}
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
		if def.HasDataField("location_id") {
			condition = condition.Cond("location_id", rootLocation.ID)
		}
	}

	//fetch
	permissionCond := permission.GetListCondition(ctx, userID, contentType, rootContent)
	condition = condition.And(permissionCond)
	list, count, err := List(ctx, contentType, condition)
	return list, count, err
}

// ListWithUser fetches a list of content which the user has read permission to.
// Note: If you have parent condition, use SubList, because this will not optimize 'under' policy and parent paramter
func ListWithUser(ctx context.Context, userID int, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error) {
	//permission condition
	permissionCondition := permission.GetListCondition(ctx, userID, contentType, nil)
	condition = condition.And(permissionCondition)

	//fetch
	list, count, err := List(ctx, contentType, condition)
	return list, count, err
}

//List without considering permission
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

func OutputField(ctx context.Context, content contenttype.ContentTyper, field string) (interface{}, error) {
	fieldMap := content.Definition().FieldMap
	if fieldDef, ok := fieldMap[field]; ok {
		//todo: handle error here
		result := outputField(ctx, fieldDef, content.Value(field))
		return result, nil
	} else {
		return nil, errors.New("Field not found")
	}
}

func outputField(ctx context.Context, fieldDef definition.FieldDef, value interface{}) interface{} {
	handler := fieldtype.GethHandler(fieldDef)
	if handler != nil {
		if washer, ok := handler.(fieldtype.Outputer); ok {
			washedValue := washer.Output(ctx, DefaultQuerier, value)
			return washedValue
		} else {
			return value
		}
	} else {
		return value
	}
}

//Output converts content into output format.(eg. add text to select, name in relation, etc)
func Output(ctx context.Context, content contenttype.ContentTyper) (contenttype.ContentMap, error) {
	if content != nil {
		def, _ := definition.GetDefinition(content.ContentType())
		contentMap, err := contenttype.ContentToMap(content)
		if err != nil {
			return nil, err
		}
		for identifier, fieldDef := range def.FieldMap {
			value := content.Value(identifier)
			washedValue := outputField(ctx, fieldDef, value)
			contentMap[identifier] = washedValue
		}
		return contentMap, nil
	}
	return nil, nil
}

//OutputList converts contents into output format, see Output for single content.
func OutputList(ctx context.Context, contentList []contenttype.ContentTyper) ([]contenttype.ContentMap, error) {
	result := []contenttype.ContentMap{}
	for _, content := range contentList {
		contentMap, err := Output(ctx, content)
		if err != nil {
			return nil, err
		}
		result = append(result, contentMap)
	}
	return result, nil
}
