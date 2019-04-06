//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"dm/model"
	"dm/util"
	"errors"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/queries"
)

type SQLDriver struct{}

//Open opens db with verification
func (m *SQLDriver) Open() (*sql.DB, error) {
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

//Query to fill in contentTyper. Use reference in content parameter
//It fill in with nil if nothing found(no error returned in this case)
func (*RMDB) GetByID(contentType string, id int, content model.ContentTyper) error {
	db, err := new(SQLDriver).Open()
	if err != nil {
		return errors.New("Error when connecting db: " + err.Error())
	}

	contentTypeDef := model.ContentTypeDefinition[contentType]
	tableName := contentTypeDef.TableName
	sql := `SELECT * FROM dm_location l, ` + tableName + ` c
					WHERE l.content_id=c.id
								AND l.content_type= ` + contentType + `
							  AND l.id=?`

	util.Debug("db", sql)
	err = queries.Raw(sql, id).Bind(context.Background(), db, content)

	if err != nil {
		return errors.New("Error when query table: " + tableName + " id: " + strconv.Itoa(id) + " " + err.Error())
	}
	return nil
}
