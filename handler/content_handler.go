//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package handler

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/fieldtype"
	"dm/util"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Contenter interface {
	Publish()

	Create()

	Edit()

	Delete()
}

type ContentHandler struct {
}

func (content ContentHandler) CreateLocation(parentID int) {
	location := entity.Location{ParentID: parentID}
	location.Store()
}

//Validate and Return a validation result
func (handler *ContentHandler) Validate(contentType string, inputs map[string]interface{}) (bool, ValidationResult) {
	definition := contenttype.GetContentDefinition(contentType)
	//todo: check there is no extra field in the inputs
	//todo: check max length
	//check required
	fieldsDef := definition.Fields
	result := ValidationResult{}
	for identifier, fieldDef := range fieldsDef {
		fieldHandler := fieldtype.GetHandler(fieldsDef[identifier].FieldType)
		_, fieldExists := inputs[identifier]
		if fieldDef.Required &&
			(!fieldExists || (fieldExists && fieldHandler != nil && fieldHandler.IsEmpty(inputs[identifier]))) {
			fieldResult := FieldValidationResult{Identifier: identifier, Detail: "1"}
			result.Fields = append(result.Fields, fieldResult)
		}
	}
	if len(result.Fields) > 0 {
		return false, result
	}
	//Validate field
	for identifier, input := range inputs {
		fieldHanlder := fieldtype.GetHandler(fieldsDef[identifier].FieldType)
		if valid, detail := fieldHanlder.Validate(input); !valid {
			fieldResult := FieldValidationResult{Identifier: identifier, Detail: detail}
			result.Fields = append(result.Fields, fieldResult)
		}
	}

	//todo: add more custom validation based on type
	return true, ValidationResult{}
}

func (content ContentHandler) Draft(contentType string, parentID int) error {
	//create empty
	now := int(time.Now().Unix())
	article := entity.Article{ContentCommon: entity.ContentCommon{Published: now, Modified: now}}
	err := article.Store()
	if err != nil {
		return errors.Wrap(err, "[Handler.Draft]Error when creating article with parent id:"+
			strconv.Itoa(parentID)+". No location crreated.")
	}
	//Save location
	location := entity.Location{ParentID: -parentID,
		ContentType: contentType,
		ContentID:   article.CID,
		UID:         util.GenerateUID()}
	err = location.Store()
	if err != nil {
		return errors.Wrap(err, "[Handler.Draft]Error when creating location with content type -"+
			contentType+", content id -"+strconv.Itoa(article.CID)+", parent id - "+strconv.Itoa(parentID))
	}
	return nil
}

//Publish a draft
func (content ContentHandler) Publish() {

}

//Create a content(same behavior as Draft&Publish but store published version directly)
func (handler *ContentHandler) Create(parentID int, contentType string, inputs map[string]interface{}) (bool, ValidationResult, error) {
	//Validate
	valid, validationResult := handler.Validate(contentType, inputs)

	if !valid {
		return false, validationResult, nil
	}

	content := entity.NewInstance(contentType)

	contentDefinition := contenttype.GetContentDefinition(contentType)
	fieldsDefinition := contentDefinition.Fields

	//todo: check if all the inputs are needed.

	//todo: check all kind of validation

	for identifier, input := range inputs {
		fieldType := fieldsDefinition[identifier].FieldType
		fieldtypeHandler := fieldtype.GetHandler(fieldType)
		fieldValue := fieldtypeHandler.ToStorage(input)
		err := content.SetValue(identifier, fieldValue)
		if err != nil {
			return false, ValidationResult{}, errors.Wrap(err, "[Create]Can not set input to "+identifier)
		}
	}

	now := int(time.Now().Unix())
	content.SetValue("published", now)
	content.SetValue("modified", now)
	content.SetValue("remote_id", util.GenerateUID())

	err := content.Store()
	if err != nil {
		return false, ValidationResult{}, err
	}
	//todo: add commit and rollback for the whole saving

	//Save location
	location := entity.Location{ParentID: parentID,
		ContentID:   content.Value("cid").(int),
		ContentType: contentType,
		UID:         util.GenerateUID()}
	//todo: set name based on rules. Now it's all based on title.
	contentName := content.Value("title").(fieldtype.TextField).Data
	location.Name = contentName
	err = location.Store()
	if err != nil {
		return false, ValidationResult{}, err
	}

	//todo: update other things in location like main_id, hierarchy

	return true, ValidationResult{}, nil
}

//Update content.
//The inputs doesn't need to include all required fields. However if it's there,
// it will check if it's required&empty
func (content ContentHandler) Update(id int, inputs map[string]interface{}) {

}

//Delete content
func (content ContentHandler) Delete(id int, toTrash bool) {

}

//Format of fields: eg. title:"test", modified: 12121
func (handler *ContentHandler) store(parentID int, contentType string, fields map[string]interface{}) {
	handler.Validate(contentType, fields)
}
