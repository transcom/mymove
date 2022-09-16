package query

import (
	"fmt"
	"os"
	"path"

	"github.com/qustavo/dotsql"
)

// getSQLScriptFilePath returns the path for the SQL scripts directory to be used
// by the QueryFetcher package
func getSQLScriptFilePath(fileName string) (string, error) {
	sqlDir := "pkg/services/query/sql_scripts"

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	queryLocation := path.Join(cwd, sqlDir, fileName+".sql")

	if queryLocation == cwd {
		return "", fmt.Errorf("The SQL scripts directory was not found at %s", sqlDir)
	}

	return queryLocation, nil
}

// GetQueryStringFromFile accepts a query name,
// returns SQL query as a string
func GetQueryString(queryName string) (string, error) {
	queryPath, err := getSQLScriptFilePath(queryName)
	if err != nil {
		panic(err)
	}

	dot, err := dotsql.LoadFromFile(queryPath)

	if err != nil {
		return "", err
	}
	queries := dot.QueryMap()
	query := queries[queryName]
	if query == "" {
		// return a new error
		return "", fmt.Errorf("The query is empty or does not exist. queryPath: %s", queryPath)
	}

	return query, nil
}
