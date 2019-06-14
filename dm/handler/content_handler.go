//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package handler

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"context"
	"database/sql"
	"dm/dm/contenttype"
	"dm/dm/db"
	"dm/dm/fieldtype"
	"dm/dm/util"
	"dm/dm/util/debug"
	"strconv"
	"strings"
	"time"

	. "dm/dm/db"

	"github.com/pkg/errors"
)

type Contenter interface {
	Publish()

	Create()

	Edit()

	Delete()
}

type ContentHandler struct {
	Context context.Context
}

var ErrorNoPermission = errors.New("The user doesn't have access to the action.")

//Validate and Return a validation result
func (ch *ContentHandler) Validate(contentType string, inputs map[string]interface{}) (bool, ValidationResult) {
	definition := contenttype.GetContentDefinition(contentType)
	//todo: check max length
	//todo: check all kind of validation
	fieldsDef := definition.Fields

	result := ValidationResult{}

	//Check if there are more fields than defined
	for identifier, _ := range inputs {
		_, exist := fieldsDef[identifier]
		if !exist {
			result.Fields = append(result.Fields, FieldValidationResult{Identifier: identifier, Detail: "2"}) //not needed
		}
	}
	if !result.Passed() {
		return false, result
	}

	//check required
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

	contentTypeHanlder := GetContentTypeHandler(contentType)
	if contentTypeHanlder != nil {
		debug.Debug(ch.Context, "Validating from content type handler", "contenthandler.validate")
		contentTypeHanlder.Validate(inputs, &result)
	}
	return result.Passed(), result
}

func GenerateName(content contenttype.ContentTyper) string {
	if content.ContentType() == "user" {
		return content.Value("firstname").(fieldtype.TextField).Data + " " + content.Value("lastname").(fieldtype.TextField).Data
	}
	return content.Value("title").(fieldtype.TextField).Data //todo: make it patter based.
}

//Store content. Note it doesn't rollback - please rollback in invoking part if error happens.
//If it's no-location content, ingore the parentID.
func (ch *ContentHandler) storeCreatedContent(content contenttype.ContentTyper, tx *sql.Tx, parentID ...int) error {
	if content.GetCID() == 0 {
		debug.Debug(ch.Context, "Content is new.", "contenthandler.StoreCreatedContent")
	}
	err := content.Store(tx)
	if err != nil {
		return err
	}

	debug.Debug(ch.Context, "Content is saved. id :"+strconv.Itoa(content.GetCID())+". ", "contenthandler.StoreCreatedContent")

	//todo: deal with relations

	contentDefinition := content.Definition()
	if contentDefinition.HasLocation {
		if len(parentID) == 0 {
			return errors.New("Need parent location id.")
		}
		parentIDInt := parentID[0]
		parent, err := contenttype.GetLocationByID(parentIDInt)
		if err != nil {
			return errors.Wrap(err, "Can not get parent location with "+strconv.Itoa(parentIDInt))
		}

		//Save location
		location := contenttype.Location{}
		location.ParentID = parentIDInt
		location.ContentID = content.GetCID()
		location.ContentType = content.ContentType()
		location.UID = util.GenerateUID()
		contentName := GenerateName(content)
		location.IdentifierPath = parent.IdentifierPath + "/" + util.NameToIdentifier(contentName)
		location.Depth = parent.Depth + 1
		location.Section = parent.Section
		location.Name = contentName

		err = location.Store(tx)

		if err != nil {
			return errors.Wrap(err, "Transaction failed in location when saving new location.")
		}
		location.Hierarchy = parent.Hierarchy + "/" + strconv.Itoa(location.ID)
		location.MainID = location.ID
		err = location.Store(tx)
		if err != nil {
			return errors.Wrap(err, "Transaction failed in location when saving location for main_id and hierarchy.")
		}
		debug.Debug(ch.Context, "Location is saved. location id :"+strconv.Itoa(location.ID)+". ", "contenthandler.StoreCreatedContent")
	}
	return nil
}

