//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"dm/contenttype"
	"dm/query"
	"dm/util"
	"strconv"

	_ "github.com/go-sql-driver/mysql" //todo: move this to loader
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
)

// Implement DBEntitier
type RMDB struct{}

//Query by ID
func (rmdb *RMDB) GetByID(contentType string, id int, content interface{}) error {
	return rmdb.GetByFields(contentType, query.Cond("location.id", id), content) //todo: use table name as parameter
}

//Query to fill in contentTyper. Use reference in content parameter.
//It fill in with nil if nothing found(no error returned in this case)
//  var content contenttype.Article
//  rmdb.GetByFields("article", map[string]interface{}{"id": 12}, content)
//
func (*RMDB) GetByFields(contentType string, condition query.Condition, content interface{}) error {
	db, err := DB()
	if err != nil {
		return errors.Wrap(err, "[RMDB.GetByFields]Error when connecting db.")
	}

	contentTypeDef := contenttype.GetContentDefinition(contentType)
	tableName := contentTypeDef.TableName

	//get condition string for fields
	conditions, values := BuildCondition(condition)
	//todo: get columns from either config or entities
	columns := []string{"id", "parent_id", "main_id",
		"hierarchy", "content_type",
		"content_id", "language",
		"name", "is_hidden", "is_invisible",
		"priority", "uid", "section", "p"}
	locationColumns := ""
	for _, column := range columns {
		locationColumns += `location.` + column + ` AS "location.` + column + `",`
	}
	locationColumns = locationColumns[:len(locationColumns)-1]

	relationQuery := ` ,
                    GROUP_CONCAT( JSON_OBJECT( 'identifier', relation.identifier,
                                      'to_type', relation.to_type,
                                      'from_location', relation.from_location,
                                      'description',relation.description,
                                      'data' ,relation.data )
                         ORDER BY relation.priority ) as relations`

	sql := `SELECT c.*, ` + locationColumns + relationQuery + `
                   FROM ( ` + tableName + ` c
                     INNER JOIN dm_location location
                        ON location.content_type = '` + contentType + `' AND location.content_id=c.id )
                     LEFT JOIN dm_relation relation
                        ON c.id=relation.to_content_id AND relation.to_type='` + contentType + `'
                     WHERE ` + conditions + `
                     GROUP BY location.id`

	util.Debug("db", sql)
	err = queries.Raw(sql, values...).Bind(context.Background(), db, content)

	if err != nil {
		message := "[RMDB.GetByFields]Error when query. sql - " + sql
		return errors.Wrap(err, message)
	}
	return nil
}

//todo: support limit.
func (*RMDB) GetEnity(tablename string, condition query.Condition, entity interface{}) error {
	conditions, values := BuildCondition(condition)
	sql := "SELECT * FROM " + tablename + " WHERE " + conditions
	util.Debug("db", sql)
	db, err := DB()
	if err != nil {
		return errors.Wrap(err, "[RMDB.GetEntity]Error when connecting db.")
	}
	err = queries.Raw(sql, values...).Bind(context.Background(), db, entity)
	if err != nil {
		return errors.Wrap(err, "[RMDB.GetEntity]Error when query.")
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
		return 0, errors.Wrap(err, "[RBDM.Insert] Error when getting db connection.")
	}
	result, err := db.ExecContext(context.Background(), sql, valueParameters...)
	if err != nil {
		return 0, errors.Wrap(err, "[RBDM.Insert]Error when executing. sql - "+sql)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "[RBDM.Insert]Error when inserting. sql - "+sql)
	}
	util.Debug("db", "Insert results in id: "+strconv.FormatInt(id, 10))

	return int(id), nil
}

//Generic update an entity
func (RMDB) Update(tablename string, values map[string]interface{}, condition query.Condition) error {
	sql := "UPDATE " + tablename + " SET "
	var valueParameters []interface{}
	for name, value := range values {
		if name != "id" {
			sql += name + "=?,"
			valueParameters = append(valueParameters, value)
		}
	}
	sql = sql[:len(sql)-1]
	conditionString, conditionValues := BuildCondition(condition)
	valueParameters = append(valueParameters, conditionValues...)
	sql += " WHERE " + conditionString
	db, err := DB()
	if err != nil {
		return errors.Wrap(err, "[RMDB.Update]Error when getting db connection.")
	}
	util.Debug("db", sql)
	//todo: use transaction
	result, err := db.ExecContext(context.Background(), sql, valueParameters...)
	if err != nil {
		return errors.Wrap(err, "[RMDB.Update]Error when updating. sql - "+sql)
	}
	resultRows, _ := result.RowsAffected()
	util.Debug("db", "Updated rows:"+strconv.FormatInt(resultRows, 10))
	return nil
}

//Delete based on condition
func (*RMDB) Delete(name string, condition interface{}) {

}

var dbObject = RMDB{}

func DBHanlder() RMDB {
	return dbObject
}
