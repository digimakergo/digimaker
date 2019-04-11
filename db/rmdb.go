//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"dm/model"
	"dm/query"
	"dm/util"
	"errors"
	"strconv"

	_ "github.com/go-sql-driver/mysql" //todo: move this to loader
	"github.com/volatiletech/sqlboiler/queries"
)

// Implement DBEntitier
type RMDB struct{}

//Query by ID
func (rmdb *RMDB) GetByID(contentType string, id int, content model.ContentTyper) error {
	return rmdb.GetByFields(contentType, query.Cond("id", id), content)
}

//Query to fill in contentTyper. Use reference in content parameter.
//It fill in with nil if nothing found(no error returned in this case)
//  var content model.Article
//  rmdb.GetByFields("article", map[string]interface{}{"id": 12}, content)
//
func (*RMDB) GetByFields(contentType string, condition query.Condition, content model.ContentTyper) error {
	db, err := DB()
	if err != nil {
		return errors.New("Error when connecting db: " + err.Error())
	}

	contentTypeDef := model.ContentTypeDefinition[contentType]
	tableName := contentTypeDef.TableName

	//get condition string for fields
	conditions, values := BuildCondition(condition)
	sql := `SELECT * FROM dm_location location, ` + tableName + ` c
                   WHERE location.content_id=c.id
                         AND location.content_type= '` + contentType + `'
                         AND ` + conditions

	util.Debug("db", sql)
	err = queries.Raw(sql, values...).Bind(context.Background(), db, content)

	if err != nil {
		message := "Error when query table: " + tableName + " " + err.Error() + "sql: " + sql
		util.Error(message)
		return errors.New(message)
	}
	return nil
}

//Fetch multiple enities
func (*RMDB) GetEntities() {

}

func (RMDB) Insert(tablename string, values map[string]interface{}) (int, error) {
	sql := "INSERT INTO " + tablename + " ("
	valuesString := "VALUES("
	var valueParameters []interface{}
	if len(values) > 0 {
		for name, value := range values {
			if name != "id" {
				sql += name + ","
				valuesString += "?,"
				valueParameters = append(valueParameters, value)
			}
		}
		sql = sql[:len(sql)-1]
		valuesString = valuesString[:len(valuesString)-1]
	}
	sql += ")"
	valuesString += ")"
	sql = sql + " " + valuesString
	util.Debug("db", sql)
	db, err := DB()
	if err != nil {
		return 0, err //todo: use new error type
	}
	result, err := db.ExecContext(context.Background(), sql, valueParameters...)
	if err != nil {
		return 0, err //todo: use new error type
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err //todo: use new error type
	}
	util.Debug("db", "Insert results in id: "+strconv.FormatInt(id, 10))

	return int(id), err
}

//Generic update an entity
func (RMDB) Update(tablename string, values map[string]interface{}) error {
	sql := "UPDATE " + tablename + " SET "
	if id, ok := values["id"].(int); ok {
		var valueParameters []interface{}
		for name, value := range values {
			if name != "id" {
				sql += name + "=?,"
				valueParameters = append(valueParameters, value)
			}
		}
		sql = sql[:len(sql)-1]
		sql += " WHERE id=" + strconv.Itoa(id)
		db, err := DB()
		if err != nil {
			return err
		}
		util.Debug("db", sql)
		//todo: use transaction
		result, err := db.ExecContext(context.Background(), sql, valueParameters...)
		resultRows, _ := result.RowsAffected()
		util.Debug("db", "Affected rows:"+strconv.FormatInt(resultRows, 10))
		if err != nil {
			return err //todo: use new error type
		}
	} else {
		return errors.New("There is no id")
	}
	return nil
}

//Update multiple enities
//todo: make a generic condition format/struct
// type Condition struct{}
//
// Cond( "id",GT, 10 )
// Cond( AND( "id",GT, 10, "modified", GT, 12012 ) )
//
func (*RMDB) UpdateAll(name string, condition interface{}) {

}

//Delete a entity
func (*RMDB) Delete(entity model.Entitier) {

}

//Delete based on condition
func (*RMDB) DeleteAll(name string, condition interface{}) {

}

var dbObject = RMDB{}

func DBHanlder() RMDB {
	return dbObject
}
