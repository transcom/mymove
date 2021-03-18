package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from command line with go run ./cmd/fetch-from-syncada-via-sftp/ --local-file-path <localFilePath> --syncada-file-name <syncadaFileName>

func checkConfig(v *viper.Viper, logger logger) error {
	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	err = cli.CheckLogging(v)
	if err != nil {
		return err
	}

	err = cli.CheckSyncadaSFTP(v)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {
	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Syncada SFTP Config
	cli.InitSyncadaSFTPFlags(flag)

	flag.String("last-read-time", "", "Files older than this time will not be fetched.")
	flag.String("directory", "", "syncada path")

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

	logger, err := logging.Config(logging.WithEnvironment(v.GetString(cli.DbEnvFlag)), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	db, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		logger.Fatal("connecting to DB", zap.Error(err))
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.Fatal("could not close database", zap.Error(closeErr))
		}
	}()

	// TODO: Do we need to get connection back so we can close it too?
	sftpClient, err := cli.InitSyncadaSFTP(v, logger)
	if err != nil {
		logger.Fatal("couldn't initialize sftp session", zap.Error(err))
	}
	defer func() {
		if closeErr := sftpClient.Close(); closeErr != nil {
			logger.Fatal("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	//2021-03-16T18:25:36Z
	t, err := time.Parse(time.RFC3339, v.GetString("last-read-time"))
	if err != nil {
		logger.Error("couldn't parse time", zap.Error(err))
	}
	//t := time.Now().Add(-1 * time.Hour)
	logger.Info("lastRead", zap.String("a", t.String()))
	syncadaSFTPSession := invoice.InitNewSyncadaSFTPReaderSession(sftpClient)
	data, _, err := syncadaSFTPSession.ReadFromSyncadaViaSFTP(v.GetString("directory"), t)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(data)
}
