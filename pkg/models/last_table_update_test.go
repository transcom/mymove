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
   <getTableResponseElement xmlns="http://ReturnTablePackage/">
	  <output>
		 <TRDM>
			<status>
			   <rowCount>28740</rowCount>
			   <statusCode>%v</statusCode>
			   <dateTime>2020-01-27T19:12:25.326Z</dateTime>
			</status>
		 </TRDM>
	  </output>
	  <attachment>
		 <xop:Include href="cid:fefe5d81-468c-4639-a543-e758a3cbceea-2@ReturnTablePackage" xmlns:xop="http://www.w3.org/2004/08/xop/include"/>
	  </attachment>
   </getTableResponseElement>
</soap:Body>
</soap:Envelope> `

type TRDMTestSuite struct {
	*testingsuite.PopTestSuite
}

const (
	physicalName  = "fakePhysicalName"
	returnContent = false
)

// TODO: Replace tableName and all references
func soapResponseForGetLastTableUpdate(statusCode string) *gosoap.Response {
	return &gosoap.Response{
		Body: []byte(fmt.Sprintf(getLastTableUpdateTemplate, statusCode)),
	}
}

func (suite *TRDMTestSuite) TestTRDMGetLastTableUpdateFake() {

	tests := []struct {
		name          string
		statusCode    string
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
			).Return(soapResponseForGetLastTableUpdate(test.statusCode), soapError)

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
