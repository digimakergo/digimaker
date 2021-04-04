//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/volatiletech/sqlboiler/queries"

	_ "github.com/go-sql-driver/mysql" //todo: move this to loader
	"github.com/pkg/errors"
)

type Datamap map[string]interface{}
type DatamapList []Datamap

//cache all columns which is fetched from database.
var tableColumns = map[string][]string{}

// Implement DBEntitier
type MysqlHandler struct {
}

//get table columns with Cache
func (handler MysqlHandler) GetColumns(table string) []string {
	if _, ok := tableColumns[table]; !ok {
		sql := "SELECT COLUMN_NAME AS 'column' FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='" + config_database + "' AND TABLE_NAME=?"
		db, err := DB()

		if err != nil {
			log.Error("No connection when fetching meta columns", "")
		}
		columns := []struct {
			Column string `boil:"column"`
		}{}
		err = queries.Raw(sql, table).Bind(context.Background(), db, &columns)
		if err != nil {
			log.Error(err.Error(), "")
		}
		result := []string{}
		for _, column := range columns {
			result = append(result, column.Column)
		}
		tableColumns[table] = result
	}

	return tableColumns[table]
}

//
func (handler MysqlHandler) WithContent(query Query, contentType string, option ContentOption) (Query, error) {
	def, err := definition.GetDefinition(contentType)
	if err != nil {
		return query, err
	}

	contentQuery := query.Queries[0]
	hasAlias := contentQuery.Alias == "" //without alias

	contentFields := contentQuery.Select

	contentAlias := "c"
	if !hasAlias {
		contentAlias = contentQuery.Alias
	}
	contentPrefix := contentAlias + "."

	//add prefix to content fields
	fieldPrefix := ""
	if !hasAlias {
		fieldPrefix = contentQuery.Alias + "."
	}

	if len(contentFields) == 0 {
		fields := handler.GetColumns(def.TableName)

		fieldsWithPrefix := []string{}
		for _, field := range fields {
			fieldStr := contentPrefix + field + " AS '" + fieldPrefix + field + "'"
			fieldsWithPrefix = append(fieldsWithPrefix, fieldStr)
		}
		contentFields = fieldsWithPrefix
	}

	//add cid to content fields
	contentFields = append(contentFields, contentPrefix+"id AS '"+fieldPrefix+"cid'")
	contentQuery.Select = contentFields

	//add location join if needed
	if def.HasLocation {
		//add location onequery
		locationAlias := "l"
		if !hasAlias {
			locationAlias = contentQuery.Alias + "_" + locationAlias
		}

		//fields
		columns := handler.GetColumns("dm_location")
		renamedColumns := []string{}
		for _, column := range columns {
			renamedColumns = append(renamedColumns, locationAlias+"."+column+" AS '"+fieldPrefix+"location."+column+"'")
		}

		locationQuery := SingleQuery{
			Table:    "dm_location",
			Alias:    locationAlias,
			Select:   renamedColumns,
			Relation: Cond(locationAlias+".content_type ==", "'"+contentType+"'").Cond(locationAlias+".content_id ==", contentAlias+".id"),
		}
		query.Add(locationQuery)
	}

	//add author join if needed
	if option.WithAuthor {
		authorAlias := contentAlias + "_l_author"
		authorQuery := SingleQuery{
			Table:    "dm_location",
			Alias:    authorAlias,
			Select:   []string{authorAlias + ".name AS " + "'" + fieldPrefix + "author_name'"},
			Relation: Cond(contentAlias+".author ==", authorAlias+".content_id").Cond(authorAlias+".content_type ==", "'user'"),
		}
		query.Add(authorQuery)
	}

	//use c for one content alias for conditions&tables
	if hasAlias {
		contentQuery.Alias = "c"
	}
	//change content query
	query.Queries[0] = contentQuery
	return query, nil
}

//Implement BuildQuery
func (handler MysqlHandler) BuildQuery(query Query, count bool) (string, []interface{}, error) {
	tables := []string{}
	fields := []string{}

	condition := EmptyCond()
	for _, item := range query.Queries {
		table := item.Table
		tables = append(tables, table+" "+item.Alias)

		fields = append(fields, item.Select...)
		condition = condition.And(item.Condition).And(item.Relation) //todo: deal relations using on
		//todo: add alias if alias is not empty
	}

	//select
	fieldStr := ""
	if len(fields) > 0 {
		fieldStr = strings.Join(fields, ",")
	} else {
		fieldStr = "*"
	}

	//tables
	tableStr := strings.Join(tables, ",")

	//condition
	conditionStr, values := BuildCondition(condition) //todo: return error error

	//group by
	groupbyStr := ""
	groupby := query.Groupby
	if len(groupby) > 0 {
		groupbyStr = " GROUP BY " + strings.Join(groupby, ",")
	}

	//sort by
	sortby, err := handler.getSortBy(query.SortArr)
	if err != nil {
		return "", nil, err
	}

	//limit
	limit, err := handler.getLimit(query.LimitArr)
	if err != nil {
		return "", nil, err
	}

	sqlStr := ""
	if !count {
		sqlStr = `SELECT ` + fieldStr + ` FROM ` + tableStr + " WHERE " + conditionStr + groupbyStr + sortby + limit
	} else {
		sqlStr = `SELECT COUNT(*) AS count FROM ` + tableStr + " WHERE " + conditionStr + groupbyStr
	}

	return sqlStr, values, nil
}

//Get sort by sql based on sortby pattern(eg.[]string{"name asc", "id desc"})
func (r MysqlHandler) getSortBy(sortby []string) (string, error) {
	//sort by
	sortbyArr := []string{}
	for _, item := range sortby {
		if strings.TrimSpace(item) != "" {
			itemArr := util.Split(item, " ")
			sortByField := itemArr[0]
			locationColumns := definition.LocationColumns
			if len(locationColumns) > 0 && util.Contains(locationColumns, sortByField) {
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
		sortbyStr = " ORDER BY " + strings.Join(sortbyArr, ",")
		sortbyStr = util.StripSQLPhrase(sortbyStr)
	}
	return sortbyStr, nil
}

func (r MysqlHandler) getLimit(limit []int) (string, error) {
	//limit
	limitStr := ""
	if len(limit) > 0 {
		if len(limit) != 2 {
			return "", errors.New("limit should be array with only 2 int. There are: " + strconv.Itoa(len(limit)))
		}
		limitStr = " LIMIT " + strconv.Itoa(limit[0]) + "," + strconv.Itoa(limit[1])
	}
	return limitStr, nil
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
func (MysqlHandler) Delete(ctx context.Context, tableName string, condition Condition, transation ...*sql.Tx) error {
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

func init() {
	RegisterHandler(MysqlHandler{})
}
