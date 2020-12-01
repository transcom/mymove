package certs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
)

// InitDoDEntrustCertificates initializes the DoD Certificates
func InitDoDEntrustCertificates(v *viper.Viper, logger Logger) ([]tls.Certificate, *x509.CertPool, error) {

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

	entrustL1KCertString := v.GetString("entrust-l1k-cert")
	entrustL1KCert := cli.ParseCertificates(entrustL1KCertString)
	if len(entrustL1KCert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing entrust L1K cert", cli.EntrustL1KCertFlag)
	}

	entrustG2CertString := v.GetString("entrust-g2-cert")
	entrustG2Cert := cli.ParseCertificates(entrustG2CertString)
	if len(entrustG2Cert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing entrust G2 cert", cli.EntrustG2CertFlag)
	}

	entrustCerts := strings.Join(append(append(make([]string, 0), entrustL1KCert...), entrustG2Cert...), "\n")

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(entrustCerts))
	if !ok {
		panic("failed to parse root certificate")
	}

	return []tls.Certificate{keyPair}, certPool, nil

}
