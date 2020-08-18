package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

func initDbWebhookNotifyFlags(flag *pflag.FlagSet) {

	flag.SortFlags = false
}

func dbWebhookNotify(cmd *cobra.Command, args []string) error {
	v := viper.New()
	// Temp url to use when making requests
	basePath := "/support/v1/webhook-notify"

	err := ParseFlags(cmd, v, args)
	if err != nil {
		return err
	}

	db, logger, err := InitRootConfig(v)
	if err != nil {
		logger.Error("invalid configuration", zap.Error(err))
		return err
	}

	// Read notifications
	notifications := []models.WebhookNotification{}
	err = db.All(&notifications)
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

	// Make the API call
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
	fmt.Println(body)
	return nil
}
