package trdm_test

import (
	"fmt"

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
	bodyXML := fmt.Sprintf(`<soap:Body wsu:Id="%v"
xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
<ret:getLastTableUpdateRequestElement>
	<ret:physicalName>TRNSPRTN_ACNT</ret:physicalName>
</ret:getLastTableUpdateRequestElement>
</soap:Body>`, bodyID)
	_, err = trdm.GenerateSignedHeader(cert, key, bodyID, []byte(bodyXML))
	suite.NoError(err)
}
