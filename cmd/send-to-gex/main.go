package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
)

func initFlags(flag *pflag.FlagSet) {

	flag.String(cli.GEXBasicAuthUsernameFlag, "", "GEX api auth username")
	flag.String(cli.GEXBasicAuthPasswordFlag, "", "GEX api auth password")
	flag.String(cli.GEXURLFlag, "", "URL for sending an HTTP POST request to GEX")

	flag.StringSlice(cli.DoDCAPackageFlag, []string{}, "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
	flag.String(cli.MoveMilDoDCACertFlag, "", "The DoD CA certificate used to sign the move.mil TLS certificate.")
	flag.String(cli.MoveMilDoDTLSCertFlag, "", "The DoD-signed TLS certificate for various move.mil services.")
	flag.String(cli.MoveMilDoDTLSKeyFlag, "", "The private key for the DoD-signed TLS certificate for various move.mil services.")

	flag.String("edi", "", "The filepath to an edi file to send to GEX")
	flag.String("transaction-name", "test", "The required name sent in the url of the gex api request")

	// Don't sort flags
	flag.SortFlags = false
}

func checkConfig(v *viper.Viper) error {

	if err := cli.CheckGEX(v); err != nil {
		return err
	}

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	if ediFile := v.GetString("edi"); ediFile == "" {
		return errors.New("EDI file path must not be empty")
	}

	return nil
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

// Call this from command line with go run cmd/send-to-gex/main.go --edi <filepath>
func main() {
	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	flag := pflag.CommandLine
	initFlags(flag)
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

	checkConfigErr := checkConfig(v)
	if checkConfigErr != nil {
		quit(logger, flag, checkConfigErr)
	}

	ediFile := v.GetString("edi")

	file, err := os.Open(filepath.Clean(ediFile))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Fatalf("Failed to close file due to %v", closeErr)
		}
	}()

	edi, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	ediString := string(edi[:])
	// make sure edi ends in new line
	ediString = strings.TrimSpace(ediString) + "\n"

	logger.Println(ediString)

	certLogger, _, err := logging.Config(logging.WithEnvironment("development"), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	certificates, rootCAs, err := certs.InitDoDCertificates(v, certLogger)
	if certificates == nil || rootCAs == nil || err != nil {
		log.Fatal("Error in getting tls certs", err)
	}

	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}

	logger.Println("Sending to GEX ...")
	resp, err := invoice.NewGexSenderHTTP(
		v.GetString("gex-url"),
		true,
		tlsConfig,
		v.GetString("gex-basic-auth-username"),
		v.GetString("gex-basic-auth-password"),
	).SendToGex(services.GEXChannelInvoice, ediString, v.GetString("transaction-name"))

	if err != nil {
		log.Fatalf("Gex Sender encountered an error: %v", err)
	}

	statusCode := 0

	if resp == nil {
		log.Fatal("Gex Sender had no response")
	} else {
		statusCode = resp.StatusCode
	}

	fmt.Printf("status code: %v, error: %v \n", statusCode, err)
}
