package invoice

import (
	"testing"

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
	suite.T().Run("addes bx header line", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := GHCPaymentRequestInvoiceGenerator{}.Generate(paymentRequest)
		suite.NoError(err)
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
		result, err := GHCPaymentRequestInvoiceGenerator{}.Generate(paymentRequest)
		suite.NoError(err)
		suite.IsType(&edisegment.N9{}, result.Header[1])
		n9 := result.Header[1].(*edisegment.N9)
		suite.Equal("CN", n9.ReferenceIdentificationQualifier)
		suite.Equal(paymentRequest.PaymentRequestNumber, n9.ReferenceIdentification)
	})

	suite.T().Run("does not error out creating EDI from Invoice858", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := GHCPaymentRequestInvoiceGenerator{}.Generate(paymentRequest)
		suite.NoError(err)
		_, err = result.EDIString()
		suite.NoError(err)
	})
}
