package support

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"

	movetaskorderclient "github.com/transcom/mymove/pkg/gen/supportclient/move_task_order"
)

// InitCreateMTOFlags initializes flags.
func InitCreateMTOFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Path to the file with the payment request JSON payload")

	flag.SortFlags = false
}

// CheckCreateMTOConfig checks the args.
func CheckCreateMTOConfig(v *viper.Viper, args []string) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		return err
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		return errors.New("create-move-task-order expects a file to be passed in")
	}

	return nil
}

// CreateMTO sends a CreateMoveTaskOrder request to the support endpoint
func CreateMTO(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger - remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := CheckCreateMTOConfig(v, args)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed in
	filename := v.GetString(utils.FilenameFlag)
	var createMTOParams movetaskorderclient.CreateMoveTaskOrderParams
	err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &createMTOParams)
	if err != nil {
		logger.Fatal(err)
	}
	createMTOParams.SetTimeout(time.Second * 30)

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
	resp, err := gateway.MoveTaskOrder.CreateMoveTaskOrder(&createMTOParams)
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
