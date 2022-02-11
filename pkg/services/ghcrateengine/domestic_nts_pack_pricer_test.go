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
	dnpkTestEscalationCompounded   = 1.0407
	dnpkTestIsPeakPeriod           = true
	dnpkTestWeight                 = unit.Pound(2100)
	dnpkTestServicesScheduleOrigin = 1
	dnpkTestContractYearName       = "DNPK Test Year"
	dnpkTestBasePriceCents         = unit.Cents(6333)
	dnpkTestFactor                 = 1.35
	dnpkTestPriceCents             = unit.Cents(186848) // dnpkTestBasePriceCents * (dnpkTestWeight / 100) * dnpkTestEscalationCompounded * dnpkTestFactor
)

var dnpkTestRequestedPickupDate = time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 5, 5, 5, 5, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticNTSPackPricer() {
	pricer := NewDomesticNTSPackPricer()

	suite.Run("success using PaymentServiceItemParams", func() {
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		paymentServiceItem := suite.setupDomesticNTSPackServiceItem()

		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dnpkTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dnpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dnpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dnpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dnpkTestBasePriceCents)},
			{Key: models.ServiceItemParamNameNTSPackingFactor, Value: FormatFloat(dnpkTestFactor, 2)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid parameters to PriceUsingParams", func() {
		paymentServiceItem := suite.setupDomesticNTSPackServiceItem()

		// Setting each param's type to something incorrect should trigger an error.

		// WeightBilled
		paymentServiceItem.PaymentServiceItemParams[3].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameWeightBilled))

		// ServicesScheduleOrigin
		paymentServiceItem.PaymentServiceItemParams[2].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to an int", models.ServiceItemParamNameServicesScheduleOrigin))

		// ReferenceDate
		paymentServiceItem.PaymentServiceItemParams[1].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a time", models.ServiceItemParamNameReferenceDate))

		// ContractCode
		paymentServiceItem.PaymentServiceItemParams[0].ServiceItemParamKey.Type = models.ServiceItemParamTypeBoolean
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("trying to convert %s to a string", models.ServiceItemParamNameContractCode))
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticNTSPackServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDNPK,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dnpkTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dnpkTestServicesScheduleOrigin),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dnpkTestWeight.Int()),
			},
		},
	)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticNTSPackPrices(schedule int, isPeakPeriod bool, priceCents unit.Cents, market models.Market, factor float64, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	packService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDPK,
			},
		})

	otherPrice := models.ReDomesticOtherPrice{
		ContractID:   contractYear.Contract.ID,
		ServiceID:    packService.ID,
		IsPeakPeriod: isPeakPeriod,
		Schedule:     schedule,
		PriceCents:   priceCents,
	}

	suite.MustSave(&otherPrice)

	ntsPackService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDNPK,
			},
		})

	shipmentTypePrice := models.ReShipmentTypePrice{
		ContractID: contractYear.Contract.ID,
		ServiceID:  ntsPackService.ID,
		Market:     market,
		Factor:     factor,
	}

	suite.MustSave(&shipmentTypePrice)
}
