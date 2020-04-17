package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/transcom/mymove/pkg/gen/supportmessages"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	mtoShipment "github.com/transcom/mymove/pkg/gen/supportclient/mto_shipment"
)

func initPatchMTOShipmentStatusFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Name of the file being passed in")
	flag.String(ETagFlag, "", "ETag for the mto shipment being updated")

	flag.SortFlags = false
}

func checkPatchMTOShipmentStatusConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(ETagFlag) == "" {
		logger.Fatal(errors.New("support-patch-mto-shipment-status expects an etag"))
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		logger.Fatal(errors.New("support-patch-mto-shipment-status expects a file to be passed in"))
	}

	return nil
}

func patchMTOShipmentStatus(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkPatchMTOShipmentStatusConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	supportGateway, cacStore, errCreateClient := CreateSupportClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Decode json from file that was passed into MTOShipment
	filename := v.GetString(FilenameFlag)

	var reader *bufio.Reader
	if filename != "" {
		file, fileErr := os.Open(filepath.Clean(filename))
		if fileErr != nil {
			logger.Fatal(fileErr)
		}
		reader = bufio.NewReader(file)
	}

	if len(args) > 0 && containsDash(args) {
		reader = bufio.NewReader(os.Stdin)
	}

	jsonDecoder := json.NewDecoder(reader)
	var shipment supportmessages.MTOShipment

	err = jsonDecoder.Decode(&shipment)

	if err != nil {
		return fmt.Errorf("decoding data failed: %w", err)
	}

	params := mtoShipment.PatchMTOShipmentStatusParams{
		ShipmentID: shipment.ID,
		IfMatch:    v.GetString(ETagFlag),
		Body: &supportmessages.PatchMTOShipmentStatusPayload{
			Status:          shipment.Status,
			RejectionReason: shipment.RejectionReason},
	}

	params.SetTimeout(time.Second * 30)

	resp, errPatchMTOShipmentStatus := supportGateway.MtoShipment.PatchMTOShipmentStatus(&params)
	if errPatchMTOShipmentStatus != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errPatchMTOShipmentStatus.Error())
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
