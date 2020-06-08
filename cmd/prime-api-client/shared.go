package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	openapi "github.com/go-openapi/runtime"
)

const (
	// FilenameFlag is the name of the file being passed in
	FilenameFlag string = "filename"
	// ETagFlag is the etag for the mto shipment being updated
	ETagFlag string = "etag"
	// PaymentRequestID is the ID of payment request to use
	PaymentRequestID string = "paymentRequestID"
)

func containsDash(args []string) bool {
	for _, arg := range args {
		if arg == "-" {
			return true
		}
	}
	return false
}

// decodeJSONFileToPayload takes a filename, or stdin and decodes the file into
// the supplied json payload.
// If the filename is not supplied, the isStdin bool should be set to true to use stdin.
// If the file contains parameters that do not exist in the payload struct, it will fail with an error
// Otherwise it will populate the payload
func decodeJSONFileToPayload(filename string, isStdin bool, payload interface{}) error {
	var reader *bufio.Reader
	if filename != "" {
		file, err := os.Open(filepath.Clean(filename))
		if err != nil {
			return fmt.Errorf("File open failed: %w", err)
		}
		reader = bufio.NewReader(file)
	} else if isStdin { // Uses std in if "-"" is provided instead
		reader = bufio.NewReader(os.Stdin)
	} else {
		return errors.New("no file input was found")
	}

	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.DisallowUnknownFields()

	// Read the json into the mto payload
	err := jsonDecoder.Decode(payload)
	if err != nil {
		return fmt.Errorf("File decode failed: %w", err)
	}

	return nil
}

// handleGatewayError handles errors returned by the gateway
func handleGatewayError(err error, logger *log.Logger) error {
	if _, ok := err.(*openapi.APIError); ok {
		// If you see an error like "unknown error (status 422)", it means
		// we hit a completely unhandled error that we should handle.
		// We should be enabling said error in the endpoint in swagger.
		// 422 for example is an Unprocessable Entity and is returned by the swagger
		// validation before it even hits the handler.
		apiErr := err.(*openapi.APIError).Response.(openapi.ClientResponse)
		logger.Fatal(fmt.Sprintf("%s: %s", err, apiErr.Message()))

	} else if typedErr, ok := err.(*url.Error); ok {
		// If the server is not running you are likely to see a connection error
		// This catches the error and prints a useful message.
		logger.Fatal(fmt.Sprintf("%s operation to %s failed, check if server is running : %s", typedErr.Op, typedErr.URL, typedErr.Err.Error()))
	}
	// If it is a handled error, we should be able to pull out the payload here
	data, _ := json.Marshal(err)
	fmt.Printf("%s", data)
	return nil
}
