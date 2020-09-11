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

const testDateFormat = "20060102"
const testTimeFormat = "1504"

func (suite *GHCInvoiceSuite) TestGenerateGHCInvoiceStartEndSegments() {
	currentTime := time.Now()
	generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}
	basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(dateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilledActual,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip3,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "2424",
		},
	}

	suite.T().Run("adds isa start segment", func(t *testing.T) {
		paymentServiceItem := testdatagen.MakeMultiplePaymentServiceItemParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
		)

		result, err := generator.Generate(paymentServiceItem.PaymentRequest, false)
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
		paymentServiceItem := testdatagen.MakeMultiplePaymentServiceItemParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
		)

		result, err := generator.Generate(paymentServiceItem.PaymentRequest, false)
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
		paymentServiceItem := testdatagen.MakeMultiplePaymentServiceItemParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
		)

		result, err := generator.Generate(paymentServiceItem.PaymentRequest, false)
		suite.FatalNoError(err)
		suite.Equal(1, result.GE.NumberOfTransactionSetsIncluded)
		suite.Equal(int64(100001251), result.GE.GroupControlNumber)
	})

	suite.T().Run("adds iea end segment", func(t *testing.T) {
		paymentServiceItem := testdatagen.MakeMultiplePaymentServiceItemParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
		)

		result, err := generator.Generate(paymentServiceItem.PaymentRequest, false)
		suite.FatalNoError(err)
		suite.Equal(1, result.IEA.NumberOfIncludedFunctionalGroups)
		suite.Equal(int64(100001272), result.IEA.InterchangeControlNumber)
	})
}

