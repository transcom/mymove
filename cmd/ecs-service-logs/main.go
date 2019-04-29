//
// ecs-service-logs is a simple program to print ECS Service logs to stdout.
//
//	Usage of ecs-service-logs:
//	      --aws-profile string               The aws-vault profile
//	      --aws-region string                The AWS Region (default "us-west-2")
//	      --aws-vault-keychain-name string   The aws-vault keychain name
//	      --cluster string                   The cluster name
//	      --environment string               The environment name
//	      --limit int                        The log limit.  If 1 and above, will limit the results returned by each log stream. (default -1)
//	      --service string                   The service name
//	      --status string                    The task status: RUNNING, STOPPED
//	      --verbose                          Print section lines
//
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/99designs/aws-vault/prompt"
	"github.com/99designs/aws-vault/vault"
	"github.com/99designs/keyring"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// The ECS ARN format is changing as explained
// https://aws.amazon.com/ecs/faqs/#Transition_to_new_ARN_and_ID_format
var regexpTaskArnNew = regexp.MustCompile("^arn:aws:ecs:([^:]+?):([^:]+?):task/([^/]+?)/(.+)$")
var regexpTaskArnOld = regexp.MustCompile("^arn:aws:ecs:([^:]+?):([^:]+?):task/(.+)$")

// We need to use regex to extract tasks ids from service events,
// because stopped tasks are only returned by ecs.ListTasks for up to an hour after stopped.
//	- https://docs.aws.amazon.com/sdk-for-go/api/service/ecs/#ECS.ListTasks
var regexpServiceEventStoppedTask = regexp.MustCompile("^[(]service ([0-9a-zA-Z_-]+)[)] has stopped (\\d+) running tasks:\\s+(.+)[.]")
var regexpServiceEventStoppedTaskID = regexp.MustCompile("[(]task ([0-9a-z-]+)[)]")

const (
	flagAWSRegion       string = "aws-region"
	flagAWSProfile      string = "aws-profile"
	flagAWSSessionToken string = "aws-session-token"

	flagAWSVaultKeychainName string = "aws-vault-keychain-name"
	flagAWSVaultProfile      string = "aws-vault"

	flagCluster                string = "cluster"
	flagService                string = "service"
	flagEnvironment            string = "environment"
	flagLogLevel               string = "level"
	flagTaskDefinitionFamily   string = "ecs-task-def-family"
	flagTaskDefinitionRevision string = "ecs-task-def-revision"
	flagGitBranch              string = "git-branch"
	flagGitCommit              string = "git-commit"
	flagPageSize               string = "page-size"
	flagLimit                  string = "limit"
	flagStatus                 string = "status"
	flagVerbose                string = "verbose"

	defaultAWSRegion string = "us-west-2"

	logLevel                  string = "level"
	logTaskDefinitionFamily   string = "ecs_task_def_family"
	logTaskDefinitionRevision string = "ecs_task_def_revision"
	logGitBranch              string = "git_branch"
	logGitCommit              string = "git_commit"
)

var environments = []string{"prod", "staging", "experimental"}
var ecsTaskStatuses = []string{"RUNNING", "STOPPED"}
var logLevels = []string{"debug", "info", "warn", "error", "panic", "fatal"}

func parseTaskID(taskArn string) string {

	// Each match will include a slice of strings starting with
	// (0) the full match, then
	// (1) the region,
	// (2) the account name,
	// (3) (the cluster name if a new arn), and then
	// (4) the task id.

	if matches := regexpTaskArnNew.FindStringSubmatch(taskArn); len(matches) > 0 {
		return matches[4] // returns the task id that was parsed from the new format
	}

	if matches := regexpTaskArnOld.FindStringSubmatch(taskArn); len(matches) > 0 {
		return matches[3] // returns the task id that was parse from the old format
	}

	return ""
}

