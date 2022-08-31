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
	mtoserviceitem "github.com/transcom/mymove/pkg/gen/supportclient/mto_service_item"
)

// InitUpdateMTOServiceItemStatusFlags declares which flags are enabled
func InitUpdateMTOServiceItemStatusFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Name of the file being passed in")

	flag.SortFlags = false
}

func checkUpdateMTOServiceItemStatusConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		logger.Fatal(errors.New("support-update-mto-service-item-status expects a file to be passed in"))
	}

	return nil
}

// UpdateMTOServiceItemStatus creates a gateway and sends the request to the endpoint
func UpdateMTOServiceItemStatus(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkUpdateMTOServiceItemStatusConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed into MTO Service item
	filename := v.GetString(utils.FilenameFlag)
	var updateServiceItemParams mtoserviceitem.UpdateMTOServiceItemStatusParams
	err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &updateServiceItemParams)
	if err != nil {
		logger.Fatal(err)
	}
	updateServiceItemParams.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	supportGateway, cacStore, errCreateClient := utils.CreateSupportClient(v)
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
	resp, err := supportGateway.MtoServiceItem.UpdateMTOServiceItemStatus(&updateServiceItemParams)
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
