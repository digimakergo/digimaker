//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package handler

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/
import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/xc/digimaker/core/contenttype"
	"github.com/xc/digimaker/core/db"
	"github.com/xc/digimaker/core/fieldtype"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/permission"
	"github.com/xc/digimaker/core/util"

	. "github.com/xc/digimaker/core/db"

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

// Validate validates and returns a validation result.
func (ch *ContentHandler) Validate(contentType string, fieldsDef map[string]contenttype.FieldDef, inputs map[string]interface{}) (bool, ValidationResult) {
	//todo: check max length
	//todo: check all kind of validation
	result := ValidationResult{Fields: map[string]string{}}

	//check required
	for identifier, fieldDef := range fieldsDef {
		input, fieldExists := inputs[identifier]
		fieldResult := ""
		//validat both required and others together.
		isEmpty := fieldtype.IsEmptyInput(input)
		if fieldDef.Required {
			if fieldExists && isEmpty || !fieldExists {
				fieldResult = "1"
			} else {
				if fieldExists {
					field := fieldtype.NewField(fieldsDef[identifier].FieldType)
					field.LoadFromInput(input)
					if fieldDef.Required && field.IsEmpty() { //not empty input, but empty in this fieldtype. eg. {} or [] for json fieldtype
						fieldResult = "1"
					} else {
						if valid, res := field.Validate(input, fieldDef.Validation); !valid {
							fieldResult = res
						}
					}
				}
			}
		}
		if fieldResult != "" {
			result.Fields[identifier] = fieldResult
		}
	}

	//validate from contenttype handler
	contentTypeHanlder := GetContentTypeHandler(contentType)
	if contentTypeHanlder != nil {
		log.Debug("Validating from content type handler", "contenthandler.validate", ch.Context)
		contentTypeHanlder.Validate(inputs, &result)
	}
	return result.Passed(), result
}

//Store content. Note it doesn't rollback - please rollback in invoking part if error happens.
//If it's no-location content, ingore the parentID.
func (ch *ContentHandler) storeCreatedContent(content contenttype.ContentTyper, userId int, tx *sql.Tx, parentID ...int) error {
	if content.GetCID() == 0 {
		log.Debug("Content is new.", "contenthandler.StoreCreatedContent", ch.Context)
	}
	err := content.Store(tx)
	if err != nil {
		return err
	}

	log.Debug("Content is saved. id :"+strconv.Itoa(content.GetCID())+". ", "contenthandler.StoreCreatedContent", ch.Context)

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
		location := content.GetLocation()
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

		//todo: set location to the content
		log.Debug("Location is saved. location id :"+strconv.Itoa(location.ID)+". ", "contenthandler.StoreCreatedContent", ch.Context)
	}
	return nil
}