func parseStoppedTaskEvent(message string) []string {
	if matches := regexpServiceEventStoppedTask.FindStringSubmatch(message); len(matches) > 0 {
		if tasks := regexpServiceEventStoppedTaskID.FindAllStringSubmatch(matches[3], -1); len(tasks) > 0 {
			taskIds := make([]string, 0, len(tasks))
			for _, task := range tasks {
				taskIds = append(taskIds, task[1])
			}
			return taskIds
		}
	}
	return make([]string, 0)
}

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %q", e.Region)
}

type errInvalidEnvironment struct {
	Environment string
}

func (e *errInvalidEnvironment) Error() string {
	return fmt.Sprintf("invalid environment %q, expecting one of %q", e.Environment, environments)
}

type errInvalidTaskStatus struct {
	Status string
}

func (e *errInvalidTaskStatus) Error() string {
	return fmt.Sprintf("invalid status %q, expecting one of %q", e.Status, ecsTaskStatuses)
}

type errInvalidLogLevel struct {
	Level string
}

func (e *errInvalidLogLevel) Error() string {
	return fmt.Sprintf("invalid log level %q, expecting one of %q", e.Level, logLevels)
}

type errInvalidCluster struct {
	Cluster string
}

func (e *errInvalidCluster) Error() string {
	return fmt.Sprintf("invalid cluster %q", e.Cluster)
}

type errInvalidService struct {
	Service string
}

func (e *errInvalidService) Error() string {
	return fmt.Sprintf("invalid service %q", e.Service)
}

func initFlags(flag *pflag.FlagSet) {
	flag.String(flagAWSRegion, defaultAWSRegion, "The AWS Region")
	flag.String(flagAWSProfile, "", "The aws-vault profile")
	flag.String(flagAWSVaultKeychainName, "", "The aws-vault keychain name")
	flag.StringP(flagCluster, "c", "", "The cluster name")
	flag.StringP(flagEnvironment, "e", "", "The environment name")
	flag.StringP(flagService, "s", "", "The service name")
	flag.String(flagStatus, "", "The task status: "+strings.Join(ecsTaskStatuses, ", "))
	flag.StringP(flagLogLevel, "l", "", "The log level: "+strings.Join(logLevels, ", "))
	flag.StringP(flagGitBranch, "b", "", "The git branch")
	flag.String(flagGitCommit, "", "The git commit")
	flag.StringP(flagTaskDefinitionFamily, "f", "", "The ECS task definition family.")
	flag.StringP(flagTaskDefinitionRevision, "r", "", "The ECS task definition revision.")
	flag.IntP(flagPageSize, "p", -1, "The page size or maximum number of log events to return during each API call.  The default is 10,000 log events.")
	flag.IntP(flagLimit, "n", -1, "If 1 or above, the maximum number of log events to print to stdout.")
	flag.BoolP(flagVerbose, "v", false, "Print section lines")
}

