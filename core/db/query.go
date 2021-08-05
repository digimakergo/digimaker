package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
)

//todo: potencially can support left join
//todo: create error types for database issue(eg. connection) or sync error or other type

type SingleQuery struct {
	Table     string
	Alias     string
	Select    []string
	Condition Condition
	Relation  Condition
}

type ContentOption struct {
	WithAuthor       bool
	WithRelation     []string //relation field
	WithRelationlist []string //relationlist field
}

type Query struct {
	Queries     []SingleQuery
	Groupby     []string //group by
	LimitArr    []int
	LeftQueries []SingleQuery
	SortArr     []string
	AlwaysCount bool
}

//Add a query
func (q *Query) Add(added SingleQuery) {
	queries := append(q.Queries, added)
	q.Queries = queries
}

func (q *Query) AddLeft(query SingleQuery) {
	queries := append(q.LeftQueries, query)
	q.LeftQueries = queries
}

//Join a content or table
func (q *Query) Join(name string, condition Condition, joinCondition Condition) Query {
	tableName, alias, _ := getTableAlias(name)

	oneQuery := SingleQuery{
		Table:     tableName,
		Alias:     alias,
		Select:    handler.GetColumns(tableName),
		Condition: condition,
		Relation:  joinCondition,
	}
	q.Add(oneQuery)
	return *q
}

//sort by
func (q *Query) Sortby(sortby ...string) Query {
	q.SortArr = sortby
	return *q
}

//Limit
func (q *Query) Limit(offset int, limit int) Query {
	q.LimitArr = []int{offset, limit}
	return *q
}

//get table, alias, contenttype, if content type doens't exist, use tablename
//eg. "dm_folder folder" or "dm_folder folder, dm_order order"
func getTableAlias(target string) (string, string, string) {
	nameArr := util.Split(target, " ")
	targetName := nameArr[0]
	targetAlias := ""
	if len(nameArr) > 1 {
		targetAlias = nameArr[1]
	}

	tableName := ""
	def, err := definition.GetDefinition(targetName)
	if err == nil {
		tableName = def.TableName
	} else {
		tableName = targetName
	}

	return tableName, targetAlias, targetName
}

//Build query from a condition. Returns the first target and a Query instance
func CreateQuery(targets string, condition Condition) (string, Query) {
	targetArr := util.Split(targets, ",")
	firstTarget := ""

	queries := []SingleQuery{}
	for i, target := range targetArr {
		table, alias, target := getTableAlias(target)
		sQuery := SingleQuery{
			Table:     table,
			Alias:     alias,
			Condition: condition,
		}
		if i == 0 {
			firstTarget = target
		}
		queries = append(queries, sQuery)
	}

	query := Query{
		Queries:     queries,
		SortArr:     condition.Option.Sortby,
		LimitArr:    condition.Option.LimitArr,
		AlwaysCount: condition.Option.AlwaysCount,
	}

	return firstTarget, query
}

//Bind content with a simple syntax, support innner join
//todo: might be better in another package since it better involves content model?
func BindContent(ctx context.Context, entity interface{}, targets string, condition Condition) (int, error) {
	contentType, query := CreateQuery(targets, condition)
	count, err := BindContentWithQuery(ctx, entity, contentType, query, ContentOption{WithAuthor: true})
	return count, err
}

//Count content separately
func CountContent(ctx context.Context, targets string, condition Condition) (int, error) {
	condition = condition.Limit(0, 0)
	count, err := BindContent(ctx, nil, targets, condition)
	return count, err
}

//Bind entity
func BindEntity(ctx context.Context, entity interface{}, targets string, condition Condition) (int, error) {
	_, query := CreateQuery(targets, condition)
	count, err := BindEntityWithQuery(ctx, entity, query)
	return count, err
}

//Count entities
func Count(targets string, condition Condition, ctx ...context.Context) (int, error) {
	condition = condition.Limit(0, 0) //set to prevent fetching entity
	var currentCtx context.Context
	if len(ctx) > 0 {
		currentCtx = ctx[0]
	} else {
		currentCtx = context.Background()
	}
	result := struct{}{}
	count, err := BindEntity(currentCtx, &result, targets, condition)
	return count, err
}

//Bind content with query
func BindContentWithQuery(ctx context.Context, entity interface{}, contentType string, query Query, option ContentOption) (int, error) {
	query, err := handler.WithContent(query, contentType, option)
	if err != nil {
		return -1, err
	}

	count, err := BindEntityWithQuery(ctx, entity, query)
	return count, err
}

//Bind entity with Query. The simpest is to use BindContent/BindEntity which is using condition directly.
//Condition is syntax(use function to create) easier than query(struct creating)
//If limit is <number>,0 it will ignore entity, and only count
//If limit is <number larger than 0>(eg. 5), <number> it will ignore count unless AlwaysCount is true
func BindEntityWithQuery(ctx context.Context, entity interface{}, query Query) (int, error) {
	sqlStr, values, err := handler.BuildQuery(query, false)
	if err != nil {
		return -1, err
	}
	log.Debug(sqlStr+","+fmt.Sprintln(values), "db", ctx)

	limitArr := query.LimitArr
	count := true
	fetchEntity := true

	if !query.AlwaysCount {
		if len(limitArr) == 0 {
			count = false //if no limit set, then it can be counted from result instead of query count
		} else {
			//if limit x,0, no need to fetch entity, count only
			if limitArr[1] == 0 {
				fetchEntity = false
			}

			//if limit x,10, fetch entity only, not need to count
			if limitArr[0] > 0 {
				count = false
			}
		}
	}

	countResult := -1
	if fetchEntity {
		db, err := DB()
		if err != nil {
			return -1, errors.Wrap(err, "Error when connecting db.")
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
			if err != nil {
				if err != sql.ErrNoRows {
					return -1, errors.Wrap(err, "Error when binding entity")
				}
			}
		}
	}

	if count {
		countResult, err = countWithQuery(ctx, query)
		if err != nil {
			return -1, err
		}
	}

	return countResult, nil
}

//Count only
func countWithQuery(ctx context.Context, query Query) (int, error) {
	query.Groupby = []string{}
	sqlStr, values, err := handler.BuildQuery(query, true)
	if err != nil {
		return -1, err
	}
	log.Debug(sqlStr+","+fmt.Sprintln(values), "db", ctx)

	db, err := DB()
	if err != nil {
		return -1, errors.Wrap(err, "Error when connecting db.")
	}
	entity := struct {
		Count int `boil:"count"`
	}{}

	err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, &entity)
	if err != nil {
		if err != sql.ErrNoRows {
			return -1, errors.Wrap(err, "Count error when binding entity")
		}
	}

	return entity.Count, nil
}
