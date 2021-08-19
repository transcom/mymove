package ghcrateengine

import (
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const (
	servicesScheduleOrigin = 1
	weightBilledActual     = 3600
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticPackWithServiceItemParamsBadData() {
	suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDPK)
	paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDPK,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "0",
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(servicesScheduleOrigin),
			},
		},
	)

	pricer := NewDomesticPackPricer(suite.DB())

	suite.Run("failure during pricing bubbles up", func() {
		_, _, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticPackWithServiceItemParams() {
	suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDPK)
	paymentServiceItem := suite.setupDomesticPackServiceItems()

	pricer := NewDomesticPackPricer(suite.DB())

	suite.Run("success all params for domestic pack available", func() {
		cost, displayParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Base Period Year 1"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "true"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.46"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("validation errors", func() {
		// No contract code
		_, _, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
		suite.Equal("could not find param with key ContractCode", err.Error())

		// No requested pickup date
		missingRequestedPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameRequestedPickupDate)
		_, _, err = pricer.PriceUsingParams(missingRequestedPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key RequestedPickupDate", err.Error())

		// No weight
		missingBilledActualWeight := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilledActual)
		_, _, err = pricer.PriceUsingParams(missingBilledActualWeight)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilledActual", err.Error())

		// No service schedule origin
		missingServicesScheduleOrigin := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServicesScheduleOrigin)
		_, _, err = pricer.PriceUsingParams(missingServicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("could not find param with key ServicesScheduleOrigin", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticPack() {
	suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDPK)

	pricer := NewDomesticPackPricer(suite.DB())

	suite.Run("success domestic pack cost within peak period", func() {
		cost, _, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.Run("success domestic pack cost within non-peak period", func() {
		nonPeakDate := peakStart.addDate(0, -1)
		cost, _, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)
		expectedCost := unit.Cents(4758)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.Run("failure if contract code bogus", func() {
		_, _, err := pricer.Price(
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)

		suite.Error(err)
		suite.Equal("Could not lookup Domestic Other Price: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("failure if move date is outside of contract year", func() {
		_, _, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)

		suite.Error(err)
		suite.Equal("Could not lookup contract year: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("weight below minimum", func() {
		cost, _, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unit.Pound(499),
			servicesScheduleOrigin,
		)
		suite.Equal(unit.Cents(0), cost)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})

	suite.Run("validation errors", func() {
		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No contract code
		_, _, err := pricer.Price("", requestedPickupDate, weightBilledActual, servicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No requested pickup date
		_, _, err = pricer.Price(testdatagen.DefaultContractCode, time.Time{}, weightBilledActual, servicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("RequestedPickupDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, 0, servicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())

		// No service schedule
		_, _, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, weightBilledActual, 0)
		suite.Error(err)
		suite.Equal("Service schedule is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticPackServiceItems() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDPK,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(weightBilledActual),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(servicesScheduleOrigin),
			},
		},
	)
}
