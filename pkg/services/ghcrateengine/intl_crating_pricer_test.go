package ghcrateengine

import (
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	icrtTestMarket               = models.Market("O")
	icrtTestBasePriceCents       = unit.Cents(2859)
	icrtTestEscalationCompounded = 1.11000
	icrtTestBilledCubicFeet      = unit.CubicFeet(10)
	icrtTestPriceCents           = unit.Cents(31730)
	icrtTestStandaloneCrate      = false
	icrtTestStandaloneCrateCap   = unit.Cents(1000000)
	icrtTestUncappedRequestTotal = unit.Cents(31730)
	icrtTestExternalCrate        = false
)

var icrtTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestIntlCratingPricer() {
	pricer := NewIntlCratingPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		paymentServiceItem := suite.setupIntlCratingServiceItem(icrtTestBilledCubicFeet)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(icrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(icrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(icrtTestBasePriceCents)},
			{Key: models.ServiceItemParamNameUncappedRequestTotal, Value: FormatCents(icrtTestUncappedRequestTotal)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
	suite.Run("success with truncating cubic feet", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		paymentServiceItem := suite.setupIntlCratingServiceItem(unit.CubicFeet(10.005))
		priceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(icrtTestPriceCents, priceCents)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)
		suite.NoError(err)
		suite.Equal(icrtTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid crating volume - external crate", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)
		badVolume := unit.CubicFeet(3.0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, badVolume, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, true, icrtTestMarket)
		suite.Error(err)
		suite.Contains(err.Error(), "external crates must be billed for a minimum of 4.00 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)
		twoYearsLaterPickupDate := icrtTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, twoYearsLaterPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupIntlCratingServiceItem(cubicFeet unit.CubicFeet) models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeICRT,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameCubicFeetBilled,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   cubicFeet.String(),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   icrtTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameStandaloneCrate,
				KeyType: models.ServiceItemParamTypeBoolean,
				Value:   strconv.FormatBool(false),
			},
			{
				Key:     models.ServiceItemParamNameStandaloneCrateCap,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.FormatInt(100000, 10),
			},
			{
				Key:     models.ServiceItemParamNameMarketOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   icrtTestMarket.String(),
			},
			{
				Key:     models.ServiceItemParamNameExternalCrate,
				KeyType: models.ServiceItemParamTypeBoolean,
				Value:   strconv.FormatBool(false),
			},
		}, nil, nil,
	)
}
