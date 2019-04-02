package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/99designs/aws-vault/prompt"
	"github.com/99designs/aws-vault/vault"
	"github.com/99designs/keyring"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var services = []string{"app"}
var environments = []string{"prod", "staging", "experimental"}

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %q", e.Region)
}

type errInvalidService struct {
	Service string
}

func (e *errInvalidService) Error() string {
	return fmt.Sprintf("invalid service %q, expecting one of %q", e.Service, services)
}

type errInvalidEnvironment struct {
	Environment string
}

func (e *errInvalidEnvironment) Error() string {
	return fmt.Sprintf("invalid environment %q, expecting one of %q", e.Environment, environments)
}

type errInvalidImage struct {
	Image string
}

func (e *errInvalidImage) Error() string {
	return fmt.Sprintf("invalid image %q", e.Image)
}

type errInvalidRule struct {
	Rule string
}

func (e *errInvalidRule) Error() string {
	return fmt.Sprintf("invalid rule %q", e.Rule)
}

const (
	awsRegionFlag            string = "aws-region"
	awsProfileFlag           string = "aws-profile"
	awsVaultKeychainNameFlag string = "aws-vault-keychain-name"
	chamberRetriesFlag       string = "chamber-retries"
	chamberKMSKeyAliasFlag   string = "chamber-kms-key-alias"
	chamberUsePathsFlag      string = "chamber-use-paths"
	serviceFlag              string = "service"
	environmentFlag          string = "environment"
	imageFlag                string = "image"
	ruleFlag                 string = "rule"
	eiaKeyFlag               string = "eia-key"
	eiaURLFlag               string = "eia-url"
	verboseFlag              string = "verbose"
)

func initFlags(flag *pflag.FlagSet) {

	// AWS Vault Settings
	flag.String(awsRegionFlag, "us-west-2", "The AWS Region")
	flag.String(awsProfileFlag, "", "The aws-vault profile")
	flag.String(awsVaultKeychainNameFlag, "", "The aws-vault keychain name")

	// Chamber Settings
	flag.Int(chamberRetriesFlag, 20, "Chamber Retries")
	flag.String(chamberKMSKeyAliasFlag, "alias/aws/ssm", "Chamber KMS Key Alias")
	flag.Int(chamberUsePathsFlag, 1, "Chamber Use Paths")

	// Task Definition Settings
	flag.String(serviceFlag, "", fmt.Sprintf("The service name (choose %q)", services))
	flag.String(environmentFlag, "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String(imageFlag, "", "The name of the image referenced in the task definition")
	flag.String(ruleFlag, "", "The name of the CloudWatch Event Rule targeting the Task Definition")

	// EIA Open Data API
	flag.String(eiaKeyFlag, "", "Key for Energy Information Administration (EIA) api")
	flag.String(eiaURLFlag, "", "Url for Energy Information Administration (EIA) api")

	// Script settings
	flag.BoolP(verboseFlag, "v", false, "Print section lines")
}

func checkEIAKey(v *viper.Viper) error {
	eiaKey := v.GetString(eiaKeyFlag)
	if len(eiaKey) != 32 {
		return fmt.Errorf("expected eia key to be 32 characters long; key is %d chars", len(eiaKey))
	}
	return nil
}

func checkEIAURL(v *viper.Viper) error {
	eiaURL := v.GetString(eiaURLFlag)
	if eiaURL != "https://api.eia.gov/series/" {
		return fmt.Errorf("invalid eia url %s, expecting https://api.eia.gov/series/", eiaURL)
	}
	return nil
}

func checkConfig(v *viper.Viper) error {

	regions, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.EcsServiceID)
	if !ok {
		return fmt.Errorf("could not find regions for service %q", endpoints.EcsServiceID)
	}

	region := v.GetString(awsRegionFlag)
	if len(region) == 0 {
		return errors.Wrap(&errInvalidRegion{Region: region}, fmt.Sprintf("%q is invalid", awsRegionFlag))
	}

	if _, ok := regions[region]; !ok {
		return errors.Wrap(&errInvalidRegion{Region: region}, fmt.Sprintf("%q is invalid", awsRegionFlag))
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
		return errors.Wrap(&errInvalidService{Service: serviceName}, fmt.Sprintf("%q is invalid", "service"))
	}
	validService := false
	for _, str := range services {
		if serviceName == str {
			validService = true
			break
		}
	}
	if !validService {
		return errors.Wrap(&errInvalidService{Service: serviceName}, fmt.Sprintf("%q is invalid", "service"))
	}

	environmentName := v.GetString(environmentFlag)
	if len(environmentName) == 0 {
		return errors.Wrap(&errInvalidEnvironment{Environment: environmentName}, fmt.Sprintf("%q is invalid", "environment"))
	}
	validEnvironment := false
	for _, str := range environments {
		if environmentName == str {
			validEnvironment = true
			break
		}
	}
	if !validEnvironment {
		return errors.Wrap(&errInvalidEnvironment{Environment: environmentName}, fmt.Sprintf("%q is invalid", "environment"))
	}

	image := v.GetString(imageFlag)
	if len(image) == 0 {
		return errors.Wrap(&errInvalidImage{Image: image}, fmt.Sprintf("%q is invalid", "image"))
	}

	rule := v.GetString(ruleFlag)
	if len(rule) == 0 {
		return errors.Wrap(&errInvalidRule{Rule: rule}, fmt.Sprintf("%q is invalid", "rule"))
	}

	err := checkEIAKey(v)
	if err != nil {
		return err
	}

	err = checkEIAURL(v)
	if err != nil {
		return err
	}

	return nil
}

