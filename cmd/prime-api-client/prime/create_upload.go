package prime

import (
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
	paymentrequestclient "github.com/transcom/mymove/pkg/gen/primeclient/payment_request"
)

// InitCreatePaymentRequestUploadFlags initializes flags.
func InitCreatePaymentRequestUploadFlags(flag *pflag.FlagSet) {
	flag.String(utils.FilenameFlag, "", "Path to the upload file for create-payment-request-upload payload")
	flag.String(utils.PaymentRequestIDFlag, "", "Payment Request ID to upload the proof of service document to")

	flag.SortFlags = false
}

// checkCreatePaymentRequestUploadConfig checks the args.
func checkCreatePaymentRequestUploadConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := utils.CheckRootConfig(v)
	if err != nil {
		return err
	}

	if v.GetString(utils.FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !utils.ContainsDash(args)) {
		return errors.New("create-payment-request-upload expects a file to be passed in")
	}

	// Get the paymentRequestID to use for the upload file
	if v.GetString(utils.PaymentRequestIDFlag) == "" && (len(args) < 1 || len(args) > 0) {
		return errors.New("create-payment-request-upload expects a  utils.PaymentRequestID to be passed in")
	}

	return nil
}

// CreatePaymentRequestUpload creates the payment request for an MTO
func CreatePaymentRequestUpload(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := utils.ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreatePaymentRequestUploadConfig(v, args, logger)
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
		defer cacStore.Close()
	}

	// Get the filename for the upload file to upload with command create-payment-request-upload
	filename := v.GetString(utils.FilenameFlag)

	// Get the paymentRequestID to use for the upload file
	paymentRequestID := v.GetString(utils.PaymentRequestIDFlag)

	file, fileErr := os.Open(filepath.Clean(filename))
	// G307 TODO needs review
	defer file.Close()
	if fileErr != nil {
		logger.Fatal(fileErr)
	}

	params := paymentrequestclient.CreateUploadParams{
		File:             file,
		PaymentRequestID: paymentRequestID,
	}
	params.SetTimeout(time.Second * 30)

	resp, errCreatePaymentRequestUpload := primeGateway.PaymentRequest.CreateUpload(&params)
	if errCreatePaymentRequestUpload != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errCreatePaymentRequestUpload.Error())
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
