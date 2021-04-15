package query

import (
	"context"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/pkg/errors"
)

type Querier struct {
}

func (q Querier) Fetch(ctx context.Context, contentType string, condition db.Condition) (contenttype.ContentTyper, error) {
	content := contenttype.NewInstance(contentType)
	condition = condition.Limit(0, 1)
	_, err := db.BindContent(ctx, content, contentType, condition)
	if err != nil {
		return nil, err
	}

	//finish bind
	err = q.finishBind(ctx, content)
	if err != nil {
		return nil, err
	}

	if content.GetCID() == 0 {
		return nil, nil
	}
	return content, err
}

func (q Querier) finishBind(ctx context.Context, content contenttype.ContentTyper) error {
	//Bind relation
	err := BindRelation(content)
	if err != nil {
		return err
	}

	//wash output
	q.washOutput(ctx, content)
	return nil
}

func (q Querier) washOutput(ctx context.Context, content contenttype.ContentTyper) {
	if content != nil {
		def, _ := definition.GetDefinition(content.ContentType())
		for identifier, fieldDef := range def.FieldMap {
			handler := fieldtype.GethHandler(fieldDef)
			if handler != nil {
				if washer, ok := handler.(fieldtype.Outputer); ok {
					value := content.Value(identifier)
					washedValue := washer.Ouput(ctx, q, value)
					content.SetValue(identifier, washedValue)
				}
			}
		}
	}
}

// ListWithUser fetches a list of content which the user has read permission to.
func (Querier) ListWithUser(ctx context.Context, userID int, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error) {
	//permission condition
	permissionCondition := accessCondition(userID, contentType, ctx)
	condition = condition.And(permissionCondition)

	//fetch
	list, count, err := List(ctx, contentType, condition)
	return list, count, err
}

// List fetches a list of content based on conditions. This is a database level 'list' WITHOUT permission check.
//
// For permission included, use query.ListWithUser
// For list under a tree root, use query.SubList
//
// Condition example:
//  db.Cond("l.parent_id", 4).Cond("author", 1).Cond("modified >", "2020-03-13")
//
// where content field can be used directly or with c. as prefix(eg. "c.author"), but location field need a l. prefix(eg. l.parent_id)
func (q Querier) List(ctx context.Context, contentType string, condition db.Condition) ([]contenttype.ContentTyper, int, error) {
	contentList := contenttype.NewList(contentType)
	count, err := db.BindContent(ctx, contentList, contentType, condition)
	if err != nil {
		return nil, count, err
	}
	//todo: count here. need to use ContentList type in NewList
	// if len(condition.Option.LimitArr) == 0 {
	// 	count = len(contentList)
	// }

	result := contenttype.ToList(contentType, contentList)
	//finish bind
	for _, content := range result {
		err = q.finishBind(ctx, content)
		if err != nil {
			return nil, 0, err
		}
	}

	return result, count, err
}

var DefaultQuerier = Querier{}

//FinishBind sets related data after data binding. It will be better if SQLBoiler support interface for customized  binding for struct.
func BindRelation(content contenttype.ContentTyper) error {
	contentType := content.ContentType()
	def, _ := definition.GetDefinition(contentType)
	if def.HasRelationlist() {
		relationMap := content.(contenttype.GetRelations).GetRelations()
		for identifier, fieldDef := range def.FieldMap {
			if fieldDef.FieldType == "relationlist" {
				if value, ok := relationMap[identifier]; ok {
					err := content.SetValue(identifier, value)
					if err != nil {
						return errors.Wrap(err, "Error when binding relationlist "+identifier)
					}
				}
			}
		}
	}
	return nil
}
