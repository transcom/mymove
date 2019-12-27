package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

type errInvalidImage struct {
	Image string
}

func (e *errInvalidImage) Error() string {
	return fmt.Sprintf("invalid AWS ECR image tag %q", e.Image)
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
	imageFlag              string = "image"
	variablesFileFlag      string = "variables-file"
	dryRunFlag             string = "dry-run"
)

type ECRImage struct {
	AWSRegion      string
	ImageArn       string
	ImageTag       string
	RegistryId     string
	RepositoryArn  string
	RepositoryName string
}

func NewECRImage(imageName string) *ECRImage {
	imageParts := strings.Split(imageName, ":")
	repositoryArn, imageTag := imageParts[0], imageParts[1]
	repositoryArnParts := strings.Split(repositoryArn, "/")
	repositoryName := repositoryArnParts[1]
	repositoryDomainParts := strings.Split(repositoryArnParts[0], ".")
	registryID, awsRegion := repositoryDomainParts[0], repositoryDomainParts[3]

	return &ECRImage{
		AWSRegion:      awsRegion,
		ImageArn:       imageName,
		ImageTag:       imageTag,
		RegistryId:     registryID,
		RepositoryArn:  repositoryArn,
		RepositoryName: repositoryName,
	}
}

func initTaskDefFlags(flag *pflag.FlagSet) {

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
	flag.String(serviceFlag, "app", fmt.Sprintf("The service name (choose %q)", services))
	flag.String(environmentFlag, "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String(imageFlag, "", "The name of the image referenced in the task definition")
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

	image := v.GetString(imageFlag)
	if len(image) == 0 {
		return errors.Wrap(&errInvalidImage{Image: image}, fmt.Sprintf("%q is invalid", imageFlag))
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

func buildContainerEnvironment(environmentName string, dbHost string, variablesFile string) []*ecs.KeyValuePair {

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

	return []*ecs.KeyValuePair{
		{
			Name:  aws.String("CHAMBER_KMS_KEY_ALIAS"),
			Value: aws.String(varFromCtxOrEnv("CHAMBER_KMS_KEY_ALIAS", ctx)),
		},
		{
			Name:  aws.String("CHAMBER_USE_PATHS"),
			Value: aws.String(varFromCtxOrEnv("CHAMBER_USE_PATHS", ctx)),
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

func taskDefFunction(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "could not parse flags")
	}

	flag := cmd.Flags()

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return errors.Wrap(err, "could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

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

	// Ensure the configuration works against the variables
	checkConfigErr := checkConfig(v)
	if checkConfigErr != nil {
		quit(logger, flag, checkConfigErr)
	}

	// Get the AWS configuration so we can build a session
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

	// ===== This stuff should come from the CLI or function =====
	// Take inputs from lambda
	environmentName := v.GetString(environmentFlag)
	serviceName := v.GetString(serviceFlag)
	imageName := v.GetString(imageFlag)
	// =====

	// Confirm the image exists
	ecrImage := NewECRImage(imageName)
	imageIdentifier := ecr.ImageIdentifier{}
	imageIdentifier.SetImageTag(ecrImage.ImageTag)
	errImageIdentifierValidate := imageIdentifier.Validate()
	if errImageIdentifierValidate != nil {
		quit(logger, nil, errors.Wrapf(errImageIdentifierValidate, "image identifier tag invalid %s", ecrImage.ImageTag))
	}

	_, err = serviceECR.DescribeImages(&ecr.DescribeImagesInput{
		ImageIds: []*ecr.ImageIdentifier{
			{
				ImageTag: aws.String(ecrImage.ImageTag),
			},
		},
		RegistryId:     aws.String(ecrImage.RegistryId),
		RepositoryName: aws.String(ecrImage.RepositoryName),
	})
	if err != nil {
		quit(logger, nil, errors.Wrapf(err, "unable retrieving image from %s", imageName))
	}

	// Get the current task definition
	var currentTaskDef ecs.TaskDefinition
	scheduledTask := false
	var commandName, subCommandName string
	if scheduledTask {
		commandName = "milmove-tasks"
		subCommandName = "send-post-move-survey"
		ruleName := fmt.Sprintf("%s-%s", subCommandName, environmentName)
		targetsOutput, err := serviceCloudWatchEvents.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
			Rule: aws.String(ruleName),
		})
		if err != nil {
			quit(logger, nil, errors.Wrap(err, "error retrieving targets for rule"))
		}

		currentTarget := targetsOutput.Targets[0]

		// Get the current task definition
		currentTaskDefArnStr := *currentTarget.EcsParameters.TaskDefinitionArn
		if verbose {
			logger.Println(fmt.Sprintf("Current Task Def Arn: %s", currentTaskDefArnStr))
		}
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
		currentTaskDef = *currentDescribeTaskDefinitionOutput.TaskDefinition
	} else {
		commandName = "milmove"
		subCommandName = "serve"
	}

	// Register the new task definition
	variablesFile := v.GetString(variablesFileFlag)
	newTaskDefInput, err := renderTaskDefinition(
		ecrImage,
		serviceRDS,
		serviceName,
		environmentName,
		commandName,
		subCommandName,
		variablesFile,
		currentTaskDef.ExecutionRoleArn,
		currentTaskDef.TaskRoleArn)
	if err != nil {
		quit(logger, nil, err)
	}

	if verbose {
		logger.Println(newTaskDefInput.String())
	}

	if dryRun {
		logger.Println("Dry run: ECS Task Definition not registered! CloudWatch Target Not Updated!")
		return nil
	}

	return nil
}

func renderTaskDefinition(ecrImage *ECRImage, serviceRDS *rds.RDS, serviceName, environmentName, commandName, subCommandName, variablesFile, executionRoleArn, taskRoleArn string) (*ecs.RegisterTaskDefinitionInput, error) {

	// Get the database host using the instance identifier
	dbInstanceIdentifier := fmt.Sprintf("%s-%s", serviceName, environmentName)
	dbInstancesOutput, err := serviceRDS.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error retrieving database definition for %s", dbInstanceIdentifier)
	}
	dbHost := *dbInstancesOutput.DBInstances[0].Endpoint.Address

	// Name the container definition and verify it exists
	containerDefName := fmt.Sprintf("%s-%s-%s", ecrImage.RepositoryName, subCommandName, environmentName)

	// AWS Logs Group is related to the cluster and should not be changed
	awsLogsGroup := fmt.Sprintf("ecs-%s-%s", ecrImage.RepositoryName, environmentName)

	// Entrypoint
	entryPoint := []string{
		fmt.Sprintf("/bin/%s", commandName),
	}

	newTaskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:        aws.String(containerDefName),
				Image:       aws.String(ecrImage.ImageArn),
				Essential:   aws.Bool(true),
				EntryPoint:  aws.StringSlice(entryPoint),
				Command:     []*string{},
				Environment: buildContainerEnvironment(environmentName, dbHost, variablesFile),
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: map[string]*string{
						"awslogs-group":         aws.String(awsLogsGroup),
						"awslogs-region":        aws.String(ecrImage.AWSRegion),
						"awslogs-stream-prefix": aws.String(ecrImage.RepositoryName),
					},
				},
			},
		},
		Cpu:                     aws.String("512"),
		ExecutionRoleArn:        aws.String(executionRoleArn),
		Family:                  aws.String(fmt.Sprintf("%s-%s", serviceName, environmentName)),
		Memory:                  aws.String("2048"),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String("FARGATE")},
		TaskRoleArn:             aws.String(taskRoleArn),
	}
	return &newTaskDefInput, nil
}
