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
	ddpTestServiceArea = "006"
	ddpTestWeight      = 3700
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticDestinationWithServiceItemParamsBadData() {
	pricer := NewDomesticDestinationPricer()

	suite.Run("failure during pricing bubbles up", func() {
		suite.setUpDomesticDestinationData()
		paymentServiceItem := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDP,
			[]factory.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   factory.DefaultContractCode,
				},
				{
					Key:     models.ServiceItemParamNameReferenceDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameServiceAreaDest,
					KeyType: models.ServiceItemParamTypeString,
					Value:   ddpTestServiceArea,
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   "0",
				},
			}, nil, nil,
		)

		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticDestinationWithServiceItemParams() {
	pricer := NewDomesticDestinationPricer()

	suite.Run("success all params for destination available", func() {
		suite.setUpDomesticDestinationData()
		paymentServiceItem := suite.setupDomesticDestinationServiceItems()

		cost, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		expectedCost := unit.Cents(5624)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.Run("validation errors", func() {
		suite.setUpDomesticDestinationData()
		paymentServiceItem := suite.setupDomesticDestinationServiceItems()

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
		missingServiceAreaDest := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServiceAreaDest)
		_, _, err = pricer.PriceUsingParams(suite.AppContextForTest(), missingServiceAreaDest)
		suite.Error(err)
		suite.Equal("could not find param with key ServiceAreaDest", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticDestination() {
	pricer := NewDomesticDestinationPricer()

	suite.Run("success destination cost within peak period", func() {
		suite.setUpDomesticDestinationData()
		isPPM := false
		cost, displayParams, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
			isPPM,
		)
		expectedCost := unit.Cents(5624)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Base Year 5"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "true"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.46"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success destination cost within non-peak period", func() {
		suite.setUpDomesticDestinationData()
		isPPM := false
		nonPeakDate := peakStart.addDate(0, -1)
		cost, displayParams, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
			isPPM,
		)
		expectedCost := unit.Cents(4884)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: "Base Year 5"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "false"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.27"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("successfully finds dom destination price for ppm with weight < 500 lbs with Price method", func() {
		suite.setUpDomesticDestinationData()
		suite.setupDomesticDestinationServiceItems()
		isPPM := true
		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// the PPM price for weights < 500 should be prorated from a base of 500
		basePriceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, unit.Pound(500), ddpTestServiceArea, isPPM)
		suite.NoError(err)

		halfPriceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, unit.Pound(250), ddpTestServiceArea, isPPM)
		suite.NoError(err)
		suite.Equal(basePriceCents/2, halfPriceCents)

		fifthPriceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, unit.Pound(100), ddpTestServiceArea, isPPM)
		suite.NoError(err)
		suite.Equal(basePriceCents/5, fifthPriceCents)
	})

	suite.Run("successfully finds dom destination price for ppm with weight < 500 lbs with PriceUsingParams method", func() {
		suite.setUpDomesticDestinationData()
		paymentServiceItem := suite.setupDomesticDestinationServiceItems()
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
			{Key: models.ServiceItemParamNameContractYearName, Value: "Base Year 5"},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: "1.04070"},
			{Key: models.ServiceItemParamNameIsPeak, Value: "true"},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: "1.46"},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("failure if contract code bogus", func() {
		suite.setUpDomesticDestinationData()
		isPPM := false
		_, _, err := pricer.Price(
			suite.AppContextForTest(),
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
			isPPM,
		)

		suite.Error(err)
		suite.Equal("could not lookup Domestic Service Area Price: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.Run("failure if move date is outside of contract year", func() {
		suite.setUpDomesticDestinationData()
		isPPM := false
		_, _, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
			isPPM,
		)

		suite.Error(err)
		suite.Contains(err.Error(), "could not calculate escalated price")

	})

	suite.Run("fail when is weight below minimum and shipment isn't a PPM", func() {
		suite.setUpDomesticDestinationData()
		isPPM := false
		cost, _, err := pricer.Price(
			suite.AppContextForTest(),
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unit.Pound(499),
			ddpTestServiceArea,
			isPPM,
		)
		suite.Equal(unit.Cents(0), cost)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})

	suite.Run("validation errors", func() {
		suite.setUpDomesticDestinationData()
		isPPM := false
		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No contract code
		_, _, err := pricer.Price(suite.AppContextForTest(), "", requestedPickupDate, ddpTestWeight, ddpTestServiceArea, isPPM)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No reference date
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, time.Time{}, ddpTestWeight, ddpTestServiceArea, isPPM)
		suite.Error(err)
		suite.Equal("ReferenceDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, 0, ddpTestServiceArea, isPPM)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())

		// No service area
		_, _, err = pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, requestedPickupDate, ddpTestWeight, "", isPPM)
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticDestinationServiceItems() models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDP,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   ddpTestServiceArea,
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(ddpTestWeight),
			},
		}, nil, nil,
	)
}

func (suite *GHCRateEngineServiceSuite) setUpDomesticDestinationData() {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Escalation:           1.0197,
				EscalationCompounded: 1.0407,
				Name:                 "Base Year 5",
			},
		})

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: ddpTestServiceArea,
			},
		})

	domesticDestinationService := factory.BuildReService(suite.DB(), []factory.Customization{
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDP,
				Name: "Dom. Destination Price",
			},
		},
	}, nil)
	domesticDestinationPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		DomesticServiceAreaID: serviceArea.ID,
		IsPeakPeriod:          true,
		ServiceID:             domesticDestinationService.ID,
	}

	domesticDestinationPeakPrice := domesticDestinationPrice
	domesticDestinationPeakPrice.PriceCents = 146
	suite.MustSave(&domesticDestinationPeakPrice)

	domesticDestinationNonPeakPrice := domesticDestinationPrice
	domesticDestinationNonPeakPrice.IsPeakPeriod = false
	domesticDestinationNonPeakPrice.PriceCents = 127
	suite.MustSave(&domesticDestinationNonPeakPrice)
}
