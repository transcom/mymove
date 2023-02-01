package certs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/server"
)

// InitDoDCertificates initializes the DoD Certificates
func InitDoDCertificates(v *viper.Viper, logger *zap.Logger) ([]tls.Certificate, *x509.CertPool, error) {

	tlsCertString := v.GetString(cli.MoveMilDoDTLSCertFlag)
	tlsCerts := cli.ParseCertificates(tlsCertString)
	if len(tlsCerts) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", cli.MoveMilDoDTLSCertFlag)
	}
	if len(tlsCerts) > 1 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s has too many certificate PEM blocks", cli.MoveMilDoDTLSCertFlag)
	}

	logger.Info("TLS certificate parsed", zap.String("env", cli.MoveMilDoDTLSCertFlag), zap.Any("count", len(tlsCerts)))

	caCertString := v.GetString(cli.MoveMilDoDCACertFlag)
	caCerts := cli.ParseCertificates(caCertString)
	if len(caCerts) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", cli.MoveMilDoDTLSCertFlag)
	}

	logger.Info("CA certificate chain parsed", zap.String("env", cli.MoveMilDoDCACertFlag), zap.Any("count", len(caCerts)))

	//Append move.mil cert with intermediate CA to create a validate certificate chain
	cert := strings.Join(append(append(make([]string, 0), tlsCerts...), caCerts...), "\n")

	key := v.GetString(cli.MoveMilDoDTLSKeyFlag)
	keyPair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, "failed to parse DOD x509 keypair for server")
	}

	logger.Info("DOD keypair created", zap.Any("certificates", len(keyPair.Certificate)))

	pathToPackage := v.GetString(cli.DoDCAPackageFlag)
	pkcs7Package, err := os.ReadFile(filepath.Clean(pathToPackage))
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

	logger.Info("Cert Pool loaded", zap.Any("env", cli.DoDCAPackageFlag))

	return []tls.Certificate{keyPair}, dodCACertPool, nil

}
