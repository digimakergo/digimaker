package handler

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	. "dm/query"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type RelationHandler struct {
}

//Add a content to current content(toContent)
func (handler *RelationHandler) Add(to contenttype.ContentTyper, from contenttype.ContentTyper, identifier string, priority int, description string) error {
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
	currentRelation := entity.Relation{}
	dbHandler.GetEnity("dm_relation", Cond("to_content_id", contentID).
		Cond("from_location", fromLocationID).
		Cond("identifier", identifier),
		&currentRelation)

	if currentRelation.ID != 0 {
		return errors.New("[relationhandler.Add]Relation existing already to " +
			strconv.Itoa(contentID) + " on " + identifier +
			" from " + strconv.Itoa(fromLocationID))
	}

	relation := entity.Relation{
		ToContentID:  contentID,
		ToType:       contentType,
		FromLocation: fromLocationID,
		Priority:     priority,
		Identifier:   identifier,
		Description:  description,
		Data:         data}
	err = relation.Store()
	if err != nil {
		errors.Wrap(err, "[relationhandler.AddTo]Saving relation error.")
	}

	return nil
}

//Generate relation data based on name pattern.
func (handler *RelationHandler) generateData(to contenttype.ContentTyper, identifier string, from contenttype.ContentTyper) (string, error) {
	fieldSetting, ok := contenttype.GetContentDefinition(to.ContentType()).Fields[identifier]
	if !ok {
		return "", errors.New("Target content doesn't have field " + identifier)
	}

	fieldDef := fieldSetting.GetDefinition()
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
