package ghcrateengine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testhelpers"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dpkTestEscalationCompounded   = 1.0407
	dpkTestIsPeakPeriod           = true
	dpkTestWeight                 = unit.Pound(2100)
	dpkTestServicesScheduleOrigin = 1
	dpkTestContractYearName       = "DPK Test Year"
	dpkTestBasePriceCents         = unit.Cents(146)
	dpkTestPriceCents             = unit.Cents(3192)
	dpkDMHFTestPriceCents         = unit.Cents(3402)
	mobileHomeFactor              = 33.51
)

var dpkTestRequestedPickupDate = time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticPackPricer() {
	pricer := NewDomesticPackPricer()
	featureFlagValues := testhelpers.MakeMobileHomeFFMap(true, true)
	var paymentServiceItem models.PaymentServiceItem
	suite.PreloadData(func() {
		paymentServiceItem = suite.setupDomesticPackServiceItem()
	})

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded, false)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, nil)
		suite.NoError(err)
		suite.Equal(dpkTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dpkTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success using PaymentServiceItemParams and with Domestic Mobile Home Factor applied", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded, true)
		paymentServiceItem.PaymentServiceItemParams[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType = models.MTOShipmentTypeMobileHome // SEt

		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, featureFlagValues)
		suite.NoError(err)
		suite.Equal(dpkDMHFTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Mobile Home Factor Test Year"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.11000"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "true"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.46"},
			{Key: models.ServiceItemParamNameMobileHomeFactor, Value: "33.51"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid parameters to PriceUsingParams", func() {
		paymentServiceItem := suite.setupDomesticPackServiceItem()

		// Setting each param's type to something incorrect should trigger an error.

		// WeightBilled
		paymentServiceItem.PaymentServiceItemParams[3].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, nil)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameWeightBilled))

		// ServicesScheduleOrigin
		paymentServiceItem.PaymentServiceItemParams[2].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, nil)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameServicesScheduleOrigin))

		// ReferenceDate
		paymentServiceItem.PaymentServiceItemParams[1].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, nil)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a time", models.ServiceItemParamNameReferenceDate))

		// ContractCode
		paymentServiceItem.PaymentServiceItemParams[0].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, nil)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a string", models.ServiceItemParamNameContractCode))
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticPackServiceItem() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDPK,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dpkTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dpkTestServicesScheduleOrigin),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dpkTestWeight.Int()),
			},
		}, nil, nil,
	)
}
