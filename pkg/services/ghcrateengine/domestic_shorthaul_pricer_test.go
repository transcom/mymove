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
	dshTestServiceArea = "006"
	dshTestWeight      = 3600
	dshTestMileage     = 1200
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticShorthaulWithServiceItemParamsBadData() {
	requestedPickup := time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat)
	suite.setUpDomesticShorthaulData()
	paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDSH,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   requestedPickup,
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip5,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dshTestMileage),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   "0",
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dshTestServiceArea,
			},
		},
	)

	pricer := NewDomesticShorthaulPricer(suite.DB())

	suite.Run("failure during pricing bubbles up", func() {
		_, rateEngineParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
		suite.Nil(rateEngineParams)
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticShorthaulWithServiceItemParams() {
	requestedPickup := time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat)
	suite.setUpDomesticShorthaulData()
	paymentServiceItem := suite.setupDomesticShorthaulServiceItems(requestedPickup)
	expectedPricingCreatedParams := suite.getExpectedDSHPricerCreatedParamsFromDBGivenParams(dshTestServiceArea, requestedPickup)

	pricer := NewDomesticShorthaulPricer(suite.DB())

	suite.Run("success all params for shorthaul available", func() {
		cost, rateEngineParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		expectedCost := unit.Cents(6563903)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		suite.validatePricerCreatedParams(expectedPricingCreatedParams, rateEngineParams)
	})

	suite.Run("validation errors", func() {

		// No contract code
		_, rateEngineParams, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
		suite.Equal("could not find param with key ContractCode", err.Error())
		suite.Nil(rateEngineParams)

		// No requested pickup date
		missingRequestedPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameRequestedPickupDate)
		_, rateEngineParams, err = pricer.PriceUsingParams(missingRequestedPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key RequestedPickupDate", err.Error())
		suite.Nil(rateEngineParams)

		// No distance
		missingDistanceZip5 := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameDistanceZip5)
		_, rateEngineParams, err = pricer.PriceUsingParams(missingDistanceZip5)
		suite.Error(err)
		suite.Equal("could not find param with key DistanceZip5", err.Error())
		suite.Nil(rateEngineParams)

		// No weight
		missingBilledActualWeight := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilledActual)
		_, rateEngineParams, err = pricer.PriceUsingParams(missingBilledActualWeight)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilledActual", err.Error())
		suite.Nil(rateEngineParams)

		// No service area
		missingServiceAreaOrigin := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServiceAreaOrigin)
		_, rateEngineParams, err = pricer.PriceUsingParams(missingServiceAreaOrigin)
		suite.Error(err)
		suite.Equal("could not find param with key ServiceAreaOrigin", err.Error())
		suite.Nil(rateEngineParams)
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticShorthaul() {
	suite.Run("success shorthaul cost within peak period", func() {
		requestedPickup := time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat)
		suite.setUpDomesticShorthaulData()

		pricer := NewDomesticShorthaulPricer(suite.DB())

		newRequestedPickup := time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)
		newExpectedPricingCreatedParams := suite.getExpectedDSHPricerCreatedParamsFromDBGivenParams(dshTestServiceArea, requestedPickup)
		cost, rateEngineParams, err := pricer.Price(
			testdatagen.DefaultContractCode,
			newRequestedPickup,
			dshTestMileage,
			dshTestWeight,
			dshTestServiceArea,
		)
		expectedCost := unit.Cents(6563903)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
		suite.validatePricerCreatedParams(newExpectedPricingCreatedParams, rateEngineParams)
	})

	suite.Run("success shorthaul cost within non-peak period", func() {
		suite.setUpDomesticShorthaulData()

		pricer := NewDomesticShorthaulPricer(suite.DB())

		nonPeakDate := peakStart.addDate(0, -1)
		newRequestedPickup := time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC)
		newExpectedPricingCreatedParams := suite.getExpectedDSHPricerCreatedParamsFromDBGivenParams(dshTestServiceArea, newRequestedPickup.Format(DateParamFormat))

		cost, rateEngineParams, err := pricer.Price(
			testdatagen.DefaultContractCode,
			newRequestedPickup,
			dshTestMileage,
			dshTestWeight,
			dshTestServiceArea,
		)
		expectedCost := unit.Cents(5709696)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
		suite.validatePricerCreatedParams(newExpectedPricingCreatedParams, rateEngineParams)
	})

	suite.Run("failure if contract code bogus", func() {
		suite.setUpDomesticShorthaulData()
		pricer := NewDomesticShorthaulPricer(suite.DB())

		_, rateEngineParams, err := pricer.Price(
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dshTestMileage,
			dshTestWeight,
			dshTestServiceArea,
		)

		suite.Error(err)
		suite.Equal("Could not lookup Domestic Service Area Price: "+models.RecordNotFoundErrorString, err.Error())
		suite.Nil(rateEngineParams)
	})

	suite.Run("failure if move date is outside of contract year", func() {
		suite.setUpDomesticShorthaulData()
		pricer := NewDomesticShorthaulPricer(suite.DB())

		_, rateEngineParams, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dshTestMileage,
			dshTestWeight,
			dshTestServiceArea,
		)

		suite.Error(err)
		suite.Equal("Could not lookup contract year: "+models.RecordNotFoundErrorString, err.Error())
		suite.Nil(rateEngineParams)
	})

	suite.Run("weight below minimum", func() {
		suite.setUpDomesticShorthaulData()
		pricer := NewDomesticShorthaulPricer(suite.DB())

		cost, rateEngineParams, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			dshTestMileage,
			unit.Pound(499),
			dshTestServiceArea,
		)
		suite.Equal(unit.Cents(0), cost)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
		suite.Nil(rateEngineParams)
	})

	suite.Run("validation errors", func() {
		suite.setUpDomesticShorthaulData()
		pricer := NewDomesticShorthaulPricer(suite.DB())

		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No contract code
		_, rateEngineParams, err := pricer.Price("", requestedPickupDate, dshTestMileage, dshTestWeight, dshTestServiceArea)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())
		suite.Nil(rateEngineParams)

		// No requested pickup date
		_, rateEngineParams, err = pricer.Price(testdatagen.DefaultContractCode, time.Time{}, dshTestMileage, dshTestWeight, dshTestServiceArea)
		suite.Error(err)
		suite.Equal("RequestedPickupDate is required", err.Error())
		suite.Nil(rateEngineParams)

		// No distance
		_, rateEngineParams, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, 0, dshTestWeight, dshTestServiceArea)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())
		suite.Nil(rateEngineParams)

		// No weight
		_, rateEngineParams, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, dshTestMileage, 0, dshTestServiceArea)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
		suite.Nil(rateEngineParams)

		// No service area
		_, rateEngineParams, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, dshTestMileage, dshTestWeight, "")
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
		suite.Nil(rateEngineParams)
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticShorthaulServiceItems(requestedPickup string) models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDSH,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   requestedPickup,
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip5,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dshTestMileage),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dshTestWeight),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaOrigin,
				KeyType: models.ServiceItemParamTypeString,
				Value:   dshTestServiceArea,
			},
		},
	)
}

