package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// AWSProfileFlag holds string identifier for command line usage
	AWSProfileFlag string = "aws-profile"
	// AWSRegionFlag holds string identifier for command line usage
	AWSRegionFlag string = "aws-region"
	// VerboseFlag holds string identifier for command line usage
	VerboseFlag string = "verbose"
	// WafForward string holds value that indicates the request was successfully forwarded from alb
	WafForward string = "waf,forward"
	// LimitFlag holds string identifier for command line usage
	LimitFlag string = "limit"

	// StatusCodeFlag holds string identifier for command line usage
	StatusCodeFlag string = "status-code"
	// EnvFlag holds string value for environment
	EnvFlag string = "ENV"
	// AlbTable holds string value for athena table name
	AlbTable string = "alb_logs"
	// AddPartitions holds string value that indicates if partitions need to be added
	AddPartitions string = "add-partitions"

	// AthenaWorkGroup holds the string value for work group to use for running queries
	AthenaWorkGroup string = "log-query"
)

// initialize flags
func initFlags(flag *pflag.FlagSet) {
	flag.BoolP(AddPartitions, "p", false, "Add partitions by month and year")
	flag.String(AWSProfileFlag, "", "The aws-vault profile")
	flag.String(AWSRegionFlag, "us-west-2", "The default aws region")
	flag.String(EnvFlag, "experimental", "Environment against which query would be executed. (staging, experimental, prod)")
	flag.IntP(LimitFlag, "n", int(10), "Limit number of query results")
	flag.IntP(StatusCodeFlag, "s", int(0), "Filter by an exact status code. Defaults to '> 499'")
	flag.BoolP(VerboseFlag, "v", false, "Show extra output for debugging")
	flag.SortFlags = false
}

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Arg parse failed")
		return
	}

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		fmt.Println("Arg binding failed")
		return
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	if !v.GetBool(VerboseFlag) {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	//Create info logger for verbose
	infoLogger := log.New(ioutil.Discard, "INFO", log.LstdFlags)

	// if verbose flag is passed then change output of info logger
	if v.GetBool(VerboseFlag) {
		infoLogger.SetOutput(os.Stdout)
	}

	verbose := cli.LogLevelIsDebug(v)
	AWSConfig, errorConfig := cli.GetAWSConfig(v, verbose)
	if errorConfig != nil {
		logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
	}
	session, errorSession := awssession.NewSession(AWSConfig)
	if errorSession != nil {
		logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
	}

	serviceAthena := athena.New(session)

	// check if table in db already exists, if no table exists then create it and add partitions
	CheckEnv(serviceAthena, logger, infoLogger, v)

	getAthenaQuery(serviceAthena, logger, infoLogger, v)

}

func getDatabaseByEnv(env string) (string, error) {
	switch env {
	case "staging":
		return "staging_alb_logs", nil
	case "experimental":
		return "experimental_alb_logs", nil
	case "prod":
		return "prod_alb_logs", nil
	default:
		return "", errors.New(fmt.Sprintf("Env is not valid: %s", env))
	}
}

func getALBLogS3BucketByEnv(env, region string) (string, error) {
	switch env {
	case "staging":
		return fmt.Sprintf("s3://transcom-ppp-aws-logs/alb/app-%s/AWSLogs/923914045601/elasticloadbalancing/%s", env, region), nil
	case "experimental":
		return fmt.Sprintf("s3://transcom-ppp-aws-logs/alb/app-%s/AWSLogs/923914045601/elasticloadbalancing/%s", env, region), nil
	case "prod":
		return fmt.Sprintf("s3://transcom-ppp-aws-logs/alb/app-%s/AWSLogs/923914045601/elasticloadbalancing/%s", env, region), nil
	default:
		return "", errors.New(fmt.Sprintf("Env is not valid: %s", env))
	}
}

func getS3BucketByEnv(env, region string) (string, error) {
	switch env {
	case "staging":
		return fmt.Sprintf("s3://transcom-ppp-aws-athena-%s-%s", env, region), nil
	case "experimental":
		return fmt.Sprintf("s3://transcom-ppp-aws-athena-%s-%s", env, region), nil
	case "prod":
		return fmt.Sprintf("s3://transcom-ppp-aws-athena-%s-%s", env, region), nil
	default:
		return "", errors.New(fmt.Sprintf("Env is not valid: %s", env))
	}
}

