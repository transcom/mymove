package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/invoice"
)

const (
	// ProcessTPPSLastReadTimeFlag is the ENV var for the last read time
	ProcessTPPSLastReadTimeFlag string = "process-tpps-last-read-time"
)

// Call this from the command line with go run ./cmd/milmove-tasks process-tpps

func checkProcessTPPSConfig(v *viper.Viper, logger *zap.Logger) error {
	logger.Debug("checking config for process-tpps")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	err = cli.CheckLogging(v)
	if err != nil {
		return err
	}

	// err = cli.CheckTPPSSFTP(v)
	// if err != nil {
	// 	return err
	// }

	// if err := cli.CheckSFTP(v); err != nil {
	// 	return err
	// }

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	return cli.CheckEntrustCert(v)
}

func initProcessTPPSFlags(flag *pflag.FlagSet) {
	// Logging Levels
	cli.InitLoggingFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	// TPPS SFTP
	// cli.InitTPPSFlags(flag)

	// Certificate
	cli.InitCertFlags(flag)

	// Entrust Certificates
	cli.InitEntrustCertFlags(flag)

	// TPPS SFTP Config
	// cli.InitTPPSSFTPFlags(flag)

	// maria not even sure I need this
	flag.String(ProcessTPPSLastReadTimeFlag, "", "Files older than this RFC3339 time will not be fetched.")
	// flag.Bool(ProcessTPPSDeleteFilesFlag, false, "If present, delete files on SFTP server that have been processed successfully")

	// Don't sort flags
	flag.SortFlags = false
}

func processTPPS(_ *cobra.Command, _ []string) error {
	v := viper.New()

	logger, _, err := logging.Config(
		logging.WithEnvironment(v.GetString(cli.LoggingEnvFlag)),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		logger.Fatal("Failed to initialized Zap logging for process-tpps")
	}
	zap.ReplaceGlobals(logger)

	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		logger.Info(fmt.Sprintf("Duration of processTPPS task:: %v", elapsedTime))
	}()

	flag := pflag.CommandLine
	initProcessTPPSFlags(flag)
	err = flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("failed to parse flags", zap.Error(err))
	}

	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatal("failed to bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err = checkProcessTPPSConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)
	dbEnv := v.GetString(cli.DbEnvFlag)
	// tppsURL := v.GetString(cli.TPPSURLFlag)
	// logger.Info(fmt.Sprintf("TPPS URL is %v", tppsURL))

	isDevOrTest := dbEnv == "experimental" || dbEnv == "development" || dbEnv == "test"
	if isDevOrTest {
		logger.Info(fmt.Sprintf("Starting in %s mode, which enables additional features", dbEnv))
	}

	certLogger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		logger.Fatal("Failed to initialize Zap logging", zap.Error(err))
	}
	certificates, rootCAs, err := certs.InitDoDEntrustCertificates(v, certLogger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Error in getting tls certs", zap.Error(err))
	}

	tppsInvoiceProcessor := invoice.NewTPPSPaidInvoiceReportProcessor()

	// Process TPPS paid invoice report
	pathTPPSPaidInvoiceReport := v.GetString(cli.SFTPTPPSPaidInvoiceReportPickupDirectory)
	err = tppsInvoiceProcessor.ProcessFile(appCtx, pathTPPSPaidInvoiceReport, "")

	if err != nil {
		logger.Error("Error reading TPPS Paid Invoice Report application advice responses", zap.Error(err))
	} else {
		logger.Info("Successfully processed TPPS Paid Invoice Report application advice responses")
	}

	return nil
}
