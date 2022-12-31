//Author xc, Created on 2019-04-01 22:00
//{COPYRIGHTS}
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/queries"

	_ "github.com/go-sql-driver/mysql" //todo: move this to loader
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
		sql := "SELECT COLUMN_NAME AS 'column' FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='" + getDBName() + "' AND TABLE_NAME=?"
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

	//add _contenttype to metadata
	contentFields = append(contentFields, "'"+contentType+"'AS _contenttype")

	contentQuery.Select = contentFields

	//Add relation join if there is
	if def.HasRelationlist() {
		relationColumns := []string{"JSON_ARRAYAGG( JSON_OBJECT( 'id', relation.id, 'identifier', relation.identifier, 'to_content_id', relation.to_content_id,'to_type', relation.to_type,'from_content_id', relation.from_content_id,'from_type', relation.from_type,'from_location', relation.from_location,'priority', relation.priority, 'uid', relation.uid, 'description',relation.description, 'data' ,relation.data ) ) AS _relations"}
		relationQuery := SingleQuery{
			Table:     "dm_relation",
			Alias:     "relation",
			Select:    relationColumns,
			Condition: Cond("relation.to_content_id ==", "c.id").Cond("relation.to_type =", contentType),
		}
		/*todo: fix mysql ONLY_FULL_GROUP_BY issue,
		currently it's setting  SET GLOBAL sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));
		*/
		query.AddLeft(relationQuery)
	}

	query.Groupby = []string{"c.id"}

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
		userDef, _ := definition.GetDefinition("user")
		authorQuery := SingleQuery{
			Table:    userDef.TableName,
			Alias:    authorAlias,
			Select:   []string{authorAlias + "._name AS " + "'" + fieldPrefix + "_author_name'"},
			Relation: Cond(contentAlias+"._author ==", authorAlias+".id"),
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

	//tables
	tableStr := strings.Join(tables, " INNER JOIN ")

	//condition
	conditionStr, values := BuildCondition(condition) //todo: return error error

	//left joins
	leftJoinStr := ""
	if !count {
		for _, item := range query.LeftQueries {
			leftConditionStr, leftValues := BuildCondition(item.Condition)
			fields = append(fields, item.Select...)
			leftJoinStr = leftJoinStr + " LEFT JOIN " + item.Table + " " + item.Alias + " ON " + leftConditionStr
			values = append(leftValues, values...) //put in front since left join in before where
		}
	}

	//select
	fieldStr := ""
	if len(fields) > 0 {
		fieldStr = strings.Join(fields, ",")
	} else {
		fieldStr = "*"
	}

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
	whereStr := ""
	if conditionStr != "" {
		whereStr = " WHERE " + conditionStr
	}
	if !count {
		sqlStr = `SELECT ` + fieldStr + ` FROM ` + tableStr + leftJoinStr + whereStr + groupbyStr + sortby + limit
	} else {
		sqlStr = `SELECT COUNT(*) AS count FROM ` + tableStr + whereStr + groupbyStr
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
		limitStr = " LIMIT " + strconv.Itoa(limit[0])
		if limit[1] == -1 {
			defaultLimit := 10000000
			limitStr += "," + strconv.Itoa(defaultLimit)
		} else {
			limitStr += "," + strconv.Itoa(limit[1])
		}
	}
	return limitStr, nil
}

func fixEmptyType(value interface{}) interface{} {
	//fix golang time doesn't support datetime null in db
	if v, ok := value.(time.Time); ok {
		emptyTime := time.Time{}
		if v == emptyTime {
			value, _ = time.Parse("2006-01-02 15:04", "1000-01-01 00:00")
		}
	}
	return value
}

func (MysqlHandler) Insert(ctx context.Context, tablename string, values map[string]interface{}, transation ...*sql.Tx) (int, error) {
	sqlStr := "INSERT INTO " + tablename + " ("
	valuesString := "VALUES("
	var valueParameters []interface{}
	if len(values) > 0 {
		for name, value := range values {
			value = fixEmptyType(value)
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
			return 0, fmt.Errorf("MysqlHandler.Insert] Error when getting db connection: %w", err)
		}
		//todo: create context to isolate queries.
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	//execution error
	if error != nil {
		return 0, fmt.Errorf("MysqlHandler.Insert]Error when executing. sql - %v: %w ", sqlStr, error)
	}
	id, err := result.LastInsertId()
	//Get id error
	if err != nil {
		return 0, fmt.Errorf("MysqlHandler.Insert]Error when inserting. sql - %v: %w ", sqlStr, err)
	}

	log.Debug("Insert results in id: "+strconv.FormatInt(id, 10), "db", ctx)

	return int(id), nil
}

//Generic update an entity
func (MysqlHandler) Update(ctx context.Context, tablename string, values map[string]interface{}, condition Condition, transation ...*sql.Tx) error {
	sqlStr := "UPDATE " + tablename + " SET "
	var valueParameters []interface{}
	for name, value := range values {
		value = fixEmptyType(value)
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
			return fmt.Errorf("MysqlHandler.Update] Error when getting db connection: %w", err)
		}
		result, error = db.ExecContext(context.Background(), sqlStr, valueParameters...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, valueParameters...)
	}
	if error != nil {
		return fmt.Errorf("[MysqlHandler.Update]Error when updating. sql - %v: %w ", sqlStr, error)
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
			return fmt.Errorf("MysqlHandler.Delete] Error when getting db connection: %w", err)
		}
		result, error = db.ExecContext(context.Background(), sqlStr, conditionValues...)
	} else {
		result, error = transation[0].ExecContext(context.Background(), sqlStr, conditionValues...)
	}
	if error != nil {
		return fmt.Errorf("[MysqlHandler.Delete]Error when deleting. sql - %v: %w", sqlStr, error)
	}

	//todo: return this because it's useful for error handling
	resultRows, _ := result.RowsAffected()
	log.Debug("Deleted rows:"+strconv.FormatInt(resultRows, 10), "db", ctx)
	return nil
}

var config_database string

func getDBName() string {
	if config_database == "" {
		config_database = viper.GetString("database.database")
	}
	return config_database
}

func init() {
	RegisterHandler(MysqlHandler{})
}
