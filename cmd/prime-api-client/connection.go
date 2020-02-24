package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	primeClient "github.com/transcom/mymove/pkg/gen/primeclient"
)

// CreateClient creates the prime api client
func CreateClient(cmd *cobra.Command, v *viper.Viper, args []string) (*primeClient.Mymove, error) {
	err := cmd.ParseFlags(args)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse args")
	}
	flags := cmd.Flags()
	err = v.BindPFlags(flags)
	if err != nil {
		return nil, errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Use command line inputs
	hostname := v.GetString(cli.HostnameFlag)
	port := v.GetInt(cli.PortFlag)
	insecure := v.GetBool(cli.InsecureFlag)

	var httpClient *http.Client

	// The client certificate comes from a smart card
	if v.GetBool(cli.CACFlag) {
		store, errStore := cli.GetCACStore(v)
		if errStore != nil {
			log.Fatal(errStore)
		}
		defer store.Close()
		cert, errTLSCert := store.TLSCertificate()
		if errTLSCert != nil {
			log.Fatal(errTLSCert)
		}

		// #nosec b/c gosec triggers on InsecureSkipVerify
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{*cert},
			InsecureSkipVerify: insecure,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient = &http.Client{
			Transport: transport,
		}
	} else if !v.GetBool(cli.CACFlag) {
		certPath := v.GetString(cli.CertPathFlag)
		keyPath := v.GetString(cli.KeyPathFlag)

		var errRuntimeClientTLS error
		httpClient, errRuntimeClientTLS = runtimeClient.TLSClient(runtimeClient.TLSClientOptions{
			Key:                keyPath,
			Certificate:        certPath,
			InsecureSkipVerify: insecure})
		if errRuntimeClientTLS != nil {
			log.Fatal(errRuntimeClientTLS)
		}
	}

	verbose := v.GetBool(cli.VerboseFlag)
	hostWithPort := fmt.Sprintf("%s:%d", hostname, port)
	myRuntime := runtimeClient.NewWithClient(hostWithPort, primeClient.DefaultBasePath, []string{"https"}, httpClient)
	myRuntime.EnableConnectionReuse()
	myRuntime.SetDebug(verbose)

	primeGateway := primeClient.New(myRuntime, nil)

	return primeGateway, nil
}