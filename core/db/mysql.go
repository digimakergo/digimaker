package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/util"
)

//cache all columns which is fetched from database.
var tableColumns = map[string][]string{}

//Implement DBHandler
type Handler struct {
}

//get table columns with Cache
func (handler Handler) GetColumns(table string) []string {
	if table == "dm_location" {
		return []string{"name", "content_type", "content_id"}
	} else {
		return []string{"author", "title", "modified"}
	}
}

//
func (handler Handler) WithContent(query Query, contentType string, option ContentFetchOption) (Query, error) {
	def, err := definition.GetDefinition(contentType)
	if err != nil {
		return query, err
	}

	contentQuery := query.Queries[0]
	hasAlias := contentQuery.Alias == "" //without alias

	contentFields := contentQuery.Select

	//add prefix to content fields
	fieldPrefix := ""
	if !hasAlias {
		fieldPrefix = contentQuery.Alias + "."
	}

	if len(contentFields) == 0 {
		fields := handler.GetColumns(def.TableName)
		if contentQuery.Alias == "" {
			contentFields = fields
		} else {
			fieldsWithPrefix := []string{}
			for _, field := range fields {
				fieldStr := fieldPrefix + field + " AS '" + fieldPrefix + field + "'"
				fieldsWithPrefix = append(fieldsWithPrefix, fieldStr)
			}
			contentFields = fieldsWithPrefix
		}
	}

	//add cid to content fields
	contentFields = append(contentFields, fieldPrefix+"id AS '"+fieldPrefix+"cid'")
	contentQuery.Select = contentFields

	contentPrefix := ""
	contentAlias := "c"
	if !hasAlias {
		contentAlias = contentQuery.Alias
		contentPrefix = contentQuery.Alias + "."
	}

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
			renamedColumns = append(renamedColumns, locationAlias+"."+column+" AS '"+contentPrefix+"l-"+column+"'")
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
			Select:   []string{authorAlias + ".name AS " + "'" + contentPrefix + "l-author_name'"},
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
func (handler Handler) BuildQuery(query Query) (string, []interface{}, error) {
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
	fieldStr := strings.Join(fields, ",")

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

	sqlStr := `SELECT ` + fieldStr + ` FROM ` + tableStr + " WHERE " + conditionStr + groupbyStr + sortby + limit

	fmt.Println(sqlStr)

	return sqlStr, values, nil
}

//Get sort by sql based on sortby pattern(eg.[]string{"name asc", "id desc"})
func (r Handler) getSortBy(sortby []string, locationColumns ...[]string) (string, error) {
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
		sortbyStr = " ORDER BY " + strings.Join(sortbyArr, ",")
		sortbyStr = util.StripSQLPhrase(sortbyStr)
	}
	return sortbyStr, nil
}

func (r Handler) getLimit(limit []int) (string, error) {
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
