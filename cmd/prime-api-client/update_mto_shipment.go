package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mtoShipment "github.com/transcom/mymove/pkg/gen/primeclient/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// TODO Add If-Match flag and check when replaces If-Unmodified-Since
// NOTE: leaving this function in here due to the above
func checkUpdateMTOShipmentConfig(v *viper.Viper, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	fileStats, err := os.Stdin.Stat()
	if err != nil {
		logger.Fatal(err)
	}

	if fileStats.Size() == 0 {
		logger.Fatal(errors.New("update-mto-shipment expects a file to be passed in"))
	}

	return nil
}

func updateMTOShipment(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	primeGateway, err := CreateClient(cmd, v, args)
	if err != nil {
		return err
	}

	err = checkUpdateMTOShipmentConfig(v, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed into MTOShipment
	jsonDecoder := json.NewDecoder(bufio.NewReader(os.Stdin))
	var shipment primemessages.MTOShipment
	err = jsonDecoder.Decode(&shipment)
	if err != nil {
		return errors.Wrap(err, "decoding data failed")
	}

	params := mtoShipment.UpdateMTOShipmentParams{
		MoveTaskOrderID: shipment.MoveTaskOrderID,
		MtoShipmentID:   shipment.ID,
		Body:            &shipment,
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