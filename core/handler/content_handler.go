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

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/permission"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/pkg/errors"
)

type InputMap map[string]interface{}

var ErrorNoPermission = errors.New("The user doesn't have access to the action.")

// Validate validates and returns a validation result.
// Validate is used in Create and Update, but can be also used separately
//  eg. when you have several steps, you want to validate one step only(only the fields in that step).
func Validate(contentType string, fieldsDef map[string]definition.FieldDef, inputs InputMap, checkAllRequired bool) (bool, ValidationResult) {
	//todo: check max length
	//todo: check all kind of validation
	result := ValidationResult{Fields: map[string]string{}}

	//check required
	for identifier, fieldDef := range fieldsDef {
		input, fieldExists := inputs[identifier]
		fieldResult := ""
		//validat both required and others together.
		if fieldExists {
			isEmpty := fieldtype.IsEmptyInput(input)
			if fieldDef.Required && isEmpty {
				fieldResult = "1"
			} else {
				fieldtypeDef := fieldtype.GetFieldtype(fieldDef.FieldType)
				handler := fieldtypeDef.NewHandler(fieldDef)
				if handler != nil {
					_, err := handler.LoadInput(input, "")
					if _, ok := err.(fieldtype.ValidationError); ok {
						fieldResult = err.Error()
					}
					if _, ok := err.(fieldtype.EmptyError); ok {
						fieldResult = "1"
					}
				}
			}
		} else if fieldDef.Required && checkAllRequired {
			fieldResult = "1"
		}
		if fieldResult != "" {
			result.Fields[identifier] = fieldResult
		}
	}

	return result.Passed(), result
}

//Store content. Note it doesn't rollback - please rollback in invoking part if error happens.
//If it's no-location content, ingore the parentID.
func storeCreatedContent(ctx context.Context, content contenttype.ContentTyper, userId int, tx *sql.Tx, parentID int) error {
	if content.GetCID() == 0 {
		log.Debug("Content is new.", "contenthandler.StoreCreatedContent", ctx)
	}

	contentDefinition := content.Definition()
	if !contentDefinition.HasLocation {
		content.SetValue("location_id", parentID)
	}

	for identifier, fieldDef := range contentDefinition.FieldMap {
		handler := fieldtype.GethHandler(fieldDef)
		if storeHandler, ok := handler.(fieldtype.StoreHandler); ok {
			err := storeHandler.Store(ctx, content.Value(identifier), content.ContentType(), content.GetCID(), tx)
			if err != nil {
				return err
			}
		}
	}

	err := content.Store(ctx, tx)
	if err != nil {
		return err
	}

	log.Debug("Content is saved. id :"+strconv.Itoa(content.GetCID())+". ", "contenthandler.StoreCreatedContent", ctx)

	//todo: deal with relations

	if contentDefinition.HasLocation {
		parent, err := contenttype.GetLocationByID(parentID)
		if err != nil {
			return errors.Wrap(err, "Can not get parent location with "+strconv.Itoa(parentID))
		}

		//Save location
		location := content.GetLocation()
		location.ParentID = parentID
		location.ContentID = content.GetCID()
		location.ContentType = content.ContentType()
		location.UID = util.GenerateUID()
		contentName := GenerateName(content)
		location.IdentifierPath = parent.IdentifierPath + "/" + util.NameToIdentifier(contentName)
		location.Depth = parent.Depth + 1
		location.Section = parent.Section
		location.Name = contentName

		err = location.Store(ctx, tx)

		if err != nil {
			return errors.Wrap(err, "Transaction failed in location when saving new location.")
		}
		location.Hierarchy = parent.Hierarchy + "/" + strconv.Itoa(location.ID)
		location.MainID = location.ID
		err = location.Store(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "Transaction failed in location when saving location for main_id and hierarchy.")
		}

		//todo: set location to the content
		log.Debug("Location is saved. location id :"+strconv.Itoa(location.ID)+". ", "contenthandler.StoreCreatedContent", ctx)
	}
	return nil
}

