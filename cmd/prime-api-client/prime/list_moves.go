package prime

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"

	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

// InitListMovesFlags declares which flags are enabled
func InitListMovesFlags(flag *pflag.FlagSet) {
	flag.String(utils.SinceFlag, "", "Timestamp for filtering moves. Returns moves updated since this time.")
	flag.SortFlags = false
}

func checkListMovesConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	return nil
}

// ListMoves creates a gateway and sends the request to the endpoint
func ListMoves(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkListMovesConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Get the since param, if any
	var params mto.ListMovesParams
	since := v.GetString(utils.SinceFlag)
	if since != "" {
		sinceDateTime, sinceErr := strfmt.ParseDateTime(since)
		if sinceErr != nil {
			logger.Fatal(err)
		}
		params.SetSince(&sinceDateTime)
	}

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

	params.SetTimeout(time.Second * 30)
	resp, err := primeGateway.MoveTaskOrder.ListMoves(&params)
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
