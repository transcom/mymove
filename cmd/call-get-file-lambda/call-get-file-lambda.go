package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"go.uber.org/zap"
)

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckStorage(v)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// Storage config
	cli.InitStorageFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

type Event struct {
	Key string `json:"key"`
}
type Response struct {
	PresignedURL string `json:"preSignedURL"`
	StatusCode   int    `json:"statusCode"`
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	var session *awssession.Session
	var lambdaClient *lambda.Lambda
	if v.GetString(cli.StorageBackendFlag) == "s3" {
		c, errorConfig := cli.GetAWSConfig(v, v.GetBool(cli.VerboseFlag))
		if errorConfig != nil {
			logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		session = s
		lambdaClient = lambda.New(session, c)
	}
	if lambdaClient == nil {
		fmt.Println("lambdaClient must not be nil")
		os.Exit(0)
	}
	request := Event{Key: "/user/d6aab501-dd85-4126-b71a-246fc50ec263/uploads/c17771af-2878-4aaf-923d-8faf1cd58ce"}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error marshalling request: ", request)
	}
	result, err := lambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("get-file"), Payload: payload})
	if err != nil {
		fmt.Println("Error calling get-file")
		os.Exit(0)
	}
	var resp Response

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		fmt.Println("Error unmarshalling get-file response")
		os.Exit(0)
	}

	// If the status code is NOT 200, the call failed
	if resp.StatusCode != 200 {
		fmt.Println("Error getting items, StatusCode: " + strconv.Itoa(resp.StatusCode))
		os.Exit(0)
	}

	fmt.Println(resp)
}
