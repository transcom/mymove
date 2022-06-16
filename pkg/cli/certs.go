package cli

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

// ErrInvalidPKCS7 is an Invalid PKCS7 error
type ErrInvalidPKCS7 struct {
	Path string
}

// Error is the error method
func (e *ErrInvalidPKCS7) Error() string {
	return fmt.Sprintf("invalid DER encoded PKCS7 package: %s", e.Path)
}

// InitCertFlags initializes the Certificate Flags
func InitCertFlags(flag *pflag.FlagSet) {
	flag.String(DevlocalCAFlag, "", "Path to PEM-encoded devlocal CA certificate, enabled in development and test builds")
	flag.StringSlice(DoDCAPackageFlag, []string{}, "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
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
		return errors.Wrap(&ErrInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is missing", pathToPackage))
	}

	return nil
}

// ParseCertificates takes a certificate and parses it into an slice of individual certificates
func ParseCertificates(str string) []string {

	certFormat := "-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----"

	// https://tools.ietf.org/html/rfc7468#section-2
	//	- https://stackoverflow.com/questions/20173472/does-go-regexps-any-charcter-match-newline
	re := regexp.MustCompile(`(?s)([-]{5}BEGIN CERTIFICATE[-]{5})(\s*)(.+?)(\s*)([-]{5}END CERTIFICATE[-]{5})`)
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
