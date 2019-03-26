package models

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"context"
	"dm/db"
	"dm/models/orm"
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

type Content struct {
	*orm.Location `json:"location"`
	Fields        map[string]Field //can we remove the fields and article.title directly?
}

//Create draft of a content. parent_id will be -1 in this case
func (c *Content) Create() error {
	db, err := db.Open()
	if err != nil {
		return nil
	}

	//Convert data
	for identifier, field := range c.Fields {
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

func (content Content) CreateLocation() {

}

func (content Content) Store() error {
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
