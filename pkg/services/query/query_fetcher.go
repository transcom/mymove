package query

import (
	"fmt"
	"os"
	"path"

	"github.com/qustavo/dotsql"
)

// GetQueryStringFromFile accepts a query name,
// returns SQL query as a string
func GetQueryString(queryName string) (string, error) {
	sqlDir := "services/query/sql_scripts"
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	queryPath := path.Join(cwd, "..", "..", sqlDir, queryName+".sql")
	dot, err := dotsql.LoadFromFile(queryPath)

	if err != nil {
		return "", err
	}
	queries := dot.QueryMap()
	query := queries[queryName]
	if query == "" {
		// return a new error
		return "", fmt.Errorf("query is empty or does not exist")
	}

	return query, nil
}
