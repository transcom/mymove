package trdm_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/trdm"
	"github.com/transcom/mymove/pkg/trdm/trdmmocks"
)

const getLastTableUpdateTemplate = `
   <getLastTableUpdateResponseElement xmlns="http://ReturnTablePackage/">
	  <lastUpdate>%v</lastUpdate>
	  <status>
		 <statusCode>%v</statusCode>
		 <dateTime>2020-01-27T20:18:34.226Z</dateTime>
	  </status>
   </getLastTableUpdateResponseElement>
`

const (
	physicalName = "fakePhysicalName"
)

func soapResponseForGetLastTableUpdate(lastUpdate string, statusCode string) *gosoap.Response {
	return &gosoap.Response{
		Body: []byte(fmt.Sprintf(getLastTableUpdateTemplate, lastUpdate, statusCode)),
	}
}

func (suite *TRDMSuite) TestTRDMGetLastTableUpdateFake() {
	tests := []struct {
		name          string
		lastUpdate    string
		statusCode    string
		responseError bool
		shouldError   bool
	}{
		{"No update", time.Now().Format(time.RFC3339), "Successful", false, false},
		{"Should not fetch update", time.Now().Format(time.RFC3339), "Failure", false, false},
		{"There is an update", time.Now().Add(-72 * time.Hour).Format(time.RFC3339), "Successful", false, false},
	}
	for _, test := range tests {
		suite.Run("fake call to TRDM: "+test.name, func() {
			var soapError error
			if test.responseError {
				soapError = errors.New("some error")
			}

			testSoapClient := &trdmmocks.SoapCaller{}
			testSoapClient.On("Call",
				mock.Anything,
				mock.Anything,
			).Return(soapResponseForGetLastTableUpdate(test.lastUpdate, test.statusCode), soapError)
			privatekey, keyErr := rsa.GenerateKey(rand.Reader, 2048)
			if keyErr != nil {
				suite.Error(keyErr)
			}
			// ! Real public key TODO: Remove
			publicPem, rest := pem.Decode([]byte(tlsPublicKey))
			suite.Empty(rest)
			certificate, err := x509.ParseCertificate([]byte(publicPem.Bytes))
			suite.NoError(err)
			lastTableUpdate := trdm.NewTRDMGetLastTableUpdate(physicalName, certificate, privatekey, testSoapClient)
			err = lastTableUpdate.GetLastTableUpdate(suite.AppContextForTest(), physicalName)
			suite.NoError(err)
		})
	}
}

func (suite *TRDMSuite) TestFetchLOARecordsByTime() {
	// Get initial TAC codes count
	initialCodes, err := trdm.FetchLOARecordsByTime(suite.AppContextForTest(), time.Now())
	intialLoaCodesLength := len(initialCodes)
	suite.NoError(err)

	// Creates a test TAC code record in the DB
	factory.BuildDefaultLineOfAccounting(suite.DB())

	// Fetch All TAC Records
	// !! A second time.Now() statement is intentional based on the SQL query.
	codes, err := trdm.FetchLOARecordsByTime(suite.AppContextForTest(), time.Now())
	suite.NoError(err)
	// Compare new TAC Code count to initial count
	finalCodesLength := len(codes)

	suite.NoError(err)
	suite.NotEqual(finalCodesLength, intialLoaCodesLength)
}

func (suite *TRDMSuite) TestFetchTACRecordsByTime() {
	// Get initial TAC codes count
	initialCodes, err := trdm.FetchTACRecordsByTime(suite.AppContextForTest(), time.Now())
	initialTacCodeLength := len(initialCodes)
	suite.NoError(err)

	// Creates a test TAC code record in the DB
	factory.BuildDefaultTransportationAccountingCode(suite.DB())

	// Fetch All TAC Records
	codes, err := trdm.FetchTACRecordsByTime(suite.AppContextForTest(), time.Now())

	// Compare new TAC Code count to initial count
	finalCodesLength := len(codes)

	suite.NoError(err)
	suite.NotEqual(finalCodesLength, initialTacCodeLength)

}
