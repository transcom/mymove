package main

import (
	"context"
	"log"
	"os"
	"path"
	"strings"

	"github.com/transcom/mymove/pkg/storage"

	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"go.uber.org/zap"
)

type Event struct {
	Key string `json:"name"`
}

const (
	StorageKey  string = "storage-key"
	OutFileName string = "outfile-name"
)

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckStorage(v)
	if err != nil {
		return err
	}
	if v.GetString(StorageKey) == "" {
		return errors.New("missing storage key")
	}
	if v.GetString(OutFileName) == "" {
		v.Set(OutFileName, path.Base(v.GetString(StorageKey)))
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// Storage config
	cli.InitStorageFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	flag.String(StorageKey, "", "file storage key")
	flag.String(OutFileName, "", "local filename")
	// Don't sort flags
	flag.SortFlags = false
}

func HandleRequest(ctx context.Context, name Event) (string, error) {

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
	// Have content type in db, but for now going to avoid connecting to db, so just retrieve from bucket
	contentType, err := storer.ContentType(v.GetString(StorageKey))
	if err != nil {
		logger.Fatal("can't get content type", zap.Error(err))
	}

	f, err := storer.PresignedURL(v.GetString(StorageKey), contentType)
	if err != nil {
		logger.Fatal("can't get generate presigned url", zap.Error(err))
	}
	return f, nil
}

func main() {

	HandleRequest(context.TODO(), Event{Key: StorageKey})
	//TODO lambda, lambda, lambda
	//lambda.Start(HandleRequest)
}