func checkConfig(v *viper.Viper) error {
	clusterName := v.GetString("cluster")

	if len(clusterName) == 0 {
		return &errInvalidCluster{Cluster: clusterName}
	}

	if awsVaultProfile := v.GetString(flagAWSVaultProfile); len(awsVaultProfile) > 0 {
		sessionToken := v.GetString(flagAWSSessionToken)
		if len(sessionToken) == 0 {
			return fmt.Errorf("in aws-vault session, but missing aws-session-token")
		}
	} else {
		regions, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.EcsServiceID)
		if !ok {
			return fmt.Errorf("could not find regions for service %q", endpoints.EcsServiceID)
		}

		region := v.GetString(flagAWSRegion)
		if len(region) == 0 {
			return errors.Wrap(&errInvalidRegion{Region: region}, fmt.Sprintf("%q is invalid", flagAWSRegion))
		}

		if _, ok := regions[region]; !ok {
			return errors.Wrap(&errInvalidRegion{Region: region}, fmt.Sprintf("%q is invalid", flagAWSRegion))
		}
	}

	logLevel := strings.ToLower(v.GetString(flagLogLevel))
	if len(logLevel) > 0 {
		valid := false
		for _, str := range logLevels {
			if logLevel == str {
				valid = true
				break
			}
		}

		if !valid {
			return errors.Wrap(&errInvalidLogLevel{Level: logLevel}, fmt.Sprintf("%q is invalid", flagLogLevel))
		}
	}

	status := strings.ToUpper(v.GetString(flagStatus))
	if len(status) > 0 {

		valid := false
		for _, str := range ecsTaskStatuses {
			if status == str {
				valid = true
				break
			}
		}

		if !valid {
			return errors.Wrap(&errInvalidTaskStatus{Status: status}, fmt.Sprintf("%q is invalid", flagStatus))
		}

		if status == "STOPPED" {
			environment := v.GetString(flagEnvironment)
			if len(environment) == 0 {
				return errors.New("when status is set to STOPPED then environment must be set")
			}
			valid := false
			for _, str := range environments {
				if environment == str {
					valid = true
					break
				}
			}
			if !valid {
				return errors.Wrap(&errInvalidEnvironment{Environment: environment}, fmt.Sprintf("%q is invalid", flagEnvironment))
			}

			if serviceName := v.GetString(flagService); len(serviceName) == 0 {
				return &errInvalidService{Service: serviceName}
			}
		}
	}

	return nil
}

// Job is struct linking a task id to a given CloudWatch Log Stream.
type Job struct {
	TaskID        string
	LogGroupName  string
	LogStreamName string
	Limit         int
}

// getAWSCredentials uses aws-vault to return AWS credentials
func getAWSCredentials(keychainName string, keychainProfile string) (*credentials.Credentials, error) {

	// Open the keyring which holds the credentials
	ring, _ := keyring.Open(keyring.Config{
		ServiceName:              "aws-vault",
		AllowedBackends:          []keyring.BackendType{keyring.KeychainBackend},
		KeychainName:             keychainName,
		KeychainTrustApplication: true,
	})

	// Prepare options for the vault before creating the provider
	vConfig, err := vault.LoadConfigFromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load AWS config from environment")
	}
	vOptions := vault.VaultOptions{
		Config:    vConfig,
		MfaPrompt: prompt.Method("terminal"),
	}
	vOptions = vOptions.ApplyDefaults()
	err = vOptions.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to validate aws-vault options")
	}

	// Get a new provider to retrieve the credentials
	provider, err := vault.NewVaultProvider(ring, keychainProfile, vOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create aws-vault provider")
	}
	credVals, err := provider.Retrieve()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve aws credentials from aws-vault")
	}
	return credentials.NewStaticCredentialsFromCreds(credVals), nil
}

