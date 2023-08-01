package models_test

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/ghcmocks"
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
		returnContent bool
		responseError bool
		shouldError   bool
	}{
		{"", "AFCT", true, false, false},
		{"", "FakeTable", false, false, false},
		{"", "nano", false, false, true},
		{"", "vi", false, true, true},
		{"", "not Sure", false, false, false},
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

			lastTableUpdate := models.NewTRDMGetLastTableUpdate(physicalName, returnContent, nil) //! REPLACE nil with soapClient
			err := lastTableUpdate.GetLastTableUpdate(suite.AppContextForTest(), "ACFT", true)

			if err != nil {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