// Create creates a content(same behavior as Draft&Publish but store published version directly)
func (ch *ContentHandler) Create(contentType string, inputs map[string]interface{}, userId int, parentID int) (contenttype.ContentTyper, ValidationResult, error) {

	parent, _ := querier.FetchByID(parentID)
	if parent == nil {
		return nil, ValidationResult{}, errors.New("parent doesn't exist. parent id: " + strconv.Itoa(parentID))
	}
	realData := map[string]interface{}{
		"under":       strings.Split(parent.GetLocation().Hierarchy, "/"),
		"contenttype": contentType,
	} //todo: support more conditions
	if !permission.HasAccessTo(userId, "content/create", realData, ch.Context) {
		return nil, ValidationResult{}, errors.New("User doesn't have access to create")
	}

	//todo: add validation callback.
	contentDefinition, _ := contenttype.GetDefinition(contentType)
	fieldsDefinition := contentDefinition.FieldMap

	//Validate
	valid, validationResult := ch.Validate(contentType, fieldsDefinition, inputs)
	if !valid {
		return nil, validationResult, nil
	}

	//Create empty content instance and set value
	content := contenttype.NewInstance(contentType)
	for identifier, input := range inputs {
		if fieldDef, ok := fieldsDefinition[identifier]; ok {
			fieldType := fieldDef.FieldType
			field := fieldtype.NewField(fieldType)
			field.LoadFromInput(input)
			err := content.SetValue(identifier, field)
			if err != nil {
				return nil, ValidationResult{}, errors.Wrap(err, "Can not set input to "+identifier)
			}
		}
	}
	now := int(time.Now().Unix())
	content.SetValue("published", now)
	content.SetValue("modified", now)
	content.SetValue("author", userId)
	content.SetValue("cuid", util.GenerateUID())

	log.StartTiming(ch.Context, "contenthandler_create.database")
	log.Debug("Validation passed. Start saving content.", "contenthandler.Create", ch.Context)
	//Create transaction
	database, err := db.DB()
	if err != nil {
		return nil, ValidationResult{}, errors.New("Can't get db connection.")
	}
	tx, err := database.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, ValidationResult{}, errors.New("Can't get transaction.")
	}

	//call content type handler
	contentTypeHandler := GetContentTypeHandler(contentType)
	if contentTypeHandler != nil {
		log.Debug("Calling handler for "+contentType, "contenthandler.create", ch.Context)
		err := contentTypeHandler.New(content, tx, parentID)
		if err != nil {
			tx.Rollback()
			log.Error("Error from callback: "+err.Error(), "contenthandler.create", ch.Context)
			return nil, ValidationResult{}, err
		}
	}

	//Save content and location if needed
	versionIfNeeded := 1
	if contentDefinition.HasVersion {
		content.SetValue("version", versionIfNeeded)
	}

	err = ch.storeCreatedContent(content, userId, tx, parentID)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error(), "contenthandler.Create", ch.Context)
		return nil, ValidationResult{}, errors.Wrap(err, "Create content error")
	}

	//Save version if needed
	if contentDefinition.HasVersion {
		log.Debug("creating version", "contenthandler.Create", ch.Context)
		_, err = ch.CreateVersion(content, versionIfNeeded, tx)
		if err != nil {
			log.Error(err.Error(), "contenthandler.Create", ch.Context)
			return nil, ValidationResult{}, errors.Wrap(err, "Create version error.")
		}
		log.Debug("Created version: "+strconv.Itoa(versionIfNeeded), "contenthandler.Create", ch.Context)
	}

	//Invoke callback
	matchData := map[string]interface{}{"parent_id": parentID,
		"content_type": contentType}
	if contentDefinition.HasLocation {
		hierachy := content.GetLocation().Hierarchy
		matchData["under"] = strings.Split(hierachy, "/")
	}

	err = ch.InvokeCallback("create", true, matchData, content, tx, inputs)
	if err != nil {
		tx.Rollback()
		return nil, ValidationResult{}, errors.Wrap(err, "Invoking callback error.")
	}

	//Commit all operations
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error("Commit error: "+err.Error(), "contenthandler.create", ch.Context)
		return nil, ValidationResult{}, errors.Wrap(err, "Commit error.")
	}

	log.Debug("Data committed.", "contenthandler.create", ch.Context)

	log.TickTiming(ch.Context, "contenthandler_create.database")
	return content, ValidationResult{}, nil
}

//InvokeCallback invokes callbacks based on condition match result
//see content_handler.json/yaml for conditions.
func (ch ContentHandler) InvokeCallback(event string, stopOnError bool, matchData map[string]interface{}, content contenttype.ContentTyper, params ...interface{}) error {
	operationHandlerList, matchInfo := GetOperationHandlerByCondition(event, matchData)
	count := len(operationHandlerList)
	if count > 0 {
		identifierList := []string{}
		for _, operationHandler := range operationHandlerList {
			identifierList = append(identifierList, operationHandler.Identifier)
		}
		log.Debug("Matched callbacks: "+strings.Join(identifierList, ","), "contenthandler.invoke_callback", ch.Context)
	} else {
		log.Debug("No callbacks matched.", "contenthandler.invoke_callback", ch.Context)
	}

	for _, info := range matchInfo {
		log.Debug(info, "callback_match", ch.Context)
	}
	for i, operationHandler := range operationHandlerList {
		log.Debug(strconv.Itoa(i+1)+"/"+strconv.Itoa(count)+
			" Invoking operation "+operationHandler.Identifier+" on "+event, "contehandler.invoke_callback", ch.Context)
		err := operationHandler.Execute(event, content, params...)
		if err != nil && stopOnError {
			log.Error("Error when invoking operation handler "+operationHandler.Identifier+":"+err.Error(), "contehandler.invoke_callback", ch.Context)
			return err
		}
		log.Debug("Invoking ended on "+operationHandler.Identifier, "contehandler.invoke_callback", ch.Context)
	}
	return nil
}

