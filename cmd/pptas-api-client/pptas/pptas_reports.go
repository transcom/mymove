package pptas

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

	"github.com/transcom/mymove/cmd/pptas-api-client/utils"
	moves "github.com/transcom/mymove/pkg/gen/pptasclient/moves"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
)

// InitPPTASReportsFlags declares which flags are enabled
func InitPPTASReportsFlags(flag *pflag.FlagSet) {
	flag.String(utils.SinceFlag, "", "Timestamp for filtering moves. Returns moves updated since this time.")
	flag.SortFlags = false
}

func checkPPTASReportsConfig(v *viper.Viper, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	return nil
}

// PPTASReports creates a gateway and sends the request to the endpoint
func PPTASReports(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkPPTASReportsConfig(v, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Get the since param, if any
	var params moves.PptasReportsParams
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
	var payload pptasmessages.PPTASReports
	// loop until we either time out or get a successful response
	for {
		resp, err := primeGateway.Moves.PptasReports(&params)
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
