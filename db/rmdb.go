package db

import (
	"context"
	"database/sql"
	"dm/model"
	"dm/model/entity"
	"dm/util"
	"errors"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

//Open opens db
func Open1() (*sql.DB, error) {
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

// Implement DBEntitier
type RMDB struct {
}

func (*RMDB) GetByID(contentType string, id int) model.ContentTyper {
	db, err := Open1()
	if err != nil {
		panic("Error: " + err.Error())
	}

	location, err := entity.Locations(Where("id=?", strconv.Itoa(id))).One(context.Background(), db)

	article, err := entity.Articles(Where("id=?", strconv.Itoa(location.ContentID))).One(context.Background(), db)
	article.Location = location

	if err != nil {
		panic("Error 2: " + err.Error())
	}
	return article
}
