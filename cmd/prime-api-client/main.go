package main

import (
	// "crypto/sha256"
	"crypto/tls"
	// "encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-input"

	"github.com/transcom/nom/pkg/pkcs11"

	primeClient "github.com/transcom/mymove/pkg/gen/primeclient"
	mto "github.com/transcom/mymove/pkg/gen/primeclient/move_task_order"
)

type errInvalidPath struct {
	Path string
}

func (e *errInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %q", e.Path)
}

type errInvalidLabel struct {
	Cert string
	Key  string
}

func (e *errInvalidLabel) Error() string {
	return fmt.Sprintf("invalid cert label %q or key label %q", e.Cert, e.Key)
}

const (
	// CACFlag indicates that a CAC should be used
	CACFlag string = "cac"
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
	flag.Bool(CACFlag, false, "Use a CAC for authentication")
	flag.String(PKCS11ModuleFlag, "/usr/local/lib/pkcs11/opensc-pkcs11.so", "Smart card: Path to the PKCS11 module to use")
	flag.String(TokenLabelFlag, "", "Smart card: name of the token to use")
	flag.String(CertLabelFlag, "Certificate for PIV Authentication", "Smart card: label of the public cert")
	flag.String(KeyLabelFlag, "PIV AUTH key", "Smart card: label of the private key")
	flag.String(CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(HostnameFlag, "primelocal", "The hostname to connect to")
	flag.Int(PortFlag, 9443, "The port to connect to")
	flag.Bool(InsecureFlag, false, "Skip TLS verification and validation")
	flag.BoolP(VerboseFlag, "v", false, "Show extra output for debugging")
	flag.SortFlags = false
}

func checkConfig(v *viper.Viper, logger *log.Logger) error {

	if v.GetBool(CACFlag) {
		pkcs11ModulePath := v.GetString(PKCS11ModuleFlag)
		if pkcs11ModulePath == "" {
			return fmt.Errorf("%q is invalid: %w", PKCS11ModuleFlag, &errInvalidPath{Path: pkcs11ModulePath})
		} else if _, err := os.Stat(pkcs11ModulePath); err != nil {
			return fmt.Errorf("%q is invalid: %w", PKCS11ModuleFlag, &errInvalidPath{Path: pkcs11ModulePath})
		}

		certLabel := v.GetString(CertLabelFlag)
		keyLabel := v.GetString(KeyLabelFlag)
		if certLabel == "" || keyLabel == "" {
			return fmt.Errorf("%q or %q is invalid: %w", CertLabelFlag, KeyLabelFlag, &errInvalidLabel{Cert: certLabel, Key: keyLabel})
		}
	} else {
		certPath := v.GetString(CertPathFlag)
		if certPath == "" {
			return fmt.Errorf("%q is invalid: %w", CertPathFlag, &errInvalidPath{Path: certPath})
		} else if _, err := os.Stat(certPath); err != nil {
			return fmt.Errorf("%q is invalid: %w", CertPathFlag, &errInvalidPath{Path: certPath})
		}

		keyPath := v.GetString(KeyPathFlag)
		if keyPath == "" {
			return fmt.Errorf("%q is invalid: %w", KeyPathFlag, &errInvalidPath{Path: keyPath})
		} else if _, err := os.Stat(keyPath); err != nil {
			return fmt.Errorf("%q is invalid: %w", KeyPathFlag, &errInvalidPath{Path: keyPath})
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
	if v.GetBool(CACFlag) {
		pkcs11ModulePath := v.GetString(PKCS11ModuleFlag)
		tokenLabel := v.GetString(TokenLabelFlag)
		certLabel := v.GetString(CertLabelFlag)
		keyLabel := v.GetString(KeyLabelFlag)
		pkcsConfig := pkcs11.Config{
			Module:           pkcs11ModulePath,
			CertificateLabel: certLabel,
			PrivateKeyLabel:  keyLabel,
			TokenLabel:       tokenLabel,
		}

		store, errPKCS11New := pkcs11.New(pkcsConfig)
		if errPKCS11New != nil {
			log.Fatal(errPKCS11New)
		}
		defer store.Close()

		inputUI := &input.UI{
			Writer: os.Stdout,
			Reader: os.Stdin,
		}

		pin, errUIAsk := inputUI.Ask("PIN", &input.Options{
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
					return errors.New("Invalid PIN format")
				}
				return nil
			},
		})
		if errUIAsk != nil {
			os.Exit(1)
		}

		errLogin := store.Login(pin)
		if errLogin != nil {
			log.Fatal(errLogin)
		}

		cert, errTLSCert := store.TLSCertificate()
		if errTLSCert != nil {
			panic(errTLSCert)
		}

		// Get the fingerprint
		// hash := sha256.Sum256(cert.Certificate[0])
		// hashString := hex.EncodeToString(hash[:])
		// fmt.Println(hashString)

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
	} else if !v.GetBool(CACFlag) {
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
