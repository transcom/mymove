package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	nameFlag       string = "name"
	taskDefARNFlag string = "task-def-arn"
	putTargetFlag  string = "put-target"
)

var names = []string{
	"connect-to-gex-via-sftp",
	"post-file-to-gex",
	"process-edis",
	"process-tpps",
	"save-ghc-fuel-price-data",
	"send-payment-reminder",
}

type errInvalidName struct {
	Name string
}

func (e *errInvalidName) Error() string {
	return fmt.Sprintf("invalid name %q, expecting one of %q", e.Name, names)
}

type errInvalidTaskDefARN struct {
	TaskDefARN string
}

func (e *errInvalidTaskDefARN) Error() string {
	return fmt.Sprintf("invalid AWS Task Def ARN %q", e.TaskDefARN)
}

func initPutTargetFlags(flag *pflag.FlagSet) {

	// AWS Account
	flag.String(awsAccountIDFlag, "", "The AWS Account ID")

	// AWS Flags
	cli.InitAWSFlags(flag)

	// Put Targets Settings
	flag.String(environmentFlag, "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String(nameFlag, "", "The name of the rule")
	flag.String(taskDefARNFlag, "", "The Task Definition ARN")

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Dry Run or Put target
	flag.Bool(dryRunFlag, false, "Execute as a dry-run without modifying AWS.")
	flag.Bool(putTargetFlag, false, "Execute and put target in AWS.")

	// Don't sort flags
	flag.SortFlags = false
}

func checkPutTargetsConfig(v *viper.Viper) error {

	awsAccountID := v.GetString(awsAccountIDFlag)
	if len(awsAccountID) == 0 {
		return fmt.Errorf("%q is invalid: %w", awsAccountIDFlag, &errInvalidAccountID{AwsAccountID: awsAccountID})
	}

	_, err := cli.CheckAWSRegion(v)
	if err != nil {
		return fmt.Errorf("'%q' is invalid: %w", cli.AWSRegionFlag, err)
	}

	environmentName := v.GetString(environmentFlag)
	if len(environmentName) == 0 {
		return fmt.Errorf("%q is invalid: %w", environmentFlag, &errInvalidEnvironment{Environment: environmentName})
	}
	validEnvironment := false
	for _, str := range environments {
		if environmentName == str {
			validEnvironment = true
			break
		}
	}
	if !validEnvironment {
		return fmt.Errorf("%q is invalid: %w", environmentFlag, &errInvalidEnvironment{Environment: environmentName})
	}

	name := v.GetString(nameFlag)
	if len(name) == 0 {
		return fmt.Errorf("%q is invalid: %w", nameFlag, &errInvalidName{Name: name})
	}
	validName := false
	for _, str := range names {
		if name == str {
			validName = true
			break
		}
	}
	if !validName {
		return fmt.Errorf("%q is invalid: %w", nameFlag, &errInvalidName{Name: name})
	}

	taskDefARN := v.GetString(taskDefARNFlag)
	if len(taskDefARN) == 0 || !arn.IsARN(taskDefARN) {
		return fmt.Errorf("%q is invalid: %w", taskDefARNFlag, &errInvalidTaskDefARN{TaskDefARN: taskDefARN})
	}

	return nil
}

func putTargetFunction(cmd *cobra.Command, args []string) error {

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	err := cmd.ParseFlags(args)
	if err != nil {
		return fmt.Errorf("could not parse flags: %w", err)
	}

	flag := cmd.Flags()

	v := viper.New()
	errBindPFlags := v.BindPFlags(flag)
	if errBindPFlags != nil {
		quit(logger, flag, errBindPFlags)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	verbose := cli.LogLevelIsDebug(v)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(io.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Ensure the configuration works against the variables
	checkConfigErr := checkPutTargetsConfig(v)
	if checkConfigErr != nil {
		quit(logger, flag, checkConfigErr)
	}

	// Get the AWS configuration so we can build a session
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(v.GetString(cli.AWSRegionFlag)),
	)
	if err != nil {
		logger.Fatal("error loading default aws config", zap.Error(err))
	}

	// Create the Services
	serviceCloudWatchEvents := cloudwatchevents.NewFromConfig(cfg)

	// Get the current task definition (for rollback)
	taskDefARN := v.GetString(taskDefARNFlag)
	name := v.GetString(nameFlag)
	ruleName := fmt.Sprintf("%s-%s", name, v.GetString(environmentFlag))
	targetsOutput, err := serviceCloudWatchEvents.ListTargetsByRule(
		context.Background(),
		&cloudwatchevents.ListTargetsByRuleInput{
			Rule: aws.String(ruleName),
		})
	if err != nil {
		quit(logger, nil, fmt.Errorf("error retrieving targets for rule: %w", err))
	}

	currentTarget := targetsOutput.Targets[0]

	// Update the task event target with the new task ECS parameters
	putTargetsInput := cloudwatchevents.PutTargetsInput{
		Rule: aws.String(ruleName),
		Targets: []types.Target{
			{
				Id:      currentTarget.Id,
				Arn:     currentTarget.Arn,
				RoleArn: currentTarget.RoleArn,
				EcsParameters: &types.EcsParameters{
					LaunchType:           types.LaunchTypeFargate,
					NetworkConfiguration: currentTarget.EcsParameters.NetworkConfiguration,
					TaskCount:            aws.Int32(1),
					TaskDefinitionArn:    aws.String(taskDefARN),
					PropagateTags:        types.PropagateTagsTaskDefinition,
				},
			},
		},
	}

	if v.GetBool(dryRunFlag) {
		// Format the new task def as JSON for viewing
		jsonErr := json.NewEncoder(logger.Writer()).Encode(putTargetsInput)
		if jsonErr != nil {
			quit(logger, nil, err)
		}
	} else if v.GetBool(putTargetFlag) {
		putTargetsOutput, err := serviceCloudWatchEvents.PutTargets(
			context.Background(),
			&putTargetsInput)
		if err != nil {
			quit(logger, nil, fmt.Errorf("error unable to put new target: %w", err))
		}
		logger.Println(putTargetsOutput)
	} else {
		quit(logger, flag, fmt.Errorf("Please provide either %q or %q flags when running", dryRunFlag, putTargetFlag))
	}

	return nil
}
