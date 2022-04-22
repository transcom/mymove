package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	openapi "github.com/go-openapi/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// FilenameFlag is the name of the file being passed in
	FilenameFlag string = "filename"
	// IDFlag is the UUID of the object being retrieved
	IDFlag string = "id"
	// SinceFlag is the datetime for the `since` filter for fetching moves
	SinceFlag string = "since"
	// ETagFlag is the etag for the mto shipment being updated
	ETagFlag string = "etag"
	// PaymentRequestIDFlag is the payment request ID
	PaymentRequestIDFlag string = "paymentRequestID"
	// CertPathFlag is the path to the certificate to use for TLS
	CertPathFlag string = "certpath"
	// KeyPathFlag is the path to the key to use for TLS
	KeyPathFlag string = "keypath"
	// HostnameFlag is the hostname to connect to
	HostnameFlag string = "hostname"
	// PortFlag is the port to connect to
	PortFlag string = "port"
	// InsecureFlag indicates that TLS verification and validation can be skipped
	InsecureFlag string = "insecure"
	// WaitFlag is how long to wait for the server to respond. The
	// string is parsed by https://pkg.go.dev/time#ParseDuration
	WaitFlag string = "wait"
)

// ParseFlags parses the command line flags
func ParseFlags(cmd *cobra.Command, v *viper.Viper, args []string) error {

	errParseFlags := cmd.ParseFlags(args)
	if errParseFlags != nil {
		return fmt.Errorf("Could not parse args: %w", errParseFlags)
	}
	flags := cmd.Flags()
	errBindPFlags := v.BindPFlags(flags)
	if errBindPFlags != nil {
		return fmt.Errorf("Could not bind flags: %w", errBindPFlags)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	return nil
}

// ContainsDash returns true if the original command included an empty dash
func ContainsDash(args []string) bool {
	for _, arg := range args {
		if arg == "-" {
			return true
		}
	}
	return false
}

// CheckRootConfig checks the validity of the prime api flags
func CheckRootConfig(v *viper.Viper) error {
	err := cli.CheckCAC(v)
	if err != nil {
		return err
	}

	err = cli.CheckLogging(v)
	if err != nil {
		return err
	}

	if (v.GetString(CertPathFlag) != "" && v.GetString(KeyPathFlag) == "") || (v.GetString(CertPathFlag) == "" && v.GetString(KeyPathFlag) != "") {
		return fmt.Errorf("Both TLS certificate and key paths must be provided")
	}

	return nil
}

// DecodeJSONFileToPayload takes a filename, or stdin and decodes the file into
// the supplied json payload.
// If the filename is not supplied, the isStdin bool should be set to true to use stdin.
// If the file contains parameters that do not exist in the payload struct, it will fail with an error
// Otherwise it will populate the payload
func DecodeJSONFileToPayload(filename string, isStdin bool, payload interface{}) error {
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

// HandleGatewayError handles errors returned by the gateway
func HandleGatewayError(err error, logger *log.Logger) error {
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
