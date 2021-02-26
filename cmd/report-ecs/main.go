package main

/*
 * Report information about running ECS Clusters and Services
 *
 * For example to get all the current platform versions run:
 *
 *	go run ../cmd/report-ecs/main.go   | jq -r .clusters[].services[].platformVersion | sort | uniq
 *
 */

import (
	"encoding/json"
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

// Service is a struct describing an ECS Service.
type Service struct {
	Arn             string `json:"arn"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	PlatformVersion string `json:"platformVersion"`
	RunningCount    int64  `json:"runningCount"`
	PendingCount    int64  `json:"pendingCount"`
	DesiredCount    int64  `json:"desiredCount"`
	TaskDefinition  string `json:"taskDefinition"`
}

// Cluster is a struct describing an ECS Cluster.
type Cluster struct {
	Arn      string    `json:"arn"`
	Services []Service `json:"services"`
}

// Report is a struct containing the full report on an ECS Service
type Report struct {
	Clusters []Cluster `json:"clusters"`
}

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

func initFlags(flag *pflag.FlagSet) {
	flag.String("aws-region", "us-west-2", "AWS region used inspecting ECS")
	flag.String("aws-profile", "", "The aws-vault profile")
	flag.String("aws-vault-keychain-name", "", "The aws-vault keychain name")
	flag.BoolP("verbose", "v", false, "Show extra output for debugging")
	flag.SortFlags = false
}

func checkRegion(v *viper.Viper) error {

	regions, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.GuarddutyServiceID)
	if !ok {
		return fmt.Errorf("could not find regions for service %s", endpoints.GuarddutyServiceID)
	}

	r := v.GetString("aws-region")
	if len(r) == 0 {
		return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-region"))
	}

	if _, ok := regions[r]; !ok {
		return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-region"))
	}

	return nil
}

func checkConfig(v *viper.Viper) error {

	err := checkRegion(v)
	if err != nil {
		return errors.Wrap(err, "Region check failed")
	}

	return nil
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
	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)
	initFlags(flag)

	parseErr := flag.Parse(os.Args[1:])
	if parseErr != nil {
		logger.Println(parseErr)
		logger.Println("flag parsing failed")
		flag.PrintDefaults()
		os.Exit(1)
	}

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		logger.Println(bindErr)
		logger.Println("flag binding failed")
		flag.PrintDefaults()
		os.Exit(1)
	}

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	if !v.GetBool("verbose") {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Check the config and exit with usage details if there is a problem
	errConfig := checkConfig(v)
	if errConfig != nil {
		logger.Println(errConfig)
		logger.Println("Usage of report-ecs:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	keychainName := v.GetString("aws-vault-keychain-name")
	keychainProfile := v.GetString("aws-profile")
	creds, errCreds := getAWSCredentials(keychainName, keychainProfile)
	if errCreds != nil {
		logger.Fatal(errors.Wrap(errCreds, fmt.Sprintf("Unable to get AWS credentials from the keychain %s and profile %s", keychainName, keychainProfile)))
	}

	// Define services
	sess := awssession.Must(awssession.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(v.GetBool("verbose")),
		Credentials:                   creds,
		Region:                        aws.String(v.GetString("aws-region")),
	}))

	// GuardDuty has the findings
	ecsService := ecs.New(sess)

	clusterArns := make([]*string, 0)
	var nextToken *string
	for {
		listClustersOutput, err := ecsService.ListClusters(&ecs.ListClustersInput{NextToken: nextToken})
		if err != nil {
			logger.Fatal(errors.Wrap(err, "error listing clusters"))
		}
		clusterArns = append(clusterArns, listClustersOutput.ClusterArns...)
		if listClustersOutput.NextToken == nil {
			break
		}
	}

	clusters := make([]Cluster, 0)
	for _, clusterArn := range clusterArns {
		cluster := Cluster{Arn: *clusterArn}
		serviceArns := make([]*string, 0)
		nextToken = nil
		for {
			listServicesOutput, err := ecsService.ListServices(&ecs.ListServicesInput{
				Cluster:   clusterArn,
				NextToken: nextToken,
			})
			if err != nil {
				logger.Fatal(errors.Wrap(err, "error listing services"))
			}
			serviceArns = append(serviceArns, listServicesOutput.ServiceArns...)
			if listServicesOutput.NextToken == nil {
				break
			}
		}
		services := make([]Service, 0)
		for i := 0; i < len(serviceArns); i += 10 {
			describeServicesInput := &ecs.DescribeServicesInput{
				Cluster: clusterArn,
			}
			if i+10 < len(serviceArns) {
				describeServicesInput.Services = serviceArns[i : i+10]
			} else {
				describeServicesInput.Services = serviceArns[i:]
			}
			describeServicesOutput, err := ecsService.DescribeServices(describeServicesInput)
			if err != nil {
				logger.Fatal(errors.Wrap(err, "error listing services"))
			}
			for _, service := range describeServicesOutput.Services {
				services = append(services, Service{
					Arn:             aws.StringValue(service.ServiceArn),
					Name:            aws.StringValue(service.ServiceName),
					Status:          aws.StringValue(service.Status),
					PlatformVersion: aws.StringValue(service.PlatformVersion),
					RunningCount:    aws.Int64Value(service.RunningCount),
					PendingCount:    aws.Int64Value(service.PendingCount),
					DesiredCount:    aws.Int64Value(service.DesiredCount),
					TaskDefinition:  aws.StringValue(service.TaskDefinition),
				})
			}
		}
		cluster.Services = services
		clusters = append(clusters, cluster)
	}

	report := &Report{
		Clusters: clusters,
	}

	b, err := json.Marshal(report)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "error marshalling report"))
	}

	fmt.Println(string(b))
}