func (suite *GHCRateEngineServiceSuite) removeOnePaymentServiceItem(paymentServiceItemParams models.PaymentServiceItemParams, nameToRemove models.ServiceItemParamName) models.PaymentServiceItemParams {
	var params models.PaymentServiceItemParams
	for _, param := range paymentServiceItemParams {
		if param.ServiceItemParamKey.Key != nameToRemove {
			params = append(params, param)
		}
	}
	return params
}

func (suite *GHCRateEngineServiceSuite) setUpDomesticShorthaulData() {
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
				ServiceArea: dshTestServiceArea,
			},
		})

	domesticShorthaulService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DSH",
				Name: "Dom. Shorthaul",
			},
		})

	domesticShorthaulPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		DomesticServiceAreaID: serviceArea.ID,
		IsPeakPeriod:          true,
		ServiceID:             domesticShorthaulService.ID,
	}

	domesticShorthaulPeakPrice := domesticShorthaulPrice
	domesticShorthaulPeakPrice.PriceCents = 146
	suite.MustSave(&domesticShorthaulPeakPrice)

	domesticShorthaulNonpeakPrice := domesticShorthaulPrice
	domesticShorthaulNonpeakPrice.IsPeakPeriod = false
	domesticShorthaulNonpeakPrice.PriceCents = 127
	suite.MustSave(&domesticShorthaulNonpeakPrice)
}

func (suite *GHCRateEngineServiceSuite) getExpectedDSHPricerCreatedParamsFromDBGivenParams(serviceArea string, requestedPickUp string) services.PricingDisplayParams {
	var err error

	var requestedPickUpDate time.Time
	requestedPickUpDate, err = time.Parse(DateParamFormat, requestedPickUp)
	suite.NoError(err)

	isPeakPeriod := IsPeakPeriod(requestedPickUpDate)

	var domServiceAreaPrice models.ReDomesticServiceAreaPrice
	domServiceAreaPrice, err = fetchDomServiceAreaPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDSH, serviceArea, isPeakPeriod)
	suite.NoError(err)

	var contractYear models.ReContractYear
	contractYear, err = fetchContractYear(suite.DB(), domServiceAreaPrice.ContractID, requestedPickUpDate)
	suite.NoError(err)

	var pricingRateEngineParams = services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: contractYear.Name,
		},
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: FormatCents(domServiceAreaPrice.PriceCents),
		},
		{
			Key:   models.ServiceItemParamNameIsPeak,
			Value: strconv.FormatBool(isPeakPeriod),
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: FormatEscalation(contractYear.EscalationCompounded),
		},
	}

	return pricingRateEngineParams
}