// CheckEnv checks the environment for querying
func CheckEnv(serviceAthena *athena.Athena, logger, infoLogger *log.Logger, v *viper.Viper) {

	//get values from flags
	awsRegion := v.GetString(AWSRegionFlag)
	env := v.GetString(EnvFlag)
	//get dbName name based on env
	dbName, err := getDatabaseByEnv(env)
	if err != nil {
		logger.Fatalf("Invalid env value %s", err.Error())
	}
	infoLogger.Println(dbName)

	//get bucket constaining alb logs
	logBucket, err := getALBLogS3BucketByEnv(env, awsRegion)
	if err != nil {
		logger.Fatalf("Could not parse alb log s3 bucket %s", err.Error())
	}
	infoLogger.Println(logBucket)
	// set default value to false
	isTableFound := false

	var startQueryExecutionInput athena.StartQueryExecutionInput

	// set athena workgroup to execute queries
	startQueryExecutionInput.SetWorkGroup(AthenaWorkGroup)

	// get list of all tables in DB
	query := fmt.Sprintf(`SHOW TABLES`)

	startQueryExecutionInput.SetQueryString(query)
	var queryExecutionContext athena.QueryExecutionContext

	queryExecutionContext.SetDatabase(dbName)
	startQueryExecutionInput.SetQueryExecutionContext(&queryExecutionContext)

	var resultConfiguration athena.ResultConfiguration

	resultConfiguration.SetOutputLocation(logBucket)
	startQueryExecutionInput.SetResultConfiguration(&resultConfiguration)

	result, err := serviceAthena.StartQueryExecution(&startQueryExecutionInput)
	if err != nil {
		logger.Fatal(err)
		return
	}

	infoLogger.Printf("Start Query Execution: %s", result.GoString())

	var queryExecutionInput athena.GetQueryExecutionInput
	queryExecutionInput.SetQueryExecutionId(*result.QueryExecutionId)

	var queryExecutionOutput *athena.GetQueryExecutionOutput
	duration := time.Duration(1) * time.Second

	for {
		queryExecutionOutput, err = serviceAthena.GetQueryExecution(&queryExecutionInput)
		if err != nil {
			logger.Fatalf("Query execution failed with error: %v", err)
			return
		}
		if *queryExecutionOutput.QueryExecution.Status.State != "RUNNING" {
			infoLogger.Print("waiting")
			break
		}
		time.Sleep(duration)
	}
	if *queryExecutionOutput.QueryExecution.Status.State == "SUCCEEDED" {

		var ip athena.GetQueryResultsInput
		ip.SetQueryExecutionId(*result.QueryExecutionId)

		op, getResultsErr := serviceAthena.GetQueryResults(&ip)
		if getResultsErr != nil {
			logger.Fatalf("Get Query results failed with error: %v", getResultsErr)
			return
		}

		// check to make sure alb_logs table exists in db
		for _, row := range op.ResultSet.Rows {
			for _, data := range row.Data {
				if *data.VarCharValue == AlbTable {
					isTableFound = true
					infoLogger.Println(fmt.Sprintf("Table found: %s", *data.VarCharValue))
					break
				}
			}
		}

	} else {
		infoLogger.Printf("Query status: %s", *queryExecutionOutput.QueryExecution.Status.State)
		logger.Println(*queryExecutionOutput.QueryExecution.Status.State)
	}

	if !isTableFound {
		logger.Println("Table not found, creating table....")
		errCreateLogTable := createLogTable(serviceAthena, logger, infoLogger, dbName, logBucket)
		if errCreateLogTable != nil {
			logger.Fatalf("Failed to create table: %v", errCreateLogTable)
		}
		logger.Println("creating monthly partitions....")
		errCreatePartitions := createPartitions(serviceAthena, logger, infoLogger, dbName, logBucket)
		if errCreatePartitions != nil {
			logger.Fatalf("Failed to create partitions: %v", errCreatePartitions)
		}
	} else if v.GetBool(AddPartitions) { // If add partitions flag is set to true then try adding partitions
		logger.Println("creating monthly partitions....")
		errCreatePartitions := createPartitions(serviceAthena, logger, infoLogger, dbName, logBucket)
		if errCreatePartitions != nil {
			logger.Fatalf("Failed to create partitions: %v", errCreatePartitions)
		}
		os.Exit(0)
	}
}