// Create creates a content(same behavior as Draft&Publish but store published version directly)
func Create(ctx context.Context, userID int, contentType string, inputs InputMap, parentID int) (contenttype.ContentTyper, ValidationResult, error) {
	parent, _ := query.FetchByID(ctx, parentID)
	if parent == nil {
		return nil, ValidationResult{}, errors.New("parent doesn't exist. parent id: " + strconv.Itoa(parentID))
	}

	contentDefinition, _ := definition.GetDefinition(contentType)
	fieldsDefinition := contentDefinition.FieldMap

	if !permission.CanCreate(ctx, parent, contentType, userID) {
		return nil, ValidationResult{}, errors.New("User doesn't have access to create")
	}

	//Validate
	valid, validationResult := Validate(contentType, fieldsDefinition, inputs, true)
	if !valid {
		return nil, validationResult, nil
	}

	//validate from contenttype handler
	contentTypeHandler := GetContentTypeHandler(contentType)
	if validator, ok := contentTypeHandler.(ContentTypeHandlerValidate); ok {
		log.Debug("Validating from content type handler", "contenthandler.validate", ctx)
		valid, validationResult = validator.ValidateCreate(ctx, inputs, parentID)
		if !valid {
			return nil, validationResult, nil
		}
	}

	//Create empty content instance and set value
	content := contenttype.NewInstance(contentType)
	for identifier, fieldDef := range fieldsDefinition {
		if input, ok := inputs[identifier]; ok {
			handler := fieldtype.GethHandler(fieldDef)

			//set value
			value, _ := handler.LoadInput(input, "")

			//Invoke BeforeStore
			var err error
			if fieldtypeEvent, ok := handler.(fieldtype.Event); ok {
				log.Debug("Invoking before storing on field: "+identifier, "handler", ctx)
				value, err = fieldtypeEvent.BeforeStore(value, nil, "create")
				if err != nil {
					return nil, ValidationResult{}, err
				}
			}
			content.SetValue(identifier, value)
		}
	}

	now := int(time.Now().Unix())
	if contentDefinition.HasLocation || contentDefinition.HasDataField("published") {
		content.SetValue("published", now)
	}
	if contentDefinition.HasLocation || contentDefinition.HasDataField("modified") {
		content.SetValue("modified", now)
	}
	if contentDefinition.HasLocation || contentDefinition.HasDataField("author") {
		content.SetValue("author", userID)
	}
	if contentDefinition.HasLocation || contentDefinition.HasDataField("cuid") {
		content.SetValue("cuid", util.GenerateUID())
	}

	log.StartTiming(ctx, "contenthandler_create.database")
	log.Debug("Validation passed. Start saving content.", "contenthandler.Create", ctx)
	//Create transaction
	tx, err := db.CreateTx()
	if err != nil {
		return nil, ValidationResult{}, errors.New("Can't get transaction")
	}

	//Save content and location if needed
	versionIfNeeded := 1
	if contentDefinition.HasVersion {
		content.SetValue("version", versionIfNeeded)
	}

	err = storeCreatedContent(ctx, content, userID, tx, parentID)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error(), "contenthandler.Create", ctx)
		return nil, ValidationResult{}, errors.Wrap(err, "Create content error")
	}

	//Save version if needed
	if contentDefinition.HasVersion {
		log.Debug("creating version", "contenthandler.Create", ctx)
		_, err = CreateVersion(ctx, content, versionIfNeeded, tx)
		if err != nil {
			log.Error(err.Error(), "contenthandler.Create", ctx)
			return nil, ValidationResult{}, errors.Wrap(err, "Create version error.")
		}
		log.Debug("Created version: "+strconv.Itoa(versionIfNeeded), "contenthandler.Create", ctx)
	}

	//call content type handler
	if creater, ok := contentTypeHandler.(ContentTypeHandlerCreate); ok {
		log.Debug("Calling handler for "+contentType, "contenthandler.create", ctx)
		err := creater.Create(ctx, content, inputs, parentID)
		if err != nil {
			tx.Rollback()
			log.Error("Error from callback: "+err.Error(), "contenthandler.create", ctx)
			return nil, ValidationResult{}, err
		}
	}

	//Invoke callback
	matchData := map[string]interface{}{"parent_id": parentID,
		"content_type": contentType}
	if contentDefinition.HasLocation {
		hierachy := content.GetLocation().Hierarchy
		matchData["under"] = strings.Split(hierachy, "/")
	}

	err = InvokeCallback(ctx, "create", true, matchData, content, tx, inputs)
	if err != nil {
		tx.Rollback()
		return nil, ValidationResult{}, errors.Wrap(err, "Invoking callback error.")
	}

	//Commit all operations
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error("Commit error: "+err.Error(), "contenthandler.create", ctx)
		return nil, ValidationResult{}, errors.Wrap(err, "Commit error.")
	}

	log.Debug("Data committed.", "contenthandler.create", ctx)

	log.EndTiming(ctx, "contenthandler_create.database")
	return content, ValidationResult{}, nil
}

