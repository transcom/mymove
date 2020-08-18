package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
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

func checkDatabase(db *pop.Connection, logger Logger, runtime *WebhookRuntime) error {
	// Read notifications
	notifications := []models.WebhookNotification{}
	err := db.All(&notifications)
	if err != nil {
		logger.Error("Error:", zap.Error(err))
		return err
	}
	logger.Info("Success!", zap.Int("Num notifications found", len(notifications)))

	// Construct msg to send
	message := "There were no notifications."
	if len(notifications) > 0 {
		not := notifications[0]
		message = fmt.Sprintf("There was a %s notification, id: %s, moveTaskOrderID: %s",
			string(not.EventKey), not.ObjectID.String(), not.MoveTaskOrderID.String())
	}

	payload := &WebhookRequest{
		Object: message,
	}
	json, err := json.Marshal(payload)

	if err != nil {
		logger.Error("Error creating payload:", zap.Error(err))
		return err
	}

	// Make the API call
	// Temp url to use when making requests
	basePath := "/support/v1/webhook-notify"

	runtime.BasePath = basePath
	resp, body, err := runtime.Post(json)
	// Check for error making the request
	if err != nil {
		logger.Error("Error making request:", zap.Error(err))
		return err
	}
	// Check for error response from server
	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("Received %d response from server: %s. Body: %s", resp.StatusCode, resp.Status, body)
		err = errors.New(errmsg)
		logger.Error("db-webhook-notify:", zap.Error(err))
		return err
	}

	logger.Info("Request complete: ", zap.String("Status", resp.Status))
	fmt.Println(fmt.Sprintf("%s", body))

	return nil
}

func dbListen(db *pop.Connection, logger Logger, runtime *WebhookRuntime, periodInSeconds int) error {

	t := time.Tick(time.Duration(periodInSeconds) * time.Second)
	r := rand.New(rand.NewSource(99)) /// todo delete

	checkDatabase(db, logger, runtime)
	for range t {
		//Download the current contents of the URL and do something with it
		checkDatabase(db, logger, runtime)
		jitterTime := r.Int31n(2000)
		fmt.Printf("Checked at at %s, add jitter %d s\n", time.Now(), jitterTime)

		// add a bit of jitter
		jitter := time.Duration(jitterTime) * time.Millisecond
		time.Sleep(jitter)
	}
	return nil
}

func dbWebhookNotify(cmd *cobra.Command, args []string) error {
	v := viper.New()

	err := ParseFlags(cmd, v, args)
	if err != nil {
		return err
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		changed := ""
		if f.Changed {
			changed = "user-defined"
		}
		fmt.Println(f.Name, ": \t", f.Value, " \t", changed)
	})

	db, logger, err := InitRootConfig(v)
	if err != nil {
		logger.Error("invalid configuration", zap.Error(err))
		return err
	}

	// Create the client and open the cacStore
	runtime, cacStore, errCreateClient := CreateClient(v)

	if errCreateClient != nil {
		logger.Error("Error creating client:", zap.Error(errCreateClient))
		return errCreateClient
	}

	// Defer closing the store until after the api call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	go dbListen(db, logger, runtime, v.GetInt(PeriodFlag))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	// Creates an empty context object
	if err = db.Close(); err == nil {
		logger.Info("Db connection closed")
	} else {
		logger.Error("Db connection close failed", zap.Error(err))
	}

	log.Println("Listener exiting")
	return nil
}