func createLogTable(serviceAthena *athena.Athena, logger, infoLogger *log.Logger, dbName, s3Path string) error {
	var startQueryExecutionInput athena.StartQueryExecutionInput
	query := fmt.Sprintf(`CREATE EXTERNAL TABLE IF NOT EXISTS alb_logs(
  			type string,
  			time string,
  			elb string,
			client_ip string,
			client_port int,
			target_ip string,
			target_port int,
			request_processing_time double,
			target_processing_time double,
			response_processing_time double,
			elb_status_code int,
			target_status_code int,
			received_bytes bigint,
			sent_bytes bigint,
			request_verb string,
			request_url string,
			request_proto string,
			user_agent string,
			ssl_cipher string,
			ssl_protocol string,
			target_group_arn string,
			trace_id string,
			domain_name string,
			chosen_cert_arn string,
			matched_rule_priority string,
			request_creation_time string,
			actions_executed string,
			redirect_url string,
			lambda_error_reason string,
			new_field string)
			PARTITIONED BY (
  				year int,
  				month int
			)
		ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.RegexSerDe'
		WITH SERDEPROPERTIES ('input.regex'='([^ ]*) ([^ ]*) ([^ ]*) ([^ ]*):([0-9]*) ([^ ]*)[:-]([0-9]*) ([-.0-9]*) ([-.0-9]*) ([-.0-9]*) (|[-0-9]*) (-|[-0-9]*) ([-0-9]*) ([-0-9]*) \"([^ ]*) ([^ ]*) (- |[^ ]*)\" \"([^\"]*)\" ([A-Z0-9-]+) ([A-Za-z0-9.-]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\" \"([^\"]*)\" ([-.0-9]*) ([^ ]*) \"([^\"]*)\" \"([^\"]*)\"($| \"[^ ]*\")(.*)')
		STORED AS INPUTFORMAT 'org.apache.hadoop.mapred.TextInputFormat'
		OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
		LOCATION '%s/'`, s3Path)

	startQueryExecutionInput.SetWorkGroup(AthenaWorkGroup)
	startQueryExecutionInput.SetQueryString(query)
	var queryExecutionContext athena.QueryExecutionContext
	queryExecutionContext.SetDatabase(dbName)
	startQueryExecutionInput.SetQueryExecutionContext(&queryExecutionContext)

	var resultConfiguration athena.ResultConfiguration
	resultConfiguration.SetOutputLocation(s3Path)
	startQueryExecutionInput.SetResultConfiguration(&resultConfiguration)

	result, err := serviceAthena.StartQueryExecution(&startQueryExecutionInput)
	if err != nil {
		logger.Fatal(err)
		return err
	}

	infoLogger.Printf("Start Query Execution: %s", result.GoString())

	var queryExecutionInput athena.GetQueryExecutionInput
	queryExecutionInput.SetQueryExecutionId(*result.QueryExecutionId)

	var queryExecutionOutput *athena.GetQueryExecutionOutput
	duration := time.Duration(2) * time.Second // Pause for 2 seconds

	for {
		queryExecutionOutput, err = serviceAthena.GetQueryExecution(&queryExecutionInput)
		if err != nil {
			logger.Fatalf("Query execution failed with error: %v", err)
			return err
		}
		if *queryExecutionOutput.QueryExecution.Status.State != "RUNNING" {
			infoLogger.Print("waiting")
			break
		}
		time.Sleep(duration)
	}
	if *queryExecutionOutput.QueryExecution.Status.State == "SUCCEEDED" {
		infoLogger.Printf("Query status details: %v", *queryExecutionOutput.QueryExecution.Status)
	} else {
		infoLogger.Printf("Query status details: %v", *queryExecutionOutput.QueryExecution.Status)
	}
	return nil
}

