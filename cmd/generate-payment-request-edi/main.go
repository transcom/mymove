package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
	// "github.com/transcom/mymove/pkg/edi/invoice"
)

// Call this from command line with go run cmd/generate-payment-request-edi/main.go
// Must use a payment request that is submitted, but not yet approved for payment (that does not already have a submitted invoice)

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// // Scenario config
	flag.String("payment-request-number", "", "The payment request number to generate an EDI for")
	flag.String("payment-request-uuid", "", "The payment request UUID to generate an EDI for")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

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
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	fmt.Println("logger: ", logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// DB connection
	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	paymentRequestUUID := v.GetString("payment-request-uuid")
	var paymentRequest models.PaymentRequest
	err = dbConnection.Find(&paymentRequest, uuid.Must(uuid.FromString(paymentRequestUUID)))
	if err != nil {
		logger.Fatal(err.Error())
	}

	result, err := generator.Generate(paymentRequest, false)

	fmt.Println(result)
}

// "540e2268-6899-4b67-828d-bb3b0331ecf2"
