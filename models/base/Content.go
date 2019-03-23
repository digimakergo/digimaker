package base

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"dmcaf/models/orm"
	utils "dmcaf/utils"
	"errors"
)

type Publisher interface {
	Publish()

	Create()

	Edit()

	Delete()
}

type Content struct {
	*orm.Location
	Fields map[string]Datatype //can we remove the fields and article.title directly?
}

func (content Content) CreateLocation() {

}

func (content Content) Store() error {
	//Store fields
	fields := content.Fields
	for identifier, field := range fields {
		err := field.Store()
		if err != nil {
			//log it and return higher
			utils.Log("Storing content error, break - id: "+string(content.ID)+", field: "+identifier, "error")
			return errors.New("Can not store content")
		}
	}

	//Store Location
	return nil
}
