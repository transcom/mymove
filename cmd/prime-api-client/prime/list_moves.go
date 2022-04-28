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
	"github.com/transcom/mymove/pkg/gen/primemessages"
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

	startTime := time.Now()

	// this wait retry logic would need to be replicated to all
	// commands, so start with list moves for now
	wait := v.GetDuration(utils.WaitFlag)
	params.SetTimeout(wait)
	var payload primemessages.ListMoves
	// loop until we either time out or get a successful response
	for {
		resp, err := primeGateway.MoveTaskOrder.ListMoves(&params)
		if err != nil {
			currentTime := time.Now()
			if currentTime.Sub(startTime) > wait {
				// the request timed out, so return the error
				return utils.HandleGatewayError(err, logger)
			}
			logger.Printf("Problem with request: %s, Sleeping 1s\n", err)
			time.Sleep(1 * time.Second)
		} else {
			payload = resp.GetPayload()
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
	}

}
