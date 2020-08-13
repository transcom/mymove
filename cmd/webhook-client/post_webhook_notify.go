package main

import (
	// "bytes"
	"encoding/json"
	"errors"

	// "fmt"
	// "io/ioutil"
	// "log"
	// "net/http/httputil"
	// "os"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

//WebhookMessage post webhook notify body
type WebhookMessage struct {
	// Message sent
	// Required: true
	Message string `json:"message"`
}

const (
	// MessageFlag could be moved out to utils folder later
	MessageFlag string = "message"
)

func initPostWebhookNotifyFlags(flag *pflag.FlagSet) {
	flag.String(MessageFlag, "", "Message to send")

	flag.SortFlags = false
}

func checkPostWebhookNotifyConfig(v *viper.Viper, args []string, logger Logger) error {
	err := CheckRootConfig(v)
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

	// Create the logger
	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err = checkPostWebhookNotifyConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	message := v.GetString(MessageFlag)
	payload := &WebhookMessage{
		Message: message,
	}
	json, err := json.Marshal(payload)

	if err != nil {
		panic(err)
	}

	// Create the client and open the cacStore
	client, cacStore, errCreateClient := CreateClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}
	// Defer closing the store until after the api call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}
	runtime := WebhookRuntime{
		client:      client,
		Host:        "https://primelocal:9443",
		BasePath:    "/support/v1/webhook-notify",
		Debug:       true,
		Logger:      logger,
		ContentType: "application/json; charset=utf-8",
	}

	// Make the API call

	runtime.Post(json)

	return nil
}
