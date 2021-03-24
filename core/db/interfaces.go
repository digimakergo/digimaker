package db

import (
	"context"
	"database/sql"
)

type DBHandler interface {
	WithContent(query Query, contentType string, option ContentOption) (Query, error)

	BuildQuery(querier Query) (string, []interface{}, error)

	GetColumns(table string) []string

	Insert(ctx context.Context, tablename string, values map[string]interface{}, transation ...*sql.Tx) (int, error)
	Update(ctx context.Context, tablename string, values map[string]interface{}, condition Condition, transation ...*sql.Tx) error
	Delete(ctx context.Context, tableName string, condition Condition, transation ...*sql.Tx) error
}

var handler DBHandler

func RegisterHandler(dbBuilder DBHandler) {
	handler = dbBuilder
}

//Insert record, return id created(if it has). It only invokes "handler.Insert"
func Insert(ctx context.Context, tablename string, values map[string]interface{}, transation ...*sql.Tx) (int, error) {
	return handler.Insert(ctx, tablename, values, transation...)
}

//Update record. It only invokes "handler.Update"
func Update(ctx context.Context, tablename string, values map[string]interface{}, condition Condition, transation ...*sql.Tx) error {
	return handler.Update(ctx, tablename, values, condition, transation...)
}

//Delete record. It only invokes "handler.Delete"
func Delete(ctx context.Context, tablename string, condition Condition, transation ...*sql.Tx) error {
	return handler.Delete(ctx, tablename, condition, transation...)
}
