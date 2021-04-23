package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
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

const (
	binMilMove       string = "/bin/milmove"
	binMilMoveTasks  string = "/bin/milmove-tasks"
	binOrders        string = "/bin/orders"
	binWebhookClient string = "/bin/webhook-client"
	digestSeparator  string = "@"
	tagSeparator     string = ":"
)

// Valid services names
var services = []string{
	"app",
	"app-client-tls",
	"app-migrations",
	"app-tasks",
	"app-webhook-client",
	"orders",
	"orders-migrations",
}

// Services mapped to Entry Points
// This prevents using an illegal entry point against a service
var servicesToEntryPoints = map[string][]string{
	"app":            {fmt.Sprintf("%s serve", binMilMove)},
	"app-client-tls": {fmt.Sprintf("%s serve", binMilMove)},
	"app-migrations": {fmt.Sprintf("%s migrate", binMilMove)},
	"app-tasks": {
		fmt.Sprintf("%s connect-to-gex-via-sftp", binMilMoveTasks),
		fmt.Sprintf("%s post-file-to-gex", binMilMoveTasks),
		fmt.Sprintf("%s process-edis", binMilMoveTasks),
		fmt.Sprintf("%s save-ghc-fuel-price-data", binMilMoveTasks),
		fmt.Sprintf("%s send-payment-reminder", binMilMoveTasks),
		fmt.Sprintf("%s send-post-move-survey", binMilMoveTasks),
	},
	"app-webhook-client": {
		fmt.Sprintf("%s webhook-notify", binWebhookClient),
	},
	"orders":            {fmt.Sprintf("%s serve", binOrders)},
	"orders-migrations": {fmt.Sprintf("%s migrate", binOrders)},
}

// Services mapped to App Ports
// This ensures app ports are correct for a service that requires port mappings
var servicesToAppPorts = map[string]int64{
	"app":            int64(8443),
	"app-client-tls": int64(9443),
	"orders":         int64(9443),
}

type errInvalidService struct {
	Service string
}

func (e *errInvalidService) Error() string {
	return fmt.Sprintf("invalid AWS ECS service %q, expecting one of %q", e.Service, services)
}

type errInvalidImage struct {
	Image string
}

func (e *errInvalidImage) Error() string {
	return fmt.Sprintf("invalid AWS ECR image tag %q", e.Image)
}

type errInvalidEntryPoint struct {
	EntryPoint string
}

func (e *errInvalidEntryPoint) Error() string {
	return fmt.Sprintf("invalid entry point in the /bin folder %q", e.EntryPoint)
}

type errInvalidFile struct {
	File string
}

func (e *errInvalidFile) Error() string {
	return fmt.Sprintf("invalid file path %q", e.File)
}

const (
	serviceFlag       string = "service"
	imageURIFlag      string = "image"
	variablesFileFlag string = "variables-file"
	entryPointFlag    string = "entrypoint"
	cpuFlag           string = "cpu"
	memFlag           string = "memory"
	registerFlag      string = "register"
)

// ECRImage represents an ECR Image tag broken into its constituent parts
type ECRImage struct {
	AWSRegion        string
	Digest           string
	ImageURIByDigest *string
	ImageURIByTag    *string
	ImageURI         string
	RegistryID       string
	RepositoryURI    string
	RepositoryName   string
	Tag              string
}

// NewECRImage returns a new ECR Image object
func NewECRImage(imageURI string) (*ECRImage, error) {
	var digestURI, tagURI *string
	digest, tag := "", ""
	var imageParts []string

	if strings.Contains(imageURI, digestSeparator) {
		digestURI = &imageURI
		imageParts = strings.Split(imageURI, digestSeparator)
		digest = imageParts[1]
	} else if strings.Contains(imageURI, tagSeparator) {
		tagURI = &imageURI
		imageParts = strings.Split(imageURI, tagSeparator)
		tag = imageParts[1]
	} else {
		return nil, fmt.Errorf("invalid URI, requires either a @digest or a :tag in the URI %v", imageURI)
	}

	if len(imageParts) != 2 {
		return nil, fmt.Errorf("image URI, url parsing failed: %v", imageURI)
	}
	repositoryURI := imageParts[0]
	repositoryURIParts := strings.Split(repositoryURI, "/")
	repositoryName := repositoryURIParts[1]
	repositoryDomainParts := strings.Split(repositoryURIParts[0], ".")
	registryID, awsRegion := repositoryDomainParts[0], repositoryDomainParts[3]

	return &ECRImage{
		AWSRegion:        awsRegion,
		Digest:           digest,
		ImageURI:         imageURI,
		ImageURIByTag:    tagURI,
		ImageURIByDigest: digestURI,
		RegistryID:       registryID,
		RepositoryURI:    repositoryURI,
		RepositoryName:   repositoryName,
		Tag:              tag,
	}, nil
}

