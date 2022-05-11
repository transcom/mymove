package support

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	paymentrequestclient "github.com/transcom/mymove/pkg/gen/supportclient/payment_request"
)

// InitRecalculatePaymentRequestFlags declares which flags are enabled
func InitRecalculatePaymentRequestFlags(flag *pflag.FlagSet) {
	flag.String(utils.IDFlag, "", "UUID of the payment request to recalculate")

	flag.SortFlags = false
}

func checkRecalculatePaymentRequestConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	if uuid.FromStringOrNil(v.GetString(utils.IDFlag)) == uuid.Nil {
		logger.Fatal(errors.New("support-recalculate-payment-request expects a valid UUID to be passed in"))
	}

	return nil
}

// RecalculatePaymentRequest creates a gateway and sends the request to the endpoint
func RecalculatePaymentRequest(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//  Create the logger
	//  Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkRecalculatePaymentRequestConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Get the UUID that was passed in
	paymentRequestID := v.GetString(utils.IDFlag)
	var recalculatePaymentRequestParams paymentrequestclient.RecalculatePaymentRequestParams
	recalculatePaymentRequestParams.PaymentRequestID = strfmt.UUID(paymentRequestID)
	recalculatePaymentRequestParams.SetTimeout(time.Second * 30)

	// Create the client and open the cacStore
	supportGateway, cacStore, errCreateClient := utils.CreateSupportClient(v)
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

	resp, err := supportGateway.PaymentRequest.RecalculatePaymentRequest(&recalculatePaymentRequestParams)
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