//Create a content(same behavior as Draft&Publish but store published version directly)
func (ch *ContentHandler) Create(contentType string, inputs map[string]interface{}, parentID ...int) (bool, ValidationResult, error) {
	//todo: permission check.

	//Validate
	valid, validationResult := ch.Validate(contentType, inputs)
	if !valid {
		return false, validationResult, nil
	}

	//todo: add validation callback.
	contentDefinition := contenttype.GetContentDefinition(contentType)
	fieldsDefinition := contentDefinition.Fields

	//Create empty content instance and set value
	content := contenttype.NewInstance(contentType)
	for identifier, input := range inputs {
		fieldType := fieldsDefinition[identifier].FieldType
		fieldtypeHandler := fieldtype.GetHandler(fieldType)
		fieldValue := fieldtypeHandler.ToStorage(input)
		err := content.SetValue(identifier, fieldValue)
		if err != nil {
			return false, ValidationResult{}, errors.Wrap(err, "Can not set input to "+identifier)
		}
	}
	now := int(time.Now().Unix())
	content.SetValue("published", now)
	content.SetValue("modified", now)
	content.SetValue("cuid", util.GenerateUID())

	debug.StartTiming(ch.Context, "database", "contenthandler.create")
	debug.Debug(ch.Context, "Validation passed. Start saving content.", "contenthandler.Create")
	//Create transaction
	database, err := db.DB()
	if err != nil {
		return false, ValidationResult{}, errors.New("Can't get db connection.")
	}
	tx, err := database.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return false, ValidationResult{}, errors.New("Can't get transaction.")
	}

	//call content type handler
	contentTypeHandler := GetContentTypeHandler(contentType)
	if contentTypeHandler != nil {
		debug.Debug(ch.Context, "Calling handler for "+contentType, "contenthandler.create")
		err := contentTypeHandler.New(content, tx, parentID...)
		if err != nil {
			tx.Rollback()
			debug.Error(ch.Context, "Error from callback: "+err.Error(), "contenthandler.create")
			return false, ValidationResult{}, err
		}
	}

	//Save content and location if needed
	versionIfNeeded := 1
	if contentDefinition.HasVersion {
		content.SetValue("version", versionIfNeeded)
	}
	err = ch.storeCreatedContent(content, tx, parentID...)
	if err != nil {
		tx.Rollback()
		debug.Error(ch.Context, err.Error(), "contenthandler.Create")
		return false, ValidationResult{}, errors.Wrap(err, "Create content error")
	}

	//Save version if needed
	if contentDefinition.HasVersion {
		debug.Debug(ch.Context, "creating version", "contenthandler.Create")
		_, err = ch.CreateVersion(content, versionIfNeeded, tx)
		if err != nil {
			debug.Error(ch.Context, err.Error(), "contenthandler.Create")
			return false, ValidationResult{}, errors.Wrap(err, "Create version error.")
		}
		debug.Debug(ch.Context, "Created version: "+strconv.Itoa(versionIfNeeded), "contenthandler.Create")
	}

	//Invoke callback
	matchData := map[string]interface{}{"parent_id": parentID[0], //todo: maybe set parent id as mandatory parameter in Create
		"content_type": contentType}
	if contentDefinition.HasLocation {
		hierachy := content.GetLocation().Hierarchy
		matchData["under"] = strings.Split(hierachy, "/")
	}
	err = ch.InvokeCallback("create", true, matchData, content, tx)
	if err != nil {
		tx.Rollback()
		return false, ValidationResult{}, errors.Wrap(err, "Invoking callback error.")
	}

	//Commit all operations
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		debug.Error(ch.Context, "Commit error: "+err.Error(), "contenthandler.create")
		return false, ValidationResult{}, errors.Wrap(err, "Commit error.")
	}

	debug.Debug(ch.Context, "Data committed.", "contenthandler.create")

	ch.InvokeCallback("created", false, matchData, content) //continue if there is error

	debug.Debug(ch.Context, "Create finished.", "contenthandler.create")
	debug.EndTiming(ch.Context, "database", "contenthandler.create")
	return true, ValidationResult{}, nil
}

//Invoke callbacks based on condition match result
//see content_handler.json/yaml for conditions.
func (ch ContentHandler) InvokeCallback(event string, stopOnError bool, matchData map[string]interface{}, content contenttype.ContentTyper, params ...interface{}) error {
	operationHandlerList, matchInfo := GetOperationHandlerByCondition(event, matchData)
	count := len(operationHandlerList)
	if count > 0 {
		identifierList := []string{}
		for _, operationHandler := range operationHandlerList {
			identifierList = append(identifierList, operationHandler.Identifier)
		}
		debug.Debug(ch.Context, "Matched callbacks: "+strings.Join(identifierList, ","), "contenthandler.invoke_callback")
	} else {
		debug.Debug(ch.Context, "No callbacks matched.", "contenthandler.invoke_callback")
	}

	for _, info := range matchInfo {
		debug.Debug(ch.Context, info, "callback_match")
	}
	for i, operationHandler := range operationHandlerList {
		debug.Debug(ch.Context, strconv.Itoa(i+1)+"/"+strconv.Itoa(count)+
			" Invoking operation "+operationHandler.Identifier+" on "+event, "contehandler.invoke_callback")
		err := operationHandler.Execute(event, content, params...)
		if err != nil && stopOnError {
			debug.Error(ch.Context, "Error when invoking operation handler "+operationHandler.Identifier+":"+err.Error(), "contehandler.invoke_callback")
			return err
		}
		debug.Debug(ch.Context, "Invoking ended on "+operationHandler.Identifier, "contehandler.invoke_callback")
	}
	return nil
}

