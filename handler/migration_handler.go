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
func (mh *MigrationHandler) Import(contentType string, contentData string) error {
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

	if contentDef.HasLocation {
		location := content.GetLocation()
		location.ID = 0
		location.ContentID = content.GetCID()
		err = location.Store(tx)
		util.Log("contenthandler.import", "Saving location.")
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "Can not save location")
		}
		util.Log("contenthandler.import", "Location saved. uid: "+location.UID+", new id:"+strconv.Itoa(location.ID))
	}

	//TODO: things with relations.

	//TODO: udpate main id where main_uid is not the same as uid

	tx.Commit()
	util.Log("contenthandler.import", "Committed")
	return nil
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
	location := contentMap["location"].(map[string]interface{})

	location["content_uid"] = content.Value("cuid")
	delete(location, "content_id")

	location["parent_uid"] = parent.GetLocation().UID

	delete(location, "parent_id")

	//todo: replace main_id with main_uid

	delete(location, "id")
	delete(location, "hierarchy")

	jsonObject := map[string]interface{}{content.ContentType(): contentMap}
	data, err = json.Marshal(jsonObject)
	return string(data), err
}
