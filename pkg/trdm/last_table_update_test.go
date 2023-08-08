package trdm_test

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/models/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/trdm"
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
		{"No update", time.Now().String(), "Successful", false, false},
		{"Should not fetch update", time.Now().String(), "Failure", false, false},
		{"There is an update", time.Now().Add(-72 * time.Hour).String(), "Successful", false, false},
	}
	for _, test := range tests {
		suite.Run("fake call to TRDM: "+test.name, func() {
			var soapError error
			if test.responseError {
				soapError = errors.New("some error")
			}

			testSoapClient := &mocks.SoapCaller{}
			testSoapClient.On("Call",
				mock.Anything,
				mock.Anything,
			).Return(soapResponseForGetLastTableUpdate(test.lastUpdate, test.statusCode), soapError)

			lastTableUpdate := trdm.NewTRDMGetLastTableUpdate(physicalName, testSoapClient)
			err := lastTableUpdate.GetLastTableUpdate(suite.AppContextForTest(), physicalName)

			if err != nil {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *TRDMSuite) TestFetchAllTACRecords() {

	// Get initial TAC codes count
	initialCodes, err := trdm.FetchAllTACRecords(suite.AppContextForTest())
	initialTacCodeLength := len(initialCodes)
	suite.NoError(err)

	// Creates a test TAC code record in the DB
	testdatagen.MakeDefaultTranportationAccountingCode(suite.DB())

	// Fetch All TAC Records
	codes, err := trdm.FetchAllTACRecords(suite.AppContextForTest())

	// Compare new TAC Code count to initial count
	finalCodesLength := len(codes)

	suite.NoError(err)
	suite.Equal(finalCodesLength, initialTacCodeLength)
}
