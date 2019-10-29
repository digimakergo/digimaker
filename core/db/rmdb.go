//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"dm/core/util"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" //todo: move this to loader
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
)

// Implement DBEntitier
type RMDB struct {
	Transaction *sql.Tx
}

//Query by ID
func (rmdb *RMDB) GetByID(contentType string, tableName string, id int, content interface{}) error {
	_, err := rmdb.GetByFields(contentType, tableName, Cond("location.id", id), []int{}, []string{}, content, false) //todo: use table name as parameter
	return err
}

//Query to fill in contentTyper. Use reference in content parameter.
//It fill in with nil if nothing found(no error returned in this case)
//  var content contenttype.Article
//  rmdb.GetByFields("article", map[string]interface{}{"id": 12}, {{"name","asc"}} content)
//
func (*RMDB) GetByFields(contentType string, tableName string, condition Condition, limit []int, sortby []string, content interface{}, count bool) (int, error) {
	db, err := DB()
	if err != nil {
		return -1, errors.Wrap(err, "[RMDB.GetByFields]Error when connecting db.")
	}

	//get condition string for fields
	conditions, values := BuildCondition(condition)
	//todo: get columns from either config or entities
	columns := []string{"id", "parent_id", "main_id",
		"hierarchy", "depth", "content_type",
		"content_id", "language",
		"identifier_path",
		"name", "is_hidden", "is_invisible",
		"priority", "uid", "section", "p"}
	locationColumns := ""
	for _, column := range columns {
		locationColumns += `location.` + column + ` AS "location.` + column + `",`
	}
	locationColumns = locationColumns[:len(locationColumns)-1]

	relationQuery := ` ,
                    CONCAT( '[', GROUP_CONCAT( JSON_OBJECT( 'identifier', relation.identifier,
                                      'to_content_id', relation.to_content_id,
                                      'to_type', relation.to_type,
                                      'from_content_id', relation.from_content_id,
                                      'from_type', relation.from_type,
                                      'from_location', relation.from_location,
                                      'priority', relation.priority,
                                      'uid', relation.uid,
                                      'description',relation.description,
                                      'data' ,relation.data )
                         ORDER BY relation.priority ), ']') as relations`

	//limit
	limitStr := ""
	if len(limit) > 0 {
		if len(limit) != 2 {
			return -1, errors.New("limit should be array with only 2 int. There are: " + strconv.Itoa(len(limit)))
		}
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}

	//sort by
	sortbyArr := []string{}
	for _, item := range sortby {
		if strings.TrimSpace(item) != "" {
			itemArr := util.Split(item)
			sortByField := itemArr[0]
			sortByOrder := "ASC"

			if len(itemArr) == 2 {
				sortByOrder = strings.ToUpper(itemArr[1])
				if sortByOrder != "ASC" && sortByOrder != "DESC" {
					return -1, errors.New("Invalid sorting string: " + sortByOrder)
				}
			}
			sortbyItem := sortByField + " " + sortByOrder
			sortbyArr = append(sortbyArr, sortbyItem)
		}
	}

	sortbyStr := ""
	if len(sortbyArr) > 0 {
		sortbyStr = "ORDER BY " + strings.Join(sortbyArr, ",")
		sortbyStr = util.StripSQLPhrase(sortbyStr)
	}

	sqlStr := `SELECT content.*, content.id AS cid, ` + locationColumns + relationQuery + `
                   FROM ( ` + tableName + ` content
                     INNER JOIN dm_location location
                        ON location.content_type = '` + contentType + `' AND location.content_id=content.id )
                     LEFT JOIN dm_relation relation
                        ON content.id=relation.to_content_id AND relation.to_type='` + contentType + `'
                     WHERE ` + conditions + `
                     GROUP BY location.id
										 ` + sortbyStr + " " + limitStr

	util.Debug("db", sqlStr)
	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, content)

	if err != nil {
		if err == sql.ErrNoRows {
			util.Warning("db", err.Error())
		} else {
			message := "[RMDB.GetByFields]Error when query. sql - " + sqlStr
			return -1, errors.Wrap(err, message)
		}
	}

	//count if there is
	countResult := 0
	if count {
		countSqlStr := `SELECT COUNT(*) AS count
									 FROM ( ` + tableName + ` content
										 INNER JOIN dm_location location
												ON location.content_type = '` + contentType + `' AND location.content_id=content.id )
										 WHERE ` + conditions

		fmt.Println(countSqlStr)
		rows, err := queries.Raw(countSqlStr, values...).QueryContext(context.Background(), db)
		if err != nil {
			message := "[RMDB.GetByFields]Error when query count. sql - " + countSqlStr
			return -1, errors.Wrap(err, message)
		}
		rows.Next()
		rows.Scan(&countResult)
		rows.Close()
	}

	return countResult, nil
}

