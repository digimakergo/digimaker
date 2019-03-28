package handler

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"context"
	"dm/db"
	"dm/model/entity"
	util "dm/util"
	"errors"

	"github.com/volatiletech/sqlboiler/boil"
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

//Create draft of a content. parent_id will be -1 in this case
func (handler *ContentHandler) Create() error {
	db, err := db.Open()
	if err != nil {
		return nil
	}

	//Convert data
	for identifier := range handler.Content.Fields() {
		field = handler.Content.Field(identifier)
		err := field.SetStoreData(c, identifier)
		if err != nil {
			return errors.New("Store data error. Did not store any. Field: " + identifier + ". Detail: " + err.Error())
		}
	}

	//Save content

	//Save location
	err = c.Location.Insert(context.Background(), db, boil.Infer()) //todo: use a generic way instead of sqlboil.
	if err != nil {
		return err
	}
	return nil
}

func (content ContentHandler) CreateLocation() {

}

func (content ContentHandler) Store() error {
	//Store fields
	fields := content.Fields
	for identifier, field := range fields {
		_, err := field.GetStoredData()
		if err != nil {
			//log it and return higher
			util.Log("Storing content error, break - id: "+string(content.ID)+", field: "+identifier, "error")
			return errors.New("Can not store content")
		}
	}

	//Store Location
	return nil
}

func (content Content) Publish() {

}
