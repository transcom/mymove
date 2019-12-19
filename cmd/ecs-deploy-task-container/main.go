package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

var services = []string{"app"}
var environments = []string{"prod", "staging", "experimental"}

// Commands should be the name of the binary found in the /bin directory in the container
var commands = []string{"milmove-tasks save-fuel-price-data", "milmove-tasks send-post-move-survey"}

type errInvalidAccountID struct {
	AwsAccountID string
}

func (e *errInvalidAccountID) Error() string {
	return fmt.Sprintf("invalid AWS account ID %q", e.AwsAccountID)
}

type errInvalidService struct {
	Service string
}

func (e *errInvalidService) Error() string {
	return fmt.Sprintf("invalid AWS ECS service %q, expecting one of %q", e.Service, services)
}

type errInvalidEnvironment struct {
	Environment string
}

func (e *errInvalidEnvironment) Error() string {
	return fmt.Sprintf("invalid MilMove environment %q, expecting one of %q", e.Environment, environments)
}

type errinvalidRepositoryName struct {
	RepositoryName string
}

func (e *errinvalidRepositoryName) Error() string {
	return fmt.Sprintf("invalid AWS ECR respository name %q", e.RepositoryName)
}

type errinvalidImageTag struct {
	ImageTag string
}

func (e *errinvalidImageTag) Error() string {
	return fmt.Sprintf("invalid AWS ECR image tag %q", e.ImageTag)
}

type errInvalidCommand struct {
	Command string
}

func (e *errInvalidCommand) Error() string {
	return fmt.Sprintf("invalid command in the /bin folder %q", e.Command)
}

type errInvalidFile struct {
	File string
}

func (e *errInvalidFile) Error() string {
	return fmt.Sprintf("invalid file path %q", e.File)
}

const (
	awsAccountIDFlag       string = "aws-account-id"
	chamberBinaryFlag      string = "chamber-binary"
	chamberRetriesFlag     string = "chamber-retries"
	chamberKMSKeyAliasFlag string = "chamber-kms-key-alias"
	chamberUsePathsFlag    string = "chamber-use-paths"
	serviceFlag            string = "service"
	environmentFlag        string = "environment"
	repositoryNameFlag     string = "repository-name"
	imageTagFlag           string = "image-tag"
	commandFlag            string = "command"
	commandArgsFlag        string = "command-args"
	variablesFileFlag      string = "variables-file"
	dryRunFlag             string = "dry-run"
)

