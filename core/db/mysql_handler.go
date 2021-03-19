//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"

	"github.com/digimakergo/digimaker/core/definition"

	_ "github.com/go-sql-driver/mysql" //todo: move this to loader
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
)

type Datamap map[string]interface{}
type DatamapList []Datamap

// Implement DBEntitier
type MysqlHandler struct {
}

//Query to fill in contentTyper. Use reference in content parameter.
//It fill in with nil if nothing found(no error returned in this case)
//
//todo: possible to have more joins between content/entities(relations or others), or ingegrate with ORM
func (r *MysqlHandler) GetContent(ctx context.Context, content interface{}, contentType string, condition Condition, count bool) (int, error) {
	def, err := definition.GetDefinition(contentType)
	if err != nil {
		return -1, err
	}

	limit := condition.LimitArr
	sortby := condition.Sortby

	tableName := def.TableName

	countResult := 0
	db, err := DB()
	if err != nil {
		return -1, errors.Wrap(err, "[MysqlHandler.GetContent]Error when connecting db.")
	}

	//limit
	limitStr := ""
	if len(limit) > 0 {
		if len(limit) != 2 {
			return -1, errors.New("limit should be array with only 2 int. There are: " + strconv.Itoa(len(limit)))
		}
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}

	relationQuery := ""
	relationJoin := ""
	groupby := ""

	if def.HasRelationlist() {
		relationQuery = r.getRelationQuery()
		relationJoin = ` LEFT JOIN dm_relation relation ON c.id=relation.to_content_id AND relation.to_type='` + contentType + `'`
		if def.HasLocation {
			groupby = ` GROUP BY l.id, author_name`
		} else {
			groupby = ` GROUP BY c.id`
		}
	}

	authorSelect := ",location_user.name AS author_name"
	authorJoin := "LEFT JOIN dm_location location_user ON location_user.content_type='user' AND location_user.content_id=c.author"

	if def.HasLocation {
		columns := util.GetInternalSettings("location_columns")
		columnsWithPrefix := util.Iterate(columns, func(s string) string {
			return `l.` + s + ` AS "location.` + s + `"`
		})
		locationColumns := "," + strings.Join(columnsWithPrefix, ",")

		//get condition string for fields
		conditionStr, values := BuildCondition(condition, columns)
		where := ""
		if conditionStr != "" {
			where = " WHERE " + conditionStr
		}

		//sort by
		sortbyStr, err := r.getSortBy(sortby, columns)
		if err != nil {
			return -1, err
		}

		sqlStr := `SELECT c.*, c.id AS cid` + authorSelect + locationColumns + relationQuery + `
	                   FROM (` + tableName + ` c INNER JOIN dm_location l ON l.content_type = '` + contentType + `' AND l.content_id=c.id)
	                     ` + relationJoin + authorJoin + where + groupby + sortbyStr + limitStr

		log.Debug(sqlStr+","+fmt.Sprintln(values), "db", ctx)
		err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, content)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Debug(err.Error(), "GetContent", ctx)
			} else {
				message := "[MysqlHandler.GetContent]Error when query. sql - " + sqlStr
				return -1, errors.Wrap(err, message)
			}
		}

		//count if there is
		if count {
			countSqlStr := `SELECT COUNT(*) AS count
										 FROM ( ` + tableName + ` c
											 INNER JOIN dm_location l ON l.content_type = '` + contentType + `' AND l.content_id=c.id )
											 ` + where

			rows, err := queries.Raw(countSqlStr, values...).QueryContext(context.Background(), db)
			if err != nil {
				message := "[MysqlHandler.GetContent]Error when query count. sql - " + countSqlStr
				return -1, errors.Wrap(err, message)
			}
			rows.Next()
			rows.Scan(&countResult)
			rows.Close()
		}
	} else {
		//Get non-location content

		//get condition string for fields
		conditionStr, values := BuildCondition(condition)
		where := ""
		if conditionStr != "" {
			where = " WHERE " + conditionStr
		}

		//sort by
		sortbyStr, err := r.getSortBy(sortby)
		if err != nil {
			return -1, err
		}

		sqlStr := `SELECT c.*, c.id as cid, '` + contentType + `' as content_type` + authorSelect + relationQuery + `
										 FROM (` + tableName + ` c INNER JOIN dm_location location ON c.location_id = location.id )
										 ` + relationJoin + authorJoin + where + groupby + sortbyStr + limitStr

		log.Debug(sqlStr+","+fmt.Sprintln(values), "db", ctx)
		err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, content)

		if err != nil {
			if err == sql.ErrNoRows {
				// log.Warning(err.Error(), "GetByFields")
			} else {
				message := "[MysqlHandler.GetEntityContent]Error when query. sql - " + sqlStr
				return -1, errors.Wrap(err, message)
			}
		}

		//count if there is
		if count {
			countSqlStr := `SELECT COUNT(*) AS count FROM ` + tableName + ` c INNER JOIN dm_location location ON c.location_id = location.id ` + where

			rows, err := queries.Raw(countSqlStr, values...).QueryContext(context.Background(), db)
			if err != nil {
				message := "[MysqlHandler.GetEntityContent]Error when query count. sql - " + countSqlStr
				return -1, errors.Wrap(err, message)
			}
			rows.Next()
			rows.Scan(&countResult)
			rows.Close()
		}
	}

	return countResult, nil
}

