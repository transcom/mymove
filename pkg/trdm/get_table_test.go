package trdm_test

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/models/mocks"
	"github.com/transcom/mymove/pkg/trdm"
)

const getTableTemplate = `
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
            <xop:Include href="cid:fefe5d81-468c-4639-a543-e758a3cbceea-2@ReturnTablePackage"
			xmlns:xop="http://www.w3.org/2004/08/xop/include"/>
         </attachment>
      </getTableResponseElement>
`

func soapResponseForGetTable(statusCode string) *gosoap.Response {
	return &gosoap.Response{
		Body: []byte(fmt.Sprintf(getTableTemplate, statusCode)),
	}
}

func (suite *TRDMSuite) TestTRDMGetTableFake() {
	tests := []struct {
		name          string
		statusCode    string
		responseError bool
		shouldError   bool
	}{
		{"No update", "Successful", false, false},
		{"Should not fetch update", "Failure", false, false},
		{"There is an update", "Successful", false, false},
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
			).Return(soapResponseForGetTable(test.statusCode), soapError)

			getTable := trdm.NewGetTable(physicalName, testSoapClient)
			err := getTable.GetTable(suite.AppContextForTest(), physicalName)

			if err != nil {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
