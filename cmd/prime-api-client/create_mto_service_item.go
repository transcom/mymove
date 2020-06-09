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

type getTypeBody struct {
	ModelType primemessages.MTOServiceItemModelType `json:"modelType"`
}
type getType struct {
	Body getTypeBody `json:"body"`
}

type basicParams struct {
	Body primemessages.MTOServiceItemBasic `json:"body"`
}
type dOFSITParams struct {
	Body primemessages.MTOServiceItemDOFSIT `json:"body"`
}
type dDFSITParams struct {
	Body primemessages.MTOServiceItemDDFSIT `json:"body"`
}
type domesticCratingParams struct {
	Body primemessages.MTOServiceItemDomesticCrating `json:"body"`
}
type shuttleParams struct {
	Body primemessages.MTOServiceItemShuttle `json:"body"`
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

	// once decoded, we can type cast into a more specific model type
	// then decode a second time into subtype
	var serviceItemParams mtoServiceItem.CreateMTOServiceItemParams
	switch gt.Body.ModelType {
	case primemessages.MTOServiceItemModelTypeMTOServiceItemBasic:
		var params basicParams
		err = decodeJSONFileToPayload(filename, containsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDOFSIT:
		var params dOFSITParams
		err = decodeJSONFileToPayload(filename, containsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDDFSIT:
		var params dDFSITParams
		err = decodeJSONFileToPayload(filename, containsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:
		var params domesticCratingParams
		err = decodeJSONFileToPayload(filename, containsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle:
		var params shuttleParams
		err = decodeJSONFileToPayload(filename, containsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	default:
		err = fmt.Errorf("allowed modelType(): %v", []primemessages.MTOServiceItemModelType{
			primemessages.MTOServiceItemModelTypeMTOServiceItemBasic,
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
	serviceItemParams.SetTimeout(time.Second * 30)
	resp, err := primeGateway.MtoServiceItem.CreateMTOServiceItem(&serviceItemParams)
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
