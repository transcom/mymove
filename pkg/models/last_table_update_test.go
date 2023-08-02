package models_test

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/ghcmocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
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

type TRDMTestSuite struct {
	*testingsuite.PopTestSuite
}

const (
	physicalName  = "fakePhysicalName"
	returnContent = false
)

// TODO: Replace lastUpdate and all references
func soapResponseForGetLastTableUpdate(lastUpdate string) *gosoap.Response {
	return &gosoap.Response{
		Body: []byte(fmt.Sprintf(getLastTableUpdateTemplate, lastUpdate)),
	}
}

func (suite *TRDMTestSuite) TestTRDMGetLastTableUpdateFake() {

	tests := []struct {
		name          string
		lastUpdate    string
		responseError bool
		shouldError   bool
	}{
		{"", "AFCT", false, false},
		{"", "FakeTable", false, false},
		{"", "nano", false, true},
		{"", "vi", true, true},
		{"", "not Sure", false, false},
	}
	for _, test := range tests {
		suite.Run("fake call to TRDM: "+test.name, func() {
			var soapError error
			if test.responseError {
				soapError = errors.New("some error")
			}

			testSoapClient := &ghcmocks.SoapCaller{}
			testSoapClient.On("Call",
				mock.Anything,
				mock.Anything,
			).Return(soapResponseForGetLastTableUpdate(test.lastUpdate), soapError)

			lastTableUpdate := models.NewTRDMGetLastTableUpdate(physicalName, nil) //! REPLACE nil with soapClient
			err := lastTableUpdate.GetLastTableUpdate(suite.AppContextForTest(), "ACFT", true)

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