//Create a new version.
//It doesn't validate version number is increment
func (ch ContentHandler) CreateVersion(content contenttype.ContentTyper, versionNumber int, tx *sql.Tx) (int, error) {
	debug.Debug(ch.Context, "Creating version: "+strconv.Itoa(versionNumber), "contenthandler.CreateVersion")
	id := content.GetCID()
	data, err := contenttype.ContentToJson(content)
	if err != nil {
		return 0, errors.Wrap(err, "Can not create version data on content id: "+strconv.Itoa(id))
	}
	version := contenttype.Version{}
	version.ContentType = content.ContentType()
	version.ContentID = content.GetCID()
	version.Author = content.Value("author").(int)
	version.Version = versionNumber
	version.Data = data
	//version.Created = content.Value("created").(int)
	err = version.Store(tx)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "Can not save version on contetent id: "+strconv.Itoa(id))
	}
	debug.Debug(ch.Context, "Version created.", "contenthandler.CreateVersion")
	return version.ID, nil
}

func (ch ContentHandler) UpdateByContentID(contentType string, contentID int, inputs map[string]interface{}) (bool, ValidationResult, error) {
	content, err := Querier().FetchByContentID(contentType, contentID)
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Failed to get content via content id.")
	}
	if content.GetCID() == 0 {
		return false, ValidationResult{}, errors.Wrap(err, "Got empty content.")
	}

	return ch.Update(content, inputs)
}

func (ch ContentHandler) UpdateByID(id int, inputs map[string]interface{}) (bool, ValidationResult, error) {
	content, err := Querier().FetchByID(id)
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Failed to get content via id.")
	}
	if content.GetCID() == 0 {
		return false, ValidationResult{}, errors.Wrap(err, "Got empty content.")
	}

	return ch.Update(content, inputs)
}

//Update content.
//The inputs doesn't need to include all required fields. However if it's there,
// it will check if it's required&empty
func (ch ContentHandler) Update(content contenttype.ContentTyper, inputs map[string]interface{}) (bool, ValidationResult, error) {
	//Validate
	debug.Debug(ch.Context, "Validating", "contenthandler.update")
	contentType := content.ContentType()
	ch.Validate(contentType, inputs)
	//Save to new version
	contentDef := contenttype.GetContentDefinition(contentType)

	tx, err := db.CreateTx()
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Create transaction error.")
	}

	//todo: update relations

	//Set content.
	fieldsDefinition := contentDef.Fields
	for identifier, input := range inputs {
		fieldType := fieldsDefinition[identifier].FieldType
		fieldtypeHandler := fieldtype.GetHandler(fieldType)
		fieldValue := fieldtypeHandler.ToStorage(input)
		err := content.SetValue(identifier, fieldValue)
		if err != nil {
			return false, ValidationResult{}, errors.Wrap(err, "Can not set input to "+identifier)
		}
	}

	//Save new version and set to content
	if contentDef.HasVersion {
		version := content.Value("version").(int) + 1
		debug.Debug(ch.Context, "Creating new version: "+strconv.Itoa(version), "contenthandler.update")
		_, err := ch.CreateVersion(content, version, tx)
		if err != nil {
			//todo: rollback here not inside.
			debug.Debug(ch.Context, "Create new version failed: "+err.Error(), "contenthandler.update")
			return false, ValidationResult{}, errors.Wrap(err, "Can not save version.")
		}
		debug.Debug(ch.Context, "New version created", "contenthandler.update")
		content.SetValue("version", version)
	}
	now := int(time.Now().Unix())
	content.SetValue("modified", now)

	//Save update content.
	debug.Debug(ch.Context, "Saving content", "contenthandler.update")
	err = content.Store(tx)
	if err != nil {
		tx.Rollback()
		debug.Debug(ch.Context, "Saving content failed: "+err.Error(), "contenthandler.update")
		return false, ValidationResult{}, errors.Wrap(err, "Saving content error. ")
	}

	//Updated location related
	if contentDef.HasLocation {
		//todo: update all locations to this content.
		location := content.GetLocation()
		location.Name = GenerateName(content)
		//todo: location.IdentifierPath =
		debug.Debug(ch.Context, "Updating location info", "contenthandler.update")
		err := location.Store(tx)
		if err != nil {
			tx.Rollback()
			debug.Debug(ch.Context, "Updating location failed: "+err.Error(), "contenthandler.update")
			return false, ValidationResult{}, errors.Wrap(err, "Updating location info error.")
		}
		debug.Debug(ch.Context, "Location updated", "contenthandler.update")
	}

	//Invoke callback
	matchData := map[string]interface{}{"content_type": contentType}
	if content.Definition().HasLocation {
		hierachy := content.GetLocation().Hierarchy
		matchData["under"] = strings.Split(hierachy, "/")
	}
	//todo: maybe old content need to pass to callback.
	err = ch.InvokeCallback("update", true, matchData, content, tx)
	if err != nil {
		tx.Rollback()
		return false, ValidationResult{}, errors.Wrap(err, "Invoking callback error.")
	}

	debug.Debug(ch.Context, "All done. Commitinng.", "contenthandler.update")
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		debug.Error(ch.Context, "Commit error: "+err.Error(), "contenthandler.update")
		return false, ValidationResult{}, errors.Wrap(err, "Commit error.")
	}

	ch.InvokeCallback("updated", false, matchData, content)

	return true, ValidationResult{}, nil
}

