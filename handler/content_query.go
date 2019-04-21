package handler

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/query"
	"dm/util"

	"github.com/pkg/errors"
)

type contentQuery struct{}

//Fetch one content
func (c contentQuery) One(contentType string, condition query.Condition) (contenttype.ContentTyper, error) {
	content := entity.NewInstance(contentType)

}

//Fetch a list of content
func (c contentQuery) List(contentType string, condition query.Condition) {

}

//Fill all data into content which is a pointer
func (c contentQuery) Fill(contentType string, condition query.Condition, content interface{}) error {
	dbhandler := db.DBHanlder()
	err := dbhandler.GetByFields(contentType, condition, content)
	if err != nil {
		message := "[List]Content Query error"
		util.Error(message, err.Error())
		return errors.Wrap(err, message)
	}
	return nil
}

var Query contentQuery = contentQuery{}
