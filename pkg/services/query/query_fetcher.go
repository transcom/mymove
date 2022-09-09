package query

import (
	"fmt"
	"strings"

	"github.com/qustavo/dotsql"
)

// GetQueryStringFromFile accepts a file name and a query name,
// returns SQL query as a string
func GetQueryStringFromFile(filePath string, queryName string) (string, error) {
	isValidPath := strings.HasSuffix(filePath, ".sql")

	if !isValidPath {
		filePath += ".sql"
	}

	dot, err := dotsql.LoadFromFile(filePath)

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

// GetQueryString accepts a query name, and assumes that the query resides
// in a file with the same name.
// Returns SQL query as a string.
func GetQueryString(queryName string) (string, error) {
	filePath := "pkg/sql/" + queryName + ".sql"
	return GetQueryStringFromFile(filePath, queryName)
}
