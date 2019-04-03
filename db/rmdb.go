//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"dm/model"
	"dm/model/entity"
	"dm/util"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/queries"
)

type SQL struct{}

//Open opens db with verification
func (m *SQL) Open() (*sql.DB, error) {
	dbConfig := util.GetConfigSection("database")
	connString := dbConfig["username"] + ":" + dbConfig["password"] +
		"@" + dbConfig["protocal"] +
		"(" + dbConfig["host"] + ")/" +
		dbConfig["database"] //todo: fix what if there is @ or / in the password?

	db, err := sql.Open(dbConfig["type"], connString)
	if err != nil {
		errorMessage := "Can not open. error: " + err.Error() + " Conneciton string: " + connString
		util.Error(errorMessage)
		return nil, errors.New(errorMessage)
	}

	if db.Ping() != nil {
		util.Error("Can not connect with connection string: " + connString)
		return nil, err
	}

	return db, nil
}

// Implement DBEntitier
type RMDB struct{}

func (*RMDB) GetByID(contentType string, id int) model.ContentTyper {
	db, err := new(SQL).Open()
	if err != nil {
		panic("Error: " + err.Error())
	}

	//todo: make it generic
	var location entity.Location
	queries.Raw(`SELECT * FROM dm_location WHERE id=?`, id).Bind(context.Background(), db, &location)

	var article entity.Article
	queries.Raw(`SELECT * FROM dm_article WHERE id=?`, location.ContentID).Bind(context.Background(), db, &article)

	article.Location = &location

	if err != nil {
		panic("Error 2: " + err.Error())
	}
	return &article
}
