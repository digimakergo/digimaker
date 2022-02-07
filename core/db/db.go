//Author xc, Created on 2019-03-27 21:00
//{COPYRIGHTS}

//Package db provides database operation and query building api.
package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/spf13/viper"
)

var db *sql.DB

//Get DB connection cached globally
//Note: when using it, the related driver should be imported already
func DB() (*sql.DB, error) {
	if db == nil {

		host := viper.GetString("database.host")
		database := viper.GetString("database.database")
		username := viper.GetString("database.username")
		password := viper.GetString("database.password")
		protocal := viper.GetString("database.protocal")

		dbType := viper.GetString("database.type")

		connString := username + ":" + password +
			"@" + protocal +
			"(" + host + ")/" +
			database + "?parseTime=true" //todo: fix what if there is @ or / in the password?

		currentDB, err := sql.Open(dbType, connString)
		if err != nil {
			errorMessage := "[DB]Can not open. error: " + err.Error() + " Conneciton string: " + connString
			return nil, errors.New(errorMessage)
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
