package main

/*
 * Find users that triggered a GuardDuty report
 *
 * This script will look up all GuardDuty findings and for each one will
 * pull out the access key, search for that access key in CloudTrail to find
 * the AssumedRole event.  That event will then provide the ARN of the role
 * which can be looked up in CloudTrail again to find the username of the
 * person that triggered the event.
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/guardduty"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

// FindingDetail captures the data from a guard duty finding
type FindingDetail struct {
	ID             *string `json:"id"`
	CreatedAt      *string `json:"createdAt"`
	AccessKeyID    *string `json:"accessKeyID,omitempty"`
	PrincipalID    *string `json:"principalID,omitempty"`
	AssumedRoleARN *string `json:"assumeRoleARN"`
	Username       *string `json:"username"`
	IPAddress      *string `json:"ipAddress,omitempty"`
	ServiceName    *string `json:"serviceName,omitempty"`
	API            *string `json:"api,omitempty"`
	City           *string `json:"city,omitempty"`
	Country        *string `json:"country,omitempty"`
}

// PrintJSON will log JSON-formatted output
func (fd *FindingDetail) PrintJSON(logger *log.Logger) error {
	fdJSON, err := json.Marshal(fd)
	if err != nil {
		return errors.Wrap(err, "Unable to marshal FindingDetail to JSON")

	}
	logger.Println(string(fdJSON))
	return nil
}

// Print will log plain-text output
func (fd *FindingDetail) Print(logger *log.Logger) {
	template := `
Finding ID:         %s
Finding Created At: %s
Access Key ID:      %s
Principal ID:       %s
Assumed Role ARN:   %s
Username:           %s
IPv4:               %s
Service Name:       %s
API:                %s
City, Country:      %s, %s`

	logger.Println(fmt.Sprintf(template,
		aws.StringValue(fd.ID),
		aws.StringValue(fd.CreatedAt),
		aws.StringValue(fd.AccessKeyID),
		aws.StringValue(fd.PrincipalID),
		aws.StringValue(fd.AssumedRoleARN),
		aws.StringValue(fd.Username),
		aws.StringValue(fd.IPAddress),
		aws.StringValue(fd.ServiceName),
		aws.StringValue(fd.API),
		aws.StringValue(fd.City),
		aws.StringValue(fd.Country),
	))
}

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

type errInvalidOutput struct {
	Output string
}

func (e *errInvalidOutput) Error() string {
	return fmt.Sprintf("invalid output %s", e.Output)
}

type errInvalidComparison struct {
	Comparison string
}

func (e *errInvalidComparison) Error() string {
	return fmt.Sprintf("invalid comparison %s", e.Comparison)
}

func initFlags(flag *pflag.FlagSet) {

	// aws-vault
	cli.InitVaultFlags(flag)

	flag.String("aws-guardduty-region", "us-west-2", "AWS region used inspecting guardduty")
	flag.BoolP("archived", "a", false, "Show archived findings instead of current findings")
	flag.StringP("output", "o", "json", "Whether to print output as 'text' or 'json'")

	// Logging Levels
	cli.InitLoggingFlags(flag)

	flag.SortFlags = false
}

func checkRegion(v *viper.Viper) error {

	regions, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.GuarddutyServiceID)
	if !ok {
		return fmt.Errorf("could not find regions for service %s", endpoints.GuarddutyServiceID)
	}

	r := v.GetString("aws-guardduty-region")
	if len(r) == 0 {
		return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-guardduty-region"))
	}

	if _, ok := regions[r]; !ok {
		return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-guardduty-region"))
	}

	return nil
}

func checkOutput(v *viper.Viper) error {

	outputs := map[string]string{"text": "text", "json": "json"}

	o := v.GetString("output")
	if _, ok := outputs[o]; !ok {
		return errors.Wrap(&errInvalidOutput{Output: o}, fmt.Sprintf("%s is invalid", "output"))
	}

	return nil
}

func checkConfig(v *viper.Viper) error {

	if err := cli.CheckVault(v); err != nil {
		return err
	}

	err := checkRegion(v)
	if err != nil {
		return errors.Wrap(err, "Region check failed")
	}

	err = checkOutput(v)
	if err != nil {
		return errors.Wrap(err, "Output check failed")
	}

	if err := cli.CheckLogging(v); err != nil {
		return err
	}

	return nil
}

// LookupEvent searches CloudTrail for event smatching a key-value pair
func LookupEvent(key *string, value *string, serviceCloudTrail *cloudtrail.CloudTrail) (*cloudtrail.Event, error) {
	lookupAttribute := cloudtrail.LookupAttribute{
		AttributeKey:   key,
		AttributeValue: value,
	}
	maxResults := int64(1)
	lookupEventsInput := cloudtrail.LookupEventsInput{
		LookupAttributes: []*cloudtrail.LookupAttribute{&lookupAttribute},
		MaxResults:       &maxResults,
	}
	events, err := serviceCloudTrail.LookupEvents(&lookupEventsInput)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("LookupEvents failed with Attribute Key '%s' and Attribute Value '%s'", *key, *value))
	}
	if len(events.Events) != 1 {
		return nil, fmt.Errorf("Expected exactly one event, got %d", len(events.Events))
	}
	return events.Events[0], nil
}

// GetRoleAndUser tries to use an access key or principal id to find a role arn and username
func GetRoleAndUser(key *string, value *string, serviceCloudTrail *cloudtrail.CloudTrail) (*string, *string, error) {
	event, err := LookupEvent(key, value, serviceCloudTrail)
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Unable to find CloudTrail event for %s %s", *key, *value))
	}

	// The CloudTrailEvent is a JSON object of unknown format
	dataStr := aws.StringValue(event.CloudTrailEvent)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to unmarshal JSON data from CloudTrail event")
	}
	userIdentity, ok := data["userIdentity"].(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("Could not retrieve userIdentity from JSON object")
	}
	roleArn, ok := userIdentity["arn"].(string)
	if !ok {
		return nil, nil, errors.New("Could not retrieve arn from JSON object")
	}
	username, _ := userIdentity["userName"].(string)
	return &roleArn, &username, nil
}

// GetUser uses a roleArn to find a given user
func GetUser(roleArn *string, serviceCloudTrail *cloudtrail.CloudTrail) (*string, error) {
	key := "ResourceName"
	event, err := LookupEvent(&key, roleArn, serviceCloudTrail)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unable to find CloudTrail event for role arn %s", *roleArn))
	}

	// The CloudTrailEvent is a JSON object of unknown format
	username := aws.StringValue(event.Username)
	return &username, nil
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

	verbose := cli.LogLevelIsDebug(v)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Check the config and exit with usage details if there is a problem
	err := checkConfig(v)
	if err != nil {
		logger.Println(err)
		logger.Println("Usage of find-guardduty-user:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Get credentials from environment or AWS Vault
	c, errorConfig := cli.GetAWSConfig(v, verbose)
	if errorConfig != nil {
		logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
	}
	session, errorSession := awssession.NewSession(c)
	if errorSession != nil {
		logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
	}

	// GuardDuty has the findings
	serviceGuardDuty := guardduty.New(session)

	// CloudTrail has information on who caused the event
	serviceCloudTrail := cloudtrail.New(session)

	// List Detectors
	listDetectorsInput := guardduty.ListDetectorsInput{}
	detectors, err := serviceGuardDuty.ListDetectors(&listDetectorsInput)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "Unable to list Guard Duty detectors"))
	}

	// Walk through detectors
	for _, detectorID := range detectors.DetectorIds {

		// The token for paging through findings
		var listFindingsNextToken *string

		for {
			// Define the service condition to list findings
			archived := "false"
			if v.GetBool("archived") {
				archived = "true"
			}
			serviceCondition := guardduty.Condition{}
			serviceCondition.SetEq([]*string{&archived})

			findingCriteria := guardduty.FindingCriteria{
				Criterion: map[string]*guardduty.Condition{
					"service.archived": &serviceCondition,
				},
			}
			listFindingsInput := guardduty.ListFindingsInput{
				DetectorId:      detectorID,
				FindingCriteria: &findingCriteria,
				NextToken:       listFindingsNextToken,
			}

			findingList, err := serviceGuardDuty.ListFindings(&listFindingsInput)
			if err != nil {
				logger.Println(errors.Wrap(err, "Unable to list Guard Duty findings"))
				continue
			}

			// Set the next token to page for
			listFindingsNextToken = findingList.NextToken

			// If we've run out of findings then quit
			if len(findingList.FindingIds) == 0 {
				break
			}
			getFindingsInput := guardduty.GetFindingsInput{
				DetectorId: detectorID,
				FindingIds: findingList.FindingIds,
			}
			findings, err := serviceGuardDuty.GetFindings(&getFindingsInput)
			if err != nil {
				logger.Println(errors.Wrap(err, "Unable to retrieve Guard Duty findings"))
				continue
			}

			// Walk through each finding
			for _, finding := range findings.Findings {

				// Not all events are from humans and in those cases we skip
				if finding.Resource.AccessKeyDetails == nil {
					if verbose {
						logger.Println(fmt.Sprintf("\nSkipping Non User Finding ID: %s", aws.StringValue(finding.Id)))
					}
					continue
				}

				fd := FindingDetail{
					ID:          finding.Id,
					CreatedAt:   finding.CreatedAt,
					AccessKeyID: finding.Resource.AccessKeyDetails.AccessKeyId,
					PrincipalID: finding.Resource.AccessKeyDetails.PrincipalId,
					ServiceName: finding.Service.Action.AwsApiCallAction.ServiceName,
					API:         finding.Service.Action.AwsApiCallAction.Api,
					City:        finding.Service.Action.AwsApiCallAction.RemoteIpDetails.City.CityName,
					Country:     finding.Service.Action.AwsApiCallAction.RemoteIpDetails.Country.CountryName,
				}

				var roleArn *string
				var username *string

				// The IPv4 address is not always available if the service is missing
				if finding.Service != nil {
					fd.IPAddress = finding.Service.Action.AwsApiCallAction.RemoteIpDetails.IpAddressV4
				}

				// Get Assumed Role ARN and Username details based on Access Key or Principal IDs
				if fd.AccessKeyID != nil && *fd.AccessKeyID != "" && *fd.AccessKeyID != "GeneratedFindingAccessKeyId" {
					var err error
					key := "AccessKeyId"
					roleArn, username, err = GetRoleAndUser(&key, fd.AccessKeyID, serviceCloudTrail)
					if err != nil {
						logger.Println(errors.Wrap(err, "Unable to find role arn and username from access key in finding"))
						continue
					}
				} else if fd.PrincipalID != nil && *fd.PrincipalID != "" && *fd.PrincipalID != "GeneratedFindingPrincipalId" {
					var err error
					key := "ResourceName"
					roleArn, username, err = GetRoleAndUser(&key, fd.PrincipalID, serviceCloudTrail)
					if err != nil {
						logger.Println(errors.Wrap(err, "Unable to find role arn and username from principal ID in finding"))
						continue
					}
				} else {
					continue
				}

				if roleArn == nil || *roleArn == "" {
					continue
				}
				fd.AssumedRoleARN = roleArn

				// If previous queries did not return username try again using the assumed role arn
				if username == nil || *username == "" {
					var err error
					username, err = GetUser(roleArn, serviceCloudTrail)
					if err != nil {
						logger.Println(errors.Wrap(err, "Unable to find role username from role arn"))
						continue
					}
				}
				fd.Username = username

				// Print output in desired format
				if output := v.GetString("output"); output == "json" {
					err := fd.PrintJSON(logger)
					if err != nil {
						logger.Println(errors.Wrap(err, "Unable to marshal finding detail to JSON"))
					}
				} else {
					fd.Print(logger)
				}
			}

			// If the next token is nil or an empty string then there are no more results to page through
			if listFindingsNextToken == nil || *listFindingsNextToken == "" {
				break
			}
		}
	}
}
