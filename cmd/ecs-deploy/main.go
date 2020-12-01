package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	awsAccountIDFlag string = "aws-account-id"
	dryRunFlag       string = "dry-run"
	environmentFlag  string = "environment"
)

var environments = []string{"prod", "staging", "experimental", "exp", "stg", "prd"}

type errInvalidAccountID struct {
	AwsAccountID string
}

func (e *errInvalidAccountID) Error() string {
	return fmt.Sprintf("invalid AWS account ID %q", e.AwsAccountID)
}

type errInvalidEnvironment struct {
	Environment string
}

func (e *errInvalidEnvironment) Error() string {
	return fmt.Sprintf("invalid MilMove environment %q, expecting one of %q", e.Environment, environments)
}

func quit(logger *log.Logger, flag *pflag.FlagSet, err error) {
	logger.Println(err.Error())
	if flag != nil {
		logger.Println("Usage:")
		flag.PrintDefaults()
	}
	os.Exit(1)
}

func main() {

	root := cobra.Command{
		Use:   "ecs-deploy [flags]",
		Short: "ecs-deploy tool",
		Long:  "ecs-deploy tool",
	}

	putTargetCommand := &cobra.Command{
		Use:          "put-target",
		Short:        "Put ECS Scheduled Task Target",
		Long:         "Put ECS Scheduled Task Target",
		RunE:         putTargetFunction,
		SilenceUsage: true,
	}
	initPutTargetFlags(putTargetCommand.Flags())
	root.AddCommand(putTargetCommand)

	taskDefCommand := &cobra.Command{
		Use:          "task-def",
		Short:        "Generate and optionally register Task Definitions",
		Long:         "Generate and optionally register Task Definitions",
		RunE:         taskDefFunction,
		SilenceUsage: true,
	}
	initTaskDefFlags(taskDefCommand.Flags())
	root.AddCommand(taskDefCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\necs-deploy completion > /usr/local/etc/bash_completion.d/ecs-deploy",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
