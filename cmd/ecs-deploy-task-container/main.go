package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/99designs/aws-vault/prompt"
	"github.com/99designs/aws-vault/vault"
	"github.com/99designs/keyring"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
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

type errInvalidTemplate struct {
	Template string
}

func (e *errInvalidTemplate) Error() string {
	return fmt.Sprintf("invalid template %q", e.Template)
}

type errInvalidImage struct {
	Image string
}

func (e *errInvalidImage) Error() string {
	return fmt.Sprintf("invalid image %q", e.Image)
}

func initFlags(flag *pflag.FlagSet) {
	flag.String("aws-region", "us-west-2", "The AWS Region")
	flag.String("aws-profile", "", "The aws-vault profile")
	flag.String("aws-vault-keychain-name", "", "The aws-vault keychain name")
	flag.String("service", "", fmt.Sprintf("The service name (choose %q)", services))
	flag.String("environment", "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String("template", "", "The name of the template to use for rendering the task definition")
	flag.String("image", "", "The name of the image referenced in the task definition")
	flag.BoolP("verbose", "v", false, "Print section lines")
}

func checkConfig(v *viper.Viper) error {

	regions, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.EcsServiceID)
	if !ok {
		return fmt.Errorf("could not find regions for service %q", endpoints.EcsServiceID)
	}

	region := v.GetString("aws-region")
	if len(region) == 0 {
		return errors.Wrap(&errInvalidRegion{Region: region}, fmt.Sprintf("%q is invalid", "aws-region"))
	}

	if _, ok := regions[region]; !ok {
		return errors.Wrap(&errInvalidRegion{Region: region}, fmt.Sprintf("%q is invalid", "aws-region"))
	}

	serviceName := v.GetString("service")
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

	environment := v.GetString("environment")
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

	template := v.GetString("template")
	if len(template) == 0 {
		return errors.Wrap(&errInvalidTemplate{Template: template}, fmt.Sprintf("%q is invalid", "template"))
	}
	// Confirm file exists
	fileInfo, err := os.Stat(template)
	if err != nil {
		return errors.Wrap(&errInvalidTemplate{Template: template}, fmt.Sprintf("%q file does not exist", "template"))
	}
	if fileInfo.IsDir() {
		return errors.Wrap(&errInvalidTemplate{Template: template}, fmt.Sprintf("%q is a directory, not a template file", "template"))
	}
	if fileInfo.Size() == 0 {
		return errors.Wrap(&errInvalidTemplate{Template: template}, fmt.Sprintf("%q is an empty template file", "template"))
	}

	image := v.GetString("image")
	if len(image) == 0 {
		return errors.Wrap(&errInvalidImage{Image: image}, fmt.Sprintf("%q is invalid", "image"))
	}

	return nil
}

func quit(logger *log.Logger, flag *pflag.FlagSet, err error) {
	logger.Println(err.Error())
	fmt.Println("Usage of ecs-service-logs:")
	flag.PrintDefaults()
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

	if !v.GetBool("verbose") {
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

	awsRegion := v.GetString("aws-region")

	awsConfig := &aws.Config{
		Region: aws.String(awsRegion),
	}

	verbose := v.GetBool("verbose")
	keychainName := v.GetString("aws-vault-keychain-name")
	keychainProfile := v.GetString("aws-profile")

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

	serviceECS := ecs.New(sess)
	fmt.Println(serviceECS)
}
