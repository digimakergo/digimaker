//Author xc, Created on 2019-03-27 21:00
//{COPYRIGHTS}

//Package db provides database operation and query building api.
package db

import (
	"context"
	"database/sql"

	"github.com/digimakergo/digimaker/core/util"

	"github.com/pkg/errors"
)

var db *sql.DB

var config_database string

//Get DB connection cached globally
//Note: when using it, the related driver should be imported already
func DB() (*sql.DB, error) {
	if db == nil {
		dbConfig := util.GetConfigSection("database")
		config_database = dbConfig["database"]
		connString := dbConfig["username"] + ":" + dbConfig["password"] +
			"@" + dbConfig["protocal"] +
			"(" + dbConfig["host"] + ")/" +
			dbConfig["database"] //todo: fix what if there is @ or / in the password?

		currentDB, err := sql.Open(dbConfig["type"], connString)
		if err != nil {
			errorMessage := "[DB]Can not open. error: " + err.Error() + " Conneciton string: " + connString
			return nil, errors.Wrap(err, errorMessage)
		}
		db = currentDB
	}
	//ping take extra time. todo: check it in a better way.
	/*
		err := db.Ping()
		if err != nil {
			return nil, errors.Wrap(err, "[DB]Can not ping to connection. ")
		}
	*/
	return db, nil
}

//Create transaction.
//todo: maybe some pararmeters for options
func CreateTx(ctx ...context.Context) (*sql.Tx, error) {
	database, err := DB()
	if err != nil {
		return nil, errors.New("Can't get db connection.")
	}
	var ctxValue context.Context
	if len(ctx) == 0 {
		ctxValue = context.Background()
	} else {
		ctxValue = ctx[0]
	}
	tx, err := database.BeginTx(ctxValue, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, errors.New("Can't get transaction.")
	}
	return tx, nil
}

type DBEntitier interface {
	GetByID(contentType string, id int) interface{}
}