func (r *MysqlHandler) getRelationQuery() string {
	relationQuery := `, JSON_ARRAYAGG( JSON_OBJECT( 'identifier', relation.identifier,
                                      'to_content_id', relation.to_content_id,
                                      'to_type', relation.to_type,
                                      'from_content_id', relation.from_content_id,
                                      'from_type', relation.from_type,
                                      'from_location', relation.from_location,
                                      'priority', relation.priority,
                                      'uid', relation.uid,
                                      'description',relation.description,
                                      'data' ,relation.data ) ) AS relations`
	return relationQuery
}

//Get sort by sql based on sortby pattern(eg.[]string{"name asc", "id desc"})
func (r *MysqlHandler) getSortBy(sortby []string, locationColumns ...[]string) (string, error) {
	//sort by
	sortbyArr := []string{}
	for _, item := range sortby {
		if strings.TrimSpace(item) != "" {
			itemArr := util.Split(item, " ")
			sortByField := itemArr[0]
			if len(locationColumns) > 0 && util.Contains(locationColumns[0], sortByField) {
				sortByField = "l." + sortByField
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
func (*MysqlHandler) Count(tablename string, condition Condition) (int, error) {
	conditions, values := BuildCondition(condition)
	sqlStr := "SELECT COUNT(*) AS count FROM " + tablename + " WHERE " + conditions
	log.Debug(sqlStr, "db")
	db, err := DB()
	if err != nil {
		return 0, errors.Wrap(err, "[MysqlHandler.Count]Error when connecting db.")
	}
	rows, err := queries.Raw(sqlStr, values...).QueryContext(context.Background(), db)
	if err != nil {
		return 0, errors.Wrap(err, "[MysqlHandler.Count]Error when querying.")
	}
	rows.Next()
	var count int
	rows.Scan(&count)
	rows.Close()
	return count, nil
}

//todo: support limit.
func (r *MysqlHandler) GetEntity(ctx context.Context, entity interface{}, tablename string, condition Condition, count bool) (int, error) {
	conditions, values := BuildCondition(condition)
	sortby := condition.Sortby
	sortbyStr, err := r.getSortBy(sortby)
	if err != nil {
		return 0, err
	}

	limit := condition.LimitArr
	limitStr := ""
	if limit != nil && len(limit) == 2 {
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}
	sqlStr := "SELECT * FROM " + tablename + " WHERE " + conditions + " " + sortbyStr + limitStr
	log.Debug(sqlStr, "db", ctx)
	db, err := DB()
	if err != nil {
		return 0, errors.Wrap(err, "[MysqlHandler.GetEntity]Error when connecting db.")
	}

	//Fill in with DatamapList or other entity
	if entityList, ok := entity.(*DatamapList); ok {
		rows, rowError := queries.Raw(sqlStr, values...).QueryContext(context.Background(), db)
		cols, _ := rows.Columns()
		err = rowError
		defer rows.Close()
		list := DatamapList{}

		for rows.Next() {
			//scan to columnpointers
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				return 0, err
			}

			//set to datamap
			datamap := Datamap{}
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				v := *val
				switch v.(type) {
				case []byte:
					datamap[colName] = string(v.([]byte))
				default:
					datamap[colName] = v
				}
			}
			list = append(list, datamap)
		}

		*entityList = list
	} else {
		err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, entity)
	}
	if err == sql.ErrNoRows {
		// log.Warning(err.Error(), "GetEntity")
	} else {
		return 0, errors.Wrap(err, "[MysqlHandler.GetEntity]Error when query.")
	}

	//count
	countResult := 0
	if count {
		countResult, err = r.Count(tablename, condition)
		if err != nil {
			return 0, err
		}
	}
	return countResult, nil
}

func (MysqlHandler) Insert(ctx context.Context, tablename string, values map[string]interface{}, transation ...*sql.Tx) (int, error) {
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
	log.Debug(sqlStr, "db", ctx)

	var result sql.Result
	var error error
	//execute using and without using transaction
	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return 0, errors.Wrap(err, "MysqlHandler.Insert] Error when getting db connection.")
		}
		//todo: create context to isolate queries.
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	//execution error
	if error != nil {
		return 0, errors.Wrap(error, "MysqlHandler.Insert]Error when executing. sql - "+sqlStr)
	}
	id, err := result.LastInsertId()
	//Get id error
	if err != nil {
		return 0, errors.Wrap(err, "MysqlHandler.Insert]Error when inserting. sql - "+sqlStr)
	}

	log.Debug("Insert results in id: "+strconv.FormatInt(id, 10), "db", ctx)

	return int(id), nil
}