//InvokeCallback invokes callbacks based on condition match result
//see content_handler.json/yaml for conditions.
func InvokeCallback(ctx context.Context, event string, stopOnError bool, matchData map[string]interface{}, content contenttype.ContentTyper, params ...interface{}) error {
	operationHandlerList, matchInfo := GetOperationHandlerByCondition(event, matchData)
	count := len(operationHandlerList)
	if count > 0 {
		identifierList := []string{}
		for _, operationHandler := range operationHandlerList {
			identifierList = append(identifierList, operationHandler.Identifier)
		}
		log.Debug("Matched callbacks: "+strings.Join(identifierList, ","), "contenthandler.invoke_callback", ctx)
	} else {
		log.Debug("No callbacks matched.", "contenthandler.invoke_callback", ctx)
	}

	for _, info := range matchInfo {
		log.Debug(info, "callback_match", ctx)
	}
	for i, operationHandler := range operationHandlerList {
		log.Debug(strconv.Itoa(i+1)+"/"+strconv.Itoa(count)+
			" Invoking operation "+operationHandler.Identifier+" on "+event, "contehandler.invoke_callback", ctx)
		err := operationHandler.Execute(ctx, event, content, params...)
		if err != nil && stopOnError {
			log.Error("Error when invoking operation handler "+operationHandler.Identifier+":"+err.Error(), "contehandler.invoke_callback", ctx)
			return err
		}
		log.Debug("Invoking ended on "+operationHandler.Identifier, "contehandler.invoke_callback", ctx)
	}
	return nil
}

//CreateVersion creates a new version.
//It doesn't validate version number is increment
func CreateVersion(ctx context.Context, content contenttype.ContentTyper, versionNumber int, tx *sql.Tx) (int, error) {
	log.Debug("Creating version: "+strconv.Itoa(versionNumber), "contenthandler.CreateVersion", ctx)
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
	err = version.Store(ctx, tx)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "Can not save version on contetent id: "+strconv.Itoa(id))
	}
	log.Debug("Version created.", "contenthandler.CreateVersion", ctx)
	return version.ID, nil
}

func UpdateByContentID(ctx context.Context, contentType string, contentID int, inputs InputMap, userId int) (bool, ValidationResult, error) {
	content, err := query.FetchByCID(ctx, contentType, contentID)
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Failed to get content via content id.")
	}
	if content.GetCID() == 0 {
		return false, ValidationResult{}, errors.Wrap(err, "Got empty content.")
	}

	return Update(ctx, content, inputs, userId)
}

func UpdateByID(ctx context.Context, id int, inputs InputMap, userId int) (bool, ValidationResult, error) {
	content, err := query.FetchByID(ctx, id)
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Failed to get content via id.")
	}
	if content.GetCID() == 0 {
		return false, ValidationResult{}, errors.Wrap(err, "Got empty content.")
	}

	return Update(ctx, content, inputs, userId)
}