//CreateVersion creates a new version.
//It doesn't validate version number is increment
func (ch ContentHandler) CreateVersion(content contenttype.ContentTyper, versionNumber int, tx *sql.Tx) (int, error) {
	log.Debug("Creating version: "+strconv.Itoa(versionNumber), "contenthandler.CreateVersion", ch.Context)
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
	log.Debug("Version created.", "contenthandler.CreateVersion", ch.Context)
	return version.ID, nil
}

func (ch ContentHandler) UpdateByContentID(contentType string, contentID int, inputs map[string]interface{}, userId int) (bool, ValidationResult, error) {
	content, err := Querier().FetchByContentID(contentType, contentID)
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Failed to get content via content id.")
	}
	if content.GetCID() == 0 {
		return false, ValidationResult{}, errors.Wrap(err, "Got empty content.")
	}

	return ch.Update(content, inputs, userId)
}

func (ch ContentHandler) UpdateByID(id int, inputs map[string]interface{}, userId int) (bool, ValidationResult, error) {
	content, err := Querier().FetchByID(id)
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Failed to get content via id.")
	}
	if content.GetCID() == 0 {
		return false, ValidationResult{}, errors.Wrap(err, "Got empty content.")
	}

	return ch.Update(content, inputs, userId)
}

//Update content.
//The inputs doesn't need to include all required fields. However if it's there,
// it will check if it's required&empty
func (ch ContentHandler) Update(content contenttype.ContentTyper, inputs map[string]interface{}, userId int) (bool, ValidationResult, error) {
	accessRealData := map[string]interface{}{
		"contenttype": content.ContentType(),
		"under":       content.GetLocation().ID,
	}
	if content.Value("author").(int) == userId {
		accessRealData["author"] = "self"
	}
	if !permission.HasAccessTo(userId, "content/update", accessRealData, ch.Context) {
		return false, ValidationResult{}, errors.New("User " + strconv.Itoa(userId) + " doesn't have access to update")
	}

	//Validate
	log.Debug("Validating", "contenthandler.update", ch.Context)
	contentType := content.ContentType()
	contentDef, _ := contenttype.GetDefinition(contentType)

	valid, result := ch.Validate(contentType, contentDef.FieldMap, inputs)
	if !valid {
		return valid, result, nil
	}

	//Save to new version
	tx, err := db.CreateTx()
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Create transaction error.")
	}

	//todo: update relations

	//Set content.

	fieldsDefinition := contentDef.FieldMap
	for identifier, fieldDefinition := range fieldsDefinition {
		if input, ok := inputs[identifier]; ok {
			fieldType := fieldDefinition.FieldType
			field := fieldtype.NewField(fieldType)
			field.LoadFromInput(input)
			err := content.SetValue(identifier, field)
			if err != nil {
				return false, ValidationResult{}, errors.Wrap(err, "Can not set input to "+identifier)
			}
		}
	}

	//Save new version and set to content
	if contentDef.HasVersion {
		version := content.Value("version").(int) + 1
		log.Debug("Creating new version: "+strconv.Itoa(version), "contenthandler.update", ch.Context)
		_, err := ch.CreateVersion(content, version, tx)
		if err != nil {
			//todo: rollback here not inside.
			log.Debug("Create new version failed: "+err.Error(), "contenthandler.update", ch.Context)
			return false, ValidationResult{}, errors.Wrap(err, "Can not save version.")
		}
		log.Debug("New version created", "contenthandler.update", ch.Context)
		content.SetValue("version", version)
	}
	now := int(time.Now().Unix())
	content.SetValue("modified", now)

	//Save update content.
	log.Debug("Saving content", "contenthandler.update", ch.Context)
	err = content.Store(tx)
	if err != nil {
		tx.Rollback()
		log.Debug("Saving content failed: "+err.Error(), "contenthandler.update", ch.Context)
		return false, ValidationResult{}, errors.Wrap(err, "Saving content error. ")
	}

	//Updated location related
	if contentDef.HasLocation {
		//todo: update all locations to this content.
		location := content.GetLocation()
		location.Name = GenerateName(content)
		//todo: location.IdentifierPath =
		log.Debug("Updating location info", "contenthandler.update", ch.Context)
		err := location.Store(tx)
		if err != nil {
			tx.Rollback()
			log.Debug("Updating location failed: "+err.Error(), "contenthandler.update", ch.Context)
			return false, ValidationResult{}, errors.Wrap(err, "Updating location info error.")
		}
		log.Debug("Location updated", "contenthandler.update", ch.Context)
	}

	//Invoke callback
	matchData := map[string]interface{}{"content_type": contentType}
	if content.Definition().HasLocation {
		hierachy := content.GetLocation().Hierarchy
		matchData["under"] = strings.Split(hierachy, "/")
	}

	//todo: maybe old content need to pass to callback.
	err = ch.InvokeCallback("update", true, matchData, content, tx, inputs)
	if err != nil {
		tx.Rollback()
		return false, ValidationResult{}, errors.Wrap(err, "Invoking callback error.")
	}

	log.Debug("All done. Commitinng.", "contenthandler.update", ch.Context)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error("Commit error: "+err.Error(), "contenthandler.update", ch.Context)
		return false, ValidationResult{}, errors.Wrap(err, "Commit error.")
	}

	return true, ValidationResult{}, nil
}

