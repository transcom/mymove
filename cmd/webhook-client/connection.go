package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//"pault.ag/go/pksigner"
	//"github.com/transcom/mymove/pkg/cli"
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
	Logger      Logger
	Host        string
	BasePath    string
	ContentType string
	Insecure    bool
	// certpath    string
	// keypath     string
}

// NewWebhookRuntime creates and returns a runtime client
func NewWebhookRuntime() (*WebhookRuntime, error) {
	wr := WebhookRuntime{}

	return &wr, nil
}

// SetupClient sets up either CAC or cert, key client
func (wr *WebhookRuntime) SetupClient(isCacEnabled bool, cert *tls.Certificate) (*WebhookRuntime, error) {

	// if isCacEnabled {
	// 	fmt.Print("ARE WE WORKING?")
	// 	cert, store, err := GetCacCertificate()
	// }
	// CAC setup
	// tlsConfig := &tls.Config{
	// 	Certificates:       []tls.Certificate{*cert},
	// 	InsecureSkipVerify: insecure,
	// 	MinVersion:         tls.VersionTLS12,
	// 	MaxVersion:         tls.VersionTLS12,
	// }
	// transport := &http.Transport{
	// 	TLSClientConfig: tlsConfig,
	// }
	// httpClient = &http.Client{
	// 	Transport: transport,
	// }

	// regular setup
	// #nosec b/c gosec triggers on InsecureSkipVerify
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{*cert},
				InsecureSkipVerify: wr.Insecure,
			},
		},
		Timeout: time.Second * 30,
	}

	wr.client = httpClient
	wr.Host = "https://primelocal:9443"
	wr.BasePath = "/support/v1/webhook-notify"
	wr.ContentType = "application/json; charset=utf-8"

	return wr, nil
}

// Post WebhookRuntime comment goes here
func (wr *WebhookRuntime) Post(data []byte) error {
	json := bytes.NewBuffer(data)
	r, err := wr.client.Post(wr.Host+wr.BasePath, wr.ContentType, json)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	if wr.Debug {
		Debug(httputil.DumpResponse(r, true))
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Print response body to stdout
	fmt.Printf("%s\n", body)

	return nil
}

// GetCacCertificate returns cert to use for tls
// func GetCacCertificate(v *viper.Viper) (*tls.Certificate, *pksigner.Store, error) {
// 	var store *pksigner.Store

// 	var errGetCACStore error
// 	store, errGetCACStore = cli.GetCACStore(v)
// 	if errGetCACStore != nil {
// 		log.Fatal(errGetCACStore)
// 	}
// 	cert, errTLSCert := store.TLSCertificate()
// 	if errTLSCert != nil {
// 		log.Fatal(errTLSCert)
// 	}

// 	// must explicitly state what signature algorithms we allow as of Go 1.14 to disable RSA-PSS signatures
// 	cert.SupportedSignatureAlgorithms = []tls.SignatureScheme{tls.PKCS1WithSHA256}

// 	return cert, store, nil
// }

// CreateClient creates the webhook client
func CreateClient(v *viper.Viper) (*WebhookRuntime, error) {

	insecure := v.GetBool(InsecureFlag)

	//var httpClient *http.Client
	//var runtimeClient *WebhookRuntime

	// The client certificate comes from a smart card
	//var store *pksigner.Store
	// if v.GetBool(cli.CACFlag) {
	// 	cert, store, err := GetCacCertificate(v)

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	runtimeClient, err := NewWebhookRuntime()
	// 	runtimeClient.cacCert = cert
	// 	rc := runtimeClient.SetupClient(true, nil, nil)
	// 	// #nosec b/c gosec triggers on InsecureSkipVerify
	// 	// tlsConfig := &tls.Config{
	// 	// 	Certificates:       []tls.Certificate{*cert},
	// 	// 	InsecureSkipVerify: insecure,
	// 	// 	MinVersion:         tls.VersionTLS12,
	// 	// 	MaxVersion:         tls.VersionTLS12,
	// 	// }
	// 	// transport := &http.Transport{
	// 	// 	TLSClientConfig: tlsConfig,
	// 	// }
	// 	// httpClient = &http.Client{
	// 	// 	Transport: transport,
	// 	// }
	// //}
	// else if !v.GetBool(cli.CACFlag) {
	certPath := v.GetString(CertPathFlag)
	keyPath := v.GetString(KeyPathFlag)
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)

	if err != nil {
		log.Fatal(err)
	}
	runtimeClient, err := NewWebhookRuntime()

	if err != nil {
		log.Fatal(err)
	}
	runtimeClient.Insecure = insecure
	runtimeClient.Debug = false
	rc, err := runtimeClient.SetupClient(false, &cert)
	if err != nil {
		log.Fatal(err)
	}
	// #nosec b/c gosec triggers on InsecureSkipVerify
	// httpClient = &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{
	// 			Certificates:       []tls.Certificate{cert},
	// 			InsecureSkipVerify: insecure,
	// 		},
	// 	},
	// 	Timeout: time.Second * 30,
	// }

	// }
	// return rc, store, nil
	return rc, nil
}
