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
	s3Region := "us-west-2"
	s3Bucket := "transcom-ppp-app-devlocal-us-west-2"
	s3KeyNameSpace := "matthewkrump"
	//dbEnv, err := env.Env("DB_ENV")
	//if err != nil {
	//	log.Fatalf("DB_ENV missing %v", err)
	//}
	//s3Region, err := env.Env("S3_REGION")
	//if err != nil {
	//	log.Fatalf("S3_REGION missing %v", err)
	//}
	//s3Bucket, err := env.Env("S3_BUCKET")
	//if err != nil {
	//	log.Fatalf("S3_BUCKET missing %v", err)
	//}
	//s3KeyNameSpace, err := env.Env("S3_KEY_NAME_SPACE")
	//if err != nil {
	//	log.Fatalf("S3_KEY_NAME_SPACE missing %v", err)
	//}
	logger, err := logging.Config(dbEnv, verbose)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	var session *awssession.Session
	//TODO this is wack
	c, errorConfig := cli.GetAWSConfig(v, verbose)
	log.Printf("aws config: %+v", c)
	if errorConfig != nil {
		logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
	}
	s, errorSession := awssession.NewSession(c)
	if errorSession != nil {
		logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
	}
	session = s
	log.Printf("aws session: %+v", session)

	logger.Info("Using s3 storage backend",
		zap.String("region", s3Region),
		zap.String("key", s3KeyNameSpace))
	storer := storage.NewS3(s3Bucket, s3KeyNameSpace, logger, session)

	//contentType, err := storer.ContentType(event.Key)
	//if err != nil {
	//	logger.Fatal("can't get content type", zap.Error(err))
	//}

	f, err := storer.PresignedURL(event.Key, "image/png")
	if err != nil {
		logger.Fatal("can't get generate presigned url", zap.Error(err))
	}
	return f, nil
}

func main() {
	lambda.Start(HandleRequest)
}
