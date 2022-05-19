package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/gobuffalo/pop/v6"
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
	flag.String(utils.CertPathFlag, "", "Path to the public cert")
	flag.String(utils.KeyPathFlag, "", "Path to the private key")
	flag.String(utils.HostnameFlag, cli.HTTPPrimeServerNameLocal, "The hostname to connect to")
	flag.Int(utils.PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(utils.InsecureFlag, false, "Skip TLS verification and validation")
}

// InitRootConfig checks the validity of the api flags and initializes a db connection.
func InitRootConfig(v *viper.Viper) (*pop.Connection, *zap.Logger, error) {

	// LOGGER SETUP
	// Get the db env to configure the logger level
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, _, err := logging.Config(
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

	var canUseFileCert bool
	if v.GetString(utils.CertPathFlag) != "" && v.GetString(utils.KeyPathFlag) == "" {
		return nil, logger, fmt.Errorf(
			"Must set the %v parameter if %v is set",
			utils.KeyPathFlag,
			utils.CertPathFlag,
		)
	} else if v.GetString(utils.CertPathFlag) == "" && v.GetString(utils.KeyPathFlag) != "" {
		return nil, logger, fmt.Errorf(
			"Must set the %v parameter if %v is set",
			utils.CertPathFlag,
			utils.KeyPathFlag,
		)
	} else if v.GetString(utils.CertPathFlag) != "" && v.GetString(utils.KeyPathFlag) != "" {
		canUseFileCert = true
	}

	var canUseEnvCert bool
	recipientMTLSCertEnvName := strings.ToUpper(strings.ReplaceAll(utils.RecipientMTLSCert, "-", "_"))
	recipientMTLSKeyEnvName := strings.ToUpper(strings.ReplaceAll(utils.RecipientMTLSKey, "-", "_"))
	if v.GetString(utils.RecipientMTLSCert) != "" && v.GetString(utils.RecipientMTLSKey) == "" {
		return nil, logger, fmt.Errorf(
			"Must set the %v environment variable if %v is set",
			recipientMTLSKeyEnvName,
			recipientMTLSCertEnvName,
		)
	} else if v.GetString(utils.RecipientMTLSCert) == "" && v.GetString(utils.RecipientMTLSKey) != "" {
		return nil, logger, fmt.Errorf(
			"Must set the %v environment variable if %v is set",
			recipientMTLSCertEnvName,
			recipientMTLSKeyEnvName,
		)
	} else if v.GetString(utils.RecipientMTLSCert) != "" && v.GetString(utils.RecipientMTLSKey) != "" {
		canUseEnvCert = true
	}

	if !canUseFileCert && !canUseEnvCert {
		return nil, logger, fmt.Errorf(
			"Must provide %v & %v parameters or set the %v & %v environment variables",
			utils.CertPathFlag,
			utils.KeyPathFlag,
			recipientMTLSCertEnvName,
			recipientMTLSKeyEnvName,
		)
	}
	if canUseFileCert && canUseEnvCert {
		logger.Info("A certificate is configured to be loaded from both the filesystem and environment; defaulting to the filesystem certificate")
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) {
		c := &aws.Config{
			Region: aws.String(v.GetString(cli.AWSRegionFlag)),
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal("error creating aws session", zap.Error(errorSession))
		}
		session = s
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		if session != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly,
			// because the session conflates the environment, shared, and container metadata config
			// within NewSession.  With stscreds, we use the Secure Token Service,
			// to assume the given role (that has rds db connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)
			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
			dbCreds = stscreds.NewCredentials(session, dbIamRole)
		}
	}

	// DB CONNECTION CHECK
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
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

	webhookNotifyCommand := &cobra.Command{
		Use:   "webhook-notify",
		Short: "Webhook Notify",
		Long: `
	Webhook Notify launches the engine for webhook notifications.
	This repeatedly checks the webhook_notification and webhook_subscription tables and
	sends the notifications every minute.`,
		RunE:         webhookNotify,
		SilenceUsage: true,
	}
	initWebhookNotifyFlags(webhookNotifyCommand.Flags())
	root.AddCommand(webhookNotifyCommand)

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