//Update content.
//The inputs doesn't need to include all required fields. However if it's there,
// it will check if it's required&empty
func Update(ctx context.Context, content contenttype.ContentTyper, inputs InputMap, userId int) (bool, ValidationResult, error) {
	contentType := content.ContentType()
	contentDef, _ := definition.GetDefinition(contentType)
	fieldsDefinition := contentDef.FieldMap

	//permission check
	if !permission.CanUpdate(ctx, content, userId) {
		return false, ValidationResult{}, errors.New("User " + strconv.Itoa(userId) + " doesn't have access to update")
	}
	//todo: think about merging 'content/update' and 'content/fields_update'
	allowedFields, err := permission.GetUpdateFields(ctx, content, userId)
	if err != nil {
		return false, ValidationResult{}, err
	}
	if len(allowedFields) > 0 && allowedFields[0] != "*" {
		for field, _ := range inputs {
			if _, ok := fieldsDefinition[field]; ok && !util.Contains(allowedFields, field) {
				return false, ValidationResult{}, errors.New("User doesn't have permission to update field " + field + ".")
			}
		}
	}

	//Validate
	log.Debug("Validating", "contenthandler.update", ctx)
	valid, validationResult := Validate(contentType, fieldsDefinition, inputs, false)
	if !valid {
		return false, validationResult, nil
	}

	//validate from contenttype handler
	contentTypeHandler := GetContentTypeHandler(contentType)
	if validator, ok := contentTypeHandler.(ContentTypeHandlerValidate); ok {
		log.Debug("Validating from update handler", "contenthandler.validate", ctx)
		valid, validationResult = validator.ValidateUpdate(ctx, inputs, content)
		if !valid {
			return false, validationResult, nil
		}
	}

	//todo: update relations
	//Set content.
	for identifier, fieldDef := range fieldsDefinition {
		if input, ok := inputs[identifier]; ok {
			//get field from loaded content
			handler := fieldtype.GethHandler(fieldDef)
			value, _ := handler.LoadInput(input, "")
			existing := content.Value(identifier)
			//Invoke BeforeSaving
			if fieldWithEvent, ok := handler.(fieldtype.Event); ok {
				value, err = fieldWithEvent.BeforeStore(value, existing, "")
				if err != nil {
					return false, ValidationResult{}, err
				}
			}
			content.SetValue(identifier, value)
		}
	}

	//Save to new version
	tx, err := db.CreateTx()
	if err != nil {
		return false, ValidationResult{}, errors.Wrap(err, "Create transaction error.")
	}

	//Save new version and set to content
	if contentDef.HasVersion {
		version := content.Value("version").(int) + 1
		log.Debug("Creating new version: "+strconv.Itoa(version), "contenthandler.update", ctx)
		_, err := CreateVersion(ctx, content, version, tx)
		if err != nil {
			//todo: rollback here not inside.
			log.Debug("Create new version failed: "+err.Error(), "contenthandler.update", ctx)
			return false, ValidationResult{}, errors.Wrap(err, "Can not save version.")
		}
		log.Debug("New version created", "contenthandler.update", ctx)
		content.SetValue("version", version)
	}

	if contentDef.HasLocation || contentDef.HasDataField("modified") {
		now := int(time.Now().Unix())
		content.SetValue("modified", now)
	}

	//Save update content.
	log.Debug("Saving content", "contenthandler.update", ctx)

	for identifier, fieldDef := range contentDef.FieldMap {
		handler := fieldtype.GethHandler(fieldDef)
		if storeHandler, ok := handler.(fieldtype.StoreHandler); ok {
			err := storeHandler.Store(ctx, content.Value(identifier), content.ContentType(), content.GetCID(), tx)
			if err != nil {
				tx.Rollback()
				return false, ValidationResult{}, err
			}
		}
	}

	err = content.Store(ctx, tx)
	if err != nil {
		tx.Rollback()
		log.Debug("Saving content failed: "+err.Error(), "contenthandler.update", ctx)
		return false, ValidationResult{}, errors.Wrap(err, "Saving content error. ")
	}

	//Updated location related
	if contentDef.HasLocation {
		//todo: update all locations to this content.
		location := content.GetLocation()
		location.Name = GenerateName(content)
		//todo: location.IdentifierPath =
		log.Debug("Updating location info", "contenthandler.update", ctx)
		err := location.Store(ctx, tx)
		if err != nil {
			tx.Rollback()
			log.Debug("Updating location failed: "+err.Error(), "contenthandler.update", ctx)
			return false, ValidationResult{}, errors.Wrap(err, "Updating location info error.")
		}
		log.Debug("Location updated", "contenthandler.update", ctx)
	}

	//call content type handler
	if updater, ok := contentTypeHandler.(ContentTypeHandlerUpdate); ok {
		log.Debug("Calling handler for "+contentType, "contenthandler.update", ctx)
		err := updater.Update(ctx, content, inputs)
		if err != nil {
			tx.Rollback()
			log.Error("Error from callback: "+err.Error(), "contenthandler.update", ctx)
			return false, ValidationResult{}, err
		}
	}

	//Invoke callback
	matchData := map[string]interface{}{"content_type": contentType}
	if content.Definition().HasLocation {
		hierachy := content.GetLocation().Hierarchy
		matchData["under"] = strings.Split(hierachy, "/")
	}

	//todo: maybe old content need to pass to callback.
	err = InvokeCallback(ctx, "update", true, matchData, content, tx, inputs)
	if err != nil {
		tx.Rollback()
		return false, ValidationResult{}, errors.Wrap(err, "Invoking callback error.")
	}

	log.Debug("All done. Commitinng.", "contenthandler.update", ctx)
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Error("Commit error: "+err.Error(), "contenthandler.update", ctx)
		return false, ValidationResult{}, errors.Wrap(err, "Commit error.")
	}

	return true, ValidationResult{}, nil
}

