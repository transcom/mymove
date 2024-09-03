package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	processor "github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from command line with go run ./cmd/simulate-process-tpps/
// This binary will be explicitly for simulation and testing purposes for the scenario of
// payment request numbers 1077-4079-3, 1208-5962-1, 8801-2773-2, 8801-2773-3
// Those payment request numbers must exist in the payment_requests table in order for
// this binary to be used properly

func checkConfig(v *viper.Viper, logger *zap.Logger) error {

	err := cli.CheckDatabase(v, logger)
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

	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err.Error())
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// DB connection
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	// Create new TPPS Paid Invoice Report Processor
	processor := processor.NewTPPSPaidInvoiceReportProcessor()

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	testTPPSPaidInvoiceReportFilePath := "pkg/services/invoice/fixtures/tpps_paid_invoice_report_testfile_tpps_pickup_dir.csv"

	err = processor.ProcessFile(appCtx, testTPPSPaidInvoiceReportFilePath, "")
	if err != nil {
		appCtx.Logger().Error("Error while processing TPPS Paid Invoice Report file", zap.String("path", testTPPSPaidInvoiceReportFilePath), zap.Error(err))
	}

	testTPPSPaidInvoiceData := `
1077-4079-3	2024-08-05	2024-08-05	421.87	DUPK	DUPK	10340	0.0311	321.57	1077-4079-cabd6371	2
1077-4079-3	2024-08-05	2024-08-05	421.87	DDP	DDP	10340	0.0097	100.3	1077-4079-a4e717fd	1
1208-5962-1	2024-08-05	2024-08-05	557	MS	MS	1	557	557	1208-5962-e0fb5863	1
8801-2773-2	2024-08-05	2024-08-05	2748.04	DOP	DOP	1	77.02	77.02	8801-2773-f2bb471e	1
8801-2773-2	2024-08-05	2024-08-05	2748.04	DPK	DPK	1	2671.02	2671.02	8801-2773-fdaee177	2
8801-2773-3	2024-08-05	2024-08-05	1397.74	DDP	DDP	1	91.31	91.31	8801-2773-2e54e07d	2
8801-2773-3	2024-08-05	2024-08-05	1397.74	DLH	DLH	1	1052.84	1052.84	8801-2773-27961d7f	1
8801-2773-3	2024-08-05	2024-08-05	1397.74	FSC	FSC	1	6.66	6.66	8801-2773-f9e0672c	3
8801-2773-3	2024-08-05	2024-08-05	1397.74	DUPK	DUPK	1	246.93	246.93	8801-2773-c6c78cf9	4

`

	appCtx.Logger().Info("The tpps_paid_invoice_reports table should now have the following data: ")
	appCtx.Logger().Info(testTPPSPaidInvoiceData)
}
