package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

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

// WebhookRuntime comment here
type WebhookRuntime struct {
	client      *http.Client
	Debug       bool
	Host        string
	BasePath    string
	ContentType string
	Insecure    bool
}

// NewWebhookRuntime creates and returns a runtime client
func NewWebhookRuntime(hostWithPort string, contentType string, insecure bool, debug bool) *WebhookRuntime {
	wr := WebhookRuntime{
		Host:        "https://" + hostWithPort,
		ContentType: contentType,
		Insecure:    insecure,
		Debug:       debug,
	}

	return &wr
}

// SetupClient sets up either CAC or cert, key client
func (wr *WebhookRuntime) SetupClient(cert *tls.Certificate) (*WebhookRuntime, error) {

	// Set up the httpClient with tls certificate

	// #nosec b/c gosec triggers on InsecureSkipVerify
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

// Post WebhookRuntime comment goes here
func (wr *WebhookRuntime) Post(data []byte) error {
	json := bytes.NewBuffer(data)
	// Create the POST request
	req, err := http.NewRequest(
		"POST",
		wr.Host+wr.BasePath,
		json,
	)
	req.Header.Set("Content-type", wr.ContentType)

	if err != nil {
		return err
	}

	// Print out the request when debug mode is on
	if wr.Debug {
		Debug(httputil.DumpRequest(req, true))
	}

	// Send request and capture the response
	resp, err := wr.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// Print out the response when debug mode is on
	if wr.Debug {
		Debug(httputil.DumpResponse(resp, true))
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	// Print response body to stdout
	fmt.Printf("%s\n", body)

	return nil
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

	insecure := v.GetBool(InsecureFlag)
	verbose := v.GetBool(cli.VerboseFlag)

	hostname := v.GetString(HostnameFlag)
	port := v.GetInt(PortFlag)
	hostWithPort := fmt.Sprintf("%s:%d", hostname, port)

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

		certPath := v.GetString(CertPathFlag)
		keyPath := v.GetString(KeyPathFlag)
		loadCert, err = tls.LoadX509KeyPair(certPath, keyPath)

		if err != nil {
			return nil, nil, err
		}

		cert = &loadCert
	}

	runtimeClient := NewWebhookRuntime(hostWithPort, contentType, insecure, verbose)

	rc, err = runtimeClient.SetupClient(cert)

	if err != nil {
		return nil, nil, err
	}

	return rc, store, nil
}