//Delete content by location id
func (ch ContentHandler) DeleteByID(id int, userId int, toTrash bool) error {
	content, err := Querier().FetchByID(id)
	//todo: check how many. if more than 1, delete current only(and set main_id if needed)
	if err != nil {
		return errors.New("[handler.delete]Content doesn't exist with id: " + strconv.Itoa(id))
	}
	err = ch.DeleteByContent(content, userId, toTrash)
	return err
}

//Delete content, relations and location.
//Note: this is only for when there is 1 location.
//  You need to judge if there are more than one locations before invoking this.
func (ch ContentHandler) DeleteByContent(content contenttype.ContentTyper, userId int, toTrash bool) error {
	//todo: check delete children. There should be more consideration if there are more children.

	accessRealData := map[string]interface{}{
		"contenttype": content.ContentType(),
		"under":       strings.Split(content.GetLocation().Hierarchy, "/"),
	}
	if content.Value("author").(int) == userId {
		accessRealData["author"] = "self"
	}
	if !permission.HasAccessTo(userId, "content/delete", accessRealData, ch.Context) {
		return errors.New("User " + strconv.Itoa(userId) + " Doesn't have access to delete")
	}

	//Delete location
	location := content.GetLocation()
	if location.CountLocations() > 1 {
		return errors.New("There are more than 1 location. Remove location first.")
	} else {
		database, err := db.DB()
		if err != nil {
			log.Error(err.Error(), "", ch.Context)
			return errors.New("[handler.deleteByContent]Can not create connection.")
		}
		tx, err := database.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			tx.Rollback()
			message := "[handler.deleteByContent]Can not create transaction."
			log.Error(message+err.Error(), "", ch.Context)
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
				log.Error(message+err.Error(), "", ch.Context)
				return errors.New(message)
			}
		}

		//Delete location
		err = content.GetLocation().Delete(tx)
		if err != nil {
			tx.Rollback()
		} else {
			//delete versions
			if content.Definition().HasVersion {
				dbHanlder := db.DBHanlder()
				dbHanlder.Delete("dm_version", db.Cond("content_type", content.ContentType()).
					Cond("content_id", content.GetCID()), tx)
			}

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
			log.Error(message+err.Error(), "", ch.Context)
			return errors.New(message)
		}

		//invoke callback
		err = ch.InvokeCallback("deleted", false, matchData, content)

	}
	return nil
}

func (handler *ContentHandler) UpdateRelation(content contenttype.ContentTyper) {

}

//Generate name based on name_pattern from definition.
func GenerateName(content contenttype.ContentTyper) string {
	pattern := content.Definition().NamePattern

	//Get variables
	vars := util.GetStrVar(pattern)
	values := map[string]string{}
	for i := range vars {
		varName := vars[i]
		field := content.Value(varName)
		switch field.(type) {
		//support Text for now. todo: support all fields.
		//todo: support created, modified time
		case fieldtype.Text:
			values[varName] = field.(fieldtype.Text).FieldValue().(string)
		}
	}

	//Replace variable with value
	result := util.ReplaceStrVar(pattern, values)
	return result
}