func initFlags(flag *pflag.FlagSet) {

	// AWS Account
	flag.String(awsAccountIDFlag, "", "The AWS Account ID")

	// AWS Flags
	cli.InitAWSFlags(flag)

	// Vault Flags
	cli.InitVaultFlags(flag)

	// Chamber Settings
	flag.String(chamberBinaryFlag, "/bin/chamber", "Chamber Binary")
	flag.Int(chamberRetriesFlag, 20, "Chamber Retries")
	flag.String(chamberKMSKeyAliasFlag, "alias/aws/ssm", "Chamber KMS Key Alias")
	flag.Int(chamberUsePathsFlag, 1, "Chamber Use Paths")

	// Task Definition Settings
	flag.String(serviceFlag, "", fmt.Sprintf("The service name (choose %q)", services))
	flag.String(environmentFlag, "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String(repositoryNameFlag, "", fmt.Sprintf("The name of the repository where the tagged image resides"))
	flag.String(imageTagFlag, "", "The name of the image tag referenced in the task definition")
	flag.String(commandFlag, "", fmt.Sprintf("The name of the command to run inside the docker container (choose %q)", commands))
	flag.String(commandArgsFlag, "", "The space separated arguments for the command")
	flag.String(variablesFileFlag, "", "A file containing variables for the task definiton")

	// Verbose
	cli.InitVerboseFlags(flag)
	flag.Bool(dryRunFlag, false, "Execute as a dry-run without modifying AWS.")

	// Don't sort flags
	flag.SortFlags = false
}

func checkConfig(v *viper.Viper) error {

	awsAccountID := v.GetString(awsAccountIDFlag)
	if len(awsAccountID) == 0 {
		return errors.Wrap(&errInvalidAccountID{AwsAccountID: awsAccountID}, fmt.Sprintf("%q is invalid", awsAccountIDFlag))
	}

	region, err := cli.CheckAWSRegion(v)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid", cli.AWSRegionFlag))
	}

	if err := cli.CheckAWSRegionForService(region, cloudwatchevents.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, ecs.ServiceName))
	}

	if err := cli.CheckAWSRegionForService(region, ecs.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, ecs.ServiceName))
	}

	if err := cli.CheckAWSRegionForService(region, ecr.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, ecs.ServiceName))
	}

	if err := cli.CheckAWSRegionForService(region, rds.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, ecs.ServiceName))
	}

	if err := cli.CheckVault(v); err != nil {
		return err
	}

	chamberRetries := v.GetInt(chamberRetriesFlag)
	if chamberRetries < 1 && chamberRetries > 20 {
		return errors.New("Chamber Retries must be greater than or equal to 1 and less than or equal to 20")
	}

	chamberKMSKeyAlias := v.GetString(chamberKMSKeyAliasFlag)
	if len(chamberKMSKeyAlias) == 0 {
		return errors.New("Chamber KMS Key Alias must be set")
	}

	chamberUsePaths := v.GetInt(chamberUsePathsFlag)
	if chamberUsePaths < 1 && chamberUsePaths > 20 {
		return errors.New("Chamber Use Paths must be greater than or equal to 1 and less than or equal to 20")
	}

	serviceName := v.GetString(serviceFlag)
	if len(serviceName) == 0 {
		return errors.Wrap(&errInvalidService{Service: serviceName}, fmt.Sprintf("%q is invalid", serviceFlag))
	}
	validService := false
	for _, str := range services {
		if serviceName == str {
			validService = true
			break
		}
	}
	if !validService {
		return errors.Wrap(&errInvalidService{Service: serviceName}, fmt.Sprintf("%q is invalid", serviceFlag))
	}

	environmentName := v.GetString(environmentFlag)
	if len(environmentName) == 0 {
		return errors.Wrap(&errInvalidEnvironment{Environment: environmentName}, fmt.Sprintf("%q is invalid", environmentFlag))
	}
	validEnvironment := false
	for _, str := range environments {
		if environmentName == str {
			validEnvironment = true
			break
		}
	}
	if !validEnvironment {
		return errors.Wrap(&errInvalidEnvironment{Environment: environmentName}, fmt.Sprintf("%q is invalid", environmentFlag))
	}

	repositoryName := v.GetString(repositoryNameFlag)
	if len(repositoryName) == 0 {
		return errors.Wrap(&errinvalidRepositoryName{RepositoryName: repositoryName}, fmt.Sprintf("%q is invalid", repositoryNameFlag))
	}

	imageTag := v.GetString(imageTagFlag)
	if len(imageTag) == 0 {
		return errors.Wrap(&errinvalidImageTag{ImageTag: imageTag}, fmt.Sprintf("%q is invalid", imageTagFlag))
	}

	commandName := v.GetString(commandFlag)
	if len(commandName) == 0 {
		return errors.Wrap(&errInvalidCommand{Command: commandName}, fmt.Sprintf("%q is invalid", commandFlag))
	}
	validRule := false
	for _, str := range commands {
		if commandName == str {
			validRule = true
			break
		}
	}
	if !validRule {
		return errors.Wrap(&errInvalidCommand{Command: commandName}, fmt.Sprintf("%q is invalid", commandFlag))
	}

	if variablesFile := v.GetString(variablesFileFlag); len(variablesFile) > 0 {
		if _, err := os.Stat(variablesFile); err != nil {
			return errors.Wrap(&errInvalidFile{File: variablesFile}, fmt.Sprintf("%q is invalid", variablesFileFlag))
		}
	}

	return nil
}

