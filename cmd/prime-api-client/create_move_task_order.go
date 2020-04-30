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

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	movetaskorderclient "github.com/transcom/mymove/pkg/gen/supportclient/move_task_order"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
)

// initCreateMTOFlags initializes flags.
func initCreateMTOFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Path to the file with the payment request JSON payload")

	flag.SortFlags = false
}

// checkCreateMTOConfig checks the args.
func checkCreateMTOConfig(v *viper.Viper, args []string) error {
	err := CheckRootConfig(v)
	if err != nil {
		return err
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		return errors.New("create-payment-request expects a file to be passed in")
	}

	return nil
}

// createMTO creates the payment request for an MTO
func createMTO(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger - remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreateMTOConfig(v, args)
	if err != nil {
		logger.Fatal(err)
	}

	// cac and api gateway
	gateway, cacStore, errCreateClient := CreateSupportClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Decode json from file that was passed into create-payment-request
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
	var mto supportmessages.MoveTaskOrder

	// Read the json into the mto payload
	err = jsonDecoder.Decode(&mto)
	if err != nil {
		return fmt.Errorf("decoding data failed: %w", err)
	}

	// Create the http request params
	params := movetaskorderclient.CreateMoveTaskOrderParams{
		Body: &mto,
	}
	params.SetTimeout(time.Second * 30)

	// Make the request to the server
	resp, err := gateway.MoveTaskOrder.CreateMoveTaskOrder(&params)
	if err != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		data, _ := json.Marshal(err)
		fmt.Printf("%s", data)
		return nil
	}

	// Get the response payload and convert to json for output
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
