package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	webhookOperations "github.com/transcom/mymove/pkg/gen/supportclient/webhook"
	"github.com/transcom/mymove/pkg/models"
)

func initDbWebhookNotifyFlags(flag *pflag.FlagSet) {

	flag.SortFlags = false
}

func dbWebhookNotify(cmd *cobra.Command, args []string) error {
	v := viper.New()

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

	//#TODO: To remove dependency on gen/supportclient,
	// replicate the functionality without using webhookOperations
	newNotification := webhookOperations.PostWebhookNotifyBody{
		Message: message,
	}

	//#TODO: To remove dependency on gen/supportclient,
	// replicate the functionality without using webhookOperations
	notifyParams := webhookOperations.NewPostWebhookNotifyParams()

	notifyParams.WithMessage(newNotification)
	notifyParams.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	supportGateway, cacStore, errCreateClient := CreateClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the api call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Make the API call
	resp, err := supportGateway.Webhook.PostWebhookNotify(notifyParams)
	if err != nil {
		return HandleGatewayError(err, logger)
	}

	payload := resp.GetPayload()

	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			logger.Error("Error", zap.Error(errJSONMarshall))
			return errJSONMarshall
		}
		logger.Info("Success!", zap.String("Payload", string(payload)))
	} else {
		logger.Error("Error:", zap.String("Error", resp.Error()))
		return errors.New(resp.Error())
	}

	return nil
}
