package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/transcom/mymove/pkg/storage"

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

var v *viper.Viper

func init() {
	// Had to move into init for now otherwise tried to set the flags on each call to handler which
	// causes a panic since can only set once
	v = viper.New()
	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
}

type Event struct {
	Key string `json:"key"`
}

type Response struct {
	PresignedURL string `json:"preSignedURL"`
	StatusCode   int    `json:"statusCode"`
}

func HandleRequest(ctx context.Context, event Event) (Response, error) {
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	eventJson, err := json.Marshal(event)
	logger.Info(string(eventJson))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	var session *awssession.Session
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
	}

	// Create a connection to the DB
	storer := storage.InitStorage(v, session, logger)
	log.Println("storer: ", storer)
	// Have content type in db, but for now going to avoid connecting to db, so just retrieve from bucket
	contentType, err := storer.ContentType(event.Key)
	if err != nil {
		logger.Fatal("can't get content type", zap.Error(err))
	}
	log.Println("ContentType: ", contentType)

	presignedURL, err := storer.PresignedURL(event.Key, contentType)
	if err != nil {
		logger.Fatal("can't get generate presigned url", zap.Error(err))
	}
	log.Println("URL: ", presignedURL)
	return Response{PresignedURL: presignedURL, StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)

	// For local testing
	// ctx := context.Background()
	// event := Event{Key: "user/d6aab501-dd85-4126-b71a-246fc50ec263/uploads/c17771af-2878-4aaf-923d-8faf1cd58cea"}
	// HandleRequest(ctx, event)
}
