//Author xc, Created on 2019-05-12 11:27
//{COPYRIGHTS}
package handler

import (
	"dm/core/contenttype"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type ExportHandler struct{}

//Export to json
func (eh *ExportHandler) Export(content contenttype.ContentTyper, parent contenttype.ContentTyper) (string, error) {
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
			currentRelation["from_uid"] = fromLocationUID
			currentRelation["from_cuid"] = fromContentUID
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
