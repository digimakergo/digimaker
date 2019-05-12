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

//Import, based on json
func (mh *MigrationHandler) Import(contentType string, contentData string) error {
	content := entity.NewInstance(contentType)
	contentDef := contenttype.GetContentDefinition(contentType)
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

	tx.Commit()
	util.Log("contenthandler.import", "Committed")
	return nil
}

//Export to json
func (mh *MigrationHandler) Export() {

}
