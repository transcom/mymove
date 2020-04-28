package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	openapi "github.com/go-openapi/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	mto "github.com/transcom/mymove/pkg/gen/supportclient/move_task_order"
)

func initGetMTOFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Name of the file being passed in")

	flag.SortFlags = false
}

func checkGetMTOConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		logger.Fatal(errors.New("get-mto expects a file to be passed in"))
	}

	return nil
}

func getMTO(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkGetMTOConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed in
	filename := v.GetString(FilenameFlag)
	var getMTOParams mto.GetMoveTaskOrderParams
	err = decodeJSONFileToPayload(filename, containsDash(args), &getMTOParams)
	if err != nil {
		logger.Fatal(err)
	}
	getMTOParams.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	supportGateway, cacStore, errCreateClient := CreateSupportClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}
	getMTOParams.SetTimeout(time.Second * 30)

	resp, errGetMTO := supportGateway.MoveTaskOrder.GetMoveTaskOrder(&getMTOParams)
	if errGetMTO != nil {
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
