package prime

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
	mtoShipment "github.com/transcom/mymove/pkg/gen/primeclient/mto_shipment"
)

// InitCreateMTOAgentFlags declares which flags are enabled
func InitCreateMTOAgentFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Name of the file being passed in")
	flag.SortFlags = false
}

func checkCreateMTOAgentConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		logger.Fatal(errors.New("create-agents expects a file to be passed in"))
	}

	return nil
}

// CreateMTOAgent creates a gateway and sends the request to the endpoint
func CreateMTOAgent(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger - remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreateMTOAgentConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed into MTOShipment
	filename := v.GetString(utils.FilenameFlag)
	var mtoAgentPayload mtoShipment.CreateMTOAgentParams // CreateMTOAgentParams Takes data out of JSon and puts it in struct
	err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &mtoAgentPayload)
	if err != nil {
		logger.Fatal(err)
	}
	mtoAgentPayload.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	primeGateway, cacStore, errCreateClient := utils.CreatePrimeClient(v)
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
	resp, err := primeGateway.MtoShipment.CreateMTOAgent(&mtoAgentPayload) // Sends the request
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
