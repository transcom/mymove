package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from command line with go run ./cmd/send-to-syncada-via-sftp/ --local-file-path <localFilePath> --syncada-file-name <syncadaFileName>

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	flag.String("local-file-path", "", "The path where the file to be sent is located")
	flag.String("syncada-file-name", "", "The name of the file to be stored in Syncada")

	// Don't sort flags
	flag.SortFlags = false
}

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	parseErr := flag.Parse(os.Args[1:])
	if parseErr != nil {
		log.Fatal("failed to parse flags", zap.Error(parseErr))
	}

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	fmt.Println("logger: ", logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	syncadaSFTPSession, err := invoice.InitNewSyncadaSFTPSession()
	if err != nil {
		logger.Fatal("couldn't initialize sftp session", zap.Error(err))
	}

	// open local file
	localFile, err := os.Open(filepath.Clean(v.GetString("local-file-path")))
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := syncadaSFTPSession.SendToSyncadaViaSFTP(localFile, v.GetString("syncada-file-name"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d bytes copied to the Syncada SFTP server\n", bytes)
}
