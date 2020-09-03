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

	suite.T().Run("adds isa header line", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := GHCPaymentRequestInvoiceGenerator{}.Generate(paymentRequest, false)
		suite.NoError(err)
		suite.IsType(&edisegment.ISA{}, result.Header[0])
		isa := result.Header[0].(*edisegment.ISA)
		suite.Equal("00", isa.AuthorizationInformationQualifier)
		suite.Equal("0084182369", isa.AuthorizationInformation)
		suite.Equal("00", isa.SecurityInformationQualifier)
		suite.Equal("_   _", isa.SecurityInformation)
		suite.Equal("ZZ", isa.InterchangeSenderIDQualifier)
		suite.Equal("GOVDPIBS", isa.InterchangeSenderID)
		suite.Equal("12", isa.InterchangeReceiverIDQualifier)
		suite.Equal("8004171844", isa.InterchangeReceiverID)
		suite.Equal(currentTime.Format(testDateFormat), isa.InterchangeDate)
		suite.Equal(currentTime.Format(testTimeFormat), isa.InterchangeTime)
		suite.Equal("U", isa.InterchangeControlStandards)
		suite.Equal("00401", isa.InterchangeControlVersionNumber)
		suite.Equal(100001272, isa.InterchangeControlNumber)
		suite.Equal("0", isa.AcknowledgementRequested)
		suite.Equal("T", isa.UsageIndicator)
		suite.Equal("|", isa.ComponentElementSeparator)
	})

	suite.T().Run("adds bx header line", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := GHCPaymentRequestInvoiceGenerator{}.Generate(paymentRequest, false)
		suite.NoError(err)
		suite.IsType(&edisegment.BX{}, result.Header[1])
		bx := result.Header[1].(*edisegment.BX)
		suite.Equal("00", bx.TransactionSetPurposeCode)
		suite.Equal("J", bx.TransactionMethodTypeCode)
		suite.Equal("PP", bx.ShipmentMethodOfPayment)
		suite.Equal(*paymentRequest.MoveTaskOrder.ReferenceID, bx.ShipmentIdentificationNumber)
		suite.Equal("TRUS", bx.StandardCarrierAlphaCode)
		suite.Equal("4", bx.ShipmentQualifier)
	})

	suite.T().Run("adds payment request number to header", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := GHCPaymentRequestInvoiceGenerator{}.Generate(paymentRequest, false)
		suite.NoError(err)
		suite.IsType(&edisegment.N9{}, result.Header[2])
		n9 := result.Header[2].(*edisegment.N9)
		suite.Equal("CN", n9.ReferenceIdentificationQualifier)
		suite.Equal(paymentRequest.PaymentRequestNumber, n9.ReferenceIdentification)
	})

	suite.T().Run("does not error out creating EDI from Invoice858", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := GHCPaymentRequestInvoiceGenerator{}.Generate(paymentRequest, false)
		suite.NoError(err)
		_, err = result.EDIString()
		suite.NoError(err)
	})
}
