package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/cmd/webhook-client/utils"
	"github.com/transcom/mymove/cmd/webhook-client/webhook"

	"go.uber.org/zap"
)

const (
	// PeriodFlag indicates how often to check the db in seconds
	PeriodFlag string = "period"
)

// Init flags specific to this command
func initDbWebhookNotifyFlags(flag *pflag.FlagSet) {
	flag.Int(PeriodFlag, 5, "Period in secs to check for notifications")
	flag.SortFlags = false
}

// func populateTestNotifications(db *pop.Connection, logger Logger) error {
// 	subID, _ := uuid.FromString("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")
// 	subscription := testdatagen.MakeWebhookSubscription(db, testdatagen.Assertions{
// 		WebhookSubscription: models.WebhookSubscription{
// 			ID:           uuid.Must(uuid.NewV4()),
// 			SubscriberID: subID,
// 			EventKey:     "PaymentRequest.Create",
// 			CallbackURL:  "/support/v1/webhook-notify",
// 		},
// 	})
// 	_, err := db.ValidateAndSave(&subscription)
// 	return err
// }
func dbWebhookNotify(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Parse flags from command line
	err := utils.ParseFlags(cmd, v, args)
	if err != nil {
		return err
	}

	// Print all flag values
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		changed := ""
		if f.Changed {
			changed = "user-defined"
		}
		fmt.Println(f.Name, ": \t", f.Value, " \t", changed)
	})

	// Validate all arguments passed in including DB, CAC, etc...
	// Also this opens the db connection and creates a logger
	db, logger, err := InitRootConfig(v)
	if err != nil {
		logger.Error("invalid configuration", zap.Error(err))
		return err
	}

	// Create the client and open the cacStore (which has to stay open?)
	runtime, cacStore, errCreateClient := utils.CreateClient(v)

	if errCreateClient != nil {
		logger.Error("Error creating client:", zap.Error(errCreateClient))
		return errCreateClient
	}

	// Defer closing the store until after the api call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Create an DBPoller engine
	//populateTestNotifications(db, logger)
	webhookEngine := webhook.Engine{
		Connection:      db,
		Logger:          logger,
		Client:          runtime,
		PeriodInSeconds: v.GetInt(PeriodFlag),
	}

	// Start polling the db for changes
	go webhookEngine.Start()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	if err = db.Close(); err == nil {
		logger.Info("Db connection closed")
	} else {
		logger.Error("Db connection close failed", zap.Error(err))
	}

	log.Println("Listener exiting")
	return nil
}
