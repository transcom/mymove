package invoice

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ProcessEDI997Suite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *ProcessEDI997Suite) SetupTest() {
	errTruncateAll := suite.TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
}

func TestProcessEDI997Suite(t *testing.T) {
	ts := &ProcessEDI997Suite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ProcessEDI997Suite) TestParsingEDI997() {
	edi997Processor := NewEDI997Processor(suite.DB(), suite.logger)
	suite.T().Run("successfully proccesses a valid EDI997", func(t *testing.T) {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000999*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakePaymentRequestToInterchangeControlNumber(suite.DB(), testdatagen.Assertions{
			PaymentRequestToInterchangeControlNumber: models.PaymentRequestToInterchangeControlNumber{
				PaymentRequestID:         paymentRequest.ID,
				InterchangeControlNumber: 999,
				PaymentRequest:           paymentRequest,
			},
		})
		err := edi997Processor.ProcessFile("", sample997EDIString)
		suite.NoError(err)
	})

	suite.T().Run("successfully updates a payment request status after processing a valid EDI997", func(t *testing.T) {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
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
		err := edi997Processor.ProcessFile("", sample997EDIString)
		suite.NoError(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusReceivedByGex, updatedPR.Status)
	})

	suite.T().Run("doesn't update a payment request status after processing an invalid EDI997", func(t *testing.T) {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000999*0*T*|
GS*SI*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
	`
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakePaymentRequestToInterchangeControlNumber(suite.DB(), testdatagen.Assertions{
			PaymentRequestToInterchangeControlNumber: models.PaymentRequestToInterchangeControlNumber{
				PaymentRequestID:         paymentRequest.ID,
				InterchangeControlNumber: 22,
				PaymentRequest:           paymentRequest,
			},
		})
		edi997Processor.ProcessFile("", sample997EDIString)

		var updatedPR models.PaymentRequest
		err := suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusPending, updatedPR.Status)
	})

	suite.T().Run("Return an error if payment request is not found with ICN", func(t *testing.T) {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000009*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
	`
		err := edi997Processor.ProcessFile("", sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "unable to find payment request")
	})
}

func (suite *ProcessEDI997Suite) TestValidatingEDIHeader() {
	edi997Processor := NewEDI997Processor(suite.DB(), suite.logger)

	suite.T().Run("fails when there are validation errors on EDI header fields", func(t *testing.T) {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *210217*1530*U*00401*2000000000*8*A*|
GS*FA*MILMOVE*8004171844*20190903*1617*2000000000*X*004010
ST*997*0001
AK1*FA*100001251
AK2*909*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 89
AK3*ab*124
AK4*1*2*3*4*MM*bad data goes here 100
AK5*Q
AK9*P*10*1*1
SE*6*0001
ST*997*0002
AK1*FA*100001251
AK2*900*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 90
AK5*B
AK9*P*10*1*1
SE*6*0002
GE*1*220001
GS*FA*MILMOVE*8004171844*20190903*1617*22000000001*X*004010
ST*997*0001
AK1*VV*100001251
AK2*123*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 93
AK5*C
AK9*E*11*1*1
SE*6*0001
GE*1*220001
IEA*1*000000995
`

		err := edi997Processor.ProcessFile("", sample997EDIString)
		suite.Error(err, "fail to process 997")
		errString := err.Error()
		actualErrors := strings.Split(errString, "\n")
		fmt.Printf("%+v", actualErrors[11])
		testData := []struct {
			TestName         string
			ExpectedErrorMsg string
		}{
			{TestName: "Invalid ICN causes missing PR", ExpectedErrorMsg: "unable to find payment request"},
			{TestName: "Invalid ICN", ExpectedErrorMsg: "'InterchangeControlNumber' failed on the 'max' tag"},
			{TestName: "Invalid AcknowledgementRequested", ExpectedErrorMsg: "'AcknowledgementRequested' failed on the 'oneof' tag"},
			{TestName: "Invalid UsageIndicator", ExpectedErrorMsg: "'UsageIndicator' failed on the 'oneof' tag"},
			{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Invalid GroupControlNumber", ExpectedErrorMsg: "'GroupControlNumber' failed on the 'max' tag"},
			{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "'TransactionSetIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Second AK2 Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "'TransactionSetIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Second AK1 failure for Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Invalid GroupControlNumber", ExpectedErrorMsg: "'GroupControlNumber' failed on the 'max' tag"},
			{TestName: "Second AK5 Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Third (in second functionalGroupEnvelope) AK2 Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "'TransactionSetIdentifierCode' failed on the 'eq' tag"},
		}

		for i, data := range testData {
			suite.T().Run(data.TestName, func(t *testing.T) {
				suite.Contains(actualErrors[i], data.ExpectedErrorMsg)
			})
		}
	})
}
