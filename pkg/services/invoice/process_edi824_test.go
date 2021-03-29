package invoice

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ProcessEDI824Suite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *ProcessEDI824Suite) SetupTest() {
	errTruncateAll := suite.TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
}

func TestProcessEDI824Suite(t *testing.T) {
	ts := &ProcessEDI824Suite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ProcessEDI824Suite) TestParsingEDI824() {
	edi824Processor := NewEDI824Processor(suite.DB(), suite.logger)
	suite.T().Run("successfully proccesses a valid EDI824", func(t *testing.T) {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000995
`
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakePaymentRequestToInterchangeControlNumber(suite.DB(), testdatagen.Assertions{
			PaymentRequestToInterchangeControlNumber: models.PaymentRequestToInterchangeControlNumber{
				PaymentRequestID:         paymentRequest.ID,
				InterchangeControlNumber: 995,
				PaymentRequest:           paymentRequest,
			},
		})
		err := edi824Processor.ProcessFile("", sample824EDIString)
		suite.NoError(err)
	})

	suite.T().Run("successfully updates a payment request status after processing a valid EDI824", func(t *testing.T) {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000996*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000996
`
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakePaymentRequestToInterchangeControlNumber(suite.DB(), testdatagen.Assertions{
			PaymentRequestToInterchangeControlNumber: models.PaymentRequestToInterchangeControlNumber{
				PaymentRequestID:         paymentRequest.ID,
				InterchangeControlNumber: 996,
				PaymentRequest:           paymentRequest,
			},
		})
		err := edi824Processor.ProcessFile("", sample824EDIString)
		suite.NoError(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusReceivedByGex, updatedPR.Status)
	})

	suite.T().Run("doesn't update a payment request status after processing an invalid EDI824", func(t *testing.T) {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*0000005*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*00000005
`
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakePaymentRequestToInterchangeControlNumber(suite.DB(), testdatagen.Assertions{
			PaymentRequestToInterchangeControlNumber: models.PaymentRequestToInterchangeControlNumber{
				PaymentRequestID:         paymentRequest.ID,
				InterchangeControlNumber: 22,
				PaymentRequest:           paymentRequest,
			},
		})
		edi824Processor.ProcessFile("", sample824EDIString)

		var updatedPR models.PaymentRequest
		err := suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusPending, updatedPR.Status)
	})

	suite.T().Run("Return an error if payment request is not found with ICN", func(t *testing.T) {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000997*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000997
`
		err := edi824Processor.ProcessFile("", sample824EDIString)
		suite.Error(err, "fail to process 824")
		suite.Contains(err.Error(), "unable to find payment request")
	})
}

func (suite *ProcessEDI824Suite) TestValidatingEDI824() {
	edi824Processor := NewEDI824Processor(suite.DB(), suite.logger)

	suite.T().Run("fails when there are validation errors on EDI header fields", func(t *testing.T) {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *210217*1530*U*00401*2000000000*8*A*|
GS*SA*MILMOVE*8004171844*20190903*1617*2000000000*X*004010
ST*824*000000001
BGN*19**20211313
OTI*VA*MM**X*X*20211311**-1*AB
OTI*AV*ER**X*X*20211311**20000000*CD
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*007*Missing Data
SE*5*000000001
ST*824*000000002
BGN*11*1126-9404*20210217
OTI*VA*MM**X*X*20211311**-1*AB
OTI*AV*ER**X*X*20211311**20000000*CD
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*007*Missing Data
SE*5*000000002
GE*2*1
GS*AG*8004171844*MILMOVE*20210217*1544*2*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0002
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*007*Missing Data
TED*812*Missing Transaction Reference or Trace Number
TED*PPD*Previously Paid
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*2*2
IEA*1*000000001
`

		err := edi824Processor.ProcessFile("", sample824EDIString)
		suite.Error(err, "fail to process 824")
		errString := err.Error()
		actualErrors := strings.Split(errString, "\n")
		suite.Equal(actualErrors, errString)
		// testData := []struct {
		// 	TestName         string
		// 	ExpectedErrorMsg string
		// }{
		// 	{TestName: "Invalid ICN causes missing PR", ExpectedErrorMsg: "unable to find payment request"},
		// 	{TestName: "Invalid ICN", ExpectedErrorMsg: "'InterchangeControlNumber' failed on the 'max' tag"},
		// 	{TestName: "Invalid AcknowledgementRequested", ExpectedErrorMsg: "'AcknowledgementRequested' failed on the 'oneof' tag"},
		// 	{TestName: "Invalid UsageIndicator", ExpectedErrorMsg: "'UsageIndicator' failed on the 'oneof' tag"},
		// 	{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
		// 	{TestName: "Invalid GroupControlNumber", ExpectedErrorMsg: "'GroupControlNumber' failed on the 'max' tag"},
		// 	{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
		// 	{TestName: "Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "'TransactionSetIdentifierCode' failed on the 'eq' tag"},
		// 	{TestName: "Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
		// 	{TestName: "Second AK2 Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "'TransactionSetIdentifierCode' failed on the 'eq' tag"},
		// 	{TestName: "Second AK1 failure for Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
		// 	{TestName: "Invalid GroupControlNumber", ExpectedErrorMsg: "'GroupControlNumber' failed on the 'max' tag"},
		// 	{TestName: "Second AK5 Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
		// 	{TestName: "Third (in second functionalGroupEnvelope) AK2 Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "'TransactionSetIdentifierCode' failed on the 'eq' tag"},
		// }

		// for i, data := range testData {
		// 	suite.T().Run(data.TestName, func(t *testing.T) {
		// 		suite.Contains(actualErrors[i], data.ExpectedErrorMsg)
		// 	})
		// }
	})
}
