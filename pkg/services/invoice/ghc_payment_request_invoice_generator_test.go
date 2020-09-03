package invoice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"

	"go.uber.org/zap"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

type GHCInvoiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestGHCInvoiceSuite(t *testing.T) {

	ts := &GHCInvoiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("ghcinvoice")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GHCInvoiceSuite) TestGenerateGHCInvoice() {
	const testDateFormat = "060102"
	const testTimeFormat = "1504"
	currentTime := time.Now()
	generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}

	suite.T().Run("adds isa header line", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		suite.Equal("00", result.ISA.AuthorizationInformationQualifier)
		suite.Equal("0084182369", result.ISA.AuthorizationInformation)
		suite.Equal("00", result.ISA.SecurityInformationQualifier)
		suite.Equal("_   _", result.ISA.SecurityInformation)
		suite.Equal("ZZ", result.ISA.InterchangeSenderIDQualifier)
		suite.Equal("GOVDPIBS", result.ISA.InterchangeSenderID)
		suite.Equal("12", result.ISA.InterchangeReceiverIDQualifier)
		suite.Equal("8004171844", result.ISA.InterchangeReceiverID)
		suite.Equal(currentTime.Format(testDateFormat), result.ISA.InterchangeDate)
		suite.Equal(currentTime.Format(testTimeFormat), result.ISA.InterchangeTime)
		suite.Equal("U", result.ISA.InterchangeControlStandards)
		suite.Equal("00401", result.ISA.InterchangeControlVersionNumber)
		suite.Equal(int64(100001272), result.ISA.InterchangeControlNumber)
		suite.Equal(0, result.ISA.AcknowledgementRequested)
		suite.Equal("T", result.ISA.UsageIndicator)
		suite.Equal("|", result.ISA.ComponentElementSeparator)
	})

	suite.T().Run("adds bx header line", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		suite.IsType(&edisegment.BX{}, result.Header[0])
		bx := result.Header[0].(*edisegment.BX)
		suite.Equal("00", bx.TransactionSetPurposeCode)
		suite.Equal("J", bx.TransactionMethodTypeCode)
		suite.Equal("PP", bx.ShipmentMethodOfPayment)
		suite.Equal(*paymentRequest.MoveTaskOrder.ReferenceID, bx.ShipmentIdentificationNumber)
		suite.Equal("TRUS", bx.StandardCarrierAlphaCode)
		suite.Equal("4", bx.ShipmentQualifier)
	})

	suite.T().Run("adds payment request number to header", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		suite.IsType(&edisegment.N9{}, result.Header[1])
		n9 := result.Header[1].(*edisegment.N9)
		suite.Equal("CN", n9.ReferenceIdentificationQualifier)
		suite.Equal(paymentRequest.PaymentRequestNumber, n9.ReferenceIdentification)
	})

	suite.T().Run("does not error out creating EDI from Invoice858", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		_, err = result.EDIString()
		suite.NoError(err)
	})
}
