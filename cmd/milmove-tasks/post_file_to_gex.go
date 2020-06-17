package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/invoice"
)

func checkPostFileToGEXConfig(v *viper.Viper) error {

	if err := cli.CheckGEX(v); err != nil {
		return err
	}

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	if ediFile := v.GetString("gex-helloworld-file"); ediFile == "" {
		return errors.New("must have file to send")
	}

	if trasactionName := v.GetString("transaction-name"); trasactionName == "" {
		return errors.New("transaction-name is missing")
	}

	return nil
}

func initPostFileToGEXFlags(flag *pflag.FlagSet) {
	// Verbose
	cli.InitVerboseFlags(flag)

	// GEX
	cli.InitGEXFlags(flag)

	// Certificate
	cli.InitCertFlags(flag)

	flag.String("gex-helloworld-file", "", "GEX file to post")
	flag.String("transaction-name", "test", "The required name sent in the url of the gex api request")
	// flag.Parse(os.Args[1:])

	// Don't sort flags
	flag.SortFlags = false
}

func quit(logger *log.Logger, flag *pflag.FlagSet, err error) {
	if err != nil {
		logger.Println(err.Error())
	}
	logger.Println("Usage of send-to-gex:")
	if flag != nil {
		flag.PrintDefaults()
	}
	os.Exit(1)
}

// go run ./cmd/milmove-tasks post-file-to-gex --edi filepath --transaction-name transactionName --gex-url 'url'
func postFileToGEX(cmd *cobra.Command, args []string) error {
	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	flag := pflag.CommandLine
	initPostFileToGEXFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		quit(logger, flag, err)
	}

	v := viper.New()
	pflagsErr := v.BindPFlags(flag)
	if pflagsErr != nil {
		quit(logger, flag, err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	checkConfigErr := checkPostFileToGEXConfig(v)
	if checkConfigErr != nil {
		quit(logger, flag, checkConfigErr)
	}

	// dbEnv := v.GetString(cli.DbEnvFlag)

	edi := v.GetString("gex-helloworld-file")

	ediString := string(edi[:])
	// make sure edi ends in new line
	ediString = strings.TrimSpace(ediString) + "\n"

	certLogger, err := logging.Config("development", true)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	certificates, rootCAs, err := certs.InitDoDCertificates(v, certLogger)
	if certificates == nil || rootCAs == nil || err != nil {
		log.Fatal("Error in getting tls certs", err)
	}

	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs}

	logger.Println("Sending to GEX ...")
	resp, err := invoice.NewGexSenderHTTP(
		v.GetString("gex-url"),
		true,
		tlsConfig,
		v.GetString("gex-basic-auth-username"),
		v.GetString("gex-basic-auth-password"),
	).SendToGex(ediString, v.GetString("transaction-name"))

	if err != nil {
		log.Fatalf("Gex Sender encountered an error: %v", err)
	}

	if resp == nil {
		log.Fatal("Gex Sender had no response")
	}

	fmt.Printf("status code: %v, error: %v \n", resp.StatusCode, err)

	return nil
}
