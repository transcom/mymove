package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// runQueryWithParams runs query with params
func runQueryWithParams(cmd *cobra.Command, args []string) error {
	sess, v, logger, err := getSessionAndLogger(cmd, args)
	if err != nil {
		logger.Fatal(err.Error())
	}

	utcLoc, _ := time.LoadLocation("UTC")
	iQuery := &InsightQuery{Query: v.GetString(QueryFlag)}
	iQuery.LogGroupName = v.GetString(LogGroupNameFlag)
	iQuery.Limit = v.GetInt64(LimitFlag)

	startTime, err := time.ParseInLocation(time.RFC3339, v.GetString(StartTimeFlag), utcLoc)
	if err != nil {
		logger.Fatal("unable to parse start time", err.Error())
		os.Exit(1)
	}
	iQuery.StartTime = startTime.Unix()

	endTime, err := time.ParseInLocation(time.RFC3339, v.GetString(EndTimeFlag), utcLoc)
	if err != nil {
		logger.Fatal("unable to parse end time", err.Error())
		os.Exit(1)
	}
	iQuery.EndTime = endTime.Unix()

	noop := v.GetBool(RawQueryFlag)

	// don't execute the query
	if noop {
		formattedOutput, err := json.Marshal(iQuery)
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Print(string(formattedOutput))
		os.Exit(0)
	}
	results := createAndRunQuery(sess, v, &iQuery.Query, &iQuery.LogGroupName, &iQuery.StartTime, &iQuery.EndTime, &iQuery.Limit, logger)
	printQueryResults(results, logger)

	return nil
}