func main() {
	root := cobra.Command{
		Use:   "ecs-service-logs [flags]",
		Short: "Show application logs for ECS Service",
		Long:  "Show application logs for ECS Service",
	}

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\necs-service-logs completion > /usr/local/etc/bash_completion.d/ecs-service-logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	showCommand := &cobra.Command{
		Use:   "show [flags]",
		Short: "Show application logs for ECS Service",
		Long:  "Show application logs for ECS Service",
		RunE:  showFunction,
	}
	initFlags(showCommand.Flags())
	root.AddCommand(showCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}

func showFunction(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(args)
	if err != nil {
		return err
	}

	flag := cmd.Flags()

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return err
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	if !v.GetBool(flagVerbose) {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	err = checkConfig(v)
	if err != nil {
		return err
	}

	awsRegion := v.GetString("aws-region")

	awsConfig := &aws.Config{
		Region: aws.String(awsRegion),
	}

	verbose := v.GetBool(flagVerbose)

	if awsVaultProfile := v.GetString(flagAWSVaultProfile); len(awsVaultProfile) == 0 {
		keychainName := v.GetString(flagAWSVaultKeychainName)
		keychainProfile := v.GetString(flagAWSProfile)
		if len(keychainName) > 0 && len(keychainProfile) > 0 {
			creds, err := getAWSCredentials(keychainName, keychainProfile)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Unable to get AWS credentials from the keychain %s and profile %s", keychainName, keychainProfile))
			}
			awsConfig.CredentialsChainVerboseErrors = aws.Bool(verbose)
			awsConfig.Credentials = creds
		}
	}

	sess, err := awssession.NewSession(awsConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS session")
	}

	serviceECS := ecs.New(sess)

	serviceCloudWatchLogs := cloudwatchlogs.New(sess)

	clusterName := v.GetString(flagCluster)
	serviceName := v.GetString(flagService)
	status := strings.ToUpper(v.GetString(flagStatus))
	pageSize := v.GetInt(flagPageSize)
	environment := v.GetString(flagEnvironment)

	jobs := make([]Job, 0)

	if status == "STOPPED" {
		stoppedTaskIds := make([]string, 0)
		describeServicesInput := &ecs.DescribeServicesInput{
			Cluster: aws.String(clusterName),
		}
		if len(serviceName) > 0 {
			describeServicesInput.Services = []*string{aws.String(serviceName)}
		}
		describeServicesOutput, err := serviceECS.DescribeServices(describeServicesInput)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error describing services with cluster name %q", clusterName))
		}
		for _, service := range describeServicesOutput.Services {
			for _, event := range service.Events {
				message := aws.StringValue(event.Message)
				if len(message) > 0 {
					taskIds := parseStoppedTaskEvent(message)
					if len(taskIds) > 0 {
						stoppedTaskIds = append(stoppedTaskIds, taskIds...)
					}
				}
			}
		}

		// If there are no tasks returned from the query then simply exit.
		if len(stoppedTaskIds) == 0 {
			return nil
		}

		for _, taskID := range stoppedTaskIds {

			logGroupName := fmt.Sprintf("ecs-tasks-%s-%s", serviceName, environment)
			logStreamName := fmt.Sprintf("app/%s-%s/%s", serviceName, environment, taskID)

			job := Job{
				TaskID:        taskID,
				LogGroupName:  logGroupName,
				LogStreamName: logStreamName,
				Limit:         -1,
			}
			if pageSize > 0 {
				job.Limit = pageSize
			}
			jobs = append(jobs, job)
		}

	} else {
		taskArns := make([]*string, 0)
		var nextToken *string
		for {

			listTasksInput := &ecs.ListTasksInput{
				Cluster:   aws.String(clusterName),
				NextToken: nextToken,
			}
			if len(serviceName) > 0 {
				listTasksInput.ServiceName = aws.String(serviceName)
			}
			listTasksOutput, err := serviceECS.ListTasks(listTasksInput)
			if err != nil {
				return err
			}
			taskArns = append(taskArns, listTasksOutput.TaskArns...)

			if listTasksOutput.NextToken == nil {
				break
			}
			nextToken = listTasksOutput.NextToken
		}

		// If there are no tasks returned from the query then simply exit.
		if len(taskArns) == 0 {
			return nil
		}

		describeTasksOutput, err := serviceECS.DescribeTasks(&ecs.DescribeTasksInput{
			Cluster: aws.String(clusterName),
			Tasks:   taskArns,
		})
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error describing tasks in cluster %q ", clusterName))
		}

		tasks := describeTasksOutput.Tasks

		taskDefinitionArns := map[string]struct{}{}
		for _, task := range tasks {
			taskDefinitionArns[*task.TaskDefinitionArn] = struct{}{}
		}

		taskDefinitions := map[string]*ecs.TaskDefinition{}
		for taskDefinitionArn := range taskDefinitionArns {
			describeTaskDefinitionOutput, err := serviceECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
				TaskDefinition: aws.String(taskDefinitionArn),
			})
			if err != nil {
				return err
			}
			taskDefinitions[taskDefinitionArn] = describeTaskDefinitionOutput.TaskDefinition
		}

		for _, task := range tasks {

			if status != "" && status != *task.LastStatus {
				continue
			}

			taskID := parseTaskID(*task.TaskArn)

			taskDefinition, ok := taskDefinitions[*task.TaskDefinitionArn]
			if !ok {
				return fmt.Errorf("missing task definition with arn %s for task %s", *task.TaskDefinitionArn, *task.TaskArn)
			}

			for _, containerDefinition := range taskDefinition.ContainerDefinitions {

				containerName := *containerDefinition.Name

				logDriver := *containerDefinition.LogConfiguration.LogDriver
				if logDriver != "awslogs" {
					return fmt.Errorf("found log driver %s, expecting %s", logDriver, "awslogs")
				}

				logConfigurationOptions := containerDefinition.LogConfiguration.Options
				if len(logConfigurationOptions) == 0 {
					return fmt.Errorf("log configuration options is empty")
				}

				logGroupName := logConfigurationOptions["awslogs-group"]
				//logRegion := *logConfigurationOptions["awslogs-region"]
				streamPrefix := *logConfigurationOptions["awslogs-stream-prefix"]

				logStreamName := fmt.Sprintf("%s/%s/%s", streamPrefix, containerName, taskID)

				job := Job{
					TaskID:        taskID,
					LogGroupName:  *logGroupName,
					LogStreamName: logStreamName,
					Limit:         -1,
				}
				if pageSize > 0 {
					job.Limit = pageSize
				}
				jobs = append(jobs, job)
			}

		}
	}

	filters := map[string]string{}

	if gitBranch := v.GetString(flagGitBranch); len(gitBranch) > 0 {
		filters[logGitBranch] = gitBranch
	}

	if gitCommit := v.GetString(flagGitCommit); len(gitCommit) > 0 {
		filters[logGitCommit] = gitCommit
	}

	if family := v.GetString(flagTaskDefinitionFamily); len(family) > 0 {
		filters[logTaskDefinitionFamily] = family
	}

	if revision := v.GetString(flagTaskDefinitionRevision); len(revision) > 0 {
		filters[logTaskDefinitionRevision] = revision
	}

	if level := strings.ToLower(v.GetString(flagLogLevel)); len(level) > 0 {
		filters[logLevel] = level
	}

	filterPattern := ""
	if len(filters) > 0 {
		parts := make([]string, 0, len(filters))
		for k, v := range filters {
			parts = append(parts, fmt.Sprintf("($.%s = %q)", k, v))
		}
		filterPattern = "{" + strings.Join(parts, " && ") + "}"
	}

	limit := v.GetInt(flagLimit)
	count := 0
	for _, job := range jobs {

		if verbose {
			fmt.Println(fmt.Sprintf("Task %s", job.TaskID))
			fmt.Println("-----------------------------------------")
		}

		var nextToken *string
		for {
			filterLogEventsInput := &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName:   aws.String(job.LogGroupName),
				LogStreamNames: []*string{aws.String(job.LogStreamName)},
				NextToken:      nextToken,
			}
			if job.Limit >= 0 {
				if (limit > 0) && ((limit - count) < job.Limit) {
					filterLogEventsInput.Limit = aws.Int64(int64(limit - count))
				} else {
					filterLogEventsInput.Limit = aws.Int64(int64(job.Limit))
				}
			}
			if len(filterPattern) > 0 {
				filterLogEventsInput.FilterPattern = aws.String(filterPattern)
			}
			getLogEventsOutput, err := serviceCloudWatchLogs.FilterLogEvents(filterLogEventsInput)
			if err != nil {
				return errors.Wrap(err, "error retrieving log events")
			}
			for _, event := range getLogEventsOutput.Events {
				fmt.Println(*event.Message)
				count++

				// break the print loop
				if (limit > 0) && (count == limit) {
					break
				}
			}

			// if there are no more events
			if getLogEventsOutput.NextToken == nil {
				break
			}

			// break the pagination loop
			if (limit > 0) && (count == limit) {
				break
			}

			nextToken = getLogEventsOutput.NextToken
		}

		// Break the outer loop
		if (limit > 0) && (count == limit) {
			break
		}

	}

	return nil
}
