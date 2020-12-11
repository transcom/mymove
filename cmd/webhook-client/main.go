package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/pop/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/cmd/webhook-client/utils"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// initRootFlags initializes flags relating to the webhook client
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	// Logging Levels
	cli.InitLoggingFlags(flag)
	// DB Config
	cli.InitDatabaseFlags(flag)

	// Additional flags pertinent to all commands using this tool
	flag.String(utils.CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(utils.KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(utils.HostnameFlag, cli.HTTPPrimeServerNameLocal, "The hostname to connect to")
	flag.Int(utils.PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(utils.InsecureFlag, false, "Skip TLS verification and validation")
}

// InitRootConfig checks the validity of the api flags and initializes a db connection.
func InitRootConfig(v *viper.Viper) (*pop.Connection, utils.Logger, error) {

	// LOGGER SETUP
	// Get the db env to configure the logger level
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, err := logging.Config(
		logging.WithEnvironment(dbEnv),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)
	logger.Info("Checking config and initializing")

	// FLAG CHECKS
	err = cli.CheckDatabase(v, logger)
	if err != nil {
		return nil, logger, err
	}

	err = cli.CheckCAC(v)
	if err != nil {
		return nil, logger, err
	}

	err = cli.CheckLogging(v)
	if err != nil {
		return nil, logger, err
	}

	if (v.GetString(utils.CertPathFlag) != "" && v.GetString(utils.KeyPathFlag) == "") ||
		(v.GetString(utils.CertPathFlag) == "" && v.GetString(utils.KeyPathFlag) != "") {
		return nil, logger, fmt.Errorf("Both TLS certificate and key paths must be provided")
	}

	// DB CONNECTION CHECK
	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	return dbConnection, logger, nil
}

func main() {
	// Root command
	root := cobra.Command{
		Use:   "webhook-client [flags]",
		Short: "Webhook client",
		Long:  "Webhook client",
	}
	initRootFlags(root.PersistentFlags())

	// Sub-commands
	postWebhookNotifyCommand := &cobra.Command{
		Use:          "post-webhook-notify",
		Short:        "Post Webhook Notify",
		Long:         "Post Webhook Notify",
		RunE:         postWebhookNotify,
		SilenceUsage: true,
	}
	initPostWebhookNotifyFlags(postWebhookNotifyCommand.Flags())
	root.AddCommand(postWebhookNotifyCommand)

	dbWebhookNotifyCommand := &cobra.Command{
		Use:   "db-webhook-notify",
		Short: "Database Webhook Notify",
		Long: `
	Database Webhook Notify checks the webhook_notification
	table and sends the first notification it finds there.`,
		RunE:         dbWebhookNotify,
		SilenceUsage: true,
	}
	initDbWebhookNotifyFlags(dbWebhookNotifyCommand.Flags())
	root.AddCommand(dbWebhookNotifyCommand)

	dbConnectionCommand := &cobra.Command{
		Use:   "db-connection-test",
		Short: "Database Connection Test",
		Long: `
	Database Connection Test creates, updates and deletes a
	record in the webhook_notification and webhook_subscription
	tables.`,
		RunE:         dbConnection,
		SilenceUsage: true,
	}
	initDbConnectionFlags(dbConnectionCommand.Flags())
	root.AddCommand(dbConnectionCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\nwebhook-client completion > /usr/local/etc/bash_completion.d/webhook-client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
