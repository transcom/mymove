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

	mtoServiceItem "github.com/transcom/mymove/pkg/gen/primeclient/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

type getType struct {
	ModelType primemessages.MTOServiceItemModelType `json:"modelType"`
}

// initCreateMTOServiceItemFlags initializes flags.
func initCreateMTOServiceItemFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Name of the file being passed in")

	flag.SortFlags = false
}

// checkCreateMTOServiceItemConfig checks the args.
func checkCreateMTOServiceItemConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		logger.Fatal(errors.New("create-mto-service-item expects a file to be passed in"))
	}

	return nil
}

// getFileReader will get the bufio file reader.
func getFileReader(filename string, args []string, logger *log.Logger) *bufio.Reader {
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

	return reader
}

// getJSONDecoder will get a new json decoder.
func getJSONDecoder(reader *bufio.Reader) *json.Decoder {
	return json.NewDecoder(reader)
}

// decodeServiceItem will decode the json into a mto service item type.
func decodeServiceItem(subType primemessages.MTOServiceItem, jsonDecoder *json.Decoder) (primemessages.MTOServiceItem, error) {
	err := jsonDecoder.Decode(&subType)
	if err != nil {
		return nil, fmt.Errorf("second decoding data failed: %w", err)
	}
	return subType, nil
}

// CreateMTOServiceItem creates the mto service item for a MTO and/or MTOShipment
func createMTOServiceItem(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreateMTOServiceItemConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// cac and api gateway
	primeGateway, cacStore, errCreateClient := CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}
	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// reading json file so we can unmarshal
	filename := v.GetString(FilenameFlag)
	reader := getFileReader(filename, args, logger)
	jsonDecoder := getJSONDecoder(reader)
	// decode first to determine the model type
	var gt getType
	err = jsonDecoder.Decode(&gt)
	if err != nil {
		return fmt.Errorf("first decoding data failed: %w", err)
	}

	// reset decoder to read from beginning
	jsonDecoder = getJSONDecoder(getFileReader(filename, args, logger))

	// once decoded, we can type cast into a more specific model type
	// then decode a second time into subtype
	var serviceItem primemessages.MTOServiceItem
	switch gt.ModelType {
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDOFSIT:
		serviceItem, err = decodeServiceItem(&primemessages.MTOServiceItemDOFSIT{}, jsonDecoder)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDDFSIT:
		serviceItem, err = decodeServiceItem(&primemessages.MTOServiceItemDDFSIT{}, jsonDecoder)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:
		serviceItem, err = decodeServiceItem(&primemessages.MTOServiceItemDomesticCrating{}, jsonDecoder)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle:
		serviceItem, err = decodeServiceItem(&primemessages.MTOServiceItemShuttle{}, jsonDecoder)
	default:
		err = fmt.Errorf("allowed modelType(): %v", []primemessages.MTOServiceItemModelType{
			primemessages.MTOServiceItemModelTypeMTOServiceItemDDFSIT,
			primemessages.MTOServiceItemModelTypeMTOServiceItemDOFSIT,
			primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating,
			primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle,
		})
	}
	// return any decoding errors
	if err != nil {
		return err
	}

	// Let's make a request!
	params := mtoServiceItem.CreateMTOServiceItemParams{
		MoveTaskOrderID: serviceItem.MoveTaskOrderID(),
		MtoShipmentID:   serviceItem.MtoShipmentID(),
		Body:            serviceItem,
	}
	params.SetTimeout(time.Second * 30)

	resp, errCreateMTOServiceItem := primeGateway.MtoServiceItem.CreateMTOServiceItem(&params)
	if errCreateMTOServiceItem != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errCreateMTOServiceItem.Error())
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
