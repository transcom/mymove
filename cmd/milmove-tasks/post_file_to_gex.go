package main

import (
	"crypto/tls"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
)

func checkPostFileToGEXConfig(v *viper.Viper) error {

	if err := cli.CheckGEX(v); err != nil {
		return err
	}

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	if err := cli.CheckEntrustCert(v); err != nil {
		return err
	}

	if ediFile := v.GetString("gex-helloworld-file"); ediFile == "" {
		return errors.New("must have file to send")
	}

	if filename := v.GetString("filename"); filename == "" {
		return errors.New("filename is missing")
	}

	return nil
}

func initPostFileToGEXFlags(flag *pflag.FlagSet) {
	// Logging Levels
	cli.InitLoggingFlags(flag)

	// GEX
	cli.InitGEXFlags(flag)

	// Certificate
	cli.InitCertFlags(flag)

	// Entrust Certificates
	cli.InitEntrustCertFlags(flag)

	flag.String("gex-helloworld-file", "", "GEX file to post")
	flag.String("filename", "test", "The required name sent in the url of the gex api request")

	// Don't sort flags
	flag.SortFlags = false
}

// go run ./cmd/milmove-tasks post-file-to-gex --edi filepath --filename filename --gex-url 'url'
func postFileToGEX(cmd *cobra.Command, args []string) error {
	// Create the logger
	v := viper.New()

	logger, _, err := logging.Config(
		logging.WithEnvironment(v.GetString(cli.LoggingEnvFlag)),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		logger.Fatal("Failed to initialize Zap logging", zap.Error(err))
	}

	flag := pflag.CommandLine
	initPostFileToGEXFlags(flag)
	err = flag.Parse(os.Args[1:])
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	pflagsErr := v.BindPFlags(flag)
	if pflagsErr != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	checkConfigErr := checkPostFileToGEXConfig(v)
	if checkConfigErr != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// dbEnv := v.GetString(cli.DbEnvFlag)

	edi := v.GetString("gex-helloworld-file")

	ediString := string(edi[:])
	// make sure edi ends in new line
	ediString = strings.TrimSpace(ediString) + "\n"

	certLogger, _, err := logging.Config(logging.WithEnvironment("development"), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		logger.Fatal("Failed to initialize Zap logging", zap.Error(err))
	}
	certificates, rootCAs, err := certs.InitDoDEntrustCertificates(v, certLogger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Error in getting tls certs", zap.Error(err))
	}

	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}

	logger.Info(
		"Sending to GEX",
		zap.String("filename", v.GetString("filename")),
		zap.String("url", v.GetString(cli.GEXURLFlag)))

	resp, err := invoice.NewGexSenderHTTP(
		v.GetString(cli.GEXURLFlag),
		true,
		tlsConfig,
		v.GetString(cli.GEXBasicAuthUsernameFlag),
		v.GetString(cli.GEXBasicAuthPasswordFlag),
	).SendToGex(services.GEXChannelInvoice, ediString, v.GetString("filename"))

	if err != nil {
		logger.Fatal("Gex Sender encountered an error", zap.Error(err))
	}

	if resp == nil {
		logger.Fatal("Gex Sender had no response", zap.Error(err))
	} else {
		logger.Info(
			"Posted to GEX",
			zap.String("filename", v.GetString("filename")),
			zap.Int("statusCode", resp.StatusCode),
			zap.Error(err))
	}

	return nil
}
