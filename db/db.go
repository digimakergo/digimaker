//Author xc, Created on 2019-03-27 21:00
//{COPYRIGHTS}

package db

import (
	"database/sql"
	"dm/util"

	"github.com/pkg/errors"
)

var db *sql.DB

//Get DB connection cached globally
//Note: when using it, the related driver should be imported already
func DB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	} else {
		dbConfig := util.GetConfigSection("database")
		connString := dbConfig["username"] + ":" + dbConfig["password"] +
			"@" + dbConfig["protocal"] +
			"(" + dbConfig["host"] + ")/" +
			dbConfig["database"] //todo: fix what if there is @ or / in the password?

		db, err := sql.Open(dbConfig["type"], connString)
		if err != nil {
			errorMessage := "[DB]Can not open. error: " + err.Error() + " Conneciton string: " + connString
			return nil, errors.Wrap(err, errorMessage)
		}

		err = db.Ping()
		if err != nil {
			return nil, errors.Wrap(err, "[DB]Can not ping to connect with connection string: "+connString)
		}
		return db, nil
	}
}

type DBer interface {
	Open() (*sql.DB, error)
	Query(q string)
	All(q string)
	Execute(q string)
}

type DBEntitier interface {
	GetByID(contentType string, id int) interface{}
}
