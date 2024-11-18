package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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
		fmt.Sprintf("%s process-tpps", binMilMoveTasks),
		fmt.Sprintf("%s save-ghc-fuel-price-data", binMilMoveTasks),
		fmt.Sprintf("%s send-payment-reminder", binMilMoveTasks),
	},
	"app-webhook-client": {
		fmt.Sprintf("%s webhook-notify", binWebhookClient),
	},
	"orders":            {fmt.Sprintf("%s serve", binOrders)},
	"orders-migrations": {fmt.Sprintf("%s migrate", binOrders)},
}

// Services mapped to App Ports
// This ensures app ports are correct for a service that requires port mappings
var servicesToAppPorts = map[string]int32{
	"app":            int32(8443),
	"app-client-tls": int32(9443),
	"orders":         int32(9443),
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
	serviceFlag              string = "service"
	imageURIFlag             string = "image"
	variablesFileFlag        string = "variables-file"
	entryPointFlag           string = "entrypoint"
	cpuFlag                  string = "cpu"
	memFlag                  string = "memory"
	registerFlag             string = "register"
	openTelemetrySidecarFlag string = "open-telemetry-sidecar"
	otelCollectorImageFlag   string = "otel-collector-image"
	healthCheckFlag          string = "health-check"
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
func (ecrImage ECRImage) Validate(ecrClient *ecr.Client) error {
	imageIdentifier := ecrtypes.ImageIdentifier{}
	if ecrImage.ImageURIByDigest != nil {
		imageIdentifier.ImageDigest = &ecrImage.Digest
	} else if ecrImage.ImageURIByTag != nil {
		imageIdentifier.ImageTag = &ecrImage.Tag
	} else {
		return fmt.Errorf("no valid imageuri, ImageURIByTag and ImageURIByDigest are null in ecrImage: %v", ecrImage)
	}

	//check to make sure image exists
	imageList, describeImageErr := ecrClient.DescribeImages(
		context.Background(),
		&ecr.DescribeImagesInput{
			ImageIds:       []ecrtypes.ImageIdentifier{imageIdentifier},
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

	// Open Telemetry SideCar
	flag.Bool(openTelemetrySidecarFlag, false, "Include open telemetry sidecar container")
	const defaultOtelImage = "public.ecr.aws/aws-observability/aws-otel-collector:v0.29.0"
	flag.String(otelCollectorImageFlag, defaultOtelImage,
		"Image to use for open telemetry sidecar")

	// Health Check
	flag.Bool(healthCheckFlag, false, "Include health check in the task definition")

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

	_, err := cli.CheckAWSRegion(v)
	if err != nil {
		return fmt.Errorf("%q is invalid: %w", cli.AWSRegionFlag, err)
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

func buildSecrets(cfg aws.Config, awsAccountID, serviceName, environmentName string) ([]ecstypes.Secret, error) {

	var secrets []ecstypes.Secret

	ssmClient := ssm.NewFromConfig(cfg)

	resolver := ssm.NewDefaultEndpointResolver()
	endpoint, err := resolver.ResolveEndpoint(cfg.Region,
		ssm.EndpointResolverOptions{})
	if err != nil {
		return nil, err
	}
	partition := endpoint.PartitionID

	ctx := context.Background()

	paginator := ssm.NewDescribeParametersPaginator(ssmClient,
		&ssm.DescribeParametersInput{},
		func(opts *ssm.DescribeParametersPaginatorOptions) {
			opts.Limit = 50
		})

	servicePrefix := fmt.Sprintf("/%s-%s", serviceName, environmentName)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return secrets, err
		}

		for _, parameter := range page.Parameters {
			if strings.HasPrefix(*parameter.Name, servicePrefix) {
				parameterARN := arn.ARN{
					Partition: partition,
					Service:   "ssm",
					Region:    cfg.Region,
					AccountID: awsAccountID,
					Resource:  fmt.Sprintf("parameter%s", *parameter.Name),
				}
				secrets = append(secrets, ecstypes.Secret{
					Name:      aws.String(strings.ToUpper(strings.Split(*parameter.Name, "/")[2])),
					ValueFrom: aws.String(parameterARN.String()),
				})
			}
		}
	}

	return secrets, nil
}

func buildContainerEnvironment(environmentName string, dbHost string, variablesFile string) []ecstypes.KeyValuePair {

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
			log.Fatal(fmt.Errorf("file %q does not exist: %w", variablesFile, err))
		}
		// Read contents of variables file into vars
		vars, readFileErr := os.ReadFile(filepath.Clean(variablesFile))
		if readFileErr != nil {
			log.Fatal(fmt.Errorf("error reading variables file"))
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

	var ecsKVPair []ecstypes.KeyValuePair

	// Sort these for easier reading
	keys := make([]string, 0, len(envVars))
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		ecsKVPair = append(ecsKVPair, ecstypes.KeyValuePair{
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
		log.SetOutput(io.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Ensure the configuration works against the variables
	err = checkTaskDefConfig(v)
	if err != nil {
		quit(logger, flag, err)
	}

	cfg, errCfg := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(v.GetString(cli.AWSRegionFlag)),
	)
	if errCfg != nil {
		quit(logger, flag, err)
	}

	serviceCloudWatchEvents := cloudwatchevents.NewFromConfig(cfg)
	serviceECS := ecs.NewFromConfig(cfg)
	serviceECR := ecr.NewFromConfig(cfg)
	serviceRDS := rds.NewFromConfig(cfg)

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
	var portMappings []ecstypes.PortMapping
	var containerDefName string

	ctx := context.Background()

	if commandName == binMilMoveTasks {
		executionRoleArn = fmt.Sprintf("ecs-task-exec-role-%s-%s-%s", serviceNameShort, environmentName, subCommandName)
		taskRoleArn = fmt.Sprintf("ecs-task-role-%s-%s-%s", serviceNameShort, environmentName, subCommandName)
		family = fmt.Sprintf("%s-%s-%s", serviceNameShort, environmentName, subCommandName)

		awsLogsStreamPrefix = serviceName
		awsLogsGroup = fmt.Sprintf("ecs-tasks-%s-%s", serviceNameShort, environmentName)
		containerDefName = fmt.Sprintf("%s-%s-%s", serviceName, subCommandName, environmentName)

		ruleName := fmt.Sprintf("%s-%s", subCommandName, environmentName)
		_, listTargetsByRuleErr := serviceCloudWatchEvents.ListTargetsByRule(
			ctx,
			&cloudwatchevents.ListTargetsByRuleInput{
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
		portMappings = []ecstypes.PortMapping{
			{
				ContainerPort: aws.Int32(port),
				HostPort:      aws.Int32(port),
				Protocol:      ecstypes.TransportProtocolTcp,
			},
		}
	}

	// Get the database host using the instance identifier
	dbInstanceIdentifier := fmt.Sprintf("%s-%s", serviceNameShort, environmentName)
	dbInstancesOutput, err := serviceRDS.DescribeDBInstances(
		ctx,
		&rds.DescribeDBInstancesInput{
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
	secrets, err := buildSecrets(cfg, awsAccountID, serviceNameShort, environmentName)
	if err != nil {
		quit(logger, nil, err)
	}
	containerEnvironment := buildContainerEnvironment(environmentName, dbHost, variablesFile)

	// AWS does not permit supplying both a secret and an environment variable that share the same
	// name into an ECS task. In order to gracefully transition between setting values as secrets
	// into setting them as environment variables, this function serves to remove any duplicates
	// that have been transitioned into being set as environment variables.
	secrets = removeSecretsWithMatchingEnvironmentVariables(secrets, containerEnvironment)

	containerDefinitions := []ecstypes.ContainerDefinition{
		{
			Name:        aws.String(containerDefName),
			Image:       aws.String(ecrImage.ImageURI),
			Essential:   aws.Bool(true),
			EntryPoint:  entryPointList,
			Command:     []string{},
			Secrets:     secrets,
			Environment: containerEnvironment,
			Ulimits: []ecstypes.Ulimit{
				{
					Name:      ecstypes.UlimitName("nofile"),
					SoftLimit: 10000,
					HardLimit: 10000,
				},
			},
			LogConfiguration: &ecstypes.LogConfiguration{
				LogDriver: ecstypes.LogDriverAwslogs,
				Options: map[string]string{
					"awslogs-group":         awsLogsGroup,
					"awslogs-region":        awsRegion,
					"awslogs-stream-prefix": awsLogsStreamPrefix,
				},
			},
			PortMappings:           portMappings,
			ReadonlyRootFilesystem: aws.Bool(true),
			Privileged:             aws.Bool(false),
			User:                   aws.String("1042"),
		},
	}

	// if health check is enabled, add it to the container definition
	if v.GetBool(healthCheckFlag) {
		containerDefinitions[0].HealthCheck = &ecstypes.HealthCheck{
			Command: []string{
				"CMD",
				binMilMove,
				"health",
			},
			// Interval defaults to 30 seconds
			// Retries defaults to 3
			// Timeout defaults to 5 seconds
			//
			// StartPeriod is a grace period when the app starts, it
			// defaults to off and can be between 5 - 300 seconds
		}
	}

	// do not enable otel sidecar for the webhook service
	isOtelEnabledService := !strings.Contains(containerDefName, "webhook")
	if v.GetBool(openTelemetrySidecarFlag) && isOtelEnabledService {
		// put our custom config file in the AOT_CONFIG_CONTENT
		// environment variable
		//
		// ideas from the suggested
		// https://github.com/aws-observability/aws-otel-collector/blob/main/config/ecs/container-insights/otel-task-metrics-config.yaml
		// and then container specific info from
		// https://aws-otel.github.io/docs/components/ecs-metrics-receiver#full-configuration-for-task--and-container-level-metrics
		//
		aotConfigContent := `
extensions:
  health_check:

receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 127.0.0.1:4317
      http:
        endpoint: 127.0.0.1:4318
  awsxray:
    endpoint: 127.0.0.1:2000
    transport: udp
  statsd:
    endpoint: 127.0.0.1:8125
    aggregation_interval: 60s
  awsecscontainermetrics:

processors:
  batch/traces:
    timeout: 1s
    send_batch_size: 50
  batch/metrics:
    timeout: 60s
  filter:
    metrics:
      include:
        match_type: strict
        metric_names:
          - .*memory.reserved
          - .*memory.utilized
          - .*cpu.reserved
          - .*cpu.utilized
          - .*network.rate.rx
          - .*network.rate.tx
          - .*storage.read_bytes
          - .*storage.write_bytes
          - container.duration
  metricstransform:
    transforms:
      - include: ecs.task.memory.utilized
        action: update
        new_name: MemoryUtilized
      - include: ecs.task.memory.reserved
        action: update
        new_name: MemoryReserved
      - include: ecs.task.cpu.utilized
        action: update
        new_name: CpuUtilized
      - include: ecs.task.cpu.reserved
        action: update
        new_name: CpuReserved
      - include: ecs.task.network.rate.rx
        action: update
        new_name: NetworkRxBytes
      - include: ecs.task.network.rate.tx
        action: update
        new_name: NetworkTxBytes
      - include: ecs.task.storage.read_bytes
        action: update
        new_name: StorageReadBytes
      - include: ecs.task.storage.write_bytes
        action: update
        new_name: StorageWriteBytes

  resource:
    attributes:
      - key: ClusterName
        from_attribute: aws.ecs.cluster.name
        action: insert
      - key: aws.ecs.cluster.name
        action: delete
      - key: ServiceName
        from_attribute: aws.ecs.service.name
        action: insert
      - key: aws.ecs.service.name
        action: delete
      - key: TaskId
        from_attribute: aws.ecs.task.id
        action: insert
      - key: aws.ecs.task.id
        action: delete
      - key: TaskDefinitionFamily
        from_attribute: aws.ecs.task.family
        action: insert
      - key: aws.ecs.task.family
        action: delete
      - key: TaskARN
        from_attribute: aws.ecs.task.arn
        action: insert
      - key: aws.ecs.task.arn
        action: delete
      - key: DockerName
        from_attribute: aws.ecs.docker.name
        action: insert
      - key: aws.ecs.docker.name
        action: delete
      - key: TaskDefinitionRevision
        from_attribute: aws.ecs.task.version
        action: insert
      - key: aws.ecs.task.version
        action: delete
      - key: PullStartedAt
        from_attribute: aws.ecs.task.pull_started_at
        action: insert
      - key: aws.ecs.task.pull_started_at
        action: delete
      - key: PullStoppedAt
        from_attribute: aws.ecs.task.pull_stopped_at
        action: insert
      - key: aws.ecs.task.pull_stopped_at
        action: delete
      - key: AvailabilityZone
        from_attribute: cloud.zone
        action: insert
      - key: cloud.zone
        action: delete
      - key: LaunchType
        from_attribute: aws.ecs.task.launch_type
        action: insert
      - key: aws.ecs.task.launch_type
        action: delete
      - key: Region
        from_attribute: cloud.region
        action: insert
      - key: cloud.region
        action: delete
      - key: AccountId
        from_attribute: cloud.account.id
        action: insert
      - key: cloud.account.id
        action: delete
      - key: DockerId
        from_attribute: container.id
        action: insert
      - key: container.id
        action: delete
      - key: ContainerName
        from_attribute: container.name
        action: insert
      - key: container.name
        action: delete
      - key: Image
        from_attribute: container.image.name
        action: insert
      - key: container.image.name
        action: delete
      - key: ImageId
        from_attribute: aws.ecs.container.image.id
        action: insert
      - key: aws.ecs.container.image.id
        action: delete
      - key: ExitCode
        from_attribute: aws.ecs.container.exit_code
        action: insert
      - key: aws.ecs.container.exit_code
        action: delete
      - key: CreatedAt
        from_attribute: aws.ecs.container.created_at
        action: insert
      - key: aws.ecs.container.created_at
        action: delete
      - key: StartedAt
        from_attribute: aws.ecs.container.started_at
        action: insert
      - key: aws.ecs.container.started_at
        action: delete
      - key: FinishedAt
        from_attribute: aws.ecs.container.finished_at
        action: insert
      - key: aws.ecs.container.finished_at
        action: delete
      - key: ImageTag
        from_attribute: container.image.tag
        action: insert
      - key: container.image.tag
        action: delete

exporters:
  awsxray:
    index_all_attributes: true
  awsemf/application:
    namespace: ECS/AWSOTel/Application
    log_group_name: '/aws/ecs/application/metrics'
  awsemf/performance:
    namespace: ECS/ContainerInsights
    log_group_name: '/aws/ecs/containerinsights/{ClusterName}/performance'
    log_stream_name: '{TaskId}'
    resource_to_telemetry_conversion:
      enabled: true
    dimension_rollup_option: NoDimensionRollup
    metric_declarations:
      - dimensions: [ [ ClusterName ], [ ClusterName, TaskDefinitionFamily ] ]
        metric_name_selectors:
          - MemoryUtilized
          - MemoryReserved
          - CpuUtilized
          - CpuReserved
          - NetworkRxBytes
          - NetworkTxBytes
          - StorageReadBytes
          - StorageWriteBytes
      - dimensions: [[ClusterName], [ClusterName, TaskDefinitionFamily, ContainerName]]
        metric_name_selectors: [container.*]

service:
  telemetry:
    logs:
      level: ERROR
  pipelines:
    traces:
      receivers: [otlp,awsxray]
      processors: [batch/traces]
      exporters: [awsxray]
    metrics/application:
      receivers: [otlp, statsd]
      processors: [batch/metrics]
      exporters: [awsemf/application]
    metrics/performance:
      receivers: [awsecscontainermetrics ]
      processors: [filter, metricstransform, resource]
      exporters: [ awsemf/performance ]

  extensions: [health_check]
`

		otelCollectorImage := v.GetString(otelCollectorImageFlag)
		containerDefinitions = append(containerDefinitions,
			ecstypes.ContainerDefinition{
				Name:      aws.String("otel-" + containerDefName),
				Image:     aws.String(otelCollectorImage),
				Essential: aws.Bool(true),
				Environment: []ecstypes.KeyValuePair{
					{
						Name:  aws.String("AOT_CONFIG_CONTENT"),
						Value: aws.String(aotConfigContent),
					},
				},
				LogConfiguration: &ecstypes.LogConfiguration{
					LogDriver: ecstypes.LogDriverAwslogs,
					Options: map[string]string{
						"awslogs-group":         awsLogsGroup,
						"awslogs-region":        awsRegion,
						"awslogs-stream-prefix": "otel-" + awsLogsStreamPrefix,
					},
				},
				HealthCheck: &ecstypes.HealthCheck{
					Command: []string{
						"CMD",
						"/healthcheck",
					},
				},
			},
		)
	}

	newTaskDefInput := ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    containerDefinitions,
		Cpu:                     aws.String(cpu),
		ExecutionRoleArn:        aws.String(executionRoleArn),
		Family:                  aws.String(family),
		Memory:                  aws.String(mem),
		NetworkMode:             ecstypes.NetworkModeAwsvpc,
		RequiresCompatibilities: []ecstypes.Compatibility{"FARGATE"},
		TaskRoleArn:             aws.String(taskRoleArn),
	}

	// Registration is never allowed by default and requires a flag
	if v.GetBool(dryRunFlag) {
		// Format the new task def as JSON for viewing
		jsonErr := json.NewEncoder(logger.Writer()).Encode(newTaskDefInput)
		if jsonErr != nil {
			quit(logger, nil, err)
		}
	} else if v.GetBool(registerFlag) {
		// Register the new task definition
		newTaskDefOutput, err := serviceECS.RegisterTaskDefinition(ctx, &newTaskDefInput)
		if err != nil {
			quit(logger, nil, fmt.Errorf("error registering new task definition: %w", err))
		}
		newTaskDefArn := *newTaskDefOutput.TaskDefinition.TaskDefinitionArn
		logger.Println(newTaskDefArn)
	} else {
		quit(logger, flag, fmt.Errorf("Please provide either %q or %q flags when running", dryRunFlag, registerFlag))
	}

	return nil
}

func removeSecretsWithMatchingEnvironmentVariables(secrets []ecstypes.Secret, containerEnvironment []ecstypes.KeyValuePair) []ecstypes.Secret {
	// Remove any secrets that share a name with an environment variable. Do this by creating a new
	// slice of secrets that does not any secrets that share a name with an environment variable.
	newSecrets := []ecstypes.Secret{}
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
