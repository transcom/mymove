package prime

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

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	mtoserviceitemclient "github.com/transcom/mymove/pkg/gen/primeclient/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

type getUpdateMTOServiceItemModelType struct {
	ModelType primemessages.UpdateMTOServiceItemModelType `json:"modelType"`
}
type getUpdateMTOServiceItemModelTypeBody struct {
	Body getUpdateMTOServiceItemModelType `json:"body"`
}

type sITParams struct {
	IfMatch          string
	MtoServiceItemID string
	Body             primemessages.UpdateMTOServiceItemSIT
}

// InitUpdateMTOServiceItemFlags declares which flags are enabled
func InitUpdateMTOServiceItemFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Name of the file being passed in")
	flag.SortFlags = false
}

func checkUpdateMTOServiceItemConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		logger.Fatal(errors.New("update-mto-service-item expects a file to be passed in"))
	}

	return nil
}

// UpdateMTOServiceItem creates a gateway and sends the request to the endpoint
func UpdateMTOServiceItem(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkUpdateMTOServiceItemConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Create the client and open the cacStore
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
	// Decode json from file that was passed in
	filename := v.GetString(utils.FilenameFlag)
	reader := getFileReader(filename, args, logger)
	jsonDecoder := getJSONDecoder(reader)
	// decode first to determine the model type
	var gt getUpdateMTOServiceItemModelTypeBody
	err = jsonDecoder.Decode(&gt)
	if err != nil {
		return fmt.Errorf("first decoding data failed: %w", err)
	}

	var mtoServiceItemReqParams mtoserviceitemclient.UpdateMTOServiceItemParams
	switch gt.Body.ModelType {
	case primemessages.UpdateMTOServiceItemModelTypeUpdateMTOServiceItemSIT:
		var params sITParams
		err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &params)
		mtoServiceItemReqParams.SetMtoServiceItemID(params.MtoServiceItemID)
		mtoServiceItemReqParams.SetBody(&params.Body)
		mtoServiceItemReqParams.SetIfMatch(params.IfMatch)
	default:
		err = fmt.Errorf("allowed modelType(): %v", []primemessages.UpdateMTOServiceItemModelType{
			primemessages.UpdateMTOServiceItemModelTypeUpdateMTOServiceItemSIT,
		})
	}

	if err != nil {
		logger.Fatal(err)
		return err
	}

	mtoServiceItemReqParams.SetTimeout(time.Second * 30)
	// Make the API Call
	resp, err := primeGateway.MtoServiceItem.UpdateMTOServiceItem(&mtoServiceItemReqParams)
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
