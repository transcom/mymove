package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/transcom/mymove/pkg/cli"
)

var logURL string

var logMessages []string

// a shared httpClient
var httpClient = &http.Client{
	Timeout: time.Duration(5 * time.Second),
}

// sendLog POSTs the `msg` to the the logURL endpoint, which is only
// available in the healthServer configured in cmd/milmove/serve.go
//
// The stdout/stderr logs from the health check run in AWS are not
// captured or available anywhere, so for debugging why the health
// check failed, having a logs endpoint is super helpful
func sendLog(msg string) {
	if logURL != "" {
		rdr := strings.NewReader(msg)
		// it doesn't matter if this succeeds or fails, we are doing
		// best effort delivery of the log
		_, _ = httpClient.Post(logURL, "text/plain", rdr)
	}
}

// reportError sends the `msg` to the logURL and writes to stderr so
// that the health check fails
func reportError(msg string) {
	sendLog(msg)
	// The docs for the AWS health check says
	//
	// An exit code of 0, with no stderr output, indicates
	// success, and a non-zero exit code indicates failure
	//
	// thus, print to stderr when the health check fails
	fmt.Fprintln(os.Stderr, msg)
}

// healthCheck does an HTTP GET on the provided `url`. Any response
// other than 200 OK is considered unhealthy
func healthCheck(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		werr := fmt.Errorf("Error create request for url `%s` error: %w", url, err)
		return werr
	}
	req.Header.Set("User-Agent", "milmove-ecs-health-check/1.0")
	r, err := httpClient.Do(req)
	if err != nil {
		werr := fmt.Errorf("Error checking health for url `%s` error: %w", url, err)
		return werr
	}
	if r.StatusCode != http.StatusOK {
		werr := fmt.Errorf("Health failed for url `%s`, status `%d`", url, r.StatusCode)
		return werr
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
//
// Why use a separate health listener thread in cmd/milmove/serve.go
// instead of using one of the existing listeners?
//
// We current deploy two ECS services: one for the my/office/admin
// apps and one for the prime api. It turns out that both have
// TLS_ENABLED=true so the tlsListener is enabled. Why is the TLS
// listener enabled for the prime api when it only terminates mTLS
// connections? Because the AWS ELB-HealthChecker needs an endpoint to
// connect to and we cannot configure it to use mTLS.
//
// So why not use that listener if it is enabled in both ECS services?
// The health checker would need to use SSL to connect. We could
// configure it with the proper certificate authorities, but then if
// the SSL certificate ever expired, the health check would fail
// immediately and the service would restart continuously, effectively
// being unavailable. That seems like too drastic a failure mode for
// an expired certificate.
//
// We could have the health check client ignore the ssl certificate,
// connecting insecurely, but then we would need to get RA approval.
//
// Finally, when the health check runs, its logs are not captured
// anywhere. If the health check fails, it's very difficult to
// troubleshoot. Having a separate health server means we can easily
// have an unauthenticated logs endpoint that the health check client
// can send logs to. It is not ideal or guaranteed to work, but it is
// better than nothing.
func healthFunction(cmd *cobra.Command, args []string) error {
	// Prepare to parse command line options / environment variables
	// using the viper library
	v, err := initializeViper(cmd, args)
	if err != nil {
		return err
	}

	healthEnabled := v.GetBool(cli.HealthListenerFlag)
	if !healthEnabled {
		err = errors.New("Health Check Listener Not Enabled")
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// cleanup that runs when this function ends
	defer func() {
		// if this function panics, try to log why before reporting
		// the error
		if r := recover(); r != nil {
			msg := fmt.Sprintf("health recovered from panic: %v\n", r)
			reportError(msg)
		}

		// send all log messages before exit if possible
		if len(logMessages) > 0 {
			for i := range logMessages {
				reportError(logMessages[i])
			}
		}
	}()

	// configure the logURL based on the configured HealthPort
	port := v.GetInt(cli.HealthPortFlag)
	logURL = fmt.Sprintf("http://localhost:%d/logs", port)

	// configure the health check endpoint URL based on the configured HealthPort
	url := fmt.Sprintf("http://localhost:%d/health", port)

	err = healthCheck(url)
	if err != nil {
		logMessages = append(logMessages, err.Error())
		return err
	}

	return nil
}
