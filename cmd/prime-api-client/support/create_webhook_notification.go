package support

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"

	webhookclient "github.com/transcom/mymove/pkg/gen/supportclient/webhook"
)

// InitCreateWebhookNotificationFlags initializes flags.
func InitCreateWebhookNotificationFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Path to the file with webhook notifications JSON payload")

	flag.SortFlags = false
}

// CheckCreateWebhookNotificationConfig checks the args.
func CheckCreateWebhookNotificationConfig(v *viper.Viper, args []string) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		return err
	}

	return nil
}

// CreateWebhookNotification sends a CreateWebhookNotification request to the support endpoint
func CreateWebhookNotification(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger - remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := CheckCreateWebhookNotificationConfig(v, args)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file if it was passed in
	filename := v.GetString(utils.FilenameFlag)
	var createWebhookNotificationParams webhookclient.CreateWebhookNotificationParams
	if filename != "" {
		err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &createWebhookNotificationParams)
		if err != nil {
			logger.Fatal(err)
		}
	}
	createWebhookNotificationParams.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	gateway, cacStore, errCreateClient := utils.CreateSupportClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}
	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer func() {
			if closeErr := cacStore.Close(); closeErr != nil {
				logger.Fatal(closeErr)
			}
		}()
	}

	// Make the API Call
	resp, err := gateway.Webhook.CreateWebhookNotification(&createWebhookNotificationParams)
	if err != nil {
		return utils.HandleGatewayError(err, logger)
	}

	// Get the successful response payload and convert to json for output
	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			logger.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		logger.Fatal(resp.Error())
	}

	return nil

}