//Delete content by location id
func (ch ContentHandler) DeleteByID(id int, toTrash bool) error {
	content, err := Querier().FetchByID(id)
	//todo: check how many. if more than 1, delete current only(and set main_id if needed)
	if err != nil {
		return errors.New("[handler.delete]Content doesn't exist with id: " + strconv.Itoa(id))
	}
	err = ch.DeleteByContent(content, toTrash)
	return err
}

//Delete content, relations and location.
//Note: this is only for when there is 1 location.
//  You need to judge if there are more than one locations before invoking this.
func (ch ContentHandler) DeleteByContent(content contenttype.ContentTyper, toTrash bool) error {
	//todo: check delete children. There should be more consideration if there are more children.
	//Delete location
	location := content.GetLocation()
	if location.CountLocations() > 1 {
		return errors.New("There are more than 1 location. Remove location first.")
	} else {
		database, err := db.DB()
		if err != nil {
			util.Error(err.Error())
			return errors.New("[handler.deleteByContent]Can not create connection.")
		}
		tx, err := database.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			tx.Rollback()
			message := "[handler.deleteByContent]Can not create transaction."
			util.Error(message + err.Error())
			return errors.New(message)
		}

		//Delete relation first.
		relations := content.GetRelations()
		if len(relations.Map) > 0 {
			dbHandler := db.DBHanlder()
			err = dbHandler.Delete("dm_relation", Cond("to_content_id", content.Value("cid")).Cond("to_type", content.ContentType()), tx)
			if err != nil {
				tx.Rollback()
				message := "[handler.deleteByContent]Can not delete relation."
				util.Error(message + err.Error())
				return errors.New(message)
			}
		}

		//Delete location
		err = content.GetLocation().Delete(tx)
		if err != nil {
			tx.Rollback()
		} else {
			//TODO: delete version if there is.

			//Delete content
			err = content.Delete(tx)
			if err != nil {
				tx.Rollback()
			}
		}

		//Invoke callback
		matchData := map[string]interface{}{"content_type": content.ContentType()}
		if content.Definition().HasLocation {
			hierachy := content.GetLocation().Hierarchy
			matchData["under"] = strings.Split(hierachy, "/")
		}
		err = ch.InvokeCallback("delete", true, matchData, content, tx)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "Invoking callback error.")
		}

		err = tx.Commit()
		if err != nil {
			message := "[handler.deleteByContent]Can not commit."
			util.Error(message + err.Error())
			return errors.New(message)
		}

		//invoke callback
		err = ch.InvokeCallback("deleted", false, matchData, content)

	}
	return nil
}

//Format of fields: eg. title:"test", modified: 12121
func (handler *ContentHandler) store(parentID int, contentType string, fields map[string]interface{}) {
	handler.Validate(contentType, fields)
}

func (handler *ContentHandler) UpdateRelation(content contenttype.ContentTyper) {

}
