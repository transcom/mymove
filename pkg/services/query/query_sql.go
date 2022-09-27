package query

import (
	"embed"
	"fmt"
)

//go:embed sql_scripts/*
var queries embed.FS

func GetSQLQueryByName(fileName string) (string, error) {
	sql, err := queries.ReadFile("sql_scripts/" + fileName + ".sql")

	if err != nil {
		return "", err
	}

	query := string(sql)

	if query == "" {
		// return a new error
		return "", fmt.Errorf("The query is empty or does not exist")
	}

	return query, nil
}
