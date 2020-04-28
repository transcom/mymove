package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	openapi "github.com/go-openapi/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

func initFetchMTOsFlags(flag *pflag.FlagSet) {
	flag.SortFlags = false
}

func checkFetchMTOsConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	return nil
}

func fetchMTOs(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkFetchMTOsConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	primeGateway, cacStore, errCreateClient := CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	var params mto.FetchMTOUpdatesParams
	params.SetTimeout(time.Second * 30)
	resp, err := primeGateway.MoveTaskOrder.FetchMTOUpdates(&params)
	if err != nil {
		// If you see an error like "unknown error (status 422)", it means
		// we hit a completely unhandled error that we should handle.
		// We should be enabling said error in the endpoint in swagger.
		// 422 for example is an Unprocessable Entity and is returned by the swagger
		// validation before it even hits the handler.
		if _, ok := err.(*openapi.APIError); ok {
			apiErr := err.(*openapi.APIError).Response.(openapi.ClientResponse)
			logger.Fatal(fmt.Sprintf("%s: %s", err, apiErr.Message()))
		}
		// If it is a handled error, we should be able to pull out the payload here
		data, _ := json.Marshal(err)
		fmt.Printf("%s", data)
		return nil
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