func quit(logger *log.Logger, flag *pflag.FlagSet, err error) {
	logger.Println(err.Error())
	logger.Println("Usage of ecs-service-logs:")
	if flag != nil {
		flag.PrintDefaults()
	}
	os.Exit(1)
}

func varFromCtxOrEnv(varName string, ctx map[string]string) string {
	// Return the value if it is in the context
	if i, ok := ctx[varName]; ok {
		return i
	}
	// Default to whatever exists in the environment
	return os.Getenv("DB_PORT")
}

func buildContainerEnvironment(v *viper.Viper, environmentName string, dbHost string, variablesFile string) []*ecs.KeyValuePair {

	// Construct variables from a file for the task def
	// These variables should always be preferred over env vars
	ctx := map[string]string{}
	if len(variablesFile) > 0 {
		// Read contents of variables file into vars
		vars, readFileErr := ioutil.ReadFile(variablesFile)
		if readFileErr != nil {
			log.Fatal(errors.New("error reading variables file"))
		}

		// Adds variables from file into context
		for _, x := range strings.Split(string(vars), "\n") {
			// If a line is empty or starts with #, then skip.
			if len(x) > 0 && x[0] != '#' {
				// Split each line on the first equals sign into [name, value]
				pair := strings.SplitAfterN(x, "=", 2)
				ctx[pair[0][0:len(pair[0])-1]] = pair[1]
			}
		}
	}

	chamberKMSKeyAlias := v.GetString(chamberKMSKeyAliasFlag)
	chamberUsePaths := v.GetInt(chamberUsePathsFlag)
	return []*ecs.KeyValuePair{
		{
			Name:  aws.String("CHAMBER_KMS_KEY_ALIAS"),
			Value: aws.String(chamberKMSKeyAlias),
		},
		{
			Name:  aws.String("CHAMBER_USE_PATHS"),
			Value: aws.String(strconv.Itoa(chamberUsePaths)),
		},
		{
			Name:  aws.String("DB_ENV"),
			Value: aws.String(cli.DbEnvContainer),
		},
		{
			Name:  aws.String("LOGGING_ENV"),
			Value: aws.String(cli.LoggingEnvProduction),
		},
		{
			Name:  aws.String("ENVIRONMENT"),
			Value: aws.String(environmentName),
		},
		{
			Name:  aws.String("DB_HOST"),
			Value: aws.String(dbHost),
		},
		{
			Name:  aws.String("DB_PORT"),
			Value: aws.String(varFromCtxOrEnv("DB_PORT", ctx)),
		},
		{
			Name:  aws.String("DB_USER"),
			Value: aws.String(varFromCtxOrEnv("DB_USER", ctx)),
		},
		{
			Name:  aws.String("DB_NAME"),
			Value: aws.String(varFromCtxOrEnv("DB_NAME", ctx)),
		},
		{
			Name:  aws.String("DB_SSL_MODE"),
			Value: aws.String(varFromCtxOrEnv("DB_SSL_MODE", ctx)),
		},
		{
			Name:  aws.String("DB_SSL_ROOT_CERT"),
			Value: aws.String(varFromCtxOrEnv("DB_SSL_ROOT_CERT", ctx)),
		},
		{
			Name:  aws.String("DB_IAM"),
			Value: aws.String(varFromCtxOrEnv("DB_IAM", ctx)),
		},
		{
			Name:  aws.String("DB_IAM_ROLE"),
			Value: aws.String(varFromCtxOrEnv("DB_IAM_ROLE", ctx)),
		},
		{
			Name:  aws.String("DB_REGION"),
			Value: aws.String(varFromCtxOrEnv("DB_REGION", ctx)),
		},
	}
}

