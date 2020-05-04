//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"

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
func (r *RMDB) GetByFields(contentType string, tableName string, condition Condition, limit []int, sortby []string, content interface{}, count bool) (int, error) {
	db, err := DB()
	if err != nil {
		return -1, errors.Wrap(err, "[RMDB.GetByFields]Error when connecting db.")
	}

	columns := util.GetConfigArr("internal", "location_columns")
	columnsWithPrefix := util.Iterate(columns, func(s string) string {
		return `location.` + s + ` AS "location.` + s + `"`
	})
	locationColumns := strings.Join(columnsWithPrefix, ",")

	//get condition string for fields
	conditions, values := BuildCondition(condition, columns)

	relationQuery := `,CONCAT( '[', GROUP_CONCAT( JSON_OBJECT( 'identifier', relation.identifier,
                                      'to_content_id', relation.to_content_id,
                                      'to_type', relation.to_type,
                                      'from_content_id', relation.from_content_id,
                                      'from_type', relation.from_type,
                                      'from_location', relation.from_location,
                                      'priority', relation.priority,
                                      'uid', relation.uid,
                                      'description',relation.description,
                                      'data' ,relation.data )
                         ORDER BY relation.priority ), ']') AS relations`

	//limit
	limitStr := ""
	if len(limit) > 0 {
		if len(limit) != 2 {
			return -1, errors.New("limit should be array with only 2 int. There are: " + strconv.Itoa(len(limit)))
		}
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}

	//sort by
	sortbyStr, err := r.getSortBy(sortby, columns)
	if err != nil {
		return -1, err
	}

	sqlStr := `SELECT content.*, content.id AS cid, location_user.name AS author_name, ` + locationColumns + relationQuery + `
                   FROM (` + tableName + ` content
                     INNER JOIN dm_location location ON location.content_type = '` + contentType + `' AND location.content_id=content.id)
                     LEFT JOIN dm_relation relation ON content.id=relation.to_content_id AND relation.to_type='` + contentType + `'
										 LEFT JOIN dm_location location_user ON location_user.content_type='user' AND location_user.content_id=content.author
                     WHERE ` + conditions + `
                     GROUP BY location.id, author_name
										 ` + sortbyStr + " " + limitStr

	log.Debug(sqlStr+","+fmt.Sprintln(values), "db")
	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, content)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warning(err.Error(), "GetByFields")
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

//Get sort by sql based on sortby pattern(eg.[]string{"name asc", "id desc"})
func (r *RMDB) getSortBy(sortby []string, locationColumns ...[]string) (string, error) {
	//sort by
	sortbyArr := []string{}
	for _, item := range sortby {
		if strings.TrimSpace(item) != "" {
			itemArr := util.Split(item, " ")
			sortByField := itemArr[0]
			if len(locationColumns) > 0 && util.Contains(locationColumns[0], sortByField) {
				sortByField = "location." + sortByField
			}
			sortByOrder := "ASC"

			if len(itemArr) == 2 {
				sortByOrder = strings.ToUpper(itemArr[1])
				if sortByOrder != "ASC" && sortByOrder != "DESC" {
					return "", errors.New("Invalid sorting string: " + sortByOrder)
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
	return sortbyStr, nil
}

// Count based on condition
func (*RMDB) Count(tablename string, condition Condition) (int, error) {
	conditions, values := BuildCondition(condition)
	sqlStr := "SELECT COUNT(*) AS count FROM " + tablename + " WHERE " + conditions
	log.Debug(sqlStr, "db")
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
func (r *RMDB) GetEntity(tablename string, condition Condition, sortby []string, entity interface{}) error {
	conditions, values := BuildCondition(condition)
	sortbyStr, err := r.getSortBy(sortby)
	if err != nil {
		return err
	}
	sqlStr := "SELECT * FROM " + tablename + " WHERE " + conditions + " " + sortbyStr
	log.Debug(sqlStr, "db")
	db, err := DB()
	if err != nil {
		return errors.Wrap(err, "[RMDB.GetEntity]Error when connecting db.")
	}
	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, entity)
	if err == sql.ErrNoRows {
		log.Warning(err.Error(), "GetEntity")
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
	log.Debug(sqlStr, "db")

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

	log.Debug("Insert results in id: "+strconv.FormatInt(id, 10), "db")

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

	log.Debug(sqlStr, "db")

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
	log.Debug("Updated rows:"+strconv.FormatInt(resultRows, 10), "db")
	return nil
}

//Delete based on condition
func (*RMDB) Delete(tableName string, condition Condition, transation ...*sql.Tx) error {
	conditionString, conditionValues := BuildCondition(condition)
	sqlStr := "DELETE FROM " + tableName + " WHERE " + conditionString

	log.Debug(sqlStr, "db")

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
	log.Debug("Deleted rows:"+strconv.FormatInt(resultRows, 10), "db")
	return nil
}

var dbObject = RMDB{}

func DBHanlder() RMDB {
	return dbObject
}
