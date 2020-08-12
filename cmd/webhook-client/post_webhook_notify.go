package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	//"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	//webhookOperations "github.com/transcom/mymove/pkg/gen/supportclient/webhook"
)

/*WebhookMessage post webhook notify body
swagger:model PostWebhookNotifyBody
*/
type WebhookMessage struct {

	// Message sent
	// Required: true
	Message string `json:"message"`
}

// CreateWebhookPayload returns message to send
func CreateWebhookPayload(message string) WebhookMessage {
	return WebhookMessage{
		Message: message,
	}
}

const (
	// MessageFlag could be moved out to utils folder later
	MessageFlag string = "message"
)

func initPostWebhookNotifyFlags(flag *pflag.FlagSet) {
	flag.String(MessageFlag, "", "Message to send")

	flag.SortFlags = false
}

func checkPostWebhookNotifyConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	message := v.GetString(MessageFlag)
	if len(message) == 0 {
		return errors.New("missing message, expected to be set")
	}

	return nil
}

func postWebhookNotify(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkPostWebhookNotifyConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	message := v.GetString(MessageFlag)
	//#TODO: To remove dependency on gen/supportclient,
	// replicate the functionality without using webhookOperations
	// newNotification := webhookOperations.PostWebhookNotifyBody{
	// 	Message: message,
	// }
	// //#TODO: To remove dependency on gen/supportclient,
	// // replicate the functionality without using webhookOperations
	// notifyParams := webhookOperations.NewPostWebhookNotifyParams()

	// notifyParams.WithMessage(newNotification)
	// notifyParams.SetTimeout(time.Second * 30)
	payload := CreateWebhookPayload(message)
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
	// Make the API call
	// resp, err := supportGateway.Webhook.PostWebhookNotify(notifyParams)
	r, err := client.Post("https://primelocal:9443/support/v1/webhook-notify", "application/json; charset=utf-8", bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Print response body to stdout
	fmt.Printf("%s\n", body)
	// payload := r.GetPayload()
	// if payload != nil {
	// 	payload, errJSONMarshall := json.Marshal(payload)
	// 	if errJSONMarshall != nil {
	// 		logger.Fatal(errJSONMarshall)
	// 	}
	// 	fmt.Println(string(payload))
	// } else {
	// 	logger.Fatal(resp.Error())
	// }

	return nil
}
