package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	primeClient "github.com/transcom/mymove/pkg/gen/primeclient"
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

func checkFetchMTOsConfig(v *viper.Viper, logger *log.Logger) error {

	err := cli.CheckCAC(v)
	if err != nil {
		return err
	}

	err = cli.CheckPrimeAPI(v)
	if err != nil {
		return err
	}

	err = cli.CheckVerbose(v)
	if err != nil {
		return err
	}

	return nil
}

func initFetchMTOsFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitPrimeAPIFlags(flag)
	cli.InitVerboseFlags(flag)

	flag.SortFlags = false
}

func fetchMTOs(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse args")
	}
	flags := cmd.Flags()
	v := viper.New()
	err = v.BindPFlags(flags)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	err = checkFetchMTOsConfig(v, logger)
	if err != nil {
		logger.Fatal(err)
	}

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

	var params mto.FetchMTOUpdatesParams
	params.SetTimeout(time.Second * 30)
	resp, errFetchMTOUpdates := primeGateway.MoveTaskOrder.FetchMTOUpdates(&params)
	if errFetchMTOUpdates != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		log.Fatal(errFetchMTOUpdates.Error())
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			log.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		log.Fatal(resp.Error())
	}

	return nil
}