//Move moves contents to target
//Check delete&create permission.
//note: it dosn't check if target can create sub-children(only check if it can create direct children)
func Move(ctx context.Context, contentIds []int, targetId int, userId int) error {
	target, err := query.FetchByID(ctx, targetId)
	targetLocation := target.GetLocation()
	if err != nil {
		log.Error(err.Error(), "")
		return errors.New("Target not found")
	}

	contents := []contenttype.ContentTyper{}
	for _, id := range contentIds {
		content, err := query.FetchByID(ctx, id)
		if err != nil {
			log.Error(err.Error(), "")
			return errors.New("Content id " + strconv.Itoa(id) + " is not found for this user.")
		}
		contents = append(contents, content)
	}

	//update content
	tx, err := db.CreateTx(ctx)
	if err != nil {
		return errors.New("Internal error")
	}

	for _, content := range contents {
		location := content.GetLocation()

		if !permission.CanDelete(ctx, content, userId) {
			log.Warning("No permission to delete when moving "+strconv.Itoa(location.ID), "")
			tx.Rollback() //error if no commit?
			return ErrorNoPermission
		}
		if !permission.CanCreate(ctx, target, content.ContentType(), userId) {
			log.Warning("No permission to create when moving "+strconv.Itoa(location.ID), "")
			tx.Rollback() //error if no commit?
			return ErrorNoPermission
		}

		location.ParentID = targetId
		oldHiearachy := location.Hierarchy
		newHiearachy := targetLocation.Hierarchy + "/" + strconv.Itoa(location.ID)
		location.Hierarchy = newHiearachy

		oldPath := location.IdentifierPath
		newPath := targetLocation.IdentifierPath + "/" + util.NameToIdentifier(location.Name)
		location.IdentifierPath = newPath

		location.Depth = len(strings.Split(newHiearachy, "/"))
		location.Store(ctx, tx)

		//update location
		subLocations := []contenttype.Location{}
		db.BindEntity(ctx, &subLocations, "dm_location", db.Cond("hierarchy like", oldHiearachy+"/%"))
		for _, subLocation := range subLocations {
			subContent, _ := query.FetchByID(ctx, subLocation.ID)
			if !permission.CanDelete(ctx, subContent, userId) {
				tx.Rollback()
				log.Warning("No permission to delete "+strconv.Itoa(location.ID), "")
				return ErrorNoPermission
			}

			subLocation.Hierarchy = newHiearachy + strings.TrimPrefix(subLocation.Hierarchy, oldHiearachy)
			subLocation.IdentifierPath = newPath + strings.TrimPrefix(subLocation.IdentifierPath, oldPath)
			subLocation.Depth = len(strings.Split(subLocation.Hierarchy, "/"))
			subLocation.Store(ctx, tx)
		}
	}
	tx.Commit()
	return nil
}

