package db

import (
	"database/sql"
	util "dm/util"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//Open opens db
func Open() (*sql.DB, error) {
	dbConfig, _ := util.GetConfigSectin("database")
	connString := dbConfig["username"] + ":" + dbConfig["password"] +
		"@" + dbConfig["protocal"] + "(" + dbConfig["host"] + ")/" + dbConfig["database"]

	db, err := sql.Open(dbConfig["type"], connString)
	if err != nil {
		errorMessage := "Can not connect. error" + err.Error() + " Conneciton string: " + connString
		util.LogError(errorMessage)
		return nil, errors.New(errorMessage)
	}

	if db.Ping() != nil {
		fmt.Printf("can not ping")
		return nil, err
	}

	return db, nil
}
