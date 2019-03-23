package models

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"dm/models/orm"
	utils "dm/utils"
	"errors"
)

type Contenter interface {
	Publish()

	Create()

	Edit()

	Delete()
}

type Content struct {
	*orm.Location
	Fields map[string]Field //can we remove the fields and article.title directly?
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
			utils.Log("Storing content error, break - id: "+string(content.ID)+", field: "+identifier, "error")
			return errors.New("Can not store content")
		}
	}

	//Store Location
	return nil
}

func (content Content) Publish() {

}
