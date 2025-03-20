package main

import (
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
)

func checkPaymentReminderConfig(v *viper.Viper, logger *zap.Logger) error {

	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return cli.CheckEmail(v)
}

func initPaymentReminderFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Email
	cli.InitEmailFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Command (test eamil): go run ./cmd/milmove-tasks send-payment-reminder
// Command (send email): go run ./cmd/milmove-tasks send-payment-reminder --email-backend=ses --aws-ses-domain=devlocal.dp3.us --aws-ses-region=us-gov-west-1
func sendPaymentReminder(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse args")
	}
	flags := cmd.Flags()
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
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
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkPaymentReminderConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil, nil)

	notificationSender, notificationSenderErr := notifications.InitEmail(v, logger)
	if notificationSenderErr != nil {
		logger.Fatal("notification sender sending not enabled", zap.Error(notificationSenderErr))
	}

	movePaymentReminderNotifier, err := notifications.NewPaymentReminder()
	if err != nil {
		logger.Fatal("initializing MoveReviewed", zap.Error(err))
	}

	err = notificationSender.SendNotification(appCtx, movePaymentReminderNotifier)
	if err != nil {
		logger.Fatal("Emails failed to send", zap.Error(err))
	}
	return nil
}
