//Author xc, Created on 2019-05-12 11:27
//{COPYRIGHTS}
package handler

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/util"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type MigrationHandler struct{}

//Import line by line. The imported file should be "up do down" in the tree structure.
// Options about: main location, relations, version.
//
//In the end of import the import do below:
//- update main_id where location is not the main_id.
//- update relations content/location id
//
func (mh *MigrationHandler) ImportAContent(contentType string, cuid string, parentUID string, contentData []byte) error {
	existing, _ := Querier().FetchByCUID(contentType, cuid)
	if existing != nil {
		//todo: for existing, maybe just update - need an option.
		return errors.New("It's there already. cuid: " + cuid)
	}

	content := entity.NewInstance(contentType)
	contentDef := content.Definition()
	tx, err := db.CreateTx()
	if err != nil {
		return errors.Wrap(err, "Error in getting transaction.")
	}
	json.Unmarshal(contentData, content)

	content.SetValue("cid", 0)
	util.Log("import", "Saving content first.")
	err = content.Store(tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Can not saved. Rolled back.")
	}
	util.Log("contenthandler.import", "Content saved. cuid: "+content.Value("cuid").(string)+", id: "+strconv.Itoa(content.GetCID()))

	parent, err := Querier().FetchByUID(parentUID)
	if parent == nil {
		return errors.Wrap(err, "Can not find parent.")
	}
	parentID := parent.GetLocation().ID
	if contentDef.HasLocation {
		location := content.GetLocation()
		location.ID = 0
		location.ParentID = parentID
		location.ContentID = content.GetCID()
		err = location.Store(tx)
		util.Log("contenthandler.import", "Saving location.")
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "Can not save location")
		}
		util.Log("contenthandler.import", "Location saved. uid: "+location.UID+", new id:"+strconv.Itoa(location.ID))
	} else {
		hasParent := util.Contains(content.IdentifierList(), "parent_id")
		if hasParent {
			content.SetValue("parent_id", parentID)
		}
	}

	//TODO: things with relations.

	//TODO: udpate main id where main_uid is not the same as uid

	tx.Commit()
	util.Log("contenthandler.import", "Committed")
	return nil
}

//Import one line
func (mh *MigrationHandler) ImportALine(data []byte) error {
	dataMap := map[string]interface{}{}
	json.Unmarshal(data, &dataMap)
	contentType := dataMap["content_type"].(string)
	cuid := dataMap["cuid"].(string)
	parentUID := dataMap["parent_uid"].(string)
	contentData, _ := json.Marshal(dataMap["data"])

	err := mh.ImportAContent(contentType, cuid, parentUID, contentData)
	return err
}

func (mh *MigrationHandler) RevertImport() {

}

//Verify if the target has need model(contenttype and fields) to import.
func (mh *MigrationHandler) VerifyModel() {

}

//Export to json
func (mh *MigrationHandler) Export(content contenttype.ContentTyper, parent contenttype.ContentTyper) (string, error) {
	data, err := json.Marshal(content)
	contentMap := map[string]interface{}{}
	json.Unmarshal(data, &contentMap)
	delete(contentMap, "id")
	if content.Definition().HasLocation {
		location := contentMap["location"].(map[string]interface{})
		delete(location, "content_id")
		delete(location, "parent_id")
		//todo: replace main_id with main_uid

		delete(location, "id")
		delete(location, "hierarchy")

	}

	//relations.
	relationList := content.GetRelations().List
	if len(relationList) > 0 {
		for i, relation := range relationList {
			var (
				fromLocationUID string
				fromContentUID  string
			)
			if fromLocationID := relation.FromLocation; fromLocationID != 0 {
				fromContent, err := Querier().FetchByID(fromLocationID)
				if err != nil {
					return "", errors.Wrap(err, "From location not found in content relation. from_location_id: "+strconv.Itoa(relation.FromLocation))
				}
				fromLocationUID = fromContent.GetLocation().UID
			} else if fromContentID := relation.FromContentID; fromContentID != 0 {
				fromContent, err := Querier().FetchByContentID(relation.FromType, fromContentID)
				if err != nil {
					return "", errors.Wrap(err, "From content not found in content relation. from_content_id: "+strconv.Itoa(relation.FromContentID))
				}
				fromContentUID = fromContent.Value("cuid").(string)
			}
			relationsMap := contentMap["relations"].(map[string]interface{})
			currentRelation := relationsMap["list"].([]interface{})[i].(map[string]interface{})
			currentRelation["from_location_uid"] = fromLocationUID
			currentRelation["from_content_uid"] = fromContentUID
		}
	}

	jsonObject := map[string]interface{}{
		"content_type": content.ContentType(),
		"parent_uid":   parent.GetLocation().UID,
		"cuid":         content.Value("cuid").(string),
		"data":         contentMap,
	}
	data, err = json.Marshal(jsonObject)
	return string(data), err
}