// Validate checks ecr image struct values by running validate method and making a request to aws service
func (ecrImage ECRImage) Validate(serviceECR *ecr.ECR) error {
	imageIdentifier := ecr.ImageIdentifier{}
	if ecrImage.ImageURIByDigest != nil {
		imageIdentifier.SetImageDigest(ecrImage.Digest)
	} else if ecrImage.ImageURIByTag != nil {
		imageIdentifier.SetImageTag(ecrImage.Tag)
	} else {
		return fmt.Errorf("no valid imageuri, ImageURIByTag and ImageURIByDigest are null in ecrImage: %v", ecrImage)
	}

	//check to make sure image can validate
	errImageIdentifierValidate := imageIdentifier.Validate()
	if errImageIdentifierValidate != nil {
		return fmt.Errorf("image identifier invalid %w", errImageIdentifierValidate)
	}

	//check to make sure image exists
	imageList, describeImageErr := serviceECR.DescribeImages(&ecr.DescribeImagesInput{
		ImageIds:       append([]*ecr.ImageIdentifier{}, &imageIdentifier),
		RegistryId:     aws.String(ecrImage.RegistryID),
		RepositoryName: aws.String(ecrImage.RepositoryName),
	})

	if describeImageErr != nil {
		return fmt.Errorf("unable to retrieve image: %v: Error: %w", ecrImage, describeImageErr)
	}
	if len(imageList.ImageDetails) < 1 {
		return fmt.Errorf("no images found %v", ecrImage)
	}
	return nil

}

func initTaskDefFlags(flag *pflag.FlagSet) {

	// AWS Account
	flag.String(awsAccountIDFlag, "", "The AWS Account ID")

	// AWS Flags
	cli.InitAWSFlags(flag)

	// Task Definition Settings
	flag.String(serviceFlag, "app", fmt.Sprintf("The service name (choose %q)", services))
	flag.String(environmentFlag, "", fmt.Sprintf("The environment name (choose %q)", environments))
	flag.String(imageURIFlag, "", "The URI of the container image to use in the task definition")
	flag.String(variablesFileFlag, "", "A file containing variables for the task definiton")
	flag.String(entryPointFlag, fmt.Sprintf("%s serve", binMilMove), "The entryPoint for the container")
	flag.Int(cpuFlag, int(512), "The CPU reservation")
	flag.Int(memFlag, int(2048), "The memory reservation")

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Dry Run or Registration
	flag.Bool(dryRunFlag, false, "Execute as a dry-run without modifying AWS.")
	flag.Bool(registerFlag, false, "Execute and register task defintion in AWS.")

	// Sort flags
	flag.SortFlags = true
}

