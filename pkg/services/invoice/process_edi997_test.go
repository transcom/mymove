package invoice

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ProcessEDI997Suite struct {
	*testingsuite.PopTestSuite
}

func TestProcessEDI997Suite(t *testing.T) {
	ts := &ProcessEDI997Suite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ProcessEDI997Suite) TestParsingEDI997() {
	edi997Processor := NewEDI997Processor()
	suite.Run("successfully processes a valid EDI997", func() {
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
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001251,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.NoError(err)
	})

	suite.Run("successfully processes a valid EDI997 with minimum AK4 values", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000999*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001
AK3*ab*123
AK4*3**1
AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001251,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.NoError(err)
	})

	suite.Run("successfully processes a valid EDI997 with 4 AK4 values", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000999*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001
AK3*ab*123
AK4*1**7*YE
AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001251,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.NoError(err)
	})

	suite.Run("throw error when parsing an EDI997 when an EDI997 is expected", func() {
		sample997EDIString := `
		ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
		GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
		ST*997*000000001
		BGN*11*1126-9404*20210217
		OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
		TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
		SE*5*000000001
		GE*1*1
		IEA*1*000000995
`
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.Contains(err.Error(), "unable to parse EDI997")
	})

	suite.Run("successfully updates a payment request status after processing a valid EDI997", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001252
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000995
	`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001252,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.NoError(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusTppsReceived, updatedPR.Status)
	})

	suite.Run("can handle 997 and 858 with same ICN", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001253
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000995
	`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 995,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001253,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.NoError(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.FatalNoError(err)
		suite.Equal(models.PaymentRequestStatusTppsReceived, updatedPR.Status)
	})

	suite.Run("does not error out if edi with same icn is processed for the same payment request", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001254
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000995
	`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 995,
					EDIType:                  models.EDIType997,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001254,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.NoError(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.FatalNoError(err)
		suite.Equal(models.PaymentRequestStatusTppsReceived, updatedPR.Status)
	})

	suite.Run("doesn't update a payment request status after processing an invalid EDI997", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000999*0*T*|
GS*SI*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001255
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
	`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 22,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.Error(err)
		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusSentToGex, updatedPR.Status)
	})

	suite.Run("throw an error when edi997 is missing a transaction set", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000999*0*T*|
GS*SI*8004171844*MILMOVE*20210217*152945*220001*X*004010
GE*1*220001
IEA*1*000000022
`
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.Contains(err.Error(), "Validation error(s) detected with the EDI997. EDI Errors could not be saved")
	})

	suite.Run("Return an error if payment request is not found with GCN", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000009*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001256
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
	`
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "unable to find PaymentRequest with GCN")
	})

	suite.Run("Should only find icn of edi type in AK2", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000009*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*997*0001
AK1*SI*100001257
AK2*858*0001

AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
	`
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "unable to find PaymentRequest with GCN")
	})
}

func (suite *ProcessEDI997Suite) TestValidatingEDI997() {
	edi997Processor := NewEDI997Processor()

	suite.Run("fails when there are validation errors on EDI997 fields", func() {
		sample997EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *210217*1530*U*00401*2000000000*8*A*|
GS*BS*MILMOVE*8004171844*20190903*1617*2000000000*X*004010
ST*997*0001
AK1*FA*100001251
AK2*858*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 89
AK3*ab*124
AK4*1*2*3*4*MM*bad data goes here 100
AK5*Q
AK9*P*10*1*1
SE*6*0001
ST*997*0002
AK1*BA*100001251
AK2*858*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 90
AK5*B
AK9*P*10*1*1
SE*6*0002
GE*1*220001
GS*FA*MILMOVE*8004171844*20190903*1617*22000000001*X*004010
ST*997*0001
AK1*VV*100001251
AK2*858*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 93
AK5*C
AK9*E*11*1*1
SE*6*0001
GE*1*220001
IEA*1*000000995
`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001251,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi997Processor.ProcessFile(suite.AppContextForTest(), "", sample997EDIString)
		suite.Error(err, "fail to process 997")
		errString := err.Error()
		actualErrors := strings.Split(errString, "\n")
		testData := []struct {
			TestName         string
			ExpectedErrorMsg string
		}{
			{TestName: "Invalid ICN", ExpectedErrorMsg: "'InterchangeControlNumber' failed on the 'max' tag"},
			{TestName: "Invalid AcknowledgementRequested", ExpectedErrorMsg: "'AcknowledgementRequested' failed on the 'oneof' tag"},
			{TestName: "Invalid UsageIndicator", ExpectedErrorMsg: "'UsageIndicator' failed on the 'oneof' tag"},
			{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'oneof' tag"},
			{TestName: "Invalid GroupControlNumber", ExpectedErrorMsg: "'GroupControlNumber' failed on the 'max' tag"},
			{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
			{TestName: "Invalid GroupControlNumber", ExpectedErrorMsg: "'GroupControlNumber' failed on the 'max' tag"},
			{TestName: "Second AK5 Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "'FunctionalIdentifierCode' failed on the 'eq' tag"},
		}

		for i, data := range testData {
			suite.Run(data.TestName, func() {
				suite.Contains(actualErrors[i], data.ExpectedErrorMsg)
			})
		}
	})
}
