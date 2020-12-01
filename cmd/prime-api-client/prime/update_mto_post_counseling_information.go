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
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

// InitUpdatePostCounselingInfoFlags declares which flags are enabled
func InitUpdatePostCounselingInfoFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Name of the file being passed in")

	flag.SortFlags = false
}

func checkUpdatePostCounselingInfoConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		logger.Fatal(errors.New("update-post-counseling-info expects a file to be passed in"))
	}

	return nil
}

// UpdatePostCounselingInfo creates a gateway and sends the request to the endpoint
func UpdatePostCounselingInfo(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkUpdatePostCounselingInfoConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed in
	filename := v.GetString(utils.FilenameFlag)
	var mtoParams mto.UpdateMTOPostCounselingInformationParams
	err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &mtoParams)
	if err != nil {
		logger.Fatal(err)
	}
	mtoParams.SetTimeout(time.Second * 30)

	primeGateway, cacStore, errCreateClient := utils.CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	// Make the API Call
	resp, err := primeGateway.MoveTaskOrder.UpdateMTOPostCounselingInformation(&mtoParams)
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
