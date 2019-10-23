package main

import (
	"time"

	"github.com/spf13/cobra"
)

// recentQueriesRDS gets last 100 query logs that were run against prod rds
func recentQueriesRDS(cmd *cobra.Command, args []string) error {
	sess, v, logger, err := getSessionAndLogger(cmd, args)
	if err != nil {
		logger.Fatal(err.Error())
	}
	query := `fields @message
	| parse "statement: *" as statement
	| filter statement like ""
	| filter statement not like "SELECT 1;"
	| filter statement NOT like "COMMIT"
	| filter statement NOT LIKE "BEGIN READ WRITE"`
	logGroupName := "/aws/rds/instance/app-prod/postgresql"
	endTime := time.Now().Unix()
	startTime := time.Now().Add(time.Hour * -24).Unix()
	limit := v.GetInt64(LimitFlag)
	results := createAndRunQuery(sess, v, &query, &logGroupName, &startTime, &endTime, &limit, logger)
	printQueryResults(results, logger)

	return nil
}
