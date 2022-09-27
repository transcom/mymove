package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/cmd/webhook-client/utils"
)

// WebhookRequest is the body of our request
type WebhookRequest struct {
	ID          string `json:"id"`
	EventName   string `json:"eventName"`
	TriggeredAt string `json:"triggeredAt"`
	ObjectType  string `json:"objectType"`
	Object      string `json:"object"`
}

// Flags specific to this command
const (
	// FilenameFlag is the string to send in the payload
	FilenameFlag string = "filename"
)

func initPostWebhookNotifyFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Filename of json file to send")
	flag.SortFlags = false
}

func checkPostWebhookNotifyConfig(v *viper.Viper, args []string, logger *zap.Logger) error {

	missingFilenameFlag := v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args))

	if missingFilenameFlag {
		logger.Fatal("post-webhook-notify expects --filename with json file passed in")
	}

	return nil
}

func postWebhookNotify(cmd *cobra.Command, args []string) error {
	v := viper.New()

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Validate all arguments passed in including DB, CAC, etc...
	// Also this opens the db connection and creates a logger
	_, logger, err := InitRootConfig(v)
	if err != nil {
		logger.Error("invalid configuration", zap.Error(err))
		return err
	}

	// Check the config before talking to the CAC
	err = checkPostWebhookNotifyConfig(v, args, logger)
	if err != nil {
		logger.Fatal("Error:", zap.Error(err))
	}

	// Decode the json file that was passed in
	filename := v.GetString(FilenameFlag)
	payload := &WebhookRequest{}
	err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &payload)

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
	runtime, cacStore, errCreateClient := utils.CreateClient(v)

	if errCreateClient != nil {
		logger.Error("Error creating runtime client:", zap.Error(errCreateClient))
		return errCreateClient
	}

	if cacStore != nil {
		defer func() {
			if closeErr := cacStore.Close(); closeErr != nil {
				logger.Error("Error closing CAC connection", zap.Error(closeErr))
			}
		}()
	}

	// Make the API call
	hostname := v.GetString(utils.HostnameFlag)
	port := v.GetInt(utils.PortFlag)
	// For now this is hardcoded to our support endpoint
	path := "support/v1/webhook-notify"

	url := fmt.Sprintf("https://%s:%d/%s", hostname, port, path)
	resp, body, err := runtime.Post(json, url)

	if err != nil {
		logger.Error("Error making request:", zap.Error(err))
		return err
	}
	// Check for error response from server
	if resp.StatusCode != 200 {
		errmsg := fmt.Sprintf("Failed to send. Response Status: %s. Body: %s", resp.Status, string(body))
		err = errors.New(errmsg)
		return err
	}
	logger.Info("Request complete: ", zap.String("Status", resp.Status))

	return nil
}
