package trdm_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/trdm"
)

func (suite *TRDMSuite) TestGenerateSignedHeader() {
	// Generate public cert, private key, and error
	cert, key, err := factory.Generatex509CertAndSecret()
	//cert, key, err := getRealCertAndKey()
	suite.NoError(err)
	// Pass certificate and key for header signing
	bodyID, err := trdm.GenerateSOAPURIWithPrefix("#id")
	suite.NoError(err)
	headerByte, err := trdm.GenerateSignedHeader(cert, key, bodyID)
	suite.NoError(err)
	// ! Readme is currently for debugging purposes
	readme := string(headerByte)
	suite.NotEmpty(readme)
}

/*
func getRealCertAndKey() (*x509.Certificate, *rsa.PrivateKey, error) {
	// ! Real public key from .envrc TODO: Remove
	publicPem, _ := pem.Decode([]byte(tlsPublicKey))
	cert, err := x509.ParseCertificate(publicPem.Bytes)
	if err != nil {
		return nil, nil, err
	}
	// ! Real private key from .envrc TODO: Remove
	privatePem, _ := pem.Decode([]byte(tlsPrivateKey))
	key, err := x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, err
}
*/
// TODO: Additional test to parse the header itself
