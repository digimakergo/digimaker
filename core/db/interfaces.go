package db

type DBHandler interface {
	WithContent(query Query, contentType string, option ContentOption) (Query, error)

	BuildQuery(querier Query) (string, []interface{}, error)

	GetColumns(table string) []string
}

var handler DBHandler

func RegisterHandler(dbBuilder DBHandler) {
	handler = dbBuilder
}
