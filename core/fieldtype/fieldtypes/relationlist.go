package fieldtypes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"errors"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/util"
)

type RelationListHandler struct {
	definition.FieldDef
}

//LoadFromInput load data from input before validation
func (handler RelationListHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	s := fmt.Sprint(input)
	result := contenttype.RelationList{}
	if s != "" {
		arr := strings.Split(s, ";")
		if len(arr) != 2 {
			return result, errors.New("wrong format")
		}
		arrInt, err := util.ArrayStrToInt(strings.Split(arr[0], ","))

		if err != nil {
			return result, errors.New("Not int(s)")
		}
		arrType := strings.Split(arr[1], ",")
		if len(arrType) != len(arrInt) {
			return result, errors.New("id and type are not same length")
		}

		for i, v := range arrInt {
			fromType := arrType[i]
			fromCid := v
			fromDef, err := definition.GetDefinition(fromType)
			if err != nil {
				return result, err
			}

			r := contenttype.Relation{}

			r.FromContentID = fromCid
			r.FromType = fromType
			r.Priority = len(arrInt) - i
			r.Identifier = handler.Identifier

			relationDataFields := fromDef.RelationData
			if len(relationDataFields) > 0 {
				//get content
				contents := db.DatamapList{}
				db.BindEntity(context.Background(), &contents, fromDef.TableName, db.Cond("id", fromCid))
				if len(contents) == 0 {
					return result, errors.New("No content found on " + strconv.Itoa(fromCid))
				}

				//If there is one field, use it on data, otherwise use json map as data
				if len(relationDataFields) == 1 {
					r.Data = fmt.Sprint(contents[0][relationDataFields[0]])
				} else {
					datamap := map[string]interface{}{}
					for _, field := range relationDataFields {
						datamap[field] = contents[0][field]
					}
					data, _ := json.Marshal(datamap)
					r.Data = string(data)
				}
			}

			result = append(result, r)
		}
		//todo: update data after updating from_content
	}
	return result, nil
}

func (handler RelationListHandler) DBField() string {
	return ""
}

//todo: pass the current content and value
func (handler RelationListHandler) Store(ctx context.Context, value interface{}, contentType string, cid int, transaction *sql.Tx) error {
	relations, ok := value.(contenttype.RelationList)
	if !ok {
		return errors.New("Not a relationlist type")
	}
	existingList := []contenttype.Relation{}

	currentCondition := db.Cond("to_type", contentType).Cond("to_content_id", cid).Cond("identifier", handler.Identifier)
	// //get existing
	if cid > 0 {
		_, err := db.BindEntity(ctx, &existingList, "dm_relation", currentCondition)
		if err != nil {
			return err
		}
	}

	//get to be deleted
	deleteCond := db.EmptyCond()
	for _, existing := range existingList {
		willDelete := true
		for _, relation := range relations {
			if existing.FromContentID == relation.FromContentID && existing.FromType == relation.FromType {
				willDelete = false
			}
		}
		if willDelete {
			deleteCond = deleteCond.Or(db.Cond("from_content_id", existing.FromContentID).Cond("from_type", existing.FromType))
		}
	}

	//get to be added
	toBeAdded := []contenttype.Relation{}
	toBeUpdated := []contenttype.Relation{}
	for _, relation := range relations {
		exists := false
		for _, existing := range existingList {
			if existing.FromContentID == relation.FromContentID &&
				existing.FromType == relation.FromType &&
				existing.Identifier == relation.Identifier {
				exists = true
				if existing.Priority != relation.Priority {
					toBeUpdated = append(toBeUpdated, relation)
				}
			}
		}
		if !exists {
			toBeAdded = append(toBeAdded, relation)
		}
	}

	//execute update
	if len(toBeUpdated) > 0 {
		for _, relation := range toBeUpdated {
			updateMap := map[string]interface{}{}
			updateMap["priority"] = relation.Priority
			err := db.Update(ctx, "dm_relation", updateMap,
				currentCondition.Cond("from_content_id", relation.FromContentID).Cond("from_type", relation.FromType),
				transaction)
			if err != nil {
				return fmt.Errorf("Update relationlist error: %w", err)
			}
		}
	}

	//execute delete
	if !db.IsEmptyCond(deleteCond) {
		err := db.Delete(ctx, "dm_relation", currentCondition.And(deleteCond), transaction)
		if err != nil {
			return err
		}
	}

	//execute insert
	for _, relation := range toBeAdded {
		dataMap := map[string]interface{}{}
		dataMap["identifier"] = relation.Identifier
		dataMap["to_content_id"] = cid
		dataMap["to_type"] = contentType
		dataMap["from_content_id"] = relation.FromContentID
		dataMap["from_type"] = relation.FromType
		dataMap["data"] = relation.Data
		dataMap["priority"] = relation.Priority
		_, err := db.Insert(ctx, "dm_relation", dataMap, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	fieldtype.Register(
		fieldtype.Definition{
			Name:     "relationlist",
			DataType: "contenttype.RelationList",
			NewHandler: func(fieldDef definition.FieldDef) fieldtype.Handler {
				return RelationListHandler{FieldDef: fieldDef}
			}})
}
