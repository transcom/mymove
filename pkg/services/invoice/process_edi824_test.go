package invoice

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	ediResponse824 "github.com/transcom/mymove/pkg/edi/edi824"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ProcessEDI824Suite struct {
	*testingsuite.PopTestSuite
}

func TestProcessEDI824Suite(t *testing.T) {
	ts := &ProcessEDI824Suite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ProcessEDI824Suite) TestParsingEDI824() {
	mockNotificationSender := notifications.NewStubNotificationSender("")

	edi824Processor := NewEDI824Processor(mockNotificationSender)
	suite.Run("successfully processes a valid EDI824", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sample824EDIString := fmt.Sprintf(`
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*%s*20210217
OTI*TR*BM*%s*MILMOVE*8004171844*20210217**100001251*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000995
`, *paymentRequest.MoveTaskOrder.ReferenceID, *paymentRequest.MoveTaskOrder.ReferenceID)
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
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.NoError(err)
	})

	suite.Run("throw an error when edi824 is missing an OTI segment", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sample824EDIString := fmt.Sprintf(`
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*%s*20210217
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000995
`, paymentRequest.PaymentRequestNumber)
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.Contains(err.Error(), "Validation error(s) detected with the EDI824. EDI Errors could not be saved")
	})

	suite.Run("throw an error when edi824 is missing a transaction set", func() {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
GE*1*1
IEA*1*000000995
`
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.Contains(err.Error(), "Validation error(s) detected with the EDI824. EDI Errors could not be saved")
	})

	suite.Run("throw an error when a payment request cannot be found with the OTI.GroupControlNumber", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sample824EDIString := fmt.Sprintf(`
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*%s*20210217
OTI*TR*BM*%s*MILMOVE*8004171844*20210217**100001252*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000995
`, paymentRequest.PaymentRequestNumber, *paymentRequest.MoveTaskOrder.ReferenceID)
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.Contains(err.Error(), "unable to find PaymentRequest with GCN")
	})

	suite.Run("throw an error when a the BGN02 ref identification doesn't match the PaymentRequestNumber", func() {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000995*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001253*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000995
`
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
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
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.NotNil(err)
		suite.Contains(err.Error(), fmt.Sprintf("The BGN02 Reference Identification field: 1126-9404 doesn't match the MTO reference ID %s of the associated payment request", *paymentRequest.MoveTaskOrder.ReferenceID))
	})

	suite.Run("throw error when parsing an EDI997 when an EDI824 is expected", func() {
		sample824EDIString := `
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
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.Contains(err.Error(), "unable to parse EDI824")
	})

	suite.Run("successfully updates a payment request status after processing a valid EDI824", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sample824EDIString := fmt.Sprintf(`
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000996*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*%s*20210217
OTI*TR*BM*%s*MILMOVE*8004171844*20210217**100001255*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000996
`, *paymentRequest.MoveTaskOrder.ReferenceID, *paymentRequest.MoveTaskOrder.ReferenceID)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001255,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.NoError(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusEDIError, updatedPR.Status)
	})

	suite.Run("does not update a payment request status after processing a valid EDI824 if every segment is TA", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sample824EDIString := fmt.Sprintf(`
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000996*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*%s*20210217
OTI*TA*BM*%s*MILMOVE*8004171844*20210217**100001255*0001
SE*5*000000001
GE*1*1
IEA*1*000000996
`, *paymentRequest.MoveTaskOrder.ReferenceID, *paymentRequest.MoveTaskOrder.ReferenceID)
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001255,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.NoError(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusPending, updatedPR.Status)
	})
	suite.Run("doesn't update a payment request status after processing an invalid EDI824", func() {
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

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)

		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.NotNil(err)

		var updatedPR models.PaymentRequest
		err = suite.DB().Where("id = ?", paymentRequest.ID).First(&updatedPR)
		suite.NoError(err)
		suite.Equal(models.PaymentRequestStatusPending, updatedPR.Status)
	})

	suite.Run("Save TED errors to the database", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sample824EDIString := fmt.Sprintf(`
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*00000997*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*%s*20210217
OTI*TR*BM*%s*MILMOVE*8004171844*20210217**100001252*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*K*MISSING DATA
SE*5*000000001
GE*1*1
IEA*1*000000997
`, *paymentRequest.MoveTaskOrder.ReferenceID, *paymentRequest.MoveTaskOrder.ReferenceID)
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
		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.NoError(err)

		var ediErrors models.EdiErrors
		err = suite.DB().Where("payment_request_id = ?", paymentRequest.ID).All(&ediErrors)
		suite.NoError(err)

		suite.Equal(2, len(ediErrors))
		for i, ediError := range ediErrors {
			suite.Equal("K", *ediError.Code)
			suite.Equal(models.EDIType824, ediError.EDIType)
			if i == 0 {
				suite.Equal("DOCUMENT OWNER CANNOT BE DETERMINED", *ediError.Description)
			} else {
				suite.Equal("MISSING DATA", *ediError.Description)
			}
		}
	})
}

func (suite *ProcessEDI824Suite) TestValidatingEDI824() {
	mockSender := notifications.NewStubNotificationSender("")
	edi824Processor := NewEDI824Processor(mockSender)

	suite.Run("fails when there are validation errors on the EDI", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sample824EDIString := fmt.Sprintf(`
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *210217*1530*U*00401*2000000000*8*A*|
GS*SA*MILMOVE*8004171844*20190903*1617*2000000000*X*004010
ST*824*000000001
BGN*19*%s*20211313
OTI*VA*MM*%s*X*X*20211311**100001252*AB
TED*007*Missing Data
SE*5*000000001
GE*2*1
IEA*1*000000001
`, *paymentRequest.MoveTaskOrder.ReferenceID, *paymentRequest.MoveTaskOrder.ReferenceID)
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

		err := edi824Processor.ProcessFile(suite.AppContextForTest(), "", sample824EDIString)
		suite.Error(err, "fail to process 824")
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
			{TestName: "Invalid BGN.TransactionSetPurposeCode", ExpectedErrorMsg: "'TransactionSetPurposeCode' failed on the 'eq' tag"},
			{TestName: "Invalid BGN.Date", ExpectedErrorMsg: "'Date' failed on the 'datetime' tag"},
			{TestName: "Invalid OTIs[0].ApplicationAcknowledgementCode", ExpectedErrorMsg: "'ApplicationAcknowledgementCode' failed on the 'oneof' tag"},
			{TestName: "Invalid OTIs[0].ReferenceIdentificationQualifier", ExpectedErrorMsg: "'ReferenceIdentificationQualifier' failed on the 'oneof' tag"},
			{TestName: "Invalid OTIs[0].ApplicationSendersCode", ExpectedErrorMsg: "'ApplicationSendersCode' failed on the 'min' tag"},
			{TestName: "Invalid OTIs[0].ApplicationReceiversCode", ExpectedErrorMsg: "'ApplicationReceiversCode' failed on the 'min' tag"},
			{TestName: "Invalid OTIs[0].Date", ExpectedErrorMsg: "'Date' failed on the 'datetime' tag"},
			{TestName: "Invalid OTIs[0].TransactionSetControlNumber", ExpectedErrorMsg: "'TransactionSetControlNumber' failed on the 'min' tag"},
			{TestName: "Invalid GE.NumberOfTransactionSetsIncluded", ExpectedErrorMsg: "'NumberOfTransactionSetsIncluded' failed on the 'eq' tag"},
		}

		for i, data := range testData {
			suite.Run(data.TestName, func() {
				suite.Contains(actualErrors[i], data.ExpectedErrorMsg)
			})
		}
	})
}

func (suite *ProcessEDI824Suite) TestIdentifyingOTIsAndTEDs() {
	suite.Run("fetchTransactionSetSegments can fetch all OTI and TED segments", func() {
		sample824EDIString := `
ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *210217*1530*U*00401*2000000000*8*A*|
GS*SA*MILMOVE*8004171844*20190903*1617*2000000000*X*004010
ST*824*000000001
BGN*19**20211313
OTI*TR*MM**X*X*20211311**-1*AB
TED*k*Missing Data
TED*k*Missing Data
TED*k*Missing Data
SE*5*000000001
GE*2*1
GS*SA*MILMOVE*8004171844*20190903*1617*2000000000*X*004010
ST*824*000000001
BGN*19**20211313
OTI*TE*MM**X*X*20211311**-1*AB
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*2*1
IEA*1*000000001
`
		edi824 := ediResponse824.EDI{}
		err := edi824.Parse(sample824EDIString)
		suite.NoError(err)
		otis, teds := fetchTransactionSetSegments(edi824)
		suite.Equal(2, len(otis))
		suite.Equal(6, len(teds))
	})
}
