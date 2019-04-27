package handler

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"fmt"

	"strings"

	"github.com/pkg/errors"
)

type RelationHandler struct {
}

//Add a content to current content(toContent)
func (handler *RelationHandler) AddTo(to contenttype.ContentTyper, from contenttype.ContentTyper, identifier string, priority int, description string) error {

	fieldSetting, ok := contenttype.GetContentDefinition(to.ContentType()).Fields[identifier]
	if !ok {
		return errors.New("[handler.AddTo]Target content doesn't have field " + identifier)
	}

	fieldDef := fieldSetting.GetDefinition()
	if !fieldDef.IsRelation {
		return errors.New("[handler.AddTo]field" + identifier + "is not a relation type.")
	}

	dataFields := strings.Split(fieldDef.RelationSettings.DataFields, ",")
	dataList := []interface{}{}
	//todo: validate if the fromField exist. this maybe done in bootstrap/generating datatype config

	for _, fromField := range dataFields {
		//todo: convert value
		dataList = append(dataList, from.Value(fromField))
	}
	dataPattern := fieldDef.RelationSettings.DataPattern
	data := fmt.Sprintf(dataPattern, dataList...)

	//todo: validate if it's added already.
	relation := entity.Relation{
		ToContentID:  to.Value("content_id").(int),
		ToType:       to.ContentType(),
		FromLocation: from.Value("id").(int),
		Priority:     priority,
		Identifier:   identifier,
		Description:  description,
		Data:         data}
	err := relation.Store()
	if err != nil {
		errors.Wrap(err, "[handler.AddTo]Saving relation error.")
	}

	return nil
}

//Update all contents which is related to current content(fromContent)
func (handler *RelationHandler) UpdateValues(fromContent contenttype.ContentTyper) {

}
