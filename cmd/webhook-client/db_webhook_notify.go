package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

const (
	// PeriodFlag indicates how often to check the db in seconds
	PeriodFlag string = "period"
)

func initDbWebhookNotifyFlags(flag *pflag.FlagSet) {
	flag.Int(PeriodFlag, 5, "Period in secs to check for notifications")
	flag.SortFlags = false
}

// DBPoller encapsulates the services used by the push notification engine
type DBPoller struct {
	Connection      *pop.Connection
	Logger          Logger
	Client          *WebhookRuntime
	Cmd             *cobra.Command
	PeriodInSeconds int
}

func (dbPoller *DBPoller) processPendingNotifications(notifications []models.WebhookNotification, subscriptions []models.WebhookSubscription) {
	for _, not := range notifications {

		// search for subscription
		foundSub := false
		for _, sub := range subscriptions {
			if sub.EventKey == not.EventKey {
				fmt.Println("Found sub for ", not.EventKey)
				foundSub = true
				// If found, send notification to subscription
				err := dbPoller.sendOneNotification(&not, &sub)
				if err != nil {
					dbPoller.Logger.Error("Notification should be updated as sent")
				}

				not.Status = models.WebhookNotificationSent
				dbPoller.Logger.Info("Notification should be updated as sent")
			}
		}
		if foundSub == false {
			dbPoller.Logger.Debug("No subscription found for notification event.", zap.String("eventKey", not.EventKey))
		}
	}
}

// sendPushNotication sends the notification and marks the notification as sent or failed
func (dbPoller *DBPoller) sendOneNotification(notif *models.WebhookNotification, sub *models.WebhookSubscription) error {
	logger := dbPoller.Logger

	// Construct notification to send
	message := fmt.Sprintf("There was a %s notification, id: %s, moveTaskOrderID: %s",
		string(notif.EventKey), notif.ID.String(), notif.MoveTaskOrderID.String())
	payload := &WebhookRequest{
		Object: message,
	}
	json, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error creating payload:", zap.Error(err))
		return err
	}

	// Post the notification
	url := sub.CallbackURL
	resp, body, err := dbPoller.Client.Post(json, url)

	// Check for error making the request
	if err != nil {
		logger.Error("Failed to send, error making request:", zap.Error(err))
		return err
	}
	// Check for error response from server
	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("Failed to send. Response Status: %s. Body: %s", resp.Status, body)
		err = errors.New(errmsg)
		logger.Error("db-webhook-notify: Failed to send notification", zap.Error(err))
		return err
	}

	logger.Info("Notification successfully sent: ", zap.String("Status", resp.Status), zap.String("Body", string(body)))
	return nil

}

func (dbPoller *DBPoller) checkDatabase() error {

	logger := dbPoller.Logger
	// Read all notifications
	notifications := []models.WebhookNotification{}
	err := dbPoller.Connection.Order("created_at asc").Where("status = ?", models.WebhookNotificationPending).All(&notifications)

	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	}
	logger.Debug("Success!", zap.Int("Num notifications found", len(notifications)))

	// If none, return
	if len(notifications) == 0 {
		return nil
	}

	// If there are notifications, get subscriptions
	subscriptions := []models.WebhookSubscription{}
	err = dbPoller.Connection.Where("status = ?", models.WebhookSubscriptionStatusActive).All(&subscriptions)

	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	}
	logger.Debug("Success!", zap.Int("Num subscriptions found", len(subscriptions)))

	// If none, return
	if len(notifications) == 0 {
		return nil
	}

	// If found, store in memory todo
	// process notifications
	dbPoller.processPendingNotifications(notifications, subscriptions)
	return nil
}

func dbListen(dbPoller *DBPoller) error {

	t := time.Tick(time.Duration(dbPoller.PeriodInSeconds) * time.Second)

	dbPoller.checkDatabase()
	for range t {
		//Download the current contents of the URL and do something with it
		dbPoller.checkDatabase()
	}
	return nil
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
	err := ParseFlags(cmd, v, args)
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
	runtime, cacStore, errCreateClient := CreateClient(v)

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
	dbPoller := DBPoller{
		Connection:      db,
		Logger:          logger,
		Client:          runtime,
		PeriodInSeconds: v.GetInt(PeriodFlag),
	}

	// Start polling the db for changes
	go dbListen(&dbPoller)

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
