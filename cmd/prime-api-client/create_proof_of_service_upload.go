package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/gen/primeclient/uploads"
	"github.com/transcom/mymove/pkg/gen/primemessages"
)

// initCreateProofOfServiceUploadFlags initializes flags.
func initCreateProofOfServiceUploadFlags(flag *pflag.FlagSet) {
	flag.String(FilenameFlag, "", "Path to the file with the proof of service upload JSON payload")

	flag.SortFlags = false
}

// checkCreateProofOfServiceUploadConfig checks the args.
func checkCreateProofOfServiceUploadConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		return err
	}

	if v.GetString(FilenameFlag) == "" && (len(args) < 1 || len(args) > 0 && !containsDash(args)) {
		return errors.New("proof-of-service-upload expects a file to be passed in")
	}

	return nil
}

func createAferoFile(filename string, bytes []int) (afero.file, error) {
	var fs = afero.NewMemMapFs()

	/*
	file, err := os.Open(filename)
	if err != nil {
		//suite.logger.Fatal("Error opening local file", zap.Error(err))
	}
	*/


	//suite.NotNil(fs)
	uploadFile, err := fs.Create(filename)
	if err != nil {
		//suite.logger.Fatal("Error creating afero file", zap.Error(err))
	}

	uploadFile.Write(bytes)

	io.ByteWriter(uploadFile, bytes)
	if err != nil {
		//suite.logger.Fatal("Error copying to afero file", zap.Error(err))
	}

	return outputFile, nil
}

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
	var proofOfServiceDocs primemessages.ProofOfServiceDocs
	err = jsonDecoder.Decode(&proofOfServiceDocs)
	if err != nil {
		return fmt.Errorf("decoding data failed: %w", err)
	}

	for posDoc, _ := range proofOfServiceDocs.Uploads {
		outputFile, err := suite.helperNewTempFile()
		suite.Nil(err)
		defer outputFile.Close()

		written, err := io.Copy(outputFile, download)
		suite.Nil(err)
		suite.NotEqual(0, written)

		info, err := outputFile.Stat()
		suite.Equal(fixtureFileInfo.Size(), info.Size())
	}

	params := uploads.CreateUploadParams{
		File: proofOfServiceDocs.Uploads[0].Bytes
		Body: &proofOfServiceDocs,
	}
	params.SetTimeout(time.Second * 30)

	resp, errCreatePaymentRequest := primeGateway.PaymentRequests.CreatePaymentRequest(&params)
	if errCreatePaymentRequest != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errCreatePaymentRequest.Error())
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
