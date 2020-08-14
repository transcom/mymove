package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	openapi "github.com/go-openapi/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ParseFlags parses the root command line flags
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

// HandleGatewayError handles errors returned by the gateway
func HandleGatewayError(err error, logger Logger) error {
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
