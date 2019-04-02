//Author xc, Created on 2019-03-27 21:00
//{COPYRIGHTS}

package db

import (
	"database/sql"
	"dm/model"

	_ "github.com/go-sql-driver/mysql"
)

type DBer interface {
	Open() (*sql.DB, error)
	Query(q string)
	All(q string)
	Execute(q string)
}

type DBEntitier interface {
	GetByID(contentType string, id int) model.ContentTyper
}
