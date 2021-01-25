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
	mto "github.com/transcom/mymove/pkg/gen/supportclient/move_task_order"
)

// InitHideNonFakeMoveTaskOrdersFlags declares which flags are enabled
func InitHideNonFakeMoveTaskOrdersFlags(flag *pflag.FlagSet) {
	flag.SortFlags = false
}
func checkHideMoveTaskOrderConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}
	return nil
}

// HideNonFakeMoveTaskOrders creates a gateway and sends the request to the endpoint
func HideNonFakeMoveTaskOrders(cmd *cobra.Command, args []string) error {
	v := viper.New()
	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)
	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}
	// Check the config before talking to the CAC
	err := checkHideMoveTaskOrderConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}
	// Create the client and open the cacStore
	supportGateway, cacStore, errCreateClient := utils.CreateSupportClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}
	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer func() {
			if closeErr := cacStore.Close(); closeErr != nil {
				fmt.Println(fmt.Errorf("Close store connection failed: %w", closeErr))
			}
		}()
	}
	var params mto.HideNonFakeMoveTaskOrdersParams
	params.SetTimeout(time.Second * 30)
	resp, err := supportGateway.MoveTaskOrder.HideNonFakeMoveTaskOrders(&params)
	if err != nil {
		return utils.HandleGatewayError(err, logger)
	}
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
