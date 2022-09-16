package query

import (
	"fmt"
	"path/filepath"

	"github.com/qustavo/dotsql"
)

// getSQLScriptFilePath returns the path for the SQL scripts directory to be used
// by the QueryFetcher package
func getSQLScriptFilePath(fileName string) (string, error) {
	sqlDir := "pkg/services/query/sql_scripts/"

	// Let's try and get the absolute path based on the construction of the SQL
	// directory as we know it at the root of the project, then append the file
	// name and the extension at the end. The Abs() function will return the
	// correct path here since it joins the current working directory with the
	// path that we give it to see if it is there. But this doesn't guarantee
	// that the file is actually there or if the path is unique according to
	// the documentation in the path/filepath package.
	queryLocation, err := filepath.Abs(sqlDir + fileName + ".sql")

	if err != nil {
		return "", fmt.Errorf("There was an error getting the absolute path for %s. %s", fileName, err)
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
