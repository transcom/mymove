package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	mtoShipment "github.com/transcom/mymove/pkg/gen/primeclient/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// ErrInvalidID is an error indicating that an id is invaid
type ErrInvalidID struct {
	MTOShipmentID string
	MTOID         string
}

func (e *ErrInvalidID) Error() string {
	return fmt.Sprintf("invalid mto id %q or mto shipment id %q", e.MTOID, e.MTOShipmentID)
}

// ErrInvalidTimestamp is an error indicating that the timestamp is invalid
type ErrInvalidTimestamp struct {
	Timestamp time.Time
}

func (e *ErrInvalidTimestamp) Error() string {
	return fmt.Sprintf("invalid timestamp %q", e.Timestamp)
}

const (
	// MTOShipmentIDFlag is the id of the mto shipment to be updated
	MTOShipmentIDFlag string = "mtoShipmentID"
	// MTOIDFlag is the id of the move task order whose shipment is being updated
	MTOIDFlag string = "mtoID"
	// UnmodifiedSinceFlag is the timestamp of when the mto shipment was last updated
	// TODO: refactor when if-unmodified-since is replaced if-math for this endpoint
	UnmodifiedSinceFlag string = "unmodifiedSince"
)

func checkUpdateMTOShipmentConfig(v *viper.Viper, logger *log.Logger) error {
	mtoID := v.GetString(MTOIDFlag)
	mtoShipmentID := v.GetString(MTOShipmentIDFlag)
	if mtoID == "" || mtoShipmentID == "" {
		return fmt.Errorf("%q or %q is invalid: %w", MTOIDFlag, MTOShipmentIDFlag, &ErrInvalidID{MTOID: mtoID, MTOShipmentID: mtoShipmentID})
	}

	unmodifiedSince := v.GetTime(UnmodifiedSinceFlag)
	if unmodifiedSince.IsZero() {
		return fmt.Errorf("UnmodifiedSinceFlag is invalid: %w", &ErrInvalidTimestamp{Timestamp: unmodifiedSince})
	}

	err := cli.CheckCAC(v)
	if err != nil {
		return err
	}

	err = cli.CheckPrimeAPI(v)
	if err != nil {
		return err
	}

	err = cli.CheckVerbose(v)
	if err != nil {
		return err
	}

	return nil
}

func initUpdateMTOShipmentFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitPrimeAPIFlags(flag)
	cli.InitVerboseFlags(flag)

	flag.String(MTOShipmentIDFlag, "", "ID of the MTO Shipment that is being updated")
	flag.String(MTOIDFlag, "", "ID of the MTO whose shipment is being updated")
	flag.String(UnmodifiedSinceFlag, "", "Timestamp of when mto shipment was last updated")

	flag.SortFlags = false
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

	// Use command line inputs
	mtoID := v.GetString(MTOIDFlag)
	mtoShipmentID := v.GetString(MTOShipmentIDFlag)
	unmodifiedSince := v.GetTime(UnmodifiedSinceFlag)

	// Decode json from file that was passed into MTOShipment
	jsonDecoder := json.NewDecoder(bufio.NewReader(os.Stdin))
	var shipment primemessages.MTOShipment
	err = jsonDecoder.Decode(&shipment)
	if err != nil {
		return errors.Wrap(err, "decoding data failed")
	}

	params := mtoShipment.UpdateMTOShipmentParams{
		MoveTaskOrderID:   strfmt.UUID(mtoID),
		MtoShipmentID:     strfmt.UUID(mtoShipmentID),
		IfUnmodifiedSince: strfmt.DateTime(unmodifiedSince),
		Body:              &shipment,
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