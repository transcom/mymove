package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from the command line with go run ./cmd/milmove-tasks process-tpps
func checkProcessTPPSConfig(v *viper.Viper, logger *zap.Logger) error {

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	// err = cli.CheckLogging(v)
	// if err != nil {
	// 	logger.Info("Reaching process_tpps.go line 36 in checkProcessTPPSConfig")
	// 	return err
	// }

	// if err := cli.CheckCert(v); err != nil {
	// 	logger.Info("Reaching process_tpps.go line 41 in checkProcessTPPSConfig")
	// 	return err
	// }

	// logger.Info("Reaching process_tpps.go line 45 in checkProcessTPPSConfig")
	// return cli.CheckEntrustCert(v)

	return nil
}

// initProcessTPPSFlags initializes TPPS processing flags
func initProcessTPPSFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Certificate
	// cli.InitCertFlags(flag)

	// // Entrust Certificates
	// cli.InitEntrustCertFlags(flag)

	// cli.InitTPPSFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func processTPPS(cmd *cobra.Command, args []string) error {
	flag := pflag.CommandLine
	flags := cmd.Flags()
	cli.InitDatabaseFlags(flag)

	err := cmd.ParseFlags(args)
	if err != nil {
		return fmt.Errorf("could not parse args: %w", err)
	}
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return fmt.Errorf("could not bind flags: %w", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, _, err := logging.Config(
		logging.WithEnvironment(dbEnv),
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

	// initProcessTPPSFlags(flag)
	// err = flag.Parse(os.Args[1:])
	// if err != nil {
	// 	log.Fatal("failed to parse flags", zap.Error(err))
	// }

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
	// dbEnv := v.GetString(cli.DbEnvFlag)

	// isDevOrTest := dbEnv == "experimental" || dbEnv == "development" || dbEnv == "test"
	// if isDevOrTest {
	// 	logger.Info(fmt.Sprintf("Starting in %s mode, which enables additional features", dbEnv))
	// }

	// certLogger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	// if err != nil {
	// 	logger.Fatal("Failed to initialize Zap logging", zap.Error(err))
	// }
	// certificates, rootCAs, err := certs.InitDoDEntrustCertificates(v, certLogger)
	// if certificates == nil || rootCAs == nil || err != nil {
	// 	logger.Fatal("Error in getting tls certs", zap.Error(err))
	// }

	tppsInvoiceProcessor := invoice.NewTPPSPaidInvoiceReportProcessor()

	// Process TPPS paid invoice report
	s3BucketTPPSPaidInvoiceReport := v.GetString(cli.ProcessTPPSInvoiceReportPickupDirectory)

	// Handling errors with processing a file or wanting to process specific TPPS payment file:

	// TODO have a parameter stored in s3 (customFilePathToProcess) that we could modify to have a specific date, should we need to rerun a filename from a specific day
	// the parameter value will be 'MILMOVE-enYYYYMMDD.csv' so that it's easy to look at the param value and know
	// the filepath format needed to grab files from the SFTP server (example filename = MILMOVE-en20241227.csv)

	customFilePathToProcess := "MILMOVE-enYYYYMMDD.csv" // TODO replace with the line below after param added to AWS
	// customFilePathToProcess := v.GetString(cli.TODOAddcustomFilePathToProcessParamHere)

	// The param will normally be MILMOVE-enYYYYMMDD.csv, so have a check in this function for if it's MILMOVE-enYYYYMMDD.csv
	tppsSFTPFileFormatNoCustomDate := "MILMOVE-enYYYYMMDD.csv"
	tppsFilename := ""
	logger.Info(tppsFilename)

	timezone, err := time.LoadLocation("America/New_York")
	if err != nil {
		logger.Error("Error loading timezone for process-tpps ECS task", zap.Error(err))
	}

	logger.Info(tppsFilename)
	if customFilePathToProcess == tppsSFTPFileFormatNoCustomDate {
		logger.Info("No custom filepath provided to process, processing payment file for yesterday's date.")
		// if customFilePathToProcess = MILMOVE-enYYYYMMDD.csv
		// process the filename for yesterday's date (like the TPPS lambda does)
		// the previous day's TPPS payment file should be available on external server
		yesterday := time.Now().In(timezone).AddDate(0, 0, -1)
		previousDay := yesterday.Format("20060102")
		tppsFilename = fmt.Sprintf("MILMOVE-en%s.csv", previousDay)
		previousDayFormatted := yesterday.Format("January 02, 2006")
		logger.Info(fmt.Sprintf("Starting transfer of TPPS data for %s: %s\n", previousDayFormatted, tppsFilename))
	} else {
		logger.Info("Custom filepath provided to process")
		// if customFilePathToProcess != MILMOVE-enYYYYMMDD.csv (meaning we have given an ACTUAL specific filename we want processed instead of placeholder MILMOVE-enYYYYMMDD.csv)
		// then append customFilePathToProcess to the s3 bucket path and process that INSTEAD OF
		// processing the filename for yesterday's date
		tppsFilename = customFilePathToProcess
		logger.Info(fmt.Sprintf("Starting transfer of TPPS data file: %s\n", tppsFilename))
	}

	pathTPPSPaidInvoiceReport := s3BucketTPPSPaidInvoiceReport + "/" + tppsFilename
	// temporarily adding logging here to see that s3 path was found
	logger.Info(fmt.Sprintf("Entire TPPS filepath pathTPPSPaidInvoiceReport: %s", pathTPPSPaidInvoiceReport))
	err = tppsInvoiceProcessor.ProcessFile(appCtx, pathTPPSPaidInvoiceReport, "")

	if err != nil {
		logger.Error("Error reading TPPS Paid Invoice Report application advice responses", zap.Error(err))
	} else {
		logger.Info("Successfully processed TPPS Paid Invoice Report application advice responses")
	}

	return nil
}
