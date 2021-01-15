package support

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
	"github.com/transcom/mymove/pkg/gen/supportclient/payment_request"
)

// InitGetPaymentRequestEDIFlags declares which flags are enabled
func InitGetPaymentRequestEDIFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Name of the file being passed in")

	flag.SortFlags = false
}

func checkGetPaymentRequestEDIConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		logger.Fatal(errors.New("get-mto expects a file to be passed in"))
	}

	return nil
}

// GetPaymentRequestEDI creates a gateway and sends the request to the endpoint
func GetPaymentRequestEDI(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkGetPaymentRequestEDIConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed in
	filename := v.GetString(utils.FilenameFlag)
	var getPaymentRequestEDIParams payment_request.GetPaymentRequestEDIParams
	err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &getPaymentRequestEDIParams)
	if err != nil {
		logger.Fatal(err)
	}
	getPaymentRequestEDIParams.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	supportGateway, cacStore, errCreateClient := utils.CreateSupportClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used to close a cmd line client
		//RA: Given the functions causing the lint errors are used end a local running process, it is not deemed a risk
		//RA Developer Status: Mitigated
		//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
		//RA Validator: jneuner@mitre.org
		//RA Modified Severity:
		defer cacStore.Close() // nolint:errcheck
	}
	getPaymentRequestEDIParams.SetTimeout(time.Second * 30)

	resp, err := supportGateway.PaymentRequest.GetPaymentRequestEDI(&getPaymentRequestEDIParams)
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
