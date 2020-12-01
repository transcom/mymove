package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/spf13/viper"
	"pault.ag/go/pksigner"

	"github.com/transcom/mymove/cmd/prime-api-client/utils"
	"github.com/transcom/mymove/pkg/cli"
)

// WebhookClientPoster is an interface that WebhookRuntime implements
type WebhookClientPoster interface {
	SetupClient(cert *tls.Certificate) (*WebhookRuntime, error)
	Post(data []byte, url string) (*http.Response, []byte, error)
}

// WebhookRuntime comment here
type WebhookRuntime struct {
	client      *http.Client
	Debug       bool
	ContentType string
	Insecure    bool
}

// NewWebhookRuntime creates and returns a runtime client
func NewWebhookRuntime(contentType string, insecure bool, debug bool) *WebhookRuntime {
	wr := WebhookRuntime{
		ContentType: contentType,
		Insecure:    insecure,
		Debug:       debug,
	}

	return &wr
}

// SetupClient sets up either CAC or cert, key client
func (wr *WebhookRuntime) SetupClient(cert *tls.Certificate) (*WebhookRuntime, error) {

	// Set up the httpClient with tls certificate

	//RA Summary: gosec - G402 - Look for bad TLS connection settings
	//RA: The linter is flagging this line of code because we are passing in a boolean value which can set InsecureSkipVerify to true.
	//RA: In production, the value of this flag is always false. We are, however, using
	//RA: this flag during local development to test the Prime API as further specified in the following docs:
	//RA: * https://github.com/transcom/prime_api_deliverable/wiki/Getting-Started#run-prime-api-client
	//RA: * https://github.com/transcom/mymove/wiki/How-to-Test-the-Prime-API-(Local,-Staging,-and-Experimental)#testing-locally
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: CAT III
	// #nosec G402
	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{*cert},
		InsecureSkipVerify: wr.Insecure,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
	}
	transport := &http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	httpClient := http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}

	// Add http client to our runtime client
	wr.client = &httpClient

	return wr, nil
}

// Post function of the WebhookRuntime http posts the data passed in and returns the
// response, body data, and any error
func (wr *WebhookRuntime) Post(data []byte, url string) (*http.Response, []byte, error) {
	bufferData := bytes.NewBuffer(data)
	// Create the POST request
	req, err := http.NewRequest(
		"POST",
		url,
		bufferData,
	)
	req.Header.Set("Content-type", wr.ContentType)

	if err != nil {
		return nil, nil, err
	}

	// Print out the request when debug mode is on
	if wr.Debug {
		output, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(output)) //todo switch to logger
	}

	// Send request and capture the response
	resp, err := wr.client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	// Print out the response when debug mode is on
	if wr.Debug {
		output, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(output))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, nil
}

// GetCacCertificate returns cert to use for tls
func GetCacCertificate(v *viper.Viper) (*tls.Certificate, *pksigner.Store, error) {
	var store *pksigner.Store
	var errGetCACStore error

	store, errGetCACStore = cli.GetCACStore(v)

	if errGetCACStore != nil {
		return nil, nil, errGetCACStore
	}

	cert, errTLSCert := store.TLSCertificate()
	if errTLSCert != nil {
		return nil, nil, errTLSCert
	}

	// Must explicitly state what signature algorithms we allow as of Go 1.14 to disable RSA-PSS signatures
	cert.SupportedSignatureAlgorithms = []tls.SignatureScheme{tls.PKCS1WithSHA256}

	return cert, store, nil
}

// CreateClient creates the webhook client
func CreateClient(v *viper.Viper) (*WebhookRuntime, *pksigner.Store, error) {
	var rc *WebhookRuntime
	var cert *tls.Certificate
	var err error
	var store *pksigner.Store

	insecure := v.GetBool(utils.InsecureFlag)
	verbose := v.GetBool(cli.VerboseFlag)
	contentType := "application/json; charset=utf-8"

	// Get the tls certificate
	// If using a CAC, the client cert comes from the card
	// Otherwise, use the certpath and keypath values
	if v.GetBool(cli.CACFlag) {
		cert, store, err = GetCacCertificate(v)

		if err != nil {
			return nil, nil, err
		}

	} else if !v.GetBool(cli.CACFlag) {
		var loadCert tls.Certificate

		certPath := v.GetString(utils.CertPathFlag)
		keyPath := v.GetString(utils.KeyPathFlag)
		loadCert, err = tls.LoadX509KeyPair(certPath, keyPath)

		if err != nil {
			return nil, nil, err
		}

		cert = &loadCert
	}

	runtimeClient := NewWebhookRuntime(contentType, insecure, verbose)

	rc, err = runtimeClient.SetupClient(cert)

	if err != nil {
		return nil, nil, err
	}

	return rc, store, nil
}