func createPartitions(serviceAthena *athena.Athena, logger, infoLogger *log.Logger, dbName, s3Path string) error {

	var startQueryExecutionInput athena.StartQueryExecutionInput
	startDate := time.Now().AddDate(-1, 0, 0)
	endDate := time.Now().AddDate(6, 0, 0)

	infoLogger.Println(fmt.Sprintf("Partition start date: %s", startDate.String()))
	infoLogger.Println(fmt.Sprintf("Partition end date: %s", endDate.String()))

	query := "ALTER TABLE alb_logs ADD IF NOT EXISTS"
	for startDate.Before(endDate) {
		infoLogger.Println(fmt.Sprintf("year: %d month: %d", startDate.Year(), startDate.Month()))
		partitionPath := fmt.Sprintf(s3Path+`/%d/%s`, startDate.Year(), fmt.Sprintf("%02d", startDate.Month()))
		query += fmt.Sprintf(`
		PARTITION (year=%d, month=%d)
		LOCATION '%s'`, startDate.Year(), startDate.Month(), partitionPath)

		startDate = startDate.AddDate(0, 1, 0)
	}

	startQueryExecutionInput.SetWorkGroup(AthenaWorkGroup)
	startQueryExecutionInput.SetQueryString(query)
	var queryExecutionContext athena.QueryExecutionContext
	queryExecutionContext.SetDatabase(dbName)
	startQueryExecutionInput.SetQueryExecutionContext(&queryExecutionContext)

	var resultConfiguration athena.ResultConfiguration
	resultConfiguration.SetOutputLocation(s3Path)
	startQueryExecutionInput.SetResultConfiguration(&resultConfiguration)

	result, err := serviceAthena.StartQueryExecution(&startQueryExecutionInput)
	if err != nil {
		logger.Fatal(err)
		return err
	}

	infoLogger.Printf("Start Query Execution: %s", result.GoString())

	var queryExecutionInput athena.GetQueryExecutionInput
	queryExecutionInput.SetQueryExecutionId(*result.QueryExecutionId)

	var queryExecutionOutput *athena.GetQueryExecutionOutput
	duration := time.Duration(2) * time.Second // Pause for 2 seconds

	for {
		queryExecutionOutput, err = serviceAthena.GetQueryExecution(&queryExecutionInput)
		if err != nil {
			logger.Fatalf("Query execution failed with error: %v", err)
			return err
		}
		if *queryExecutionOutput.QueryExecution.Status.State != "RUNNING" {
			infoLogger.Print("waiting")
			break
		}
		time.Sleep(duration)
	}
	if *queryExecutionOutput.QueryExecution.Status.State == "SUCCEEDED" {
		infoLogger.Printf("Query status details: %v", *queryExecutionOutput.QueryExecution.Status)
	} else {
		infoLogger.Printf("Query status details: %v", *queryExecutionOutput.QueryExecution.Status)
	}

	return nil
}

func getAthenaQuery(serviceAthena *athena.Athena, logger, infoLogger *log.Logger, v *viper.Viper) {

	//get values from env flags
	env := v.GetString(EnvFlag)
	//get dbName name based on env
	dbName, err := getDatabaseByEnv(env)
	if err != nil {
		logger.Fatalf("Invalid env value %s", err.Error())
	}
	limit := v.GetInt(LimitFlag)
	awsRegion := v.GetString(AWSRegionFlag)
	// get bucket name based on env
	s3Bucket, err := getS3BucketByEnv(env, awsRegion)
	if err != nil {
		logger.Fatalf("Could not parse s3 bucket %s", err.Error())
	}
	infoLogger.Println(s3Bucket)

	infoLogger.Println("Querying alb logs from Athena:")
	var startQueryExecutionInput athena.StartQueryExecutionInput

	filterByStatus := v.GetInt(StatusCodeFlag)
	var query string
	if filterByStatus > 0 {
		query = fmt.Sprintf(`SELECT actions_executed,
		elb_status_code,
		*,
		from_iso8601_timestamp(time) AS timestmp
		FROM alb_logs
		WHERE elb_status_code = %d
		ORDER BY  timestmp DESC
		limit %d`, filterByStatus, limit)

	} else {
		query = fmt.Sprintf(`SELECT actions_executed,
		elb_status_code,
		*,
		from_iso8601_timestamp(time) AS timestmp
		FROM alb_logs
		WHERE elb_status_code > 499
		ORDER BY  timestmp DESC
		limit %d`, limit)
	}

	startQueryExecutionInput.SetQueryString(query)
	startQueryExecutionInput.SetWorkGroup(AthenaWorkGroup)
	var queryExecutionContext athena.QueryExecutionContext
	queryExecutionContext.SetDatabase(dbName)
	startQueryExecutionInput.SetQueryExecutionContext(&queryExecutionContext)

	var resultConfiguration athena.ResultConfiguration
	resultConfiguration.SetOutputLocation(s3Bucket)

	startQueryExecutionInput.SetResultConfiguration(&resultConfiguration)

	result, err := serviceAthena.StartQueryExecution(&startQueryExecutionInput)
	if err != nil {
		logger.Fatal(err)
		return
	}

	infoLogger.Printf("Start Query Execution: %s", result.GoString())

	var queryExecutionInput athena.GetQueryExecutionInput
	queryExecutionInput.SetQueryExecutionId(*result.QueryExecutionId)

	var queryExecutionOutput *athena.GetQueryExecutionOutput
	duration := time.Duration(2) * time.Second // Pause for 2 seconds

	for {
		queryExecutionOutput, err = serviceAthena.GetQueryExecution(&queryExecutionInput)
		if err != nil {
			logger.Fatalf("Query execution failed with error: %v", err)
			return
		}
		if *queryExecutionOutput.QueryExecution.Status.State != "RUNNING" {
			infoLogger.Print("waiting")
			break
		}
		time.Sleep(duration)
	}
	if *queryExecutionOutput.QueryExecution.Status.State == "SUCCEEDED" {

		var ip athena.GetQueryResultsInput
		ip.SetQueryExecutionId(*result.QueryExecutionId)

		op, err := serviceAthena.GetQueryResults(&ip)
		if err != nil {
			logger.Fatalf("Get Query results failed with error: %v", err)
			return
		}

		res := dataFormatter(op.ResultSet.Rows, true)
		res = attachCloudwatchLinks(res, env)
		formattedOutput, err := json.Marshal(res)
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Println(string(formattedOutput))
	} else {
		infoLogger.Printf("Query status: %s", *queryExecutionOutput.QueryExecution.Status.State)
		logger.Println(*queryExecutionOutput.QueryExecution.Status)
	}
}

