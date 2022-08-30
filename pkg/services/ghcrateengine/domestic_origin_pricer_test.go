package ghcrateengine

import (
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	dopTestServiceArea = "006"
	dopTestWeight      = 3600
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticOriginWithServiceItemParamsBadData() {
	pricer := NewDomesticOriginPricer()

	suite.Run("failure during pricing bubbles up", func() {
		suite.setUpDomesticOriginData()
		paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOP,
			[]testdatagen.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   testdatagen.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameReferenceDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   "0",
				},
				{
					Key:     models.ServiceItemParamNameServiceAreaOrigin,
					KeyType: models.ServiceItemParamTypeString,
					Value:   dopTestServiceArea,
				},
			},
		)

		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticOriginWithServiceItemParams() {
	pricer := NewDomesticOriginPricer()

	suite.Run("success all params for domestic origin available", func() {
		suite.setUpDomesticOriginData()
		paymentServiceItem := suite.setupDomesticOriginServiceItems()

		cost, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		expectedCost := unit.Cents(5470)

		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Test Contract Year"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "true"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.46"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("validation errors", func() {
		suite.setUpDomesticOriginData()
		paymentServiceItem := suite.setupDomesticOriginServiceItems()

		// No contract code
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
		suite.Equal("could not find param with key ContractCode", err.Error())

		// No reference date
		missingReferenceDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameReferenceDate)
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), missingReferenceDate)
		suite.Error(err)
		suite.Equal("could not find param with key ReferenceDate", err.Error())

		// No weight
		missingBilledWeight := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilled)
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), missingBilledWeight)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilled", err.Error())

		// No service area
		missingServiceAreaOrigin := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServiceAreaOrigin)
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), missingServiceAreaOrigin)
		suite.Error(err)
		suite.Equal("could not find param with key ServiceAreaOrigin", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticOrigin() {
	pricer := NewDomesticOriginPricer()

	suite.Run("success domestic origin cost within peak period", func() {
		suite.setUpDomesticOriginData()
		isPPM := false

		cost, displayParams, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dopTestWeight,
			dopTestServiceArea,
			isPPM,
		)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Test Contract Year"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "true"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.46"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success domestic origin cost within non-peak period", func() {
		suite.setUpDomesticOriginData()
		isPPM := false

		nonPeakDate := peakStart.addDate(0, -1)
		cost, displayParams, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			dopTestWeight,
			dopTestServiceArea,
			isPPM,
		)

		expectedCost := unit.Cents(4758)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Test Contract Year"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "false"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.27"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("failure if contract code bogus", func() {
		suite.setUpDomesticOriginData()
		isPPM := false

		_, _, err := pricer.Price(
			suite.AppContextForTest(),
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dopTestWeight,
			dopTestServiceArea,
			isPPM,
		)

		suite.Error(err)
		suite.Equal("Could not lookup Domestic Service Area Price: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("failure if move date is outside of contract year", func() {
		suite.setUpDomesticOriginData()
		isPPM := false

		_, _, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dopTestWeight,
			dopTestServiceArea,
			isPPM,
		)

		suite.Error(err)
		suite.Equal("Could not lookup contract year: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("weight below minimum", func() {
		suite.setUpDomesticOriginData()
		isPPM := false

		cost, _, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unit.Pound(499),
			dopTestServiceArea,
			isPPM,
		)
		suite.Equal(unit.Cents(0), cost)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})

	suite.Run("validation errors", func() {
		suite.setUpDomesticOriginData()
		isPPM := false

		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No contract code
		_, _, err := pricer.Price(suite.AppContextForTest(), "", requestedPickupDate, dshTestWeight, dopTestServiceArea, isPPM)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No reference date
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, time.Time{}, dshTestWeight, dopTestServiceArea, isPPM)
		suite.Error(err)
		suite.Equal("ReferenceDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, 0, dopTestServiceArea, isPPM)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())

		// No service area
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, dshTestWeight, "", isPPM)
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
	})

	suite.Run("successfully finds dom destination price for ppm with weight < 500 lbs with Price method", func() {
		suite.setUpDomesticOriginData()
		suite.setupDomesticOriginServiceItems()
		isPPM := true
		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// the PPM price for weights < 500 should be prorated from a base of 500
		basePriceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, unit.Pound(500), dopTestServiceArea, isPPM)
		suite.NoError(err)

		halfPriceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, unit.Pound(250), dopTestServiceArea, isPPM)
		suite.NoError(err)
		suite.Equal(basePriceCents/2, halfPriceCents)

		fifthPriceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, unit.Pound(100), dopTestServiceArea, isPPM)
		suite.NoError(err)
		suite.Equal(basePriceCents/5, fifthPriceCents)
	})

	suite.Run("successfully finds dom destination price for ppm with weight < 500 lbs with PriceUsingParams method", func() {
		suite.setUpDomesticOriginData()
		paymentServiceItem := suite.setupDomesticOriginServiceItems()
		params := paymentServiceItem.PaymentServiceItemParams
		params[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType = models.MTOShipmentTypePPM
		weightBilledIndex := 3

		params[weightBilledIndex].Value = "500"
		basePriceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)

		params[weightBilledIndex].Value = "250"
		halfPriceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(basePriceCents/2, halfPriceCents)

		params[weightBilledIndex].Value = "100"
		fifthPriceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(basePriceCents/5, fifthPriceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Test Contract Year"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "true"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.46"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOriginServiceItems() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDOP,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dopTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dshTestWeight),
			},
		},
	)
}

func (suite *GHCRateEngineServiceSuite) setUpDomesticOriginData() {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Escalation:           1.0197,
				EscalationCompounded: 1.0407,
			},
		})

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: dopTestServiceArea,
			},
		})

	domesticOriginService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDOP,
				Name: "Dom. Origin Price",
			},
		})

	domesticOriginPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		DomesticServiceAreaID: serviceArea.ID,
		IsPeakPeriod:          true,
		ServiceID:             domesticOriginService.ID,
	}

	domesticOriginPeakPrice := domesticOriginPrice
	domesticOriginPeakPrice.PriceCents = 146
	suite.MustSave(&domesticOriginPeakPrice)

	domesticOriginNonpeakPrice := domesticOriginPrice
	domesticOriginNonpeakPrice.IsPeakPeriod = false
	domesticOriginNonpeakPrice.PriceCents = 127
	suite.MustSave(&domesticOriginNonpeakPrice)
}
