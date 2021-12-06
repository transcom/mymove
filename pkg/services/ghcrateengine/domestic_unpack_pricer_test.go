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
	servicesScheduleDest = 1
	unpackWeightBilled   = 3600
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
					Key:     models.ServiceItemParamNameWeightBilled,
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

		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticUnpackWithServiceItemParams() {
	pricer := NewDomesticUnpackPricer()

	suite.Run("success all params for domestic unpack available", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)
		paymentServiceItem := suite.setupDomesticUnpackServiceItems()

		cost, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
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
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
		suite.Equal("could not find param with key ContractCode", err.Error())

		// No requested pickup date
		missingRequestedPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameRequestedPickupDate)
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), missingRequestedPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key RequestedPickupDate", err.Error())

		// No weight
		missingBilledWeight := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilled)
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), missingBilledWeight)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilled", err.Error())

		// No services schedule destination
		missingServicesScheduleDest := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServicesScheduleDest)
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), missingServicesScheduleDest)
		suite.Error(err)
		suite.Equal("could not find param with key ServicesScheduleDest", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticUnpack() {
	pricer := NewDomesticUnpackPricer()

	suite.Run("success domestic unpack cost within peak period", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		cost, _, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilled,
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
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilled,
			servicesScheduleDest,
		)
		expectedCost := unit.Cents(4758)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.Run("failure if contract code bogus", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		_, _, err := pricer.Price(
			suite.AppContextForTest(),
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilled,
			servicesScheduleDest,
		)

		suite.Error(err)
		suite.Equal("Could not lookup domestic other price: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("failure if move date is outside of contract year", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		_, _, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unpackWeightBilled,
			servicesScheduleDest,
		)

		suite.Error(err)
		suite.Equal("Could not lookup contract year: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("weight below minimum", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		cost, _, err := pricer.Price(
			suite.AppContextForTest(),
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
		_, _, err := pricer.Price(suite.AppContextForTest(), "", requestedPickupDate, unpackWeightBilled, servicesScheduleDest)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No requested pickup date
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, time.Time{}, unpackWeightBilled, servicesScheduleDest)
		suite.Error(err)
		suite.Equal("RequestedPickupDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, 0, servicesScheduleDest)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())

		// No services schedule
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, unpackWeightBilled, 0)
		suite.Error(err)
		suite.Equal("Services schedule is required", err.Error())
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
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(unpackWeightBilled),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(servicesScheduleDest),
			},
		},
	)
}
