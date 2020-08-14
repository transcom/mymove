package main

import (
	"encoding/json"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
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
	// #TODO: move logger or pass it into setup
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
	runtime, cacStore, err := CreateClient(v)

	if err != nil {
		logger.Fatal(err.Error())
	}

	if cacStore != nil {
		defer cacStore.Close()
	}
	// Make the API call

	runtime.Post(json)

	return nil
}