func dataFormatter(rows []*athena.Row, firstRowNamed bool) []map[string]string {
	var formattedOutput []map[string]string
	var headerRow []string

	for index, row := range rows {

		// first row
		if index == 0 {
			if firstRowNamed {
				for _, data := range row.Data {
					headerRow = append(headerRow, *data.VarCharValue)
				}
				continue
			} else {
				for i := range row.Data {
					headerRow = append(headerRow, string(i))
				}
			}
		}

		var prop = make(map[string]string)
		for i, data := range row.Data {
			// if value is nil then use empty string
			if data.VarCharValue == nil {
				prop[headerRow[i]] = ""
			} else {
				prop[headerRow[i]] = *data.VarCharValue
			}
		}
		formattedOutput = append(formattedOutput, prop)
	}
	return formattedOutput
}

func attachCloudwatchLinks(rows []map[string]string, env string) []map[string]string {
	for _, row := range rows {
		actions := row["actions_executed"]
		timeStamp := row["time"]
		ip := row["client_ip"]
		layout := time.RFC3339Nano
		parsedTime, _ := time.Parse(layout, timeStamp)
		startTime := parsedTime.Add(time.Minute * -2).Format(layout)
		endTime := parsedTime.Add(time.Minute * 2).Format(layout)

		// request got past alb, generate link for cloudwatch ui
		if actions == WafForward {
			region := "us-west-2"
			cwURL := generateCustomCloudWatchLink(region, endTime, startTime, ip, env)
			row["ecs_tasks_app_logs_url"] = cwURL
		}
	}
	return rows
}

func generateCustomCloudWatchLink(region, end, start, clientIP, env string) string {
	host := "console.aws.amazon.com"
	path := "cloudwatch/home"
	dataSource := fmt.Sprintf("ecs-tasks-app-%s", env)

	cloudwatchURL := url.URL{
		Scheme: "https",
		Host:   region + "." + host,
		Path:   path,
	}

	// compose query string
	query := cloudwatchURL.Query()
	query.Set("region", region)
	cloudwatchURL.RawQuery = query.Encode()

	// set start time, end time, client ip and the source for the cloudwatch query
	fragmentParams := fmt.Sprintf("~(end~'%s~start~'%s~timeType~'ABSOLUTE~tz~'Local~editorString~'fields*20*40timestamp*2c*20*40message*0a*7c*20filter*20*60x-forwarded-for*60*20*3d*20'%s'~isLiveTail~false~source~(~'%s))~", end, start, clientIP, dataSource)
	// add additional data to query through fragment property
	cloudwatchURL.Fragment = "logs-insights:queryDetail=" + fragmentParams

	return cloudwatchURL.String()
}
