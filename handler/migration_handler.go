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
func (mh *MigrationHandler) ImportAContent(contentType string, contentData string) error {
	contentMap := map[string]interface{}{}
	json.Unmarshal([]byte(contentData), &contentMap)
	cuid := contentMap["cuid"].(string)
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
	json.Unmarshal([]byte(contentData), content)

	content.SetValue("cid", 0)
	util.Log("import", "Saving content first.")
	err = content.Store(tx)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Can not saved. Rolled back.")
	}
	util.Log("contenthandler.import", "Content saved. cuid: "+content.Value("cuid").(string)+", id: "+strconv.Itoa(content.GetCID()))

	parentUID := contentMap["parent_uid"].(string)
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
	contentType := ""
	for key, _ := range dataMap {
		contentType = key
		break
	}

	contentData, err := json.Marshal(dataMap[contentType])
	if err != nil {
		return err
	}
	err = mh.ImportAContent(contentType, string(contentData))
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
		location["content_uid"] = content.Value("cuid")
		delete(location, "content_id")
		delete(location, "parent_id")
		//todo: replace main_id with main_uid

		delete(location, "id")
		delete(location, "hierarchy")
	}

	contentMap["parent_uid"] = parent.GetLocation().UID

	jsonObject := map[string]interface{}{content.ContentType(): contentMap}
	data, err = json.Marshal(jsonObject)
	return string(data), err
}
