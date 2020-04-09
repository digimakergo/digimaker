//Author xc, Created on 2019-05-12 11:27
//{COPYRIGHTS}
package handler

import (
	"dm/core/contenttype"
	"dm/core/db"
	"dm/core/util"
	"dm/core/util/log"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type ImportHandler struct{}

//Import line by line. The imported file should be "up do down" in the tree structure.
// Options about: main location, relations, version.
//
//In the end of import the import do below:
//- update main_id where location is not the main_id.
//- update relations content/location id
//
func (mh *ImportHandler) ImportAContent(contentType string, cuid string, parentUID string, contentData []byte) error {
	existing, _ := Querier().FetchByCUID(contentType, cuid)
	if existing != nil {
		//todo: for existing, maybe just update - need an option.
		return errors.New("It's there already. cuid: " + cuid)
	}

	content := contenttype.NewInstance(contentType)
	contentDef := content.Definition()
	tx, err := db.CreateTx()
	if err != nil {
		return errors.Wrap(err, "Error in getting transaction.")
	}
	json.Unmarshal(contentData, content)

	content.SetValue("cid", 0)
	log.Debug("Saving content first.", "import")
	err = content.Store(tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Can not saved. Rolled back.")
	}
	log.Debug("Content saved. cuid: "+content.Value("cuid").(string)+", id: "+strconv.Itoa(content.GetCID()), "contenthandler.import")

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
		log.Debug("Saving location.", "contenthandler.import")
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "Can not save location")
		}
		log.Debug("Location saved. uid: "+location.UID+", new id:"+strconv.Itoa(location.ID), "contenthandler.import")
	} else {
		hasParent := util.Contains(content.IdentifierList(), "parent_id")
		if hasParent {
			content.SetValue("parent_id", parentID)
		}
	}

	//TODO: things with relations.
	//Relations should be an option, some relation identifiers can be imported(eg. useful links, related images),
	//while others can be just ignored(only keep relation but ignore the orginal content. eg. uesful articles)
	//So there should be option for this.
	//And what if related article also has relations?
	// - recursive relation can be crazy, so thinking more about senario might be better.

	//TODO: udpate main id where main_uid is not the same as uid

	tx.Commit()
	log.Debug("Committed", "contenthandler.import")
	return nil
}

//Import one line
func (mh *ImportHandler) ImportALine(data []byte) error {
	dataMap := map[string]interface{}{}
	json.Unmarshal(data, &dataMap)
	contentType := dataMap["content_type"].(string)
	cuid := dataMap["cuid"].(string)
	parentUID := dataMap["parent_uid"].(string)
	contentData, _ := json.Marshal(dataMap["data"])

	err := mh.ImportAContent(contentType, cuid, parentUID, contentData)
	return err
}

func (mh *ImportHandler) RevertImport() {

}

//Verify if the target has need model(contenttype and fields) to import.
func (mh *ImportHandler) VerifyModel() {

}