func quit(logger *log.Logger, flag *pflag.FlagSet, err error) {
	logger.Println(err.Error())
	fmt.Println("Usage of ecs-service-logs:")
	if flag != nil {
		flag.PrintDefaults()
	}
	os.Exit(1)
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
	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	if !v.GetBool(verboseFlag) {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	err := checkConfig(v)
	if err != nil {
		quit(logger, flag, err)
	}

	awsRegion := v.GetString(awsRegionFlag)

	awsConfig := &aws.Config{
		Region: aws.String(awsRegion),
	}

	verbose := v.GetBool(verboseFlag)
	keychainName := v.GetString(awsVaultKeychainNameFlag)
	keychainProfile := v.GetString(awsProfileFlag)

	if len(keychainName) > 0 && len(keychainProfile) > 0 {
		creds, err := getAWSCredentials(keychainName, keychainProfile)
		if err != nil {
			quit(logger, flag, errors.Wrap(err, fmt.Sprintf("Unable to get AWS credentials from the keychain %s and profile %s", keychainName, keychainProfile)))
		}
		awsConfig.CredentialsChainVerboseErrors = aws.Bool(verbose)
		awsConfig.Credentials = creds
	}

	sess, err := awssession.NewSession(awsConfig)
	if err != nil {
		quit(logger, flag, errors.Wrap(err, "failed to create AWS session"))
	}

	// Create the Services
	serviceCloudWatchEvents := cloudwatchevents.New(sess)
	serviceECS := ecs.New(sess)
	serviceRDS := rds.New(sess)

	// Get the current task definition (for rollback)
	ruleName := v.GetString(ruleFlag)
	targetsOutput, err := serviceCloudWatchEvents.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
		Rule: aws.String(ruleName),
	})
	if err != nil {
		quit(logger, flag, errors.Wrap(err, "error retrieving targets for rule"))
	}

	blueTaskDefArn := *targetsOutput.Targets[0].EcsParameters.TaskDefinitionArn
	fmt.Println(blueTaskDefArn)

	// aws ecs describe-task-definition --task-definition=app-scheduled-task-save_fuel_price_data-experimental:1
	blueTaskDef, err := serviceECS.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &blueTaskDefArn,
	})
	if err != nil {
		quit(logger, flag, errors.Wrapf(err, "error retrieving task definition for %s", blueTaskDefArn))
	}
	fmt.Println(blueTaskDef)

	// Get the database host using the instance identifier
	environmentName := v.GetString(environmentFlag)
	dbInstanceIdentifier := fmt.Sprintf("app-%s", environmentName)
	dbInstancesOutput, err := serviceRDS.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	})
	if err != nil {
		quit(logger, flag, errors.Wrapf(err, "error retrieving database definition for %s", dbInstanceIdentifier))
	}
	dbHost := *dbInstancesOutput.DBInstances[0].Endpoint.Address

	// Register the new task definition
	serviceName := v.GetString(serviceFlag)
	imageName := v.GetString(imageFlag)
	familyName := fmt.Sprintf("%s-%s", serviceName, environmentName)
	taskRoleArn := fmt.Sprintf("ecs-task-role-%s", familyName)
	executionRoleArn := fmt.Sprintf("ecs-task-excution-role-%s", familyName)
	containerDefName := fmt.Sprintf("app-tasks-%s-%s", ruleName, environmentName)

	// AWS Logs Group is related to the cluster and should not be changed
	awsLogsGroup := fmt.Sprintf("ecs-tasks-app-%s", environmentName)

	// Chamber Settings
	chamberRetries := v.GetInt(chamberRetriesFlag)
	chamberKMSKeyAlias := v.GetString(chamberKMSKeyAliasFlag)
	chamberUsePaths := v.GetInt(chamberUsePathsFlag)

	// Tool Settings
	eiaKey := v.GetString(eiaKeyFlag)
	eiaURL := v.GetString(eiaURLFlag)

	taskDefinitionOutput, err := serviceECS.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:      aws.String(containerDefName),
				Image:     aws.String(imageName),
				Essential: aws.Bool(false),
				EntryPoint: []*string{
					aws.String("/bin/chamber"),
					aws.String("-r"),
					aws.String(strconv.Itoa(chamberRetries)),
					aws.String("exec"),
					aws.String(familyName),
					aws.String("--"),
					aws.String(fmt.Sprintf("/bin/%s", ruleName)),
				},
				Environment: []*ecs.KeyValuePair{
					&ecs.KeyValuePair{
						Name:  aws.String("ENVIRONMENT"),
						Value: aws.String(environmentName),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("DB_HOST"),
						Value: aws.String(dbHost),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("DB_PORT"),
						Value: aws.String("5432"),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("DB_USER"),
						Value: aws.String("master"),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("DB_NAME"),
						Value: aws.String("app"),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("CHAMBER_KMS_KEY_ALIAS"),
						Value: aws.String(chamberKMSKeyAlias),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("CHAMBER_USE_PATHS"),
						Value: aws.String(strconv.Itoa(chamberUsePaths)),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("EIA_KEY"),
						Value: aws.String(eiaKey),
					},
					&ecs.KeyValuePair{
						Name:  aws.String("EIA_URL"),
						Value: aws.String(eiaURL),
					},
				},
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: map[string]*string{
						"awslogs-group":         aws.String(awsLogsGroup),
						"awslogs-region":        aws.String(awsRegion),
						"awslogs-stream-prefix": aws.String(containerDefName),
					},
				},
			},
		},
		Cpu:                     aws.String("256"),
		ExecutionRoleArn:        aws.String(executionRoleArn),
		Family:                  aws.String(familyName),
		Memory:                  aws.String("512"),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String("FARGATE")},
		TaskRoleArn:             aws.String(taskRoleArn),
	})
	if err != nil {
		quit(logger, flag, errors.Wrap(err, "error registering new task definition"))
	}
	greenTaskDefArn := *taskDefinitionOutput.TaskDefinition.TaskDefinitionArn
	fmt.Println(greenTaskDefArn)

	// aws events puts-target
}