//Delete content by content id
func DeleteByCID(ctx context.Context, cid int, contenttype string, userId int) error {
	content, err := query.FetchByCID(ctx, contenttype, cid)
	if err != nil {
		return errors.New("[handler.delete]Content doesn't exist with cid: " + strconv.Itoa(cid))
	}
	err = DeleteByContent(ctx, content, userId, false)
	return err
}

//Delete content by location id
func DeleteByID(ctx context.Context, id int, userId int, toTrash bool) error {
	content, err := query.FetchByID(ctx, id)
	//todo: check how many. if more than 1, delete current only(and set main_id if needed)
	if err != nil {
		return errors.New("[handler.delete]Content doesn't exist with id: " + strconv.Itoa(id))
	}
	err = DeleteByContent(ctx, content, userId, toTrash)
	return err
}

//Delete content, relations and location.
//Note: this is only for when there is 1 location.
//  You need to judge if there are more than one locations before invoking this.
func DeleteByContent(ctx context.Context, content contenttype.ContentTyper, userId int, toTrash bool) error {
	//todo: check delete children. There should be more consideration if there are more children.

	if !permission.CanDelete(ctx, content, userId) {
		return errors.New("User " + strconv.Itoa(userId) + " Doesn't have access to delete. cid: " + strconv.Itoa(content.GetCID()))
	}

	tx, err := db.CreateTx()
	if err != nil {
		tx.Rollback()
		message := "[handler.deleteByContent]Can not create transaction."
		log.Error(message+err.Error(), "", ctx)
		return errors.New(message)
	}

	def := content.Definition()
	if !def.HasLocation {
		err := content.Delete(ctx, tx)
		if err != nil {
			return err
		}
	} else {
		//Delete location
		location := content.GetLocation()
		if location.CountLocations() > 1 {
			return errors.New("There are more than 1 location. Remove location first.")
		} else {
			//Delete relation first. //todo: use relationlist fieldtype
			if relations, ok := content.(contenttype.GetRelations); ok {
				if len(relations.GetRelations()) > 0 {
					err = db.Delete(ctx, "dm_relation", db.Cond("to_content_id", content.Value("cid")).Cond("to_type", content.ContentType()), tx)
					if err != nil {
						tx.Rollback()
						message := "[handler.deleteByContent]Can not delete relation."
						log.Error(message+err.Error(), "", ctx)
						return errors.New(message)
					}
				}
			}

			//Delete location
			err = content.GetLocation().Delete(ctx, tx)
			if err != nil {
				tx.Rollback()
			} else {
				//delete versions
				if content.Definition().HasVersion {
					db.Delete(ctx, "dm_version", db.Cond("content_type", content.ContentType()).
						Cond("content_id", content.GetCID()), tx)
				}

				//Delete content
				err = content.Delete(ctx, tx)
				if err != nil {
					tx.Rollback()
				}
			}
		}
	}

	contentTypeHandler := GetContentTypeHandler(content.ContentType())
	if deleter, ok := contentTypeHandler.(ContentTypeHandlerDelete); ok {
		err = deleter.Delete(ctx, content)
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
	err = InvokeCallback(ctx, "delete", true, matchData, content, tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Invoking callback error.")
	}

	err = tx.Commit()
	if err != nil {
		message := "[handler.deleteByContent]Can not commit."
		log.Error(message+err.Error(), "", ctx)
		return errors.New(message)
	}

	// fieldtype callback for delete
	for identifier, field := range def.FieldMap {
		handler := fieldtype.GethHandler(field)
		if handler != nil {
			if event, ok := handler.(fieldtype.Event); ok {
				event.AfterDelete(content.Value(identifier))
			}
		}
	}

	//invoke callback
	err = InvokeCallback(ctx, "deleted", false, matchData, content)

	return nil
}

func UpdateRelation(content contenttype.ContentTyper) {

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
		//support string for now. todo: support all fields.
		//todo: support created, modified time
		case string:
			values[varName] = field.(string)
		}
	}

	//Replace variable with value
	result := util.ReplaceStrVar(pattern, values)
	return result
}
