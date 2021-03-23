package db

type DBHandler interface {
	WithContent(query Query, contentType string, option ContentFetchOption) (Query, error)

	BuildQuery(querier Query) (string, []interface{}, error)

	GetColumns(table string) []string
}
