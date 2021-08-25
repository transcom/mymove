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
	servicesScheduleDest     = 1
	unpackWeightBilledActual = 3600
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticUnpackWithServiceItemParamsBadData() {
	pricer := NewDomesticUnpackPricer()

	suite.Run("failure during pricing bubbles up", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)
		paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDUPK,
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
					Key:     models.ServiceItemParamNameServicesScheduleDest,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   strconv.Itoa(servicesScheduleDest),
				},
			},
		)

		_, _, err := pricer.PriceUsingParams(suite.TestAppContext(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticUnpackWithServiceItemParams() {
	pricer := NewDomesticUnpackPricer()

	suite.Run("success all params for domestic unpack available", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)
		paymentServiceItem := suite.setupDomesticUnpackServiceItems()

		cost, displayParams, err := pricer.PriceUsingParams(suite.TestAppContext(), paymentServiceItem.PaymentServiceItemParams)
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
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)
		paymentServiceItem := suite.setupDomesticUnpackServiceItems()

		// No contract code
		_, _, err := pricer.PriceUsingParams(suite.TestAppContext(), models.PaymentServiceItemParams{})
		suite.Error(err)
		suite.Equal("could not find param with key ContractCode", err.Error())

		// No requested pickup date
		missingRequestedPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameRequestedPickupDate)
		_, _, err = pricer.PriceUsingParams(suite.TestAppContext(), missingRequestedPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key RequestedPickupDate", err.Error())

		// No weight
		missingBilledActualWeight := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilledActual)
		_, _, err = pricer.PriceUsingParams(suite.TestAppContext(), missingBilledActualWeight)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilledActual", err.Error())

		// No service schedule destination
		missingServicesScheduleDest := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServicesScheduleDest)
		_, _, err = pricer.PriceUsingParams(suite.TestAppContext(), missingServicesScheduleDest)
		suite.Error(err)
		suite.Equal("could not find param with key ServicesScheduleDest", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticUnpack() {
	pricer := NewDomesticUnpackPricer()

	suite.Run("success domestic unpack cost within peak period", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		cost, _, err := pricer.Price(
			suite.TestAppContext(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilledActual,
			servicesScheduleDest,
		)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.Run("success domestic unpack cost within non-peak period", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		nonPeakDate := peakStart.addDate(0, -1)
		cost, _, err := pricer.Price(
			suite.TestAppContext(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilledActual,
			servicesScheduleDest,
		)
		expectedCost := unit.Cents(4758)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.Run("failure if contract code bogus", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		_, _, err := pricer.Price(
			suite.TestAppContext(),
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilledActual,
			servicesScheduleDest,
		)

		suite.Error(err)
		suite.Equal("Could not lookup Domestic Other Price: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("failure if move date is outside of contract year", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		_, _, err := pricer.Price(
			suite.TestAppContext(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilledActual,
			servicesScheduleDest,
		)

		suite.Error(err)
		suite.Equal("Could not lookup contract year: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("weight below minimum", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		cost, _, err := pricer.Price(
			suite.TestAppContext(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unit.Pound(499),
			servicesScheduleDest,
		)
		suite.Equal(unit.Cents(0), cost)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})

	suite.Run("validation errors", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No contract code
		_, _, err := pricer.Price(suite.TestAppContext(), "", requestedPickupDate, unpackWeightBilledActual, servicesScheduleDest)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No requested pickup date
		_, _, err = pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, time.Time{}, unpackWeightBilledActual, servicesScheduleDest)
		suite.Error(err)
		suite.Equal("RequestedPickupDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, requestedPickupDate, 0, servicesScheduleDest)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())

		// No service schedule
		_, _, err = pricer.Price(suite.TestAppContext(), testdatagen.DefaultContractCode, requestedPickupDate, unpackWeightBilledActual, 0)
		suite.Error(err)
		suite.Equal("Service schedule is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticUnpackServiceItems() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDUPK,
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
				Value:   strconv.Itoa(unpackWeightBilledActual),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(servicesScheduleDest),
			},
		},
	)
}
