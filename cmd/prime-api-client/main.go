package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	primeClient "github.com/transcom/mymove/pkg/gen/primeclient"
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

const (

	// CertPathFlag is the path to the certificate to use for TLS
	CertPathFlag string = "certpath"
	// KeyPathFlag is the path to the key to use for TLS
	KeyPathFlag string = "keypath"
	// HostnameFlag is the hostname to connect to
	HostnameFlag string = "host"
	// PortFlag is the port to connect to
	PortFlag string = "port"
	// Insecure flag indicates that TLS verification and validation can be skipped
	InsecureFlag string = "insecure"
	// VerboseFlag holds string identifier for command line usage
	VerboseFlag string = "verbose"
)

// initialize flags
func initFlags(flag *pflag.FlagSet) {

	cli.InitCACFlags(flag)

	flag.String(CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(HostnameFlag, "primelocal", "The hostname to connect to")
	flag.Int(PortFlag, 9443, "The port to connect to")
	flag.Bool(InsecureFlag, false, "Skip TLS verification and validation")
	flag.BoolP(VerboseFlag, "v", false, "Show extra output for debugging")
	flag.SortFlags = false
}

func checkConfig(v *viper.Viper, logger *log.Logger) error {

	if err := cli.CheckCAC(v); err != nil {
		return err
	}

	if !v.GetBool(cli.CACFlag) {
		certPath := v.GetString(CertPathFlag)
		if certPath == "" {
			return fmt.Errorf("%q is invalid: %w", CertPathFlag, &cli.ErrInvalidPath{Path: certPath})
		} else if _, err := os.Stat(certPath); err != nil {
			return fmt.Errorf("%q is invalid: %w", CertPathFlag, &cli.ErrInvalidPath{Path: certPath})
		}

		keyPath := v.GetString(KeyPathFlag)
		if keyPath == "" {
			return fmt.Errorf("%q is invalid: %w", KeyPathFlag, &cli.ErrInvalidPath{Path: keyPath})
		} else if _, err := os.Stat(keyPath); err != nil {
			return fmt.Errorf("%q is invalid: %w", KeyPathFlag, &cli.ErrInvalidPath{Path: keyPath})
		}
	}
	return nil
}

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Arg parse failed")
		return
	}

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		fmt.Println("Arg binding failed")
		return
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	verbose := v.GetBool(VerboseFlag)
	if !verbose {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Use command line inputs
	hostname := v.GetString(HostnameFlag)
	port := v.GetInt(PortFlag)
	insecure := v.GetBool(InsecureFlag)

	var httpClient *http.Client

	// The client certificate comes from a smart card
	if v.GetBool(cli.CACFlag) {
		store, errStore := cli.GetCACStore(v)
		defer store.Close()
		if errStore != nil {
			log.Fatal(errStore)
		}
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
		certPath := v.GetString(CertPathFlag)
		keyPath := v.GetString(KeyPathFlag)

		var errRuntimeClientTLS error
		httpClient, errRuntimeClientTLS = runtimeClient.TLSClient(runtimeClient.TLSClientOptions{
			Key:                keyPath,
			Certificate:        certPath,
			InsecureSkipVerify: insecure})
		if errRuntimeClientTLS != nil {
			log.Fatal(errRuntimeClientTLS)
		}
	}

	hostWithPort := fmt.Sprintf("%s:%d", hostname, port)
	myRuntime := runtimeClient.NewWithClient(hostWithPort, primeClient.DefaultBasePath, []string{"https"}, httpClient)
	myRuntime.EnableConnectionReuse()
	myRuntime.SetDebug(verbose)

	primeGateway := primeClient.New(myRuntime, nil)

	var params mto.FetchMTOUpdatesParams
	params.SetTimeout(time.Second * 30)
	resp, errFetchMTOUpdates := primeGateway.MoveTaskOrder.FetchMTOUpdates(&params)
	if errFetchMTOUpdates != nil {
		log.Fatal(errFetchMTOUpdates)
	}

	payload, errJSONMarshall := json.Marshal(resp.GetPayload())
	if errJSONMarshall != nil {
		log.Fatal(errJSONMarshall)
	}
	fmt.Println(string(payload))
}
