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
	"go.mozilla.org/pkcs7"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
)

// addToCertPoolFromPkcs7Package reads the certificates in a DER-encoded PKCS7
// package and adds those Certificates to the x509.CertPool
func addToCertPoolFromPkcs7Package(certPool *x509.CertPool, pkcs7Package []byte) error {
	p7, err := pkcs7.Parse(pkcs7Package)
	if err != nil {
		return err
	}
	for _, cert := range p7.Certificates {
		certPool.AddCert(cert)
	}
	return nil
}

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

	logger.Info(fmt.Sprintf("certificate chain from %s parsed", cli.MoveMilDoDTLSCertFlag), zap.Any("count", len(tlsCerts)))

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
	pkcs7Package, err := os.ReadFile(filepath.Clean(pathToPackage))
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, fmt.Sprintf("%s is invalid", cli.DoDCAPackageFlag))
	}

	if len(pkcs7Package) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&cli.ErrInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is an empty file", cli.DoDCAPackageFlag))
	}

	dodCACertPool := x509.NewCertPool()
	err = addToCertPoolFromPkcs7Package(dodCACertPool, pkcs7Package)
	if err != nil {
		return make([]tls.Certificate, 0), dodCACertPool, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	return []tls.Certificate{keyPair}, dodCACertPool, nil

}

func InitMutualTLSClientCAPool(v *viper.Viper, logger *zap.Logger) (*x509.CertPool, error) {
	pathToPackage := v.GetString(cli.DoDCAPackageFlag)
	pkcs7Package, err := os.ReadFile(filepath.Clean(pathToPackage))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s is invalid", cli.DoDCAPackageFlag))
	}

	if len(pkcs7Package) == 0 {
		return nil, errors.Wrap(&cli.ErrInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is an empty file", cli.DoDCAPackageFlag))
	}

	mtlsClientCAPool := x509.NewCertPool()
	err = addToCertPoolFromPkcs7Package(mtlsClientCAPool, pkcs7Package)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	additionalPackages := v.GetStringSlice(cli.MutualTLSAdditionalCAPackage)
	for i := range additionalPackages {
		pathToPackage := additionalPackages[i]
		pkcs7Package, err := os.ReadFile(filepath.Clean(pathToPackage))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("%s has invalid file: %s", cli.MutualTLSAdditionalCAPackage, pathToPackage))
		}

		if len(pkcs7Package) == 0 {
			return nil, errors.Wrap(&cli.ErrInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s has an empty file: %s", cli.MutualTLSAdditionalCAPackage, pathToPackage))
		}
		err = addToCertPoolFromPkcs7Package(mtlsClientCAPool, pkcs7Package)
		if err != nil {
			return nil, err
		}
	}
	return mtlsClientCAPool, nil
}
