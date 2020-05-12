package main

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

	mtoShipment "github.com/transcom/mymove/pkg/gen/primeclient/mto_shipment"
)

func initCreateMTOShipmentFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Name of the file being passed in")
	flag.SortFlags = false
}

func checkCreateMTOShipmentConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		logger.Fatal(errors.New("create-mto-shipment expects a file to be passed in"))
	}

	return nil
}

func createMTOShipment(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger - remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreateMTOShipmentConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed into MTOShipment
	filename := v.GetString(FilenameFlag)
	var shipmentPayload mtoShipment.CreateMTOShipmentParams
	err = decodeJSONFileToPayload(filename, containsDash(args), &shipmentPayload)
	if err != nil {
		logger.Fatal(err)
	}
	shipmentPayload.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	primeGateway, cacStore, errCreateClient := CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Make the API Call
	resp, err := primeGateway.MtoShipment.CreateMTOShipment(&shipmentPayload)
	if err != nil {
		return handleGatewayError(err, logger)
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
