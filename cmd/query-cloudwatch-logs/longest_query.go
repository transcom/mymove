package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

// TopDurationQueryResults struct is useful for capturing and formatting top duration query results
type TopDurationQueryResults struct {
	Duration  string                          `json:"duration"`
	ID        string                          `json:"id"`
	RawLogs   [][]*cloudwatchlogs.ResultField `json:"RawLogs"`
	TimeStamp string                          `json:"timestamp"`
	Message   string                          `json:"message"`
}

// longestRunningQuery gets top 10 longest running rds queries in production in the last 24 hour time period")
func longestRunningQuery(cmd *cobra.Command, args []string) error {
	sess, v, logger, err := getSessionAndLogger(cmd, args)
	if err != nil {
		logger.Fatal(err.Error())
	}

	query := `fields @timestamp, @message
	| parse "duration: * ms" as duration
	| parse "*@*:[*]" as user, app, id
	| filter duration like ""
	| sort duration desc`
	logGroupName := "/aws/rds/instance/app-prod/postgresql"
	endTime := time.Now().Unix()
	startTime := time.Now().Add(time.Hour * -24).Unix()
	limit := int64(10)
	results := createAndRunQuery(sess, v, &query, &logGroupName, &startTime, &endTime, &limit, logger)

	//loop over each duration and get raw logs so we can figure out what query is impacting the duration value
	for _, record := range results.Results {
		var connectionID, duration, message, timeStamp string

		// get connection id and time stamp so we can query again to get raw logs
		for _, recordField := range record {
			if *recordField.Field == "id" {
				connectionID = *recordField.Value
			}
			if *recordField.Field == "duration" {
				duration = *recordField.Value
			}
			if *recordField.Field == "@message" {
				message = *recordField.Value
			}
			if *recordField.Field == "@timestamp" {
				timeStamp = *recordField.Value
			}

		}
		// Query each record by connection id so we can get raw logs to figure out the query causing the slow down
		query := fmt.Sprintf(`fields @timestamp, @message
		| parse "*@*:[*]" as user, app, id
		| parse "duration: * ms" as duration
		| parse "statement: *" as statement
		| filter id = %s
		| sort @timestamp asc`, connectionID)
		limit := int64(100)

		out := createAndRunQuery(sess, v, &query, &logGroupName, &startTime, &endTime, &limit, logger)
		res := &TopDurationQueryResults{
			Duration:  duration,
			ID:        connectionID,
			RawLogs:   out.Results,
			Message:   message,
			TimeStamp: timeStamp,
		}
		formattedOutput, _ := json.Marshal(res)
		logger.Printf(string(formattedOutput))
	}
	return nil
}
