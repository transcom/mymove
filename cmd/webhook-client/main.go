package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// Logger type exports the logger for use in the command files
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

const (
	// CertPathFlag is the path to the client mTLS certificate
	CertPathFlag string = "certpath"
	// KeyPathFlag is the path to the key mTLS certificate
	KeyPathFlag string = "keypath"
	// HostnameFlag is the hostname to connect to
	HostnameFlag string = "hostname"
	// PortFlag is the port to connect to
	PortFlag string = "port"
	// InsecureFlag indicates that TLS verification and validation can be skipped
	InsecureFlag string = "insecure"
	// MessageFlag is the string to send in the payload
	MessageFlag string = "message"
	// FilenameFlag is the name of the file being passed in
	FilenameFlag string = "filename"
)

// initRootFlags initializes flags relating to the webhook client
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitVerboseFlags(flag)
	// DB Config
	cli.InitDatabaseFlags(flag)

	flag.String(CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(HostnameFlag, cli.HTTPPrimeServerNameLocal, "The hostname to connect to")
	flag.Int(PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(InsecureFlag, false, "Skip TLS verification and validation")
	flag.String(MessageFlag, "Hello World", "Message for the client to send")
	flag.String(FilenameFlag, "", "Data file passed in to the client")
}

// Debug prints helpful debugging information for requests
func Debug(data []byte, err error) {
	if err == nil {
		log.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

// InitRootConfig checks the validity of the api flags and initializes a db connection.
func InitRootConfig(v *viper.Viper) (*pop.Connection, Logger, error) {

	// LOGGER SETUP
	// Get the db env to configure the logger level
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
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

	err = cli.CheckVerbose(v)
	if err != nil {
		return nil, logger, err
	}

	if (v.GetString(CertPathFlag) != "" && v.GetString(KeyPathFlag) == "") || (v.GetString(CertPathFlag) == "" && v.GetString(KeyPathFlag) != "") {
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
