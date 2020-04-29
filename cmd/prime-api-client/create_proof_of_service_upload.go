package main

import (
	//"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	//"unsafe"

	//"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/gen/primeclient/uploads"
	//"github.com/transcom/mymove/pkg/gen/primemessages"
)

// initCreateProofOfServiceUploadFlags initializes flags.
func initCreateProofOfServiceUploadFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Path to the file with the create-payment-request-upload JSON payload")
	flag.String(PaymentRequestID, "", "Payment Request ID to upload the proof of service document to")

	flag.SortFlags = false
}

// checkCreateProofOfServiceUploadConfig checks the args.
func checkCreateProofOfServiceUploadConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		return err
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		return errors.New("create-payment-request-upload expects a file to be passed in")
	}

	return nil
}

/*
func intToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func bytesToFile(filename string, data []byte) (afero.File, error) {
	var fs = afero.NewMemMapFs()
	uploadFile, err := fs.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create afero file %w", err)
	}

	uploadFile.Write(data)
	uploadFile.Sync()

	return uploadFile, nil
}
*/

// createProofOfServiceUpload creates the payment request for an MTO
func createProofOfServiceUpload(cmd *cobra.Command, args []string) error {
	v := viper.New()

	// Create the logger
	// Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkCreateProofOfServiceUploadConfig(v, args, logger)
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

	/*
		// Decode json from file that was passed into proof-of-service-upload
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
		var upload primemessages.Upload
		err = jsonDecoder.Decode(&upload)
		if err != nil {
			return fmt.Errorf("decoding data failed: %w", err)
		}

		data := intToByteArray(*upload.Bytes)
		file, err := bytesToFile(*upload.Filename, data)
		if err != nil {
			return fmt.Errorf("failed to create file from Byte: %w", err)
		}
	*/

	// Decode json from file that was passed into proof-of-service-upload
	filename := v.GetString(FilenameFlag)
	if filename == "" {
		return fmt.Errorf("failed to open filename: %s", filename)
	}

	paymentRequestID := v.GetString(PaymentRequestID)
	if paymentRequestID == "" {
		return fmt.Errorf("paymentRequestID required: %s", paymentRequestID)
	}

	file, fileErr := os.Open(filepath.Clean(filename))
	defer file.Close()
	if fileErr != nil {
		logger.Fatal(fileErr)
	}

	params := uploads.CreateUploadParams{
		File:             file,
		PaymentRequestID: paymentRequestID,
	}
	params.SetTimeout(time.Second * 30)

	resp, errCreatePaymentRequestUpload := primeGateway.Uploads.CreateUpload(&params)
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
