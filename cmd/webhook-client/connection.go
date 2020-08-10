package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/spf13/viper"
	"pault.ag/go/pksigner"

	// #TODO: Story to remove dependency on gen/supportclient. The support endpoint
	// in support.yaml will function as our "server" for local testing. In production,
	// we need to send the payload to the Prime / an external server and will not
	// know any internals around their implementation.
	supportClient "github.com/transcom/mymove/pkg/gen/supportclient"

	"github.com/transcom/mymove/pkg/cli"
)

// CreateClient creates the support api client
func CreateClient(v *viper.Viper) (*supportClient.Mymove, *pksigner.Store, error) {

	// Use command line inputs
	hostname := v.GetString(HostnameFlag)
	port := v.GetInt(PortFlag)
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
	myRuntime := runtimeClient.NewWithClient(hostWithPort, supportClient.DefaultBasePath, []string{"https"}, httpClient)
	myRuntime.EnableConnectionReuse()
	myRuntime.SetDebug(verbose)

	supportGateway := supportClient.New(myRuntime, nil)

	return supportGateway, store, nil
}
