package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"

	primeClient "github.com/transcom/mymove/pkg/gen/primeclient"
	"github.com/transcom/nom/pkg/pkcs11"
)

const (
	// PKCS11ModuleFlag is the location of the PCKS11 module to use with the smart card
	PKCS11ModuleFlag string = "pkcs11module"
	// TokenLabel is the Token Label to use with the smart card
	TokenLabelFlag string = "tokenlabel"
	// CertLabel is the Certificate Label to use with the smart card
	CertLabelFlag string = "certlabel"
	// KeyLabel is the Key Label to use with the smart card
	KeyLabelFlag string = "keylabel"
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
	flag.String(PKCS11ModuleFlag, "/usr/local/lib/pkcs11/opensc-pkcs11.so", "Smart card: Path to the PKCS11 module to use")
	flag.String(TokenLabelFlag, "", "Smart card: name of the token to use")
	flag.String(CertLabelFlag, "Certificate for PIV Authentication", "Smart card: label of the public cert")
	flag.String(KeyLabelFlag, "PIV AUTH key", "Smart card: label of the private key")
	flag.String(CertPathFlag, "", "Smart card: label of the public cert")
	flag.String(KeyPathFlag, "", "Smart card: label of the private key")
	flag.String(HostnameFlag, "primelocal", "The hostname to connect to")
	flag.Int(PortFlag, 9443, "The port to connect to")
	flag.Bool(InsecureFlag, false, "Skip TLS verification and validation")
	flag.BoolP(VerboseFlag, "v", false, "Show extra output for debugging")
	flag.SortFlags = false
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

	if !v.GetBool(VerboseFlag) {
		// Disable any logging that isn't attached to the logger unless using the verbose flag
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)

		// Remove the flags for the logger
		logger.SetFlags(0)
	}

	// Use command line inputs
	pkcs11ModulePath := v.GetString(PKCS11ModuleFlag)
	tokenLabel := v.GetString(TokenLabelFlag)
	certLabel := v.GetString(CertLabelFlag)
	keyLabel := v.GetString(KeyLabelFlag)
	certPath := v.GetString(CertPathFlag)
	keyPath := v.GetString(KeyLabelFlag)
	hostname := v.GetString(HostnameFlag)
	port := v.GetInt(PortFlag)
	insecure := v.GetBool(InsecureFlag)

	var httpClient *http.Client

	// The client certificate comes from a smart card
	if pkcs11ModulePath != "" && certPath != "" && keyPath != "" {
		pkcsConfig := pkcs11.Config{
			Module:           pkcs11ModulePath,
			CertificateLabel: certLabel,
			PrivateKeyLabel:  keyLabel,
			TokenLabel:       tokenLabel,
		}

		store, err := pkcs11.New(pkcsConfig)
		if err != nil {
			log.Fatal(err)
		}
		defer store.Close()

		inputUI := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		pin, err := inputUI.Ask("PIN", &input.Options{
			Default:     "",
			HideOrder:   true,
			HideDefault: true,
			Required:    true,
			Loop:        true,
			Mask:        true,
			ValidateFunc: func(input string) error {
				matched, matchErr := regexp.Match("^\\d+$", []byte(input))
				if matchErr != nil {
					return matchErr
				}
				if !matched {
					return errors.New("Invalid")
				}
				return nil
			},
		})
		if err != nil {
			os.Exit(1)
		}

		err = store.Login(pin)
		if err != nil {
			log.Fatal(err)
		}

		cert, err := store.TLSCertificate()
		if err != nil {
			panic(err)
		}
		// #nosec b/c gosec triggers on InsecureSkipVerify
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{*cert},
			InsecureSkipVerify: insecure,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient = &http.Client{
			Transport: transport,
		}
	} else {
		var err error
		httpClient, err = runtimeClient.TLSClient(runtimeClient.TLSClientOptions{Key: keyPath, Certificate: certPath, InsecureSkipVerify: insecure})
		if err != nil {
			log.Fatal(err)
		}
	}

	hostWithPort := fmt.Sprintf("%s:%d", hostname, port)
	myRuntime := runtimeClient.NewWithClient(hostWithPort, primeClient.DefaultBasePath, []string{"https"}, httpClient)
	myRuntime.EnableConnectionReuse()
	myRuntime.SetDebug(true)
	primeGateway := primeClient.New(myRuntime, nil)
	fmt.Println(primeGateway)
}
