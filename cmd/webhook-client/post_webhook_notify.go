package main

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//WebhookRequest is the body of our request
type WebhookRequest struct {
	ID          string `json:"id"`
	EventName   string `json:"eventName"`
	TriggeredAt string `json:"triggeredAt"`
	ObjectType  string `json:"objectType"`
	Object      string `json:"object"`
}

func initPostWebhookNotifyFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Path to the file with the payment request JSON payload")

	flag.SortFlags = false
}

func checkPostWebhookNotifyConfig(v *viper.Viper, args []string, logger Logger) error {
	_, _, err := InitRootConfig(v)
	if err != nil {
		logger.Fatal(err.Error())
	}

	missingFilenameFlag := v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !ContainsDash(args))

	if missingFilenameFlag {
		logger.Fatal("post-webhook-notify expects --filename with json file passed in")
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

	// Decode the json file that was passed in
	filename := v.GetString(FilenameFlag)
	payload := &WebhookRequest{}
	err = DecodeJSONFileToPayload(filename, ContainsDash(args), &payload)

	if err != nil {
		logger.Error("Error opening file:", zap.Error(err))
		return err
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
