package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/transcom/mymove/pkg/cli"
)

func healthCheck(httpClient *http.Client, protocol string, host string, port int) error {
	url := fmt.Sprintf("%s://%s:%d/health", protocol, host, port)
	r, err := httpClient.Get(url)
	if err != nil {
		// The docs for the AWS health check says
		//
		// An exit code of 0, with no stderr output, indicates
		// success, and a non-zero exit code indicates failure
		//
		// thus, print to stderr when the health check fails
		fmt.Fprintf(os.Stderr, "Error checking noTLS health for url `%s`: %s", url, err)
		return err
	}
	if r.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Health NOT OK for url `%s`, status `%d`", url, r.StatusCode)
		return errors.New("Health check failed")
	}

	return nil
}

// we use a health function in the milmove binary instead of using
// cmd/health-checker because
//
//  1. We don't want to have to install another binary in the container
//     for health checks
//  2. We want to use the same command line options used to start the
//     server as when checking health
func healthFunction(cmd *cobra.Command, args []string) error {
	// cleanup that runs when this function ends
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "health recovered from panic: %v", r)
		}
	}()

	// Prepare to parse command line options / environment variables
	// using the viper library
	v, err := initializeViper(cmd, args)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	listenInterface := v.GetString(cli.InterfaceFlag)
	host := "localhost"
	if listenInterface != "" {
		host = listenInterface
	}
	noTLSEnabled := v.GetBool(cli.NoTLSListenerFlag)
	if noTLSEnabled {
		fmt.Println("Checking noTLS server health")
		port := v.GetInt(cli.NoTLSPortFlag)
		err := healthCheck(httpClient, "http", host, port)
		if err != nil {
			return err
		}
	}

	tlsEnabled := v.GetBool(cli.TLSListenerFlag)
	if tlsEnabled {
		fmt.Println("Checking TLS server health")
		port := v.GetInt(cli.TLSPortFlag)
		err := healthCheck(httpClient, "https", host, port)
		if err != nil {
			return err
		}
	}

	// To test mutualTLS, we would need a valid client cert for each
	// environment, but today those aren't included in the container,
	// so skip the health check. In the future maybe we could inject
	// them via environment variables
	mutualTLSEnabled := v.GetBool(cli.MutualTLSListenerFlag)
	if mutualTLSEnabled {
		fmt.Println("WARNING: Skipping mutualTLS server health")
	}

	return nil
}