func (suite *GHCInvoiceSuite) TestGenerateGHCInvoiceHeader() {
	currentTime := time.Now()
	basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(dateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilledActual,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip3,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "2424",
		},
	}
	generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}
	paymentServiceItem := testdatagen.MakeMultiplePaymentServiceItemParams(
		suite.DB(),
		models.ReServiceCodeDLH,
		basicPaymentServiceItemParams,
	)

	result, err := generator.Generate(paymentServiceItem.PaymentRequest, false)
	suite.FatalNoError(err)

	suite.T().Run("adds bx header segment", func(t *testing.T) {
		suite.IsType(&edisegment.BX{}, result.Header[0])
		bx := result.Header[0].(*edisegment.BX)
		suite.Equal("00", bx.TransactionSetPurposeCode)
		suite.Equal("J", bx.TransactionMethodTypeCode)
		suite.Equal("PP", bx.ShipmentMethodOfPayment)
		suite.Equal(*paymentServiceItem.PaymentRequest.MoveTaskOrder.ReferenceID, bx.ShipmentIdentificationNumber)
		suite.Equal("TRUS", bx.StandardCarrierAlphaCode)
		suite.Equal("4", bx.ShipmentQualifier)
	})

	suite.T().Run("does not error out creating EDI from Invoice858", func(t *testing.T) {
		_, err := result.EDIString()
		suite.NoError(err)
	})

	serviceMember := paymentServiceItem.PaymentRequest.MoveTaskOrder.Orders.ServiceMember
	testData := []struct {
		TestName      string
		Qualifier     string
		ExpectedValue string
	}{
		{TestName: "payment request number", Qualifier: "CN", ExpectedValue: paymentServiceItem.PaymentRequest.PaymentRequestNumber},
		{TestName: "contract code", Qualifier: "CT", ExpectedValue: "TRUSS_TEST"},
		{TestName: "service member name", Qualifier: "1W", ExpectedValue: serviceMember.ReverseNameLineFormat()},
		{TestName: "service member rank", Qualifier: "ML", ExpectedValue: string(*serviceMember.Rank)},
		{TestName: "service member branch", Qualifier: "3L", ExpectedValue: string(*serviceMember.Affiliation)},
	}

	for idx, data := range testData {
		suite.T().Run(fmt.Sprintf("adds %s to header", data.TestName), func(t *testing.T) {
			suite.IsType(&edisegment.N9{}, result.Header[idx+1])
			n9 := result.Header[idx+1].(*edisegment.N9)
			suite.Equal(data.Qualifier, n9.ReferenceIdentificationQualifier)
			suite.Equal(data.ExpectedValue, n9.ReferenceIdentification)
		})
	}

	suite.T().Run("adds actual pickup date to header", func(t *testing.T) {
		suite.IsType(&edisegment.G62{}, result.Header[6])
		g62 := result.Header[6].(*edisegment.G62)
		suite.Equal(86, g62.DateQualifier)
		suite.Equal(currentTime.Format(testDateFormat), g62.Date)
	})

	suite.T().Run("adds orders destination address", func(t *testing.T) {
		// name
		expectedDutyStation := paymentServiceItem.PaymentRequest.MoveTaskOrder.Orders.NewDutyStation
		suite.IsType(&edisegment.N1{}, result.Header[7])
		n1 := result.Header[7].(*edisegment.N1)
		suite.Equal("ST", n1.EntityIdentifierCode)
		suite.Equal(expectedDutyStation.Name, n1.Name)
		suite.Equal("10", n1.IdentificationCodeQualifier)
		suite.Equal(expectedDutyStation.TransportationOffice.Gbloc, n1.IdentificationCode)
		// street address
		address := expectedDutyStation.Address
		suite.IsType(&edisegment.N3{}, result.Header[8])
		n3 := result.Header[8].(*edisegment.N3)
		suite.Equal(address.StreetAddress1, n3.AddressInformation1)
		suite.Equal(*address.StreetAddress2, n3.AddressInformation2)
		// city state info
		suite.IsType(&edisegment.N4{}, result.Header[9])
		n4 := result.Header[9].(*edisegment.N4)
		suite.Equal(address.City, n4.CityName)
		suite.Equal(address.State, n4.StateOrProvinceCode)
		suite.Equal(address.PostalCode, n4.PostalCode)
		suite.Equal(*address.Country, n4.CountryCode)
	})

	suite.T().Run("adds orders origin address", func(t *testing.T) {
		// name
		expectedDutyStation := paymentServiceItem.PaymentRequest.MoveTaskOrder.Orders.OriginDutyStation
		suite.IsType(&edisegment.N1{}, result.Header[10])
		n1 := result.Header[10].(*edisegment.N1)
		suite.Equal("SF", n1.EntityIdentifierCode)
		suite.Equal(expectedDutyStation.Name, n1.Name)
		suite.Equal("10", n1.IdentificationCodeQualifier)
		suite.Equal(expectedDutyStation.TransportationOffice.Gbloc, n1.IdentificationCode)
		// street address
		address := expectedDutyStation.Address
		suite.IsType(&edisegment.N3{}, result.Header[11])
		n3 := result.Header[11].(*edisegment.N3)
		suite.Equal(address.StreetAddress1, n3.AddressInformation1)
		suite.Equal(*address.StreetAddress2, n3.AddressInformation2)
		// city state info
		suite.IsType(&edisegment.N4{}, result.Header[12])
		n4 := result.Header[12].(*edisegment.N4)
		suite.Equal(address.City, n4.CityName)
		suite.Equal(address.State, n4.StateOrProvinceCode)
		suite.Equal(address.PostalCode, n4.PostalCode)
		suite.Equal(*address.Country, n4.CountryCode)
	})
}

func (suite *GHCInvoiceSuite) TestGenerateGHCInvoiceBody() {
	currentTime := time.Now()
	basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(dateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilledActual,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip3,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "2424",
		},
	}
	generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}

	suite.T().Run("adds l0 service item segment", func(t *testing.T) {
		paymentServiceItem := testdatagen.MakeMultiplePaymentServiceItemParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
		)

		result, err := generator.Generate(paymentServiceItem.PaymentRequest, false)
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
