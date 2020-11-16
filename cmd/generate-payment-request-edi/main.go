package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from command line with go run ./cmd/generate-payment-request-edi/ --payment-request-number <paymentRequestNumber>
// Must use a payment request that is submitted, but not yet approved for payment (that does not already have a submitted invoice)

func checkConfig(v *viper.Viper, logger logger) error {
	paymentRequestNumber := v.GetString("payment-request-number")
	if paymentRequestNumber == "" {
		return errors.New("must provide payment-request-number")
	}

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {
	// This command's config
	flag.String("payment-request-number", "", "The payment request number to generate an EDI for")

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
	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	paymentRequestNumber := v.GetString("payment-request-number")

	var paymentRequest models.PaymentRequest
	err = dbConnection.Where("payment_request_number = ?", paymentRequestNumber).First(&paymentRequest)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(os.Stderr, "ERROR: Could not find a payment request with number %s\n", paymentRequestNumber)
		} else {
			logger.Error(err.Error())
		}
		os.Exit(1)
	}

	generator := invoice.NewGHCPaymentRequestInvoiceGenerator(dbConnection)
	edi858c, err := generator.Generate(paymentRequest, false)
	if err != nil {
		logger.Fatal(err.Error())
	}

	edi858String, err := edi858c.EDIString(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	fmt.Print(edi858String)
}
