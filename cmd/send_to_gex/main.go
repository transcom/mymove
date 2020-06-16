package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services/invoice"
)

// Call this from command line with go run cmd/send_to_gex/main.go --edi <filepath>
func main() {
	flag := pflag.CommandLine
	// EDI Invoice Config
	flag.String("gex-basic-auth-username", "", "GEX api auth username")
	flag.String("gex-basic-auth-password", "", "GEX api auth password")
	flag.String("gex-url", "", "URL for sending an HTTP POST request to GEX")

	flag.String("dod-ca-package", "", "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
	flag.String("move-mil-dod-ca-cert", "", "The DoD CA certificate used to sign the move.mil TLS certificate.")
	flag.String("move-mil-dod-tls-cert", "", "The DoD-signed TLS certificate for various move.mil services.")
	flag.String("move-mil-dod-tls-key", "", "The private key for the DoD-signed TLS certificate for various move.mil services.")

	flag.String("edi", "", "The filepath to an edi file to send to GEX")
	flag.String("transaction-name", "test", "The required name sent in the url of the gex api request")
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, err := logging.Config("development", true)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	ediFile := v.GetString("edi")
	if ediFile == "" {
		log.Fatal("Usage: go run cmd/send_to_gex/main.go  --edi <edi filepath> --transaction-name <name>")
	}

	file, err := os.Open(ediFile) // #nosec
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	edi, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	ediString := string(edi[:])
	// make sure edi ends in new line
	ediString = strings.TrimSpace(ediString) + "\n"

	fmt.Println(ediString)

	certificates, rootCAs, err := initDODCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		log.Fatal("Error in getting tls certs", err)
	}
	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs}
	url := v.GetString("gex-url")
	if len(url) == 0 {
		log.Fatal("Not sending to GEX because no URL set. Set GEX_URL in your envrc.local.")
	}

	resp, err := invoice.NewGexSenderHTTP(
		url,
		true,
		tlsConfig,
		v.GetString("gex-basic-auth-username"),
		v.GetString("gex-basic-auth-password"),
	).SendToGex(ediString, v.GetString("transaction-name"))
	if resp == nil || err != nil {
		log.Fatal("Gex Sender had no response", err)
	}

	fmt.Println("Sending to GEX. . .")
	fmt.Printf("status code: %v, error: %v \n", resp.StatusCode, err)
}

//TODO: Infra will work to refactor and reduce duplication (also found in cmd/milmove/main.go)
func initDODCertificates(v *viper.Viper, logger *zap.Logger) ([]tls.Certificate, *x509.CertPool, error) {
	tlsCertString := v.GetString(cli.MoveMilDoDTLSCertFlag)
	tlsCerts := cli.ParseCertificates(tlsCertString)
	if len(tlsCerts) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", cli.MoveMilDoDTLSCertFlag)
	}
	if len(tlsCerts) > 1 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s has too many certificate PEM blocks", cli.MoveMilDoDTLSCertFlag)
	}

	logger.Info(fmt.Sprintf("certificate chain from %s parsed", cli.MoveMilDoDTLSCertFlag), zap.Any("count", len(tlsCerts)))
	fmt.Println(len(tlsCerts))

	caCertString := v.GetString(cli.MoveMilDoDCACertFlag)
	caCerts := cli.ParseCertificates(caCertString)
	if len(caCerts) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", cli.MoveMilDoDTLSCertFlag)
	}
	logger.Info(fmt.Sprintf("certificate chain from %s parsed", cli.MoveMilDoDCACertFlag), zap.Any("count", len(caCerts)))

	//Append move.mil cert with intermediate CA to create a validate certificate chain
	cert := strings.Join(append(append(make([]string, 0), tlsCerts...), caCerts...), "\n")

	key := v.GetString(cli.MoveMilDoDTLSKeyFlag)
	keyPair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, "failed to parse DOD x509 keypair for server")
	}

	logger.Info("DOD keypair", zap.Any("certificates", len(keyPair.Certificate)))

	pathToPackage := v.GetString(cli.DoDCAPackageFlag)
	pkcs7Package, err := ioutil.ReadFile(pathToPackage) // #nosec
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, fmt.Sprintf("%s is invalid", cli.DoDCAPackageFlag))
	}

	if len(pkcs7Package) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&cli.ErrInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is an empty file", cli.DoDCAPackageFlag))
	}

	dodCACertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		return make([]tls.Certificate, 0), dodCACertPool, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	return []tls.Certificate{keyPair}, dodCACertPool, nil

}

//TODO: Infra will refactor to reduce duplication
type errInvalidPKCS7 struct {
	Path string
}

//TODO: Infra will refactor to reduce duplication
func (e *errInvalidPKCS7) Error() string {
	return fmt.Sprintf("invalid DER encoded PKCS7 package: %s", e.Path)
}
