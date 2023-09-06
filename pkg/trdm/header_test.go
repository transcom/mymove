package trdm_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/trdm"
)

func (suite *TRDMSuite) TestGenerateSignedHeader() {
	// Generate public cert, private key, and error
	cert, key, err := factory.Generatex509CertAndSecret()
	suite.NoError(err)
	// Pass certificate and key for header signing
	_, err = trdm.GenerateSignedHeader(cert, key)
	suite.NoError(err)
}

// TODO: Additional test to parse the header itself
