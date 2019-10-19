package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

// getAllLogGroupNames lists all log group names
func getAllLogGroupNames(cmd *cobra.Command, args []string) error {
	sess, _, logger, err := getSessionAndLogger(cmd, args)
	if err != nil {
		logger.Fatal(err.Error())
	}

	svc := cloudwatchlogs.New(sess)
	out, err := svc.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	formattedOutput, _ := json.Marshal(out.LogGroups)
	logger.Println(string(formattedOutput))

	return nil
}
