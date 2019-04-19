package handler

import (
	"dm/db"
	"dm/query"
	"dm/util"

	"github.com/pkg/errors"
)

type contentQuery struct{}

func (c contentQuery) List(contentType string, condition query.Condition, content interface{}) error {
	dbhanlder := db.DBHanlder()
	err := dbhanlder.GetByFields(contentType, condition, content)
	if err != nil {
		message := "[List]Content Query error"
		util.Error(message, err.Error())
		return errors.Wrap(err, message)
	}
	return nil
}

var Query contentQuery = contentQuery{}
