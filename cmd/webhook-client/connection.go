package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pault.ag/go/pksigner"

	"github.com/transcom/mymove/pkg/cli"
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

// CreateClient creates the support api client
func CreateClient(v *viper.Viper) (*http.Client, *pksigner.Store, error) {

	// Use command line inputs
	//hostname := v.GetString(HostnameFlag)
	//port := v.GetInt(PortFlag)
	insecure := v.GetBool(InsecureFlag)

	var httpClient *http.Client

	// The client certificate comes from a smart card
	var store *pksigner.Store
	if v.GetBool(cli.CACFlag) {
		var errGetCACStore error
		store, errGetCACStore = cli.GetCACStore(v)
		if errGetCACStore != nil {
			log.Fatal(errGetCACStore)
		}
		cert, errTLSCert := store.TLSCertificate()
		if errTLSCert != nil {
			log.Fatal(errTLSCert)
		}

		// must explicitly state what signature algorithms we allow as of Go 1.14 to disable RSA-PSS signatures
		cert.SupportedSignatureAlgorithms = []tls.SignatureScheme{tls.PKCS1WithSHA256}

		// #nosec b/c gosec triggers on InsecureSkipVerify
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{*cert},
			InsecureSkipVerify: insecure,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
		}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient = &http.Client{
			Transport: transport,
		}
	} else if !v.GetBool(cli.CACFlag) {
		certPath := v.GetString(CertPathFlag)
		keyPath := v.GetString(KeyPathFlag)

		// var errRuntimeClientTLS error
		// httpClient, errRuntimeClientTLS = runtimeClient.TLSClient(runtimeClient.TLSClientOptions{
		// 	Key:                keyPath,
		// 	Certificate:        certPath,
		// 	InsecureSkipVerify: insecure})
		// if errRuntimeClientTLS != nil {
		// 	log.Fatal(errRuntimeClientTLS)
		// }
		// Create a HTTPS client and supply the created CA pool and certificate
		// Read the key pair to create certificate
		// certPath = "config/tls/devlocal-mtls.cer"
		// keyPath = "config/tls/devlocal-mtls.key"
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		cert.SupportedSignatureAlgorithms = []tls.SignatureScheme{tls.PKCS1WithSHA256}
		if err != nil {
			log.Fatal(err)
		}
		// #nosec
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates:       []tls.Certificate{cert},
					InsecureSkipVerify: insecure,
				},
			},
		}

	}

	// verbose := v.GetBool(cli.VerboseFlag)
	// hostWithPort := fmt.Sprintf("%s:%d", hostname, port)
	// myRuntime := runtimeClient.NewWithClient(hostWithPort, "/support/v1", []string{"https"}, httpClient)
	// myRuntime.EnableConnectionReuse()
	// myRuntime.SetDebug(verbose)

	// gateway := supportClient.New(myRuntime, nil)

	return httpClient, store, nil
}
