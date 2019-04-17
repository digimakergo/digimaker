package handler

import (
	"dm/db"
	"dm/query"

	"github.com/pkg/errors"
)

type contentQuery struct{}

func (c contentQuery) List(contentType string, condition query.Condition, content interface{}) error {
	dbhanlder := db.DBHanlder()
	err := dbhanlder.GetByFields(contentType, condition, content)
	if err != nil {
		errors.Wrap(err, "Content Query error")
	}
	return nil
}

var Query contentQuery = contentQuery{}
