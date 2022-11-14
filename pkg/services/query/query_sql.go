package query

import (
	"fmt"

	"github.com/transcom/mymove/pkg/assets"
)

func GetSQLQueryByName(fileName string) (string, error) {
	sql, err := assets.Asset("sql_scripts/" + fileName + ".sql")

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
