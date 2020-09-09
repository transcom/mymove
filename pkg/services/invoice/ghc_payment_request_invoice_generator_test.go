package invoice

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
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

func (suite *GHCInvoiceSuite) TestGenerateGHCInvoiceStartEndSegments() {
	const testDateFormat = "060102"
	const testTimeFormat = "1504"
	currentTime := time.Now()
	generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}

	suite.T().Run("adds isa start segment", func(t *testing.T) {
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

	suite.T().Run("adds gs start segment", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		suite.Equal("SI", result.GS.FunctionalIdentifierCode)
		suite.Equal("MYMOVE", result.GS.ApplicationSendersCode)
		suite.Equal("8004171844", result.GS.ApplicationReceiversCode)
		suite.Equal(currentTime.Format(dateFormat), result.GS.Date)
		suite.Equal(currentTime.Format(timeFormat), result.GS.Time)
		suite.Equal(int64(100001251), result.GS.GroupControlNumber)
		suite.Equal("X", result.GS.ResponsibleAgencyCode)
		suite.Equal("004010", result.GS.Version)
	})

	suite.T().Run("adds ge end segment", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		suite.Equal(1, result.GE.NumberOfTransactionSetsIncluded)
		suite.Equal(int64(100001251), result.GE.GroupControlNumber)
	})

	suite.T().Run("adds iea end segment", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		suite.Equal(1, result.IEA.NumberOfIncludedFunctionalGroups)
		suite.Equal(int64(100001272), result.IEA.InterchangeControlNumber)
	})
}

func (suite *GHCInvoiceSuite) TestGenerateGHCInvoiceHeader() {
	generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}
	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})

	result, err := generator.Generate(paymentRequest, false)
	suite.FatalNoError(err)

	suite.T().Run("adds bx header segment", func(t *testing.T) {
		suite.IsType(&edisegment.BX{}, result.Header[0])
		bx := result.Header[0].(*edisegment.BX)
		suite.Equal("00", bx.TransactionSetPurposeCode)
		suite.Equal("J", bx.TransactionMethodTypeCode)
		suite.Equal("PP", bx.ShipmentMethodOfPayment)
		suite.Equal(*paymentRequest.MoveTaskOrder.ReferenceID, bx.ShipmentIdentificationNumber)
		suite.Equal("TRUS", bx.StandardCarrierAlphaCode)
		suite.Equal("4", bx.ShipmentQualifier)
	})

	suite.T().Run("does not error out creating EDI from Invoice858", func(t *testing.T) {
		_, err := result.EDIString()
		suite.NoError(err)
	})

	serviceMember := paymentRequest.MoveTaskOrder.Orders.ServiceMember
	testData := []struct {
		TestName      string
		Position      int
		Qualifier     string
		ExpectedValue string
	}{
		{TestName: "payment request number", Position: 1, Qualifier: "CN", ExpectedValue: paymentRequest.PaymentRequestNumber},
		{TestName: "service member name", Position: 2, Qualifier: "1W", ExpectedValue: serviceMember.ReverseNameLineFormat()},
		{TestName: "service member rank", Position: 3, Qualifier: "ML", ExpectedValue: string(*serviceMember.Rank)},
		{TestName: "service member branch", Position: 4, Qualifier: "3L", ExpectedValue: string(*serviceMember.Affiliation)},
	}

	for _, data := range testData {
		suite.T().Run(fmt.Sprintf("adds %s to header", data.TestName), func(t *testing.T) {
			suite.IsType(&edisegment.N9{}, result.Header[data.Position])
			n9 := result.Header[data.Position].(*edisegment.N9)
			suite.Equal(data.Qualifier, n9.ReferenceIdentificationQualifier)
			suite.Equal(data.ExpectedValue, n9.ReferenceIdentification)
		})
	}
}

func (suite *GHCInvoiceSuite) TestGenerateGHCInvoiceBody() {
	generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}

	suite.T().Run("adds l0 service item segment", func(t *testing.T) {
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{})
		// add PSIs
		var params models.PaymentServiceItemParams

		paymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
			PaymentServiceItem: models.PaymentServiceItem{
				PaymentRequestID: paymentRequest.ID,
			},
		})

		weightServiceItemParamKey := testdatagen.MakeServiceItemParamKey(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: models.ServiceItemParamKey{
					Key:  models.ServiceItemParamNameWeightBilledActual,
					Type: models.ServiceItemParamTypeInteger,
				},
			})

		weightServiceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem:  paymentServiceItem,
				ServiceItemParamKey: weightServiceItemParamKey,
				PaymentServiceItemParam: models.PaymentServiceItemParam{
					Value: "4242",
				},
			})
		params = append(params, weightServiceItemParam)

		distanceServiceItemParamKey := testdatagen.MakeServiceItemParamKey(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: models.ServiceItemParamKey{
					Key:  models.ServiceItemParamNameDistanceZip3,
					Type: models.ServiceItemParamTypeInteger,
				},
			})

		distanceServiceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem:  paymentServiceItem,
				ServiceItemParamKey: distanceServiceItemParamKey,
				PaymentServiceItemParam: models.PaymentServiceItemParam{
					Value: "2424",
				},
			})
		params = append(params, distanceServiceItemParam)

		paymentServiceItem.PaymentServiceItemParams = params
		suite.MustSave(&paymentServiceItem)

		result, err := generator.Generate(paymentRequest, false)
		suite.FatalNoError(err)
		lastIdx := len(result.ServiceItems) - 1
		suite.IsType(&edisegment.L0{}, result.ServiceItems[lastIdx])
		l0 := result.ServiceItems[lastIdx].(*edisegment.L0)
		suite.Equal(1, l0.LadingLineItemNumber)
		suite.Equal(float64(2424), l0.BilledRatedAsQuantity)
		suite.Equal("DM", l0.BilledRatedAsQualifier)
		suite.Equal(float64(4242), l0.Weight)
		suite.Equal("B", l0.WeightQualifier)
		suite.Equal("L", l0.WeightUnitCode)
	})
}
