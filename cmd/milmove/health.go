package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// a shared httpClient
var httpClient = &http.Client{
	Timeout: time.Duration(5 * time.Second),
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
func healthFunction(cmd *cobra.Command, args []string) error {
	// Prepare to parse command line options / environment variables
	// using the viper library
	v, err := initializeViper(cmd, args)
	if err != nil {
		return err
	}
	zapConfig := logging.BuildZapConfig(logging.WithEnvironment(v.GetString(cli.LoggingEnvFlag)),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	// write to stdout of the server process so log is captured
	// https://github.com/aws/containers-roadmap/issues/1114#issuecomment-1260905963
	zapConfig.OutputPaths = []string{"/proc/1/fd/1"}
	logger, err := zapConfig.Build()
	if err != nil {
		lout, lerr := os.OpenFile("/proc/1/fd/1", os.O_WRONLY|os.O_APPEND, 0644)
		if lerr == nil {
			// try to write something
			_, _ = fmt.Fprintln(lout, "{\"level\":\"error\",\"msg\":\"Health Check logger failed to initialize\"}")
			// this might produce bogus JSON
			// again ignore errors
			_, _ = fmt.Fprintln(lout, "{\"level\":\"error\",\"msg\":\"Health Check: "+lerr.Error()+"\"}")
		}
		return err
	}

	healthEnabled := v.GetBool(cli.HealthListenerFlag)
	if !healthEnabled {
		err = errors.New("Health Check Listener Not Enabled")
		logger.Error("Health Check failed, listener not enabled", zap.Any(cli.HealthListenerFlag, healthEnabled))
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	// cleanup that runs when this function ends
	defer func() {
		// if this function panics, try to log why before reporting
		// the error
		if r := recover(); r != nil {
			msg := fmt.Sprintf("Health Check recovered from panic: %v\n", r)
			logger.Error(msg)
		}
		// try to sync any logs to the output file, ignore errors
		_ = logger.Sync()
	}()

	// configure the logURL based on the configured HealthPort
	port := v.GetInt(cli.HealthPortFlag)
	// configure the health check endpoint URL based on the configured HealthPort
	url := fmt.Sprintf("http://localhost:%d/health", port)

	err = healthCheck(url)
	if err != nil {
		logger.Error("Health Check failed", zap.Any("url", url), zap.Error(err))
		return err
	}

	return nil
}
