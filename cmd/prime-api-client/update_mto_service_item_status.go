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

	mtoserviceitem "github.com/transcom/mymove/pkg/gen/supportclient/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
)

func initUpdateMTOServiceItemStatusFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Name of the file being passed in")
	flag.String(ETagFlag, "", "ETag for the mto service item being updated")

	flag.SortFlags = false
}

func checkUpdateMTOServiceItemStatusConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(ETagFlag) == "" {
		logger.Fatal(errors.New("support-update-mto-service-item-status expects an etag"))
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		logger.Fatal(errors.New("support-update-mto-service-item-status expects a file to be passed in"))
	}

	return nil
}

func updateMTOServiceItemStatus(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkUpdateMTOServiceItemStatusConfig(v, args, logger)
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

	// Decode json from file that was passed into MTO Service item
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
	var serviceItem supportmessages.UpdateMTOServiceItemStatus
	err = jsonDecoder.Decode(&serviceItem)
	if err != nil {
		return fmt.Errorf("decoding data failed: %w", err)
	}

	params := mtoserviceitem.UpdateMTOServiceItemStatusParams{
		MtoServiceItemID: serviceItem.ID.String(),
		Body:             &serviceItem,
		IfMatch:          v.GetString(ETagFlag),
	}

	params.SetTimeout(time.Second * 30)

	resp, errUpdateMTOServiceItemStatus := supportGateway.MtoServiceItem.UpdateMTOServiceItemStatus(&params)
	if errUpdateMTOServiceItemStatus != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errUpdateMTOServiceItemStatus.Error())
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
