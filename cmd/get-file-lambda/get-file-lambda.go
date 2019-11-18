package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/transcom/mymove/pkg/storage"

	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"go.uber.org/zap"
)

type Event struct {
	Key string `json:"key"`
}

func HandleRequest(ctx context.Context, event Event) (string, error) {

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	//TODO pass via env
	verbose := true
	dbEnv := "development"
	storageBackend := "s3"
	s3Region := "us-west-2"
	s3Bucket := "transcom-ppp-app-devlocal-us-west-2"
	s3KeyNameSpace := "matthewkrump"
	logger, err := logging.Config(dbEnv, verbose)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	var session *awssession.Session
	if storageBackend == "s3" {
		c, errorConfig := cli.GetAWSConfig(v, verbose)
		if errorConfig != nil {
			logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		session = s
	}

	//var storer FileStorer
	if storageBackend == "s3" {
		logger.Info("Using s3 storage backend",
			zap.String("region", s3Region),
			zap.String("key", s3KeyNameSpace))
		if len(s3Bucket) == 0 {
			logger.Fatal("must provide aws-s3-bucket-name parameter, exiting")
		}
		if len(s3Region) == 0 {
			logger.Fatal("Must provide aws-s3-region parameter, exiting")
		}
		if len(s3KeyNameSpace) == 0 {
			logger.Fatal("Must provide aws_s3_key_namespace parameter, exiting")
		}
	}
	storer := storage.NewS3(s3KeyNameSpace, s3KeyNameSpace, logger, session)

	contentType, err := storer.ContentType(event.Key)
	if err != nil {
		logger.Fatal("can't get content type", zap.Error(err))
	}

	f, err := storer.PresignedURL(event.Key, contentType)
	if err != nil {
		logger.Fatal("can't get generate presigned url", zap.Error(err))
	}
	return f, nil
}

func main() {
	lambda.Start(HandleRequest)
}
