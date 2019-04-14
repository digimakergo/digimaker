//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package handler

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"dm/model"
	"dm/model/entity"
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
	Content *entity.Article
}

func (content ContentHandler) CreateLocation(parentID int) {
	location := entity.Location{ParentID: parentID}
	location.Store()
}

//Create draft of a content. parent_id will be -1 in this case
func (handler *ContentHandler) Create(title string, parentID int) error {
	//Save content
	now := int(time.Now().Unix())
	article := entity.Article{Author: 1, Published: now, Modified: now}
	article.Store()

	//Save location
	location := entity.Location{ParentID: parentID, ContentID: article.CID, UID: util.GenerateUID()}
	err := location.Store()
	if err != nil {
		return err
	}
	return nil
}

//Format of fields: eg. title:"test", modified: 12121
func (handler *ContentHandler) store(parentID int, contentType string, fields map[string]interface{}) {
	handler.Validate(contentType, fields)
}

//return a validation result instance
func (handler *ContentHandler) Validate(contentType string, fields map[string]interface{}) error {
	definition, err := model.GetContentDefinition(contentType)
	if err != nil {
		return errors.Wrap(err, "Error in "+contentType)
	}
	//check required
	fieldsDef := definition.Fields
	var requiredFields []string
	for identifier, fieldDef := range fieldsDef {
		if fieldDef.Required {
			if _, found := fields[identifier]; !found {
				requiredFields = append(requiredFields, identifier)
			} else {
				//check if it's empty
			}
		}
	}
	//Validate with hook

	//Store with hook
	return nil
}

func (content ContentHandler) Store() error {
	//Store Location
	return nil
}

func (content ContentHandler) Draft(contentType string, parentID int) error {
	//create empty
	now := int(time.Now().Unix())
	article := entity.Article{Author: 1, Published: now, Modified: now}
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

func (content ContentHandler) Publish() {

}
