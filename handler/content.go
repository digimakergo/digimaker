//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package handler

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"dm/model/entity"
	"dm/util"
	"time"
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

func (content ContentHandler) Store() error {
	//Store Location
	return nil
}

func (content ContentHandler) Draft(contentType string, parentID int) error {
	//create empty
	now := int(time.Now().Unix())
	article := entity.Article{Author: 1, Published: now, Modified: now}
	err := article.Store()

	//Save location
	location := entity.Location{ParentID: -parentID,
		ContentType: contentType,
		ContentID:   article.CID,
		UID:         util.GenerateUID()}
	err = location.Store()
	if err != nil {
		return err
	}
	return nil
}

func (content ContentHandler) Publish() {

}
