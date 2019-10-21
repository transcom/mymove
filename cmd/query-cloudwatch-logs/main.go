package main

/*
* Query against cloudwatch logs and insights
*
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// AWSRegionFlag holds string identifier for command line usage
	AWSRegionFlag string = "aws-region"
	// VerboseFlag holds string identifier for command line usage
	VerboseFlag string = "verbose"
	// QueryFlag holds string identifier for command line usage
	QueryFlag string = "query"
	// LogGroupNameFlag holds string identifier for command line usage
	LogGroupNameFlag string = "log-group-name"
	// StartTimeFlag holds string identifier for command line usage
	StartTimeFlag string = "start-time"
	// EndTimeFlag holds string identifier for command line usage
	EndTimeFlag string = "end-time"
	// LimitFlag holds string identifier for command line usage
	LimitFlag string = "limit"
	// RawQueryFlag holds string identifier for command line usage
	RawQueryFlag string = "raw-query"
)

// InsightQuery struct is useful for capturing all query params in one place
type InsightQuery struct {
	Query        string `json:"query"`
	LogGroupName string `json:"logGroupName"`
	EndTime      int64  `json:"endTime"`
	StartTime    int64  `json:"startTime"`
	Limit        int64  `json:"limit"`
}

func getSessionAndLogger(cmd *cobra.Command, args []string) (*awssession.Session, *viper.Viper, *log.Logger, error) {

	err := cmd.ParseFlags(args)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "Could not parse flags")
	}

	flag := cmd.Flags()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

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

	if !v.GetBool(VerboseFlag) {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	isAwsVault := v.GetString("AWS_VAULT")
	if isAwsVault != "" {
		logger.Println("aws-vault double wrap detected, please run this command without aws-vault")
		os.Exit(1)
	}

	keychainName := v.GetString(cli.VaultAWSKeychainNameFlag)
	keychainProfile := v.GetString(cli.VaultAWSProfileFlag)
	sessionDuration := v.GetDuration(cli.VaultAWSSessionDurationFlag)
	assumeRoleTTL := v.GetDuration(cli.VaultAWSAssumeRoleTTLFlag)
	creds, err := cli.GetAWSCredentialsFromKeyring(keychainName, keychainProfile, sessionDuration, assumeRoleTTL)
	if err != nil {
		logger.Fatal(errors.Wrap(err, fmt.Sprintf("Unable to get AWS credentials from the keychain %s and profile %s", keychainName, keychainProfile)))
		os.Exit(1)

	}

	// Define services
	sess := awssession.Must(awssession.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(v.GetBool(VerboseFlag)),
		Credentials:                   creds,
		Region:                        aws.String(v.GetString(AWSRegionFlag)),
	}))
	return sess, v, logger, nil
}

// function for starting a new insight query, returns a queryid as its return type
func startInsightQuery(svc cloudwatchlogs.CloudWatchLogs, query, logGroup *string, startTime, endTime, limit *int64, logger *log.Logger) string {

	queryInput := cloudwatchlogs.StartQueryInput{
		StartTime:    startTime,
		EndTime:      endTime,
		LogGroupName: logGroup,
		QueryString:  query,
		Limit:        limit,
	}

	output, err := svc.StartQuery(&queryInput)
	if err != nil {
		logger.Printf("Got error starting the query events: %s", *query)
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	return *output.QueryId
}

// function for obtaining results of a query by queryID
func getQueryResultsByID(svc cloudwatchlogs.CloudWatchLogs, queryID string, logger *log.Logger) cloudwatchlogs.GetQueryResultsOutput {
	queryResultsInput := cloudwatchlogs.GetQueryResultsInput{QueryId: &queryID}
	results, err := svc.GetQueryResults(&queryResultsInput)
	if err != nil {
		logger.Printf("failed to get query by queryId: %s", queryID)
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	return *results
}

// printQueryResults formats and prints query results output
func printQueryResults(results cloudwatchlogs.GetQueryResultsOutput, logger *log.Logger) {
	formattedOutput, err := json.Marshal(results)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Println(string(formattedOutput))
}

func createAndRunQuery(sess *awssession.Session, v *viper.Viper, query, logGroupName *string, startTime, endTime, limit *int64, logger *log.Logger) cloudwatchlogs.GetQueryResultsOutput {
	noop := v.GetBool(RawQueryFlag)

	// don't execute the query
	if noop {
		formattedOutput, err := json.Marshal(InsightQuery{Query: *query, Limit: *limit, StartTime: *startTime, EndTime: *endTime, LogGroupName: *logGroupName})
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Print(string(formattedOutput))
		os.Exit(0)
	}

	svc := cloudwatchlogs.New(sess)
	queryID := startInsightQuery(*svc, query, logGroupName, startTime, endTime, limit, logger)
	results := getQueryResultsByID(*svc, queryID, logger)

	verbose := v.GetBool(VerboseFlag)
	// keep querying until results stop running
	for *results.Status == "Running" {
		results = getQueryResultsByID(*svc, queryID, logger)
		formattedOutput, err := json.Marshal(results)
		if err != nil {
			logger.Fatal(err.Error())
		}
		if verbose {
			logger.Println("-----------------------------------")
			logger.Println(string(formattedOutput))
		}
		time.Sleep(time.Duration(4000))
	}
	return results
}

// initRootFlags initializes the flags for the root and all sub-commands
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitVaultFlags(flag)
	flag.String(AWSRegionFlag, "us-west-2", "The default aws region")
	flag.BoolP(VerboseFlag, "v", false, "Show extra output for debugging")
}

// initQueryFlags initializes the flags for the query commands
func initQueryFlags(flag *pflag.FlagSet) {
	defaultEndTime := time.Now()
	defaultStartTime := time.Now().Add(time.Hour * -24)
	defaultQuery := `fields @message
	| parse "statement: *" as statement
	| filter statement not like "SELECT 1;"
	| sort @timestamp desc`
	defaultLogGroupName := "/aws/rds/instance/app-staging/postgresql"

	// For more see https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html
	flag.StringP(QueryFlag, "q", defaultQuery, "Custom query https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html")
	flag.StringP(LogGroupNameFlag, "l", defaultLogGroupName, "The log group you would like to query against")
	flag.StringP(StartTimeFlag, "s", defaultStartTime.Format(time.RFC3339), "Start time in RFC3339 format for the query")
	flag.StringP(EndTimeFlag, "e", defaultEndTime.Format(time.RFC3339), "End time in RFC3339 format for the query")
	flag.IntP(LimitFlag, "n", int(100), "Limit number of query results")
	flag.Bool(RawQueryFlag, false, "Print query and its params without executing it")

	flag.SortFlags = false
}

func main() {
	root := cobra.Command{
		Use:   "query-cloudwatch-logs [flags]",
		Short: "Query CloudWatch Logs",
		Long:  "Query CloudWatchLogs",
	}
	initRootFlags(root.PersistentFlags())

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Query",
		Long:  "Query",
		RunE:  runQueryWithParams,
	}
	initQueryFlags(queryCmd.Flags())
	root.AddCommand(queryCmd)

	listGroupsCmd := &cobra.Command{
		Use:   "list-groups",
		Short: "List Groups",
		Long:  "List All Log Group Names",
		RunE:  getAllLogGroupNames,
	}
	root.AddCommand(listGroupsCmd)

	longestQueryCmd := &cobra.Command{
		Use:   "longest-query",
		Short: "Get Longest Queries",
		Long:  "Get 10 Longest Running DB Queries in last 24 hour time period",
		RunE:  longestRunningQuery,
	}
	initQueryFlags(longestQueryCmd.Flags())
	root.AddCommand(longestQueryCmd)

	recentProdQueriesCmd := &cobra.Command{
		Use:   "recent-query",
		Short: "Get Recent Queries",
		Long:  "Get 100 Recent Queries run against Prod RDS",
		RunE:  recentQueriesRDS,
	}
	initQueryFlags(recentProdQueriesCmd.Flags())
	root.AddCommand(recentProdQueriesCmd)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
