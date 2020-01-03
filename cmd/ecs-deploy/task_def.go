package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol/json/jsonutil"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

var services = []string{"app"}
var environments = []string{"prod", "staging", "experimental"}

// Commands should be the name of the binary found in the /bin directory in the container
//var commands = []string{"milmove-tasks save-fuel-price-data", "milmove-tasks send-post-move-survey"}

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
	awsAccountIDFlag  string = "aws-account-id"
	serviceFlag       string = "service"
	environmentFlag   string = "environment"
	imageURIFlag      string = "image"
	variablesFileFlag string = "variables-file"
)

type ECRImage struct {
	AWSRegion      string
	imageURI       string
	ImageTag       string
	RegistryID     string
	RepositoryURI  string
	RepositoryName string
}

func NewECRImage(imageURI string) *ECRImage {
	imageParts := strings.Split(imageURI, ":")
	repositoryURI, imageTag := imageParts[0], imageParts[1]
	repositoryURIParts := strings.Split(repositoryURI, "/")
	repositoryName := repositoryURIParts[1]
	repositoryDomainParts := strings.Split(repositoryURIParts[0], ".")
	registryID, awsRegion := repositoryDomainParts[0], repositoryDomainParts[3]

	return &ECRImage{
		AWSRegion:      awsRegion,
		imageURI:       imageURI,
		ImageTag:       imageTag,
		RegistryID:     registryID,
		RepositoryURI:  repositoryURI,
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

	// Task Definition Settings
	flag.String(serviceFlag, "app", fmt.Sprintf("The service name (choose %q)", services))
	flag.String(environmentFlag, "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String(imageURIFlag, "", "The URI of the container image to use in the task definition")
	flag.String(variablesFileFlag, "", "A file containing variables for the task definiton")

	// Verbose
	cli.InitVerboseFlags(flag)

	// Sort flags
	flag.SortFlags = true
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
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, cloudwatchevents.ServiceName))
	}

	if err := cli.CheckAWSRegionForService(region, ecs.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, ecs.ServiceName))
	}

	if err := cli.CheckAWSRegionForService(region, ecr.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, ecr.ServiceName))
	}

	if err := cli.CheckAWSRegionForService(region, rds.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, rds.ServiceName))
	}

	if err := cli.CheckAWSRegionForService(region, ssm.ServiceName); err != nil {
		return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", cli.AWSRegionFlag, ssm.ServiceName))
	}

	if err := cli.CheckVault(v); err != nil {
		return err
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

	image := v.GetString(imageURIFlag)
	if len(image) == 0 {
		return errors.Wrap(&errInvalidImage{Image: image}, fmt.Sprintf("%q is invalid", imageURIFlag))
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

func buildSecrets(serviceSSM *ssm.SSM, awsRegion, awsAccountID, serviceName, environmentName string) []*ecs.Secret {

	var secrets []*ecs.Secret

	params := ssm.DescribeParametersInput{
		MaxResults: aws.Int64(50),
	}

	ctx := context.Background()

	p := request.Pagination{
		NewRequest: func() (*request.Request, error) {
			req, _ := serviceSSM.DescribeParametersRequest(&params)
			req.SetContext(ctx)
			return req, nil
		},
	}

	for p.Next() {
		page := p.Page().(*ssm.DescribeParametersOutput)

		for _, parameter := range page.Parameters {
			if strings.HasPrefix(*parameter.Name, fmt.Sprintf("/%s-%s", serviceName, environmentName)) {
				secrets = append(secrets, &ecs.Secret{
					Name:      aws.String(strings.ToUpper(strings.Split(*parameter.Name, "/")[2])),
					ValueFrom: aws.String(fmt.Sprintf("arn:aws:ssm:%s:%s:parameter%s", awsRegion, awsAccountID, *parameter.Name)),
				})
			}
		}
	}

	return secrets
}

func buildContainerEnvironment(environmentName string, dbHost string, variablesFile string) []*ecs.KeyValuePair {

	envVars := map[string]string{
		"DB_ENV":      cli.DbEnvContainer,
		"DB_HOST":     dbHost,
		"ENVIRONMENT": environmentName,
		"LOGGING_ENV": cli.LoggingEnvProduction,
	}

	// Construct variables from a file for the task def
	// These variables should always be preferred over env vars
	if len(variablesFile) > 0 {
		if _, err := os.Stat(variablesFile); os.IsNotExist(err) {
			log.Fatal(fmt.Errorf("File %q does not exist: %w", variablesFile, err))
		}
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
				envVars[pair[0][0:len(pair[0])-1]] = pair[1]
			}
		}
	}

	var ecsKVPair []*ecs.KeyValuePair

	// Sort these for easier reading
	keys := make([]string, 0, len(envVars))
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		ecsKVPair = append(ecsKVPair, &ecs.KeyValuePair{
			Name:  aws.String(key),
			Value: aws.String(envVars[key]),
		})
	}
	return ecsKVPair

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
	// serviceECS := ecs.New(sess)
	serviceECR := ecr.New(sess)
	serviceRDS := rds.New(sess)
	serviceSSM := ssm.New(sess)

	// ===== Limit the variables required =====
	awsAccountID := v.GetString(awsAccountIDFlag)
	awsRegion := v.GetString(cli.AWSRegionFlag)
	environmentName := v.GetString(environmentFlag)
	serviceName := v.GetString(serviceFlag)
	imageURI := v.GetString(imageURIFlag)
	variablesFile := v.GetString(variablesFileFlag)

	// Confirm the image exists
	ecrImage := NewECRImage(imageURI)
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
		RegistryId:     aws.String(ecrImage.RegistryID),
		RepositoryName: aws.String(ecrImage.RepositoryName),
	})
	if err != nil {
		quit(logger, nil, errors.Wrapf(err, "unable retrieving image from %s", imageURI))
	}

	// Set the command name
	scheduledTask := false
	var commandName, subCommandName string
	if scheduledTask {
		commandName = "milmove-tasks"
		subCommandName = "send-post-move-survey"
		ruleName := fmt.Sprintf("%s-%s", subCommandName, environmentName)
		_, listTargetsByRuleErr := serviceCloudWatchEvents.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
			Rule: aws.String(ruleName),
		})
		if listTargetsByRuleErr != nil {
			quit(logger, nil, fmt.Errorf("error retrieving targets for rule %q: %w", ruleName, listTargetsByRuleERr))
		}
	} else {
		commandName = "milmove"
		subCommandName = "serve"
	}

	// Register the new task definition
	executionRoleArn := fmt.Sprintf("ecs-task-execution-role-%s-%s", serviceName, environmentName)
	taskRoleArn := fmt.Sprintf("ecs-task-role-%s-%s", serviceName, environmentName)

	// Get the database host using the instance identifier
	dbInstanceIdentifier := fmt.Sprintf("%s-%s", serviceName, environmentName)
	dbInstancesOutput, err := serviceRDS.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	})
	if err != nil {
		quit(logger, nil, fmt.Errorf("error retrieving database definition for %q: %w", dbInstanceIdentifier, err))
	}
	dbHost := *dbInstancesOutput.DBInstances[0].Endpoint.Address

	// Name the container definition and verify it exists
	containerDefName := fmt.Sprintf("%s-%s", serviceName, environmentName)

	// AWS Logs Group is related to the cluster and should not be changed
	awsLogsGroup := fmt.Sprintf("ecs-tasks-%s-%s", serviceName, environmentName)

	// Entrypoint
	entryPoint := []string{
		fmt.Sprintf("/bin/%s", commandName),
		subCommandName,
	}

	newTaskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:        aws.String(containerDefName),
				Image:       aws.String(ecrImage.imageURI),
				Essential:   aws.Bool(true),
				EntryPoint:  aws.StringSlice(entryPoint),
				Command:     []*string{},
				Secrets:     buildSecrets(serviceSSM, awsRegion, awsAccountID, serviceName, environmentName),
				Environment: buildContainerEnvironment(environmentName, dbHost, variablesFile),
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: map[string]*string{
						"awslogs-group":         aws.String(awsLogsGroup),
						"awslogs-region":        aws.String(awsRegion),
						"awslogs-stream-prefix": aws.String(serviceName),
					},
				},
				PortMappings: []*ecs.PortMapping{
					{
						ContainerPort: aws.Int64(8443),
						HostPort:      aws.Int64(8443),
						Protocol:      aws.String("tcp"),
					},
				},
				ReadonlyRootFilesystem: aws.Bool(true),
			},
		},
		ExecutionRoleArn:        aws.String(executionRoleArn),
		Family:                  aws.String(fmt.Sprintf("%s-%s", serviceName, environmentName)),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String("FARGATE")},
		TaskRoleArn:             aws.String(taskRoleArn),
		Cpu:                     aws.String("512"),
		Memory:                  aws.String("2048"),
	}

	newTaskDefJSON, jsonErr := jsonutil.BuildJSON(newTaskDefInput)
	if jsonErr != nil {
		quit(logger, nil, err)
	}
	logger.Println(string(newTaskDefJSON))
	// logger.Println(newTaskDefInput.GoString())

	return nil
}
