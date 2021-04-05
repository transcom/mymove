package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/services/invoice"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// Call this from command line with go run ./cmd/fetch-from-syncada-via-sftp/ --directory <syncada directory to download from> --last-read-time <time of last run>

const (
	// LastReadTimeFlag is the ENV var for the last read time
	LastReadTimeFlag string = "last-read-time"
	// DirectoryFlag is the ENV var for the directory
	DirectoryFlag string = "directory"
	// DeleteFilesFlag is the ENV var for deleting SFTP files after they've been processed
	DeleteFilesFlag string = "delete-files-after-processing"
)

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

	flag.String(LastReadTimeFlag, "", "Files older than this RFC3339 time will not be fetched.")
	flag.String(DirectoryFlag, "", "syncada path")
	flag.Bool(DeleteFilesFlag, false, "If present, delete files on SFTP server that have been processed successfully")

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

	sshClient, err := cli.InitSyncadaSSH(v, logger)
	if err != nil {
		logger.Fatal("couldn't initialize SSH client", zap.Error(err))
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			logger.Fatal("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	sftpClient, err := cli.InitSyncadaSFTP(sshClient, logger)
	if err != nil {
		logger.Fatal("couldn't initialize SFTP client", zap.Error(err))
	}
	defer func() {
		if closeErr := sftpClient.Close(); closeErr != nil {
			logger.Fatal("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	// Sample expected format: 2021-03-16T18:25:36Z
	lastReadTime := v.GetString(LastReadTimeFlag)
	var t time.Time
	if lastReadTime != "" {
		t, err = time.Parse(time.RFC3339, lastReadTime)
		if err != nil {
			logger.Error("couldn't parse time", zap.Error(err))
		}
	}
	logger.Info("lastRead", zap.String("t", t.String()))

	wrappedSFTPClient := invoice.NewSFTPClientWrapper(sftpClient)
	syncadaSFTPSession := invoice.NewSyncadaSFTPReaderSession(wrappedSFTPClient, db, logger, v.GetBool(DeleteFilesFlag))

	_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(v.GetString(DirectoryFlag), t, invoice.NewEDI997Processor(db, logger))
	if err != nil {
		logger.Error("Error reading 997 responses", zap.Error(err))
	} else {
		logger.Info("Successfully processed 997 responses")
	}

	_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(v.GetString(DirectoryFlag), t, invoice.NewEDI824Processor(db, logger))
	if err != nil {
		logger.Error("Error reading 824 responses", zap.Error(err))
	} else {
		logger.Info("Successfully processed 824 responses")
	}
}