// Count based on condition
func (*RMDB) Count(tablename string, condition Condition) (int, error) {
	conditions, values := BuildCondition(condition)
	sqlStr := "SELECT COUNT(*) AS count FROM " + tablename + " WHERE " + conditions
	util.Debug("db", sqlStr)
	db, err := DB()
	if err != nil {
		return 0, errors.Wrap(err, "[RMDB.Count]Error when connecting db.")
	}
	rows, err := queries.Raw(sqlStr, values...).QueryContext(context.Background(), db)
	if err != nil {
		return 0, errors.Wrap(err, "[RMDB.Count]Error when querying.")
	}
	rows.Next()
	var count int
	rows.Scan(&count)
	rows.Close()
	return count, nil
}

//todo: support limit.
func (*RMDB) GetEntity(tablename string, condition Condition, entity interface{}) error {
	conditions, values := BuildCondition(condition)
	sqlStr := "SELECT * FROM " + tablename + " WHERE " + conditions
	util.Debug("db", sqlStr)
	db, err := DB()
	if err != nil {
		return errors.Wrap(err, "[RMDB.GetEntity]Error when connecting db.")
	}
	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, entity)
	if err == sql.ErrNoRows {
		util.Warning("db", err.Error())
	} else {
		return errors.Wrap(err, "[RMDB.GetEntity]Error when query.")
	}
	return nil
}

//Fetch multiple enities
func (*RMDB) GetMultiEntities(tablenames []string, condition Condition, entity interface{}) {

}

func (RMDB) Insert(tablename string, values map[string]interface{}, transation ...*sql.Tx) (int, error) {
	sqlStr := "INSERT INTO " + tablename + " ("
	valuesString := "VALUES("
	var valueParameters []interface{}
	if len(values) > 0 {
		for name, value := range values {
			if name != "id" {
				sqlStr += name + ","
				valuesString += "?,"
				valueParameters = append(valueParameters, value)
			}
		}
		sqlStr = sqlStr[:len(sqlStr)-1]
		valuesString = valuesString[:len(valuesString)-1]
	}
	sqlStr += ")"
	valuesString += ")"
	sqlStr = sqlStr + " " + valuesString
	util.Debug("db", sqlStr)

	var result sql.Result
	var error error
	//execute using and without using transaction
	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return 0, errors.Wrap(err, "[RBDM.Insert] Error when getting db connection.")
		}
		//todo: create context to isolate queries.
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	//execution error
	if error != nil {
		return 0, errors.Wrap(error, "[RBDM.Insert]Error when executing. sql - "+sqlStr)
	}
	id, err := result.LastInsertId()
	//Get id error
	if err != nil {
		return 0, errors.Wrap(err, "[RBDM.Insert]Error when inserting. sql - "+sqlStr)
	}

	util.Debug("db", "Insert results in id: "+strconv.FormatInt(id, 10))

	return int(id), nil
}

//Generic update an entity
func (RMDB) Update(tablename string, values map[string]interface{}, condition Condition, transation ...*sql.Tx) error {
	sqlStr := "UPDATE " + tablename + " SET "
	var valueParameters []interface{}
	for name, value := range values {
		if name != "id" {
			sqlStr += name + "=?,"
			valueParameters = append(valueParameters, value)
		}
	}
	sqlStr = sqlStr[:len(sqlStr)-1]
	conditionString, conditionValues := BuildCondition(condition)
	valueParameters = append(valueParameters, conditionValues...)
	sqlStr += " WHERE " + conditionString

	util.Debug("db", sqlStr)

	var result sql.Result
	var error error
	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return errors.Wrap(err, "[RBDM.Update] Error when getting db connection.")
		}
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	if error != nil {
		return errors.Wrap(error, "[RMDB.Update]Error when updating. sql - "+sqlStr)
	}
	resultRows, _ := result.RowsAffected()
	util.Debug("db", "Updated rows:"+strconv.FormatInt(resultRows, 10))
	return nil
}

//Delete based on condition
func (*RMDB) Delete(tableName string, condition Condition, transation ...*sql.Tx) error {
	conditionString, conditionValues := BuildCondition(condition)
	sqlStr := "DELETE FROM " + tableName + " WHERE " + conditionString

	util.Debug("db", sqlStr)

	var result sql.Result
	var error error

	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return errors.Wrap(err, "[RBDM.Delete] Error when getting db connection.")
		}
		result, error = db.ExecContext(context.Background(), sqlStr, conditionValues...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, conditionValues...)
	}
	if error != nil {
		return errors.Wrap(error, "[RMDB.Delete]Error when deleting. sql - "+sqlStr)
	}
	resultRows, _ := result.RowsAffected()
	util.Debug("db", "Deleted rows:"+strconv.FormatInt(resultRows, 10))
	return nil
}

var dbObject = RMDB{}

func DBHanlder() RMDB {
	return dbObject
}