func main() {
	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		quit(logger, flag, err)
	}

	v := viper.New()
	pflagsErr := v.BindPFlags(flag)
	if pflagsErr != nil {
		quit(logger, flag, err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	verbose := v.GetBool(cli.VerboseFlag)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Flag for Dry Run
	dryRun := v.GetBool(dryRunFlag)

	checkConfigErr := checkConfig(v)
	if checkConfigErr != nil {
		quit(logger, flag, checkConfigErr)
	}

	awsConfig, err := cli.GetAWSConfig(v, verbose)
	if err != nil {
		quit(logger, nil, err)
	}

	sess, err := awssession.NewSession(awsConfig)
	if err != nil {
		quit(logger, nil, errors.Wrap(err, "failed to create AWS session"))
	}

	// Create the Services
	serviceCloudWatchEvents := cloudwatchevents.New(sess)
	serviceECS := ecs.New(sess)
	serviceECR := ecr.New(sess)
	serviceRDS := rds.New(sess)

	// Get the current task definition (for rollback)
	commandName := v.GetString(commandFlag)
	cmds := strings.Fields(commandName)
	subCommandName := cmds[0]
	if len(cmds) > 1 {
		subCommandName = cmds[1]
	}
	commandArgs := []string{}
	if str := v.GetString(commandArgsFlag); len(str) > 0 {
		commandArgs = strings.Split(str, " ")
	}
	ruleName := fmt.Sprintf("%s-%s", subCommandName, v.GetString(environmentFlag))
	targetsOutput, err := serviceCloudWatchEvents.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
		Rule: aws.String(ruleName),
	})
	if err != nil {
		quit(logger, nil, errors.Wrap(err, "error retrieving targets for rule"))
	}

	currentTarget := targetsOutput.Targets[0]

	// Get the current task definition
	currentTaskDefArnStr := *currentTarget.EcsParameters.TaskDefinitionArn
	logger.Println(fmt.Sprintf("Current Task Def Arn: %s", currentTaskDefArnStr))
	currentTaskDefArn, err := arn.Parse(currentTaskDefArnStr)
	if err != nil {
		quit(logger, nil, errors.Wrap(err, "Unable to parse current task definition arn"))
	}

	currentTaskDefName := strings.Split(currentTaskDefArn.Resource, ":")[0]
	currentTaskDefName = strings.Split(currentTaskDefName, "/")[1]
	currentDescribeTaskDefinitionOutput, err := serviceECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(currentTaskDefName),
	})
	if err != nil {
		quit(logger, nil, errors.Wrapf(err, "unable to parse current task arn %s", currentTaskDefArnStr))
	}
	currentTaskDef := *currentDescribeTaskDefinitionOutput.TaskDefinition

	// Confirm the image exists
	awsRegion := v.GetString(cli.AWSRegionFlag)
	imageTag := v.GetString(imageTagFlag)
	registryID := v.GetString(awsAccountIDFlag)
	repositoryName := v.GetString(repositoryNameFlag)
	imageName := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", registryID, awsRegion, repositoryName, imageTag)

	_, err = serviceECR.DescribeImages(&ecr.DescribeImagesInput{
		ImageIds: []*ecr.ImageIdentifier{
			{
				ImageTag: aws.String(imageTag),
			},
		},
		RegistryId:     aws.String(registryID),
		RepositoryName: aws.String(repositoryName),
	})
	if err != nil {
		quit(logger, nil, errors.Wrapf(err, "unable retrieving image from %s", imageName))
	}

	// Get the database host using the instance identifier
	serviceName := v.GetString(serviceFlag)
	environmentName := v.GetString(environmentFlag)
	dbInstanceIdentifier := fmt.Sprintf("%s-%s", serviceName, environmentName)
	dbInstancesOutput, err := serviceRDS.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	})
	if err != nil {
		quit(logger, nil, errors.Wrapf(err, "error retrieving database definition for %s", dbInstanceIdentifier))
	}
	dbHost := *dbInstancesOutput.DBInstances[0].Endpoint.Address

	// Name the container definition and verify it exists
	containerDefName := fmt.Sprintf("%s-tasks-%s-%s", serviceName, subCommandName, environmentName)

	// AWS Logs Group is related to the cluster and should not be changed
	awsLogsGroup := fmt.Sprintf("ecs-tasks-%s-%s", serviceName, environmentName)
	awsLogsStreamPrefix := fmt.Sprintf("%s-tasks", serviceName)

	// Chamber Settings
	chamberBinary := v.GetString(chamberBinaryFlag)
	chamberRetries := v.GetInt(chamberRetriesFlag)
	chamberStore := fmt.Sprintf("%s-%s", serviceName, environmentName)

	entryPoint := []string{
		chamberBinary,
		"-r",
		strconv.Itoa(chamberRetries),
		"exec",
		chamberStore,
		"--",
		fmt.Sprintf("/bin/%s", cmds[0]),
	}
	if len(cmds) > 1 {
		entryPoint = append(entryPoint, cmds[1:]...)
	}
	if len(commandArgs) > 0 {
		entryPoint = append(entryPoint, commandArgs...)
	}

	// Register the new task definition
	variablesFile := v.GetString(variablesFileFlag)
	newTaskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:        aws.String(containerDefName),
				Image:       aws.String(imageName),
				Essential:   aws.Bool(true),
				EntryPoint:  aws.StringSlice(entryPoint),
				Command:     []*string{},
				Environment: buildContainerEnvironment(v, environmentName, dbHost, variablesFile),
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: map[string]*string{
						"awslogs-group":         aws.String(awsLogsGroup),
						"awslogs-region":        aws.String(awsRegion),
						"awslogs-stream-prefix": aws.String(awsLogsStreamPrefix),
					},
				},
			},
		},
		Cpu:                     currentTaskDef.Cpu,
		ExecutionRoleArn:        currentTaskDef.ExecutionRoleArn,
		Family:                  currentTaskDef.Family,
		Memory:                  currentTaskDef.Memory,
		NetworkMode:             currentTaskDef.NetworkMode,
		RequiresCompatibilities: currentTaskDef.RequiresCompatibilities,
		TaskRoleArn:             currentTaskDef.TaskRoleArn,
	}
	if verbose {
		logger.Println(newTaskDefInput.String())
	}

	if dryRun {
		logger.Println("Dry run: ECS Task Definition not registered! CloudWatch Target Not Updated!")
		return
	}

	// Register the new task definition
	newTaskDefOutput, err := serviceECS.RegisterTaskDefinition(&newTaskDefInput)
	if err != nil {
		quit(logger, nil, errors.Wrap(err, "error registering new task definition"))
	}
	newTaskDefArn := *newTaskDefOutput.TaskDefinition.TaskDefinitionArn
	logger.Println(fmt.Sprintf("New Task Def Arn: %s", newTaskDefArn))

	// Update the task event target with the new task ECS parameters
	putTargetsOutput, err := serviceCloudWatchEvents.PutTargets(&cloudwatchevents.PutTargetsInput{
		Rule: aws.String(ruleName),
		Targets: []*cloudwatchevents.Target{
			{
				Id:      currentTarget.Id,
				Arn:     currentTarget.Arn,
				RoleArn: currentTarget.RoleArn,
				EcsParameters: &cloudwatchevents.EcsParameters{
					LaunchType:           aws.String("FARGATE"),
					NetworkConfiguration: currentTarget.EcsParameters.NetworkConfiguration,
					TaskCount:            aws.Int64(1),
					TaskDefinitionArn:    aws.String(newTaskDefArn),
				},
			},
		},
	})
	if err != nil {
		quit(logger, nil, errors.Wrap(err, "Unable to put new target"))
	}
	logger.Println(putTargetsOutput)
}
