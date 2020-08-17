package main

import (
	"encoding/json"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//WebhookMessage is the body of our request
type WebhookMessage struct {
	// Message sent
	// Required: true
	Message string `json:"message"`
}

func initPostWebhookNotifyFlags(flag *pflag.FlagSet) {
	flag.String(MessageFlag, "", "Message to send")

	flag.SortFlags = false
}

func checkPostWebhookNotifyConfig(v *viper.Viper, args []string, logger Logger) error {
	_, _, err := InitRootConfig(v)
	if err != nil {
		logger.Fatal(err.Error())
	}

	message := v.GetString(MessageFlag)
	if len(message) == 0 {
		return errors.New("missing message, expected to be set")
	}

	return nil
}

func postWebhookNotify(cmd *cobra.Command, args []string) error {
	v := viper.New()
	// basePath represents url where we're sending our request
	// For now this is hardcoded to our support endpoint
	basePath := "/support/v1/webhook-notify"

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	_, logger, err := InitRootConfig(v)
	if err != nil {
		logger.Fatal("Invalid configuration", zap.Error(err))
	}

	// Check the config before talking to the CAC
	err = checkPostWebhookNotifyConfig(v, args, logger)
	if err != nil {
		logger.Fatal("Error:", zap.Error(err))
	}

	message := v.GetString(MessageFlag)
	payload := &WebhookMessage{
		Message: message,
	}
	json, err := json.Marshal(payload)

	if err != nil {
		logger.Error("Error creating payload:", zap.Error(err))
		return err
	}

	// Create the client and open the cacStore
	runtime, cacStore, errCreateClient := CreateClient(v)

	if errCreateClient != nil {
		logger.Error("Error creating runtime client:", zap.Error(errCreateClient))
		return errCreateClient
	}

	if cacStore != nil {
		defer cacStore.Close()
	}

	// Make the API call
	runtime.BasePath = basePath
	resp, err := runtime.Post(json)

	if err != nil {
		logger.Error("Error making request:", zap.Error(err))
		return err
	}

	logger.Info("Request complete: ", zap.String("Status", resp.Status))

	return nil
}
