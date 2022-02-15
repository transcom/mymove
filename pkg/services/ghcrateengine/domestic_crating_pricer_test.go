package ghcrateengine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dcrtTestServiceSchedule      = 3
	dcrtTestBasePriceCents       = unit.Cents(2300)
	dcrtTestEscalationCompounded = 1.125
	dcrtTestBilledCubicFeet      = 10
	dcrtTestPriceCents           = unit.Cents(25875) // dcrtTestBasePriceCents * (dcrtTestBilledCubicFeet ) * dcrtTestEscalationCompounded
)

var dcrtTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticCratingPricer() {
	pricer := NewDomesticCratingPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)

		paymentServiceItem := suite.setupDomesticCratingServiceItem()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dcrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dcrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dcrtTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(dcrtTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid crating volume", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)
		badVolume := unit.CubicFeet(3.0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, badVolume, dcrtTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "crate must be billed for a minimum of 4 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, testdatagen.DefaultContractCode, dcrtTestEscalationCompounded)
		twoYearsLaterPickupDate := dcrtTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticCratingServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDCRT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameCubicFeetBilled,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   fmt.Sprintf("%d", int(dcrtTestBilledCubicFeet)),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dcrtTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dcrtTestServiceSchedule),
			},
		},
	)
}
