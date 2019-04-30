package cli

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/server"

	"github.com/spf13/pflag"
)

const (
	// DevlocalCAFlag is the Devlocal CA Flag
	DevlocalCAFlag string = "devlocal-ca"
	// DoDCAPackageFlag is the DoD CA Package Flag
	DoDCAPackageFlag string = "dod-ca-package"
	// MoveMilDoDCACertFlag is the Move.mil DoD CA Cert Flag
	MoveMilDoDCACertFlag string = "move-mil-dod-ca-cert"
	// MoveMilDoDTLSCertFlag is the Move.mil DoD TLS Cert Flag
	MoveMilDoDTLSCertFlag string = "move-mil-dod-tls-cert"
	// MoveMilDoDTLSKeyFlag is the Move.mil DoD TLS Key Flag
	MoveMilDoDTLSKeyFlag string = "move-mil-dod-tls-key"
)

type errInvalidPKCS7 struct {
	Path string
}

func (e *errInvalidPKCS7) Error() string {
	return fmt.Sprintf("invalid DER encoded PKCS7 package: %s", e.Path)
}

// InitCertFlags initializes the Certificate Flags
func InitCertFlags(flag *pflag.FlagSet) {
	flag.String(DevlocalCAFlag, "", "Path to PEM-encoded devlocal CA certificate, enabled in development and test builds")
	flag.String(DoDCAPackageFlag, "", "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
	flag.String(MoveMilDoDCACertFlag, "", "The DoD CA certificate used to sign the move.mil TLS certificate.")
	flag.String(MoveMilDoDTLSCertFlag, "", "The DoD-signed TLS certificate for various move.mil services.")
	flag.String(MoveMilDoDTLSKeyFlag, "", "The private key for the DoD-signed TLS certificate for various move.mil services.")
}

// CheckCert validates Cert command line flags
func CheckCert(v *viper.Viper) error {

	dbEnv := v.GetString(DbEnvFlag)
	isDevOrTest := dbEnv == "development" || dbEnv == "test"
	if devlocalCAPath := v.GetString(DevlocalCAFlag); isDevOrTest && devlocalCAPath == "" {
		return errors.Errorf("No devlocal CA path defined")
	}

	tlsCertString := v.GetString(MoveMilDoDTLSCertFlag)
	if len(tlsCertString) == 0 {
		return errors.Errorf("%s is missing", MoveMilDoDTLSCertFlag)
	}

	caCertString := v.GetString(MoveMilDoDCACertFlag)
	if len(caCertString) == 0 {
		return errors.Errorf("%s is missing", MoveMilDoDCACertFlag)
	}

	key := v.GetString(MoveMilDoDTLSKeyFlag)
	if len(key) == 0 {
		return errors.Errorf("%s is missing", MoveMilDoDTLSKeyFlag)
	}
	pathToPackage := v.GetString(DoDCAPackageFlag)
	if len(pathToPackage) == 0 {
		return errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is missing", DoDCAPackageFlag))
	}

	return nil
}

// InitDoDCertificates initializes the DoD Certificates
func InitDoDCertificates(v *viper.Viper, logger Logger) ([]tls.Certificate, *x509.CertPool, error) {

	tlsCertString := v.GetString(MoveMilDoDTLSCertFlag)
	tlsCerts := ParseCertificates(tlsCertString)
	if len(tlsCerts) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", MoveMilDoDTLSCertFlag)
	}
	if len(tlsCerts) > 1 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s has too many certificate PEM blocks", MoveMilDoDTLSCertFlag)
	}

	logger.Info(fmt.Sprintf("certitficate chain from %s parsed", MoveMilDoDTLSCertFlag), zap.Any("count", len(tlsCerts)))

	caCertString := v.GetString(MoveMilDoDCACertFlag)
	caCerts := ParseCertificates(caCertString)
	if len(caCerts) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", MoveMilDoDTLSCertFlag)
	}

	logger.Info(fmt.Sprintf("certitficate chain from %s parsed", MoveMilDoDCACertFlag), zap.Any("count", len(caCerts)))

	//Append move.mil cert with intermediate CA to create a validate certificate chain
	cert := strings.Join(append(append(make([]string, 0), tlsCerts...), caCerts...), "\n")

	key := v.GetString(MoveMilDoDTLSKeyFlag)
	keyPair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, "failed to parse DOD x509 keypair for server")
	}

	logger.Info("DOD keypair", zap.Any("certificates", len(keyPair.Certificate)))

	pathToPackage := v.GetString(DoDCAPackageFlag)
	pkcs7Package, err := ioutil.ReadFile(pathToPackage) // #nosec
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, fmt.Sprintf("%s is invalid", DoDCAPackageFlag))
	}

	if len(pkcs7Package) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is an empty file", DoDCAPackageFlag))
	}

	dodCACertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		return make([]tls.Certificate, 0), dodCACertPool, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	return []tls.Certificate{keyPair}, dodCACertPool, nil

}

// ParseCertificates takes a certificate and parses it into an slice of individual certificates
func ParseCertificates(str string) []string {

	certFormat := "-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----"

	// https://tools.ietf.org/html/rfc7468#section-2
	//	- https://stackoverflow.com/questions/20173472/does-go-regexps-any-charcter-match-newline
	re := regexp.MustCompile("(?s)([-]{5}BEGIN CERTIFICATE[-]{5})(\\s*)(.+?)(\\s*)([-]{5}END CERTIFICATE[-]{5})")
	matches := re.FindAllStringSubmatch(str, -1)

	certs := make([]string, 0, len(matches))
	for _, m := range matches {
		// each match will include a slice of strings starting with
		// (0) the full match, then
		// (1) "-----BEGIN CERTIFICATE-----",
		// (2) whitespace if any,
		// (3) base64-encoded certificate data,
		// (4) whitespace if any, and then
		// (5) -----END CERTIFICATE-----
		certs = append(certs, fmt.Sprintf(certFormat, m[3]))
	}
	return certs
}
