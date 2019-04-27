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
func (handler *RelationHandler) Add(to contenttype.ContentTyper, from contenttype.ContentTyper, identifier string, priority int, description string) error {
	//todo: validate if the fromField exist. this maybe done in bootstrap/generating datatype config
	data, err := handler.generateData(to, identifier, from)
	if err != nil {
		return errors.Wrap(err, "[hanlder.AddTo]")
	}

	//todo: validate if it's added already.
	relation := entity.Relation{
		ToContentID:  to.Value("content_id").(int),
		ToType:       to.ContentType(),
		FromLocation: from.Value("id").(int),
		Priority:     priority,
		Identifier:   identifier,
		Description:  description,
		Data:         data}
	err = relation.Store()
	if err != nil {
		errors.Wrap(err, "[handler.AddTo]Saving relation error.")
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
