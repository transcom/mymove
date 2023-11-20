package factory

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
)

func Generatex509CertAndSecret() (*x509.Certificate, *rsa.PrivateKey, error) {
	// Create an x509 template to be used for creating a new certificate for testing purposes
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1),
	}

	// Generate private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Generate public certificate
	certByte, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}
	certificate, err := x509.ParseCertificate(certByte)
	if err != nil {
		return nil, nil, err
	}

	return certificate, key, nil
}
