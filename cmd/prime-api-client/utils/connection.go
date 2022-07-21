package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/spf13/viper"
	"pault.ag/go/pksigner"

	"github.com/transcom/mymove/pkg/cli"
	primeClient "github.com/transcom/mymove/pkg/gen/primeclient"
	supportClient "github.com/transcom/mymove/pkg/gen/supportclient"
)

// CreatePrimeClientWithCACStoreParam creates the prime api client
func CreatePrimeClientWithCACStoreParam(v *viper.Viper, store *pksigner.Store) (*primeClient.Mymove, *pksigner.Store, error) {

	// Use command line inputs
	hostname := v.GetString(HostnameFlag)
	port := v.GetInt(PortFlag)
	insecure := v.GetBool(InsecureFlag)

	var httpClient *http.Client

	// The client certificate comes from a smart card
	//var store *pksigner.Store
	if v.GetBool(cli.CACFlag) {
		/* TODO how to check if logged in already??
		var errCACStoreLogin error
		store, errCACStoreLogin = cli.CACStoreLogin(v, store)
		if errCACStoreLogin != nil {
			log.Fatal(errCACStoreLogin)
		}
		*/
		cert, errTLSCert := store.TLSCertificate()
		if errTLSCert != nil {
			fmt.Printf("\n\nstore.TLSCertificate() failed with: [%s]\n\n", errTLSCert.Error())
			log.Fatal(errTLSCert)
		}

		// must explicitly state what signature algorithms we allow as of Go 1.14 to disable RSA-PSS signatures
		cert.SupportedSignatureAlgorithms = []tls.SignatureScheme{tls.PKCS1WithSHA256}

		//RA Summary: gosec - G402 - Look for bad TLS connection settings
		//RA: The linter is flagging this line of code because we are passing in a boolean value which can set InsecureSkipVerify to true.
		//RA: In production, the value of this flag is always false. We are, however, using
		//RA: this flag during local development to test the Prime API as further specified in the following docs:
		//RA: * https://github.com/transcom/prime_api_deliverable/wiki/Getting-Started#run-prime-api-client
		//RA: * https://transcom.github.io/mymove-docs/docs/dev/getting-started/How-to-Test-the-Prime-API
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: CAT III
		// #nosec G402
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

	verbose := cli.LogLevelIsDebug(v)
	hostWithPort := fmt.Sprintf("%s:%d", hostname, port)
	myRuntime := runtimeClient.NewWithClient(hostWithPort, primeClient.DefaultBasePath, []string{"https"}, httpClient)
	myRuntime.EnableConnectionReuse()
	myRuntime.SetDebug(verbose)

	primeGateway := primeClient.New(myRuntime, nil)

	return primeGateway, store, nil
}

// CreatePrimeClient creates the prime api client
func CreatePrimeClient(v *viper.Viper) (*primeClient.Mymove, *pksigner.Store, error) {

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

		//RA Summary: gosec - G402 - Look for bad TLS connection settings
		//RA: The linter is flagging this line of code because we are passing in a boolean value which can set InsecureSkipVerify to true.
		//RA: In production, the value of this flag is always false. We are, however, using
		//RA: this flag during local development to test the Prime API as further specified in the following docs:
		//RA: * https://github.com/transcom/prime_api_deliverable/wiki/Getting-Started#run-prime-api-client
		//RA: * https://transcom.github.io/mymove-docs/docs/dev/getting-started/How-to-Test-the-Prime-API
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: CAT III
		// #nosec G402
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

	verbose := cli.LogLevelIsDebug(v)
	hostWithPort := fmt.Sprintf("%s:%d", hostname, port)
	myRuntime := runtimeClient.NewWithClient(hostWithPort, primeClient.DefaultBasePath, []string{"https"}, httpClient)
	myRuntime.EnableConnectionReuse()
	myRuntime.SetDebug(verbose)

	primeGateway := primeClient.New(myRuntime, nil)

	return primeGateway, store, nil
}

// CreateSupportClient creates the support api client
func CreateSupportClient(v *viper.Viper) (*supportClient.Mymove, *pksigner.Store, error) {

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

		//RA Summary: gosec - G402 - Look for bad TLS connection settings
		//RA: The linter is flagging this line of code because we are passing in a boolean value which can set InsecureSkipVerify to true.
		//RA: In production, the value of this flag is always false. We are, however, using
		//RA: this flag during local development to test the Prime API as further specified in the following docs:
		//RA: * https://github.com/transcom/prime_api_deliverable/wiki/Getting-Started#run-prime-api-client
		//RA: * https://transcom.github.io/mymove-docs/docs/dev/getting-started/How-to-Test-the-Prime-API
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: CAT III
		// #nosec G402
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

	verbose := cli.LogLevelIsDebug(v)
	hostWithPort := fmt.Sprintf("%s:%d", hostname, port)
	myRuntime := runtimeClient.NewWithClient(hostWithPort, supportClient.DefaultBasePath, []string{"https"}, httpClient)
	myRuntime.EnableConnectionReuse()
	myRuntime.SetDebug(verbose)

	supportGateway := supportClient.New(myRuntime, nil)

	return supportGateway, store, nil
}
