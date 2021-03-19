package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	. "github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"

	"github.com/pkg/errors"
)

type RelationHandler struct {
}

func (handler *RelationHandler) AddContentID(to contenttype.ContentTyper, contentID int, contentType string, identifer string, priority int, description string) error {
	return nil
}

//Add a content to current content(toContent)
func (handler *RelationHandler) AddContent(ctx context.Context, to contenttype.ContentTyper, from contenttype.ContentTyper, identifier string, priority int, description string) error {
	//todo: validate if the fromField exist. this maybe done in bootstrap/generating datatype config
	data, err := handler.generateData(to, identifier, from)
	if err != nil {
		return errors.Wrap(err, "[relationhandler.AddTo]")
	}

	contentID := to.Value("content_id").(int)
	contentType := to.ContentType()
	fromLocationID := from.Value("id").(int)

	//Check if it's added already.
	dbHandler := db.DBHanlder()
	currentRelation := contenttype.Relation{}
	dbHandler.GetEntity(ctx, &currentRelation,
		"dm_relation", Cond("to_content_id", contentID).
			Cond("from_location", fromLocationID).
			Cond("identifier", identifier),
		false)

	if currentRelation.ID != 0 {
		return errors.New("[relationhandler.Add]Relation existing already to " +
			strconv.Itoa(contentID) + " on " + identifier +
			" from " + strconv.Itoa(fromLocationID))
	}

	relation := contenttype.Relation{
		ToContentID:  contentID,
		ToType:       contentType,
		FromLocation: fromLocationID,
		Priority:     priority,
		Identifier:   identifier,
		Description:  description,
		Data:         data}
	err = relation.Store(ctx)
	if err != nil {
		errors.Wrap(err, "[relationhandler.AddTo]Saving relation error.")
	}

	return nil
}

//Generate relation data based on name pattern.
func (handler *RelationHandler) generateData(to contenttype.ContentTyper, identifier string, from contenttype.ContentTyper) (string, error) {
	def, _ := definition.GetDefinition(to.ContentType())
	fieldSetting, ok := def.FieldMap[identifier]
	if !ok {
		return "", errors.New("Target content doesn't have field " + identifier)
	}

	fieldDef := fieldtype.GetDef(fieldSetting.Identifier)
	if !fieldDef.IsRelation {
		return "", errors.New("field" + identifier + "is not a relation type.")
	}

	dataFields := strings.Split(fieldDef.RelationSettings.DataFields, ",")
	dataList := []interface{}{}

	for _, fromField := range dataFields {
		//todo: convert value
		dataList = append(dataList, from.Value(fromField))
	}
	dataPattern := fieldDef.RelationSettings.DataPattern
	data := fmt.Sprintf(dataPattern, dataList...)
	return data, nil
}

//Update all contents which is related to current content(fromContent)
func (handler *RelationHandler) UpdateValues(fromContent contenttype.ContentTyper) {

}
