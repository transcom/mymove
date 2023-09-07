package trdm_test

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/trdm"
	"github.com/transcom/mymove/pkg/trdm/trdmmocks"
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

func getTextFile(filePath string) ([]byte, error) {
	text, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return text, nil
}

func soapResponseForGetTable(statusCode string, payload []byte) *gosoap.Response {
	return &gosoap.Response{
		Body:    []byte(fmt.Sprintf(getTableTemplate, statusCode)),
		Payload: []byte(payload),
	}
}

func (suite *TRDMSuite) TestGetTableFake() {
	loaFile, err := getTextFile("../parser/loa/fixtures/Line Of Accounting.txt")
	suite.NoError(err)
	tacFile, err := getTextFile("../parser/tac/fixtures/Transportation Account.txt")
	suite.NoError(err)
	tests := []struct {
		name          string
		physicalName  string
		statusCode    string
		payload       []byte
		responseError bool
		shouldError   bool
	}{
		{"Update Line of Accounting", "LN_OF_ACCT", "Successful", loaFile, false, false},
		{"Should not fetch update", "fakeName", "Failure", nil, false, false},
		{"Update Transportation Accounting Codes", "TRNSPRTN_ACNT", "Successful", tacFile, false, false},
	}
	for _, test := range tests {
		suite.Run("fake call to TRDM: "+test.name, func() {
			var soapError error
			if test.responseError {
				soapError = errors.New("Error running range of GetTable tests")
				suite.Error(soapError)
			}

			testSoapClient := &trdmmocks.SoapCaller{}
			testSoapClient.On("Call",
				mock.Anything,
				mock.Anything,
			).Return(soapResponseForGetTable(test.statusCode, test.payload), soapError)
			cert, key, err := factory.Generatex509CertAndSecret()
			suite.NoError(err)
			getTable := trdm.NewGetTable(test.physicalName, cert, key, testSoapClient)
			err = getTable.GetTable(suite.AppContextForTest(), test.physicalName, time.Now())
			suite.NoError(err)
		})
	}
}