func checkTaskDefConfig(v *viper.Viper) error {

	awsAccountID := v.GetString(awsAccountIDFlag)
	if len(awsAccountID) == 0 {
		return fmt.Errorf("%q is invalid: %w", awsAccountIDFlag, &errInvalidAccountID{AwsAccountID: awsAccountID})
	}

	region, err := cli.CheckAWSRegion(v)
	if err != nil {
		return fmt.Errorf("%q is invalid: %w", cli.AWSRegionFlag, err)
	}

	if err := cli.CheckAWSRegionForService(region, cloudwatchevents.ServiceName); err != nil {
		return fmt.Errorf("%q is invalid for service %s: %w", cli.AWSRegionFlag, cloudwatchevents.ServiceName, err)
	}

	if err := cli.CheckAWSRegionForService(region, ecs.ServiceName); err != nil {
		return fmt.Errorf("%q is invalid for service %s: %w", cli.AWSRegionFlag, ecs.ServiceName, err)
	}

	if err := cli.CheckAWSRegionForService(region, ecr.ServiceName); err != nil {
		return fmt.Errorf("%q is invalid for service %s: %w", cli.AWSRegionFlag, ecr.ServiceName, err)
	}

	if err := cli.CheckAWSRegionForService(region, rds.ServiceName); err != nil {
		return fmt.Errorf("%q is invalid for service %s: %w", cli.AWSRegionFlag, rds.ServiceName, err)
	}

	if err := cli.CheckAWSRegionForService(region, ssm.ServiceName); err != nil {
		return fmt.Errorf("%q is invalid for service %s: %w", cli.AWSRegionFlag, ssm.ServiceName, err)
	}

	serviceName := v.GetString(serviceFlag)
	if len(serviceName) == 0 {
		return fmt.Errorf("%q is invalid: %w", serviceFlag, &errInvalidService{Service: serviceName})
	}
	validService := false
	for _, str := range services {
		if serviceName == str {
			validService = true
			break
		}
	}
	if !validService {
		return fmt.Errorf("%q is invalid: %w", serviceFlag, &errInvalidService{Service: serviceName})
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

	image := v.GetString(imageURIFlag)
	if len(image) == 0 {
		return fmt.Errorf("%q is invalid: %w", imageURIFlag, &errInvalidImage{Image: image})
	}

	if variablesFile := v.GetString(variablesFileFlag); len(variablesFile) > 0 {
		if _, err := os.Stat(variablesFile); err != nil {
			return fmt.Errorf("%q is invalid: %w", variablesFileFlag, &errInvalidFile{File: variablesFile})
		}
	}

	entryPoint := v.GetString(entryPointFlag)
	if len(entryPointFlag) == 0 {
		return fmt.Errorf("%q is invalid: %w", entryPointFlag, &errInvalidEntryPoint{EntryPoint: entryPoint})
	}
	validEntryPoint := false
	entryPoints := servicesToEntryPoints[serviceName]
	for _, str := range entryPoints {
		if entryPoint == str {
			validEntryPoint = true
			break
		}
	}
	if !validEntryPoint {
		return fmt.Errorf("%q is invalid: %w", entryPointFlag, &errInvalidEntryPoint{EntryPoint: entryPoint})
	}

	return nil
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

	partition, _ := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), awsRegion)

	for p.Next() {
		page := p.Page().(*ssm.DescribeParametersOutput)

		for _, parameter := range page.Parameters {
			if strings.HasPrefix(*parameter.Name, fmt.Sprintf("/%s-%s", serviceName, environmentName)) {
				parameterARN := arn.ARN{
					Partition: partition.ID(),
					Service:   "ssm",
					Region:    awsRegion,
					AccountID: awsAccountID,
					Resource:  fmt.Sprintf("parameter%s", *parameter.Name),
				}
				secrets = append(secrets, &ecs.Secret{
					Name:      aws.String(strings.ToUpper(strings.Split(*parameter.Name, "/")[2])),
					ValueFrom: aws.String(parameterARN.String()),
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
		vars, readFileErr := ioutil.ReadFile(filepath.Clean(variablesFile))
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
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Ensure the configuration works against the variables
	checkConfigErr := checkTaskDefConfig(v)
	if checkConfigErr != nil {
		quit(logger, flag, checkConfigErr)
	}

	awsConfig := createAwsConfig(v.GetString(cli.AWSRegionFlag))
	sess, err := awssession.NewSession(awsConfig)
	if err != nil {
		quit(logger, nil, fmt.Errorf("failed to create AWS session: %w", err))
	}

	// Create the Services
	serviceCloudWatchEvents := cloudwatchevents.New(sess)
	serviceECS := ecs.New(sess)
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

	// Short service name needed for RDS, CloudWatch Logs, and SSM
	serviceNameParts := strings.Split(serviceName, "-")
	serviceNameShort := serviceNameParts[0]

	// Confirm the image exists
	ecrImage, errECRImage := NewECRImage(imageURI)
	if errECRImage != nil {
		quit(logger, nil, fmt.Errorf("unable to recognize image URI %q: %w", imageURI, errECRImage))
	}

	errValidateImage := ecrImage.Validate(serviceECR)
	if errValidateImage != nil {
		quit(logger, nil, fmt.Errorf("unable to validate image %v: %w", ecrImage, errValidateImage))
	}

	// Entrypoint
	entryPoint := v.GetString(entryPointFlag)
	entryPointList := strings.Split(entryPoint, " ")
	commandName := entryPointList[0]
	subCommandName := entryPointList[1]

	// Register the new task definition
	executionRoleArn := fmt.Sprintf("ecs-task-execution-role-%s-%s", serviceName, environmentName)
	taskRoleArn := fmt.Sprintf("ecs-task-role-%s-%s", serviceName, environmentName)
	family := fmt.Sprintf("%s-%s", serviceName, environmentName)

	// handle entrypoint specific logic
	var awsLogsStreamPrefix string
	var awsLogsGroup string
	var portMappings []*ecs.PortMapping
	var containerDefName string
	if commandName == binMilMoveTasks {
		executionRoleArn = fmt.Sprintf("ecs-task-exec-role-%s-%s-%s", serviceNameShort, environmentName, subCommandName)
		taskRoleArn = fmt.Sprintf("ecs-task-role-%s-%s-%s", serviceNameShort, environmentName, subCommandName)
		family = fmt.Sprintf("%s-%s-%s", serviceNameShort, environmentName, subCommandName)

		awsLogsStreamPrefix = serviceName
		awsLogsGroup = fmt.Sprintf("ecs-tasks-%s-%s", serviceNameShort, environmentName)
		containerDefName = fmt.Sprintf("%s-%s-%s", serviceName, subCommandName, environmentName)

		ruleName := fmt.Sprintf("%s-%s", subCommandName, environmentName)
		_, listTargetsByRuleErr := serviceCloudWatchEvents.ListTargetsByRule(&cloudwatchevents.ListTargetsByRuleInput{
			Rule: aws.String(ruleName),
		})
		if listTargetsByRuleErr != nil {
			quit(logger, nil, fmt.Errorf("error retrieving targets for rule %q: %w", ruleName, listTargetsByRuleErr))
		}
	} else if subCommandName == "migrate" {
		awsLogsStreamPrefix = serviceName
		awsLogsGroup = fmt.Sprintf("ecs-tasks-%s-%s", serviceNameShort, environmentName)
		containerDefName = fmt.Sprintf("%s-%s", serviceName, environmentName)

		// TODO: The execution role needs to be split from the app
		// This needs to be fixed in terraform
		executionRoleArn = fmt.Sprintf("ecs-task-execution-role-%s-%s", serviceNameShort, environmentName)
		// TODO: The task role is missing an (s) so we can't use service name
		// This is `ecs-task-role-app-migration-experimental` vs `ecs-task-role-app-migration(s)-experimental`
		// This needs to be fixed in terraform and then rolled out
		taskRoleArn = fmt.Sprintf("ecs-task-role-%s-migration-%s", serviceNameShort, environmentName)
	} else if commandName == binWebhookClient {
		awsLogsStreamPrefix = serviceName
		awsLogsGroup = fmt.Sprintf("ecs-tasks-%s-%s", serviceName, environmentName)
		containerDefName = fmt.Sprintf("%s-%s", serviceName, environmentName)
	} else {
		awsLogsStreamPrefix = serviceNameShort
		awsLogsGroup = fmt.Sprintf("ecs-tasks-%s-%s", serviceName, environmentName)
		containerDefName = fmt.Sprintf("%s-%s", serviceName, environmentName)

		// Ports
		port := servicesToAppPorts[serviceName]
		portMappings = []*ecs.PortMapping{
			{
				ContainerPort: aws.Int64(port),
				HostPort:      aws.Int64(port),
				Protocol:      aws.String("tcp"),
			},
		}
	}

	// Get the database host using the instance identifier
	dbInstanceIdentifier := fmt.Sprintf("%s-%s", serviceNameShort, environmentName)
	dbInstancesOutput, err := serviceRDS.DescribeDBInstances(&rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbInstanceIdentifier),
	})
	if err != nil {
		quit(logger, nil, fmt.Errorf("error retrieving database definition for %q: %w", dbInstanceIdentifier, err))
	}
	dbHost := *dbInstancesOutput.DBInstances[0].Endpoint.Address

	// CPU / MEM
	cpu := strconv.Itoa(v.GetInt(cpuFlag))
	mem := strconv.Itoa(v.GetInt(memFlag))

	// Create the set of secrets and environment variables that will be injected into the
	// container.
	secrets := buildSecrets(serviceSSM, awsRegion, awsAccountID, serviceNameShort, environmentName)
	containerEnvironment := buildContainerEnvironment(environmentName, dbHost, variablesFile)

	// AWS does not permit supplying both a secret and an environment variable that share the same
	// name into an ECS task. In order to gracefully transition between setting values as secrets
	// into setting them as environment variables, this function serves to remove any duplicates
	// that have been transitioned into being set as environment variables.
	secrets = removeSecretsWithMatchingEnvironmentVariables(secrets, containerEnvironment)

	newTaskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:        aws.String(containerDefName),
				Image:       aws.String(ecrImage.ImageURI),
				Essential:   aws.Bool(true),
				EntryPoint:  aws.StringSlice(entryPointList),
				Command:     []*string{},
				Secrets:     secrets,
				Environment: containerEnvironment,
				LogConfiguration: &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: map[string]*string{
						"awslogs-group":         aws.String(awsLogsGroup),
						"awslogs-region":        aws.String(awsRegion),
						"awslogs-stream-prefix": aws.String(awsLogsStreamPrefix),
					},
				},
				PortMappings:           portMappings,
				ReadonlyRootFilesystem: aws.Bool(true),
				Privileged:             aws.Bool(false),
				User:                   aws.String("1042"),
			},
		},
		Cpu:                     aws.String(cpu),
		ExecutionRoleArn:        aws.String(executionRoleArn),
		Family:                  aws.String(family),
		Memory:                  aws.String(mem),
		NetworkMode:             aws.String("awsvpc"),
		RequiresCompatibilities: []*string{aws.String("FARGATE")},
		TaskRoleArn:             aws.String(taskRoleArn),
	}

	// Registration is never allowed by default and requires a flag
	if v.GetBool(dryRunFlag) {
		// Format the new task def as JSON for viewing
		newTaskDefJSON, jsonErr := jsonutil.BuildJSON(newTaskDefInput)
		if jsonErr != nil {
			quit(logger, nil, err)
		}

		logger.Println(string(newTaskDefJSON))
	} else if v.GetBool(registerFlag) {
		// Register the new task definition
		newTaskDefOutput, err := serviceECS.RegisterTaskDefinition(&newTaskDefInput)
		if err != nil {
			quit(logger, nil, fmt.Errorf("error registering new task definition: %w", err))
		}
		newTaskDefArn := *newTaskDefOutput.TaskDefinition.TaskDefinitionArn
		logger.Println(newTaskDefArn)
	} else {
		quit(logger, flag, errors.New(fmt.Sprintf("Please provide either %q or %q flags when running", dryRunFlag, registerFlag)))
	}

	return nil
}

func createAwsConfig(awsRegionFlag string) *aws.Config {
	awsConfig := &aws.Config{
		Region: aws.String(awsRegionFlag),
	}
	return awsConfig
}

func removeSecretsWithMatchingEnvironmentVariables(secrets []*ecs.Secret, containerEnvironment []*ecs.KeyValuePair) []*ecs.Secret {
	// Remove any secrets that share a name with an environment variable. Do this by creating a new
	// slice of secrets that does not any secrets that share a name with an environment variable.
	newSecrets := []*ecs.Secret{}
	for _, secret := range secrets {
		conflictFound := false
		for _, envSetting := range containerEnvironment {
			if *secret.Name == *envSetting.Name {
				conflictFound = true
			}
		}

		if conflictFound {
			// Report any conflicts that are found.
			fmt.Fprintln(os.Stderr, "Found a secret with the same name as an environment variable. Discarding secret in favor of the environment variable:", *secret.Name)
		} else {
			// If no conflict is found, keep the secret.
			newSecrets = append(newSecrets, secret)
		}
	}

	return newSecrets
}