//Generic update an entity
func (MysqlHandler) Update(ctx context.Context, tablename string, values map[string]interface{}, condition Condition, transation ...*sql.Tx) error {
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

	log.Debug(sqlStr, "db", ctx)

	var result sql.Result
	var error error
	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return errors.Wrap(err, "MysqlHandler.Update] Error when getting db connection.")
		}
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	if error != nil {
		return errors.Wrap(error, "[MysqlHandler.Update]Error when updating. sql - "+sqlStr)
	}
	resultRows, _ := result.RowsAffected()
	log.Debug("Updated rows:"+strconv.FormatInt(resultRows, 10), "db", ctx)
	return nil
}

//Delete based on condition
func (*MysqlHandler) Delete(ctx context.Context, tableName string, condition Condition, transation ...*sql.Tx) error {
	conditionString, conditionValues := BuildCondition(condition)
	sqlStr := "DELETE FROM " + tableName + " WHERE " + conditionString

	log.Debug(sqlStr, "db", ctx)

	var result sql.Result
	var error error

	if len(transation) == 0 {
		db, err := DB()
		if err != nil {
			return errors.Wrap(err, "MysqlHandler.Delete] Error when getting db connection.")
		}
		result, error = db.ExecContext(context.Background(), sqlStr, conditionValues...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, conditionValues...)
	}
	if error != nil {
		return errors.Wrap(error, "[MysqlHandler.Delete]Error when deleting. sql - "+sqlStr)
	}

	//todo: return this because it's useful for error handling
	resultRows, _ := result.RowsAffected()
	log.Debug("Deleted rows:"+strconv.FormatInt(resultRows, 10), "db", ctx)
	return nil
}

//todo: consider to remove DBHanlder()?
var dbObject = MysqlHandler{}

func DBHanlder() MysqlHandler {
	return dbObject
}
