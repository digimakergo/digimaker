package db

import (
	"context"
	"fmt"

	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
)

//todo: potencially can support left join

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
	SortArr     []string
	AlwaysCount bool
}

//Add a query
func (q *Query) Add(added SingleQuery) {
	queries := append(q.Queries, added)
	q.Queries = queries
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
		Queries:  queries,
		SortArr:  condition.Option.Sortby,
		LimitArr: condition.Option.LimitArr,
	}

	return firstTarget, query
}

//Bind content with a simple syntax, support innner join
func BindContent(ctx context.Context, entity interface{}, targets string, condition Condition) (int, error) {
	contentType, query := CreateQuery(targets, condition)
	count, err := BindContentWithQuery(ctx, entity, contentType, query, ContentOption{WithAuthor: true})
	return count, err
}

//Bind content with query
func BindContentWithQuery(ctx context.Context, entity interface{}, contentType string, query Query, option ContentOption) (int, error) {
	query, err := handler.WithContent(query, contentType, option)
	if err != nil {
		return -1, err
	}

	count, err := BindEntity(ctx, entity, query)
	return count, err
}

//Bind entity with Query
//If limit is x,0(even 0,0) it will ignore entity, and only count
//If limit is 10(>0),y it will ignore count unless AlwaysCount is true
func BindEntity(ctx context.Context, entity interface{}, query Query) (int, error) {
	sqlStr, values, err := handler.BuildQuery(query)
	if err != nil {
		return -1, err
	}
	log.Debug(sqlStr+","+fmt.Sprintln(values), "db", ctx)

	limitArr := query.LimitArr
	count := true
	fetchEntity := true

	if len(limitArr) > 0 {
		//if limit x,0, no need to fetch entity, count only
		if limitArr[1] == 0 {
			fetchEntity = false
		}

		//if limit x,10, fetch entity only, not need to count, unless AlwaysCount is set
		if limitArr[0] > 0 && !query.AlwaysCount {
			count = false
		}
	}

	countResult := -1
	if fetchEntity {
		db, err := DB()
		if err != nil {
			return -1, errors.Wrap(err, "Error when connecting db.")
		}
		err = queries.Raw(sqlStr, values...).Bind(context.Background(), db, entity)
		if err != nil {
			return -1, errors.Wrap(err, "Error when binding entity."+err.Error())
		}
	}

	if count {
		countResult, err = Count(ctx, query)
		if err != nil {
			return -1, err
		}
	}

	return countResult, nil
}

//Count only
func Count(ctx context.Context, query Query) (int, error) {
	sqlStr, values, err := handler.BuildQuery(query)
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
		return -1, errors.Wrap(err, "Count error when binding entity."+err.Error())
	}

	return entity.Count, nil
}
