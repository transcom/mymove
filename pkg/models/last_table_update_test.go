package models_test

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const getLastTableUpdateTemplate = `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
   <getLastTableUpdateResponseElement xmlns="http://ReturnTablePackage/">
	  <lastUpdate>%v</lastUpdate>
	  <status>
		 <statusCode>Successful</statusCode>
		 <dateTime>2020-01-27T20:18:34.226Z</dateTime>
	  </status>
   </getLastTableUpdateResponseElement>
</soap:Body>
</soap:Envelope> `

const (
	physicalName = "fakePhysicalName"
)

func soapResponseForGetLastTableUpdate(lastUpdate time.Time) *gosoap.Response {
	return &gosoap.Response{
		Body: []byte(fmt.Sprintf(getLastTableUpdateTemplate, lastUpdate)),
	}
}

func (suite *ModelSuite) TestTRDMGetLastTableUpdateFake() {
	tests := []struct {
		name          string
		lastUpdate    time.Time
		responseError bool
		shouldError   bool
	}{
		{"No update", time.Now(), false, false},
		{"Should error", time.Now(), true, true},
		{"There is an update", time.Now().Add(-72 * time.Hour), false, false},
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
			).Return(soapResponseForGetLastTableUpdate(test.lastUpdate), soapError)

			lastTableUpdate := models.NewTRDMGetLastTableUpdate(physicalName, testSoapClient)
			err := lastTableUpdate.GetLastTableUpdate(suite.AppContextForTest(), physicalName)

			if err != nil {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *ModelSuite) TestFetchAllTACRecords() {

	// Get initial TAC codes count
	initialCodes, err := models.FetchAllTACRecords(suite.DB())
	initialTacCodeLength := len(initialCodes)
	suite.NoError(err)

	// Creates a test TAC code record in the DB
	testdatagen.MakeDefaultTransportationAccountingCode(suite.DB())

	// Fetch All TAC Records
	codes, err := models.FetchAllTACRecords(suite.DB())

	// Compare new TAC Code count to initial count
	finalCodesLength := len(codes)

	suite.NoError(err)
	suite.Equal(finalCodesLength, initialTacCodeLength+1)
}
