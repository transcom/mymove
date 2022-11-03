package prime

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

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	mtoServiceItem "github.com/transcom/mymove/pkg/gen/primeclient/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

type getTypeBody struct {
	ModelType primemessages.MTOServiceItemModelType `json:"modelType"`
}
type getType struct {
	Body getTypeBody `json:"body"`
}

type dOFSITParams struct {
	Body primemessages.MTOServiceItemOriginSIT `json:"body"`
}
type dDFSITParams struct {
	Body primemessages.MTOServiceItemDestSIT `json:"body"`
}
type domesticCratingParams struct {
	Body primemessages.MTOServiceItemDomesticCrating `json:"body"`
}
type shuttleParams struct {
	Body primemessages.MTOServiceItemShuttle `json:"body"`
}

// InitCreateMTOServiceItemFlags initializes flags.
func InitCreateMTOServiceItemFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Name of the file being passed in")

	flag.SortFlags = false
}

// checkCreateMTOServiceItemConfig checks the args.
func checkCreateMTOServiceItemConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
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
	if len(args) > 0 && utils.ContainsDash(args) {
		reader = bufio.NewReader(os.Stdin)
	}

	return reader
}

// getJSONDecoder will get a new json decoder.
func getJSONDecoder(reader *bufio.Reader) *json.Decoder {
	return json.NewDecoder(reader)
}

// CreateMTOServiceItem creates the mto service item for a MTO and/or MTOShipment
func CreateMTOServiceItem(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreateMTOServiceItemConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// cac and api gateway
	primeGateway, cacStore, errCreateClient := utils.CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}
	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer func() {
			if closeErr := cacStore.Close(); closeErr != nil {
				logger.Fatal(closeErr)
			}
		}()
	}

	// reading json file so we can unmarshal
	filename := v.GetString(utils.FilenameFlag)
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
	case primemessages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:
		var params dOFSITParams
		err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDestSIT:
		var params dDFSITParams
		err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:
		var params domesticCratingParams
		err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	case primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle:
		var params shuttleParams
		err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &params)
		serviceItemParams.SetBody(&params.Body)
	default:
		err = fmt.Errorf("allowed modelType(): %v", []primemessages.MTOServiceItemModelType{
			primemessages.MTOServiceItemModelTypeMTOServiceItemDestSIT,
			primemessages.MTOServiceItemModelTypeMTOServiceItemOriginSIT,
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
		return utils.HandleGatewayError(err, logger)
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
