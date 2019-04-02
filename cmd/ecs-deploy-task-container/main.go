package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

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

type errInvalidTemplateFile struct {
	TemplateFile string
}

func (e *errInvalidTemplateFile) Error() string {
	return fmt.Sprintf("invalid template %q", e.TemplateFile)
}

type errInvalidVariablesFile struct {
	VariablesFile string
}

func (e *errInvalidVariablesFile) Error() string {
	return fmt.Sprintf("invalid variables %q", e.VariablesFile)
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
	serviceFlag              string = "service"
	environmentFlag          string = "environment"
	templateFlag             string = "template"
	variablesFlag            string = "variables"
	imageFlag                string = "image"
	ruleFlag                 string = "rule"
	verboseFlag              string = "verbose"
)

func initFlags(flag *pflag.FlagSet) {
	flag.String(awsRegionFlag, "us-west-2", "The AWS Region")
	flag.String(awsProfileFlag, "", "The aws-vault profile")
	flag.String(awsVaultKeychainNameFlag, "", "The aws-vault keychain name")
	flag.String(serviceFlag, "", fmt.Sprintf("The service name (choose %q)", services))
	flag.String(environmentFlag, "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String(templateFlag, "", "The name of the template file to use for rendering the task definition")
	flag.String(variablesFlag, "", "The name of the variables file to use for rendering the task definition")
	flag.String(imageFlag, "", "The name of the image referenced in the task definition")
	flag.String(ruleFlag, "", "The name of the CloudWatch Event Rule targeting the Task Definition")
	flag.BoolP(verboseFlag, "v", false, "Print section lines")
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

	environment := v.GetString(environmentFlag)
	if len(environment) == 0 {
		return errors.Wrap(&errInvalidEnvironment{Environment: environment}, fmt.Sprintf("%q is invalid", "environment"))
	}
	validEnvironment := false
	for _, str := range environments {
		if environment == str {
			validEnvironment = true
			break
		}
	}
	if !validEnvironment {
		return errors.Wrap(&errInvalidEnvironment{Environment: environment}, fmt.Sprintf("%q is invalid", "environment"))
	}

	templateFile := v.GetString(templateFlag)
	if len(templateFile) == 0 {
		return errors.Wrap(&errInvalidTemplateFile{TemplateFile: templateFile}, fmt.Sprintf("%q is invalid", "template"))
	}
	// Confirm file exists
	templateFileInfo, err := os.Stat(templateFile)
	if err != nil {
		return errors.Wrap(&errInvalidTemplateFile{TemplateFile: templateFile}, fmt.Sprintf("%q file does not exist", "template"))
	}
	if templateFileInfo.IsDir() {
		return errors.Wrap(&errInvalidTemplateFile{TemplateFile: templateFile}, fmt.Sprintf("%q is a directory, not a file", "template"))
	}
	if templateFileInfo.Size() == 0 {
		return errors.Wrap(&errInvalidTemplateFile{TemplateFile: templateFile}, fmt.Sprintf("%q is an empty file", "template"))
	}

	variablesFile := v.GetString(variablesFlag)
	if len(variablesFile) == 0 {
		return errors.Wrap(&errInvalidVariablesFile{VariablesFile: variablesFile}, fmt.Sprintf("%q is invalid", "variables"))
	}
	// Confirm file exists
	variablesFileInfo, err := os.Stat(variablesFile)
	if err != nil {
		return errors.Wrap(&errInvalidVariablesFile{VariablesFile: variablesFile}, fmt.Sprintf("%q file does not exist", "variables"))
	}
	if variablesFileInfo.IsDir() {
		return errors.Wrap(&errInvalidVariablesFile{VariablesFile: variablesFile}, fmt.Sprintf("%q is a directory, not a file", "variables"))
	}
	if variablesFileInfo.Size() == 0 {
		return errors.Wrap(&errInvalidVariablesFile{VariablesFile: variablesFile}, fmt.Sprintf("%q is an empty file", "variables"))
	}

	image := v.GetString(imageFlag)
	if len(image) == 0 {
		return errors.Wrap(&errInvalidImage{Image: image}, fmt.Sprintf("%q is invalid", "image"))
	}

	rule := v.GetString(ruleFlag)
	if len(rule) == 0 {
		return errors.Wrap(&errInvalidRule{Rule: rule}, fmt.Sprintf("%q is invalid", "rule"))
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

func render(logger *log.Logger, templateFile string, variablesFile string, templateVars map[string]string) (*string, error) {
	// Read contents of template file into tmpl
	// #nosec because we want to read in from a file
	tmpl, err := ioutil.ReadFile(templateFile)
	if err != nil {
		quit(logger, nil, errors.Wrap(err, fmt.Sprintf("unable to read template file %s", templateFile)))
	}

	ctx := map[string]string{}

	if len(variablesFile) > 0 {
		// Read contents of variables file into vars
		// #nosec because we want to read in from a file
		vars, err := ioutil.ReadFile(variablesFile)
		if err != nil {
			quit(logger, nil, errors.Wrap(err, fmt.Sprintf("unable to read variables file %s", variablesFile)))
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

	// Adds environment vairables to context
	// os.Environ() returns a copy of strings representing the environment, in the form "key=value".
	// https://golang.org/pkg/os/#Environ
	for _, x := range os.Environ() {
		// Split each environment variable on the first equals sign into [name, value]
		pair := strings.SplitAfterN(x, "=", 2)
		// Add to context
		ctx[pair[0][0:len(pair[0])-1]] = pair[1]
	}

	// Adds template variables to context
	for k, v := range templateVars {
		ctx[k] = v
	}

	t, err := template.New("main").Option("missingkey=error").Parse(string(tmpl))
	if err != nil {
		quit(logger, nil, errors.Wrap(err, "unable to parse the template"))
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, ctx); err != nil {
		quit(logger, nil, errors.Wrap(err, "unable to render the template"))
	}

	tplStr := tpl.String()

	return &tplStr, nil
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
	environment := v.GetString(environmentFlag)
	dbInstanceIdentifier := fmt.Sprintf("app-%s", environment)
	dbInstancesOutput, err := serviceRDS.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	})
	if err != nil {
		quit(logger, flag, errors.Wrapf(err, "error retrieving database definition for %s", dbInstanceIdentifier))
	}
	dbHost := *dbInstancesOutput.DBInstances[0].Endpoint.Address

	// Build the template
	templateFile := v.GetString(templateFlag)
	variablesFile := v.GetString(variablesFlag)
	templateVars := map[string]string{
		"environment": environment,
		"image":       v.GetString(imageFlag),
		"db_host":     dbHost,
	}
	newDef, err := render(logger, templateFile, variablesFile, templateVars)
	if err != nil {
		quit(logger, flag, err)
	}
	fmt.Println(*newDef)

	// Register the new task definition
	serviceName := v.GetString(serviceFlag)
	familyName := fmt.Sprintf("%s-%s", serviceName, environmentName)
	executionTaskRoleArn := fmt.Sprintf("ecs-task-role-%s-%s", serviceName, environmentName)
	taskDefinitionOutput, err := serviceECS.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		ContainerDefinition: [], // JSON?
		ExcecutionRoleArn: aws.String(executionRoleArn),
		Family: aws.String(familyName),
		NetworkMode: aws.String("awsvpc"),
		task role arn
		requires compatibilities fargate
		execution roel arn

	})
	if err != nil {
		quit(logger, flag, errors.Wrap(err, "error registering new task definition"))
	}

	// aws events puts-target
}
