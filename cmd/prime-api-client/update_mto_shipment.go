package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mtoShipment "github.com/transcom/mymove/pkg/gen/primeclient/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

const (
	// FilenameFlag is the name of the file being passed in
	FilenameFlag string = "filename"
	// ETagFlag is the etag for the mto shipment being updated
	ETagFlag string = "etag"
)

func containsDash(args []string) bool {
	for _, arg := range args {
		if arg == "-" {
			return true
		}
	}
	return false
}

func initUpdateMTOShipmentFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Name of the file being passed in")
	flag.String(ETagFlag, "", "ETag for the mto shipment being updated")

	flag.SortFlags = false
}

func checkUpdateMTOShipmentConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(ETagFlag) == "" {
		logger.Fatal(errors.New("update-mto-shipment expects an etag"))
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		logger.Fatal(errors.New("update-mto-shipment expects a file to be passed in"))
	}

	return nil
}

func updateMTOShipment(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	primeGateway, cacStore, err := CreateClient(cmd, v, args)
	if err != nil {
		return err
	}

	if cacStore != nil {
		defer cacStore.Close()
	}

	err = checkUpdateMTOShipmentConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
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
	var shipment primemessages.MTOShipment
	err = jsonDecoder.Decode(&shipment)
	if err != nil {
		return errors.Wrap(err, "decoding data failed")
	}

	params := mtoShipment.UpdateMTOShipmentParams{
		MoveTaskOrderID: shipment.MoveTaskOrderID,
		MtoShipmentID:   shipment.ID,
		Body:            &shipment,
		IfMatch:         v.GetString(ETagFlag),
	}
	params.SetTimeout(time.Second * 30)

	resp, errUpdateMTOShipment := primeGateway.MtoShipment.UpdateMTOShipment(&params)
	if errUpdateMTOShipment != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		log.Fatal(errUpdateMTOShipment.Error())
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			log.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		log.Fatal(resp.Error())
	}

	return nil
}
