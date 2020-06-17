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

	"github.com/transcom/mymove/pkg/gen/primeclient/payment_requests"
)

// InitCreatePaymentRequestFlags initializes flags.
func InitCreatePaymentRequestFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Path to the file with the payment request JSON payload")

	flag.SortFlags = false
}

// checkCreatePaymentRequestConfig checks the args.
func checkCreatePaymentRequestConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		return err
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		return errors.New("create-payment-request expects a file to be passed in")
	}

	return nil
}

// CreatePaymentRequest creates the payment request for an MTO
func CreatePaymentRequest(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreatePaymentRequestConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Decode json from file that was passed in
	filename := v.GetString(utils.FilenameFlag)
	var paymentRequestParams payment_requests.CreatePaymentRequestParams
	err = utils.DecodeJSONFileToPayload(filename, utils.ContainsDash(args), &paymentRequestParams)
	if err != nil {
		logger.Fatal(err)
	}
	paymentRequestParams.SetTimeout(time.Second * 30)

	// cac and api gateway
	primeGateway, cacStore, errCreateClient := utils.CreatePrimeClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer cacStore.Close()
	}

	resp, err := primeGateway.PaymentRequests.CreatePaymentRequest(&paymentRequestParams)
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
