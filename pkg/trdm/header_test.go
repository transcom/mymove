package trdm_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"

	"github.com/transcom/mymove/pkg/trdm"
)

func (suite *TRDMSuite) TestGenerateSignedHeader() {
	// Ceate an x509 template to be used for creating a new certificate. This template is lacking a lot of information, but provides just enough
	// to satisfy the test's needs
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1),
	}

	// Generate private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	suite.NoError(err)

	// Generate public certificate
	certByte, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	suite.NoError(err)
	certificate, err := x509.ParseCertificate(certByte)
	suite.NoError(err)

	// Pass certificate and key for header signing
	_, err = trdm.GenerateSignedHeader(certificate, key)
	suite.NoError(err)
}

// TODO: Additional test to parse the header itself
