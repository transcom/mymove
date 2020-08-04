package main

import (
	"encoding/json"
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

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	db, logger, err := InitRootConfig(v)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	notifications := []models.WebhookNotification{}
	err = db.All(&notifications)
	if err != nil {
		fmt.Print("ERROR!\n")
		fmt.Printf("%v\n", err)
	} else {

		fmt.Printf("Success! %d notifications found.\n", len(notifications))
	}
	message := "There were no notifications."
	if len(notifications) > 0 {
		message = fmt.Sprintf("There was a %s notification", string(notifications[0].EventKey))
	}

	//	message := v.GetString(MessageFlag)
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
		logger.Fatal("Error:", zap.Error(err))
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			logger.Fatal("Error", zap.Error(errJSONMarshall))
		}
		fmt.Println("payload", string(payload))
	} else {
		logger.Fatal("Error:", zap.String("Error", resp.Error()))
	}

	return nil
}
