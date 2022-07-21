package prime

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

// InitGetMTOFlags declares which flags are enabled
func InitGetMTOFlags(flag *pflag.FlagSet) {
	flag.String(utils.IDFlag, "", "UUID of the desired move")

	flag.SortFlags = false
}

func checkGetMTOConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if uuid.FromStringOrNil(v.GetString(utils.IDFlag)) == uuid.Nil {
		logger.Fatal(errors.New("get-move-task-order expects a valid UUID to be passed in"))
	}

	return nil
}

// GetMTO creates a gateway and sends the request to the endpoint
func GetMTO(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkGetMTOConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Get the UUID that was passed in
	moveID := v.GetString(utils.IDFlag)
	var getMTOParams mto.GetMoveTaskOrderParams
	getMTOParams.MoveID = moveID
	getMTOParams.SetTimeout(time.Second * 30)

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
	getMTOParams.SetTimeout(time.Second * 30)

	resp, err := primeGateway.MoveTaskOrder.GetMoveTaskOrder(&getMTOParams)
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
