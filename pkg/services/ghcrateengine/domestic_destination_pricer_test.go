package ghcrateengine

import (
	"strconv"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const (
	ddpTestServiceArea = "006"
	ddpTestWeight      = 3600
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticDestinationWithServiceItemParamsBadData() {
	suite.setUpDomesticDestinationData()
	paymentServiceItem := testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDP,
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
				Key:     models.ServiceItemParamNameServiceAreaDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   ddpTestServiceArea,
			},
		},
	)

	pricer := NewDomesticDestinationPricer(suite.DB())

	suite.T().Run("failure during pricing bubbles up", func(t *testing.T) {
		_, _, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticDestinationWithServiceItemParams() {
	suite.setUpDomesticDestinationData()
	paymentServiceItem := suite.setupDomesticDestinationServiceItems()

	pricer := NewDomesticDestinationPricer(suite.DB())

	suite.T().Run("success all params for destination available", func(t *testing.T) {
		cost, _, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.T().Run("validation errors", func(t *testing.T) {
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

		// No service area
		missingServiceAreaDest := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServiceAreaDest)
		_, _, err = pricer.PriceUsingParams(missingServiceAreaDest)
		suite.Error(err)
		suite.Equal("could not find param with key ServiceAreaDest", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticDestination() {
	suite.setUpDomesticDestinationData()

	pricer := NewDomesticDestinationPricer(suite.DB())

	suite.T().Run("success destination cost within peak period", func(t *testing.T) {
		cost, displayParams, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
		)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)

		suite.Equal(expectedCost, cost)
		if suite.Len(displayParams, 4) {
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNameContractYearName, "Base Year 5")
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNameEscalationCompounded, "1.04070")
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNameIsPeak, "true")
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNamePriceRateOrFactor, "1.46")
		}

	})

	suite.T().Run("success destination cost within non-peak period", func(t *testing.T) {
		nonPeakDate := peakStart.addDate(0, -1)
		cost, displayParams, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
		)
		expectedCost := unit.Cents(4758)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)

		if suite.Len(displayParams, 4) {
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNameContractYearName, "Base Year 5")
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNameEscalationCompounded, "1.04070")
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNameIsPeak, "false")
			suite.HasDisplayParam(displayParams, models.ServiceItemParamNamePriceRateOrFactor, "1.27")
		}
	})

	suite.T().Run("failure if contract code bogus", func(t *testing.T) {
		_, _, err := pricer.Price(
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
		)

		suite.Error(err)
		suite.Equal("Could not lookup Domestic Service Area Price: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.T().Run("failure if move date is outside of contract year", func(t *testing.T) {
		_, _, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			ddpTestWeight,
			ddpTestServiceArea,
		)

		suite.Error(err)
		suite.Equal("Could not lookup contract year: "+models.RecordNotFoundErrorString, err.Error())
	})

	suite.T().Run("weight below minimum", func(t *testing.T) {
		cost, _, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unit.Pound(499),
			ddpTestServiceArea,
		)
		suite.Equal(unit.Cents(0), cost)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})

	suite.T().Run("validation errors", func(t *testing.T) {
		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No contract code
		_, _, err := pricer.Price("", requestedPickupDate, ddpTestWeight, ddpTestServiceArea)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No requested pickup date
		_, _, err = pricer.Price(testdatagen.DefaultContractCode, time.Time{}, ddpTestWeight, ddpTestServiceArea)
		suite.Error(err)
		suite.Equal("RequestedPickupDate is required", err.Error())

		// No weight
		_, _, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, 0, ddpTestServiceArea)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())

		// No service area
		_, _, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, ddpTestWeight, "")
		suite.Error(err)
		suite.Equal("ServiceArea is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticDestinationServiceItems() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDDP,
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
				Value:   strconv.Itoa(ddpTestWeight),
			},
			{
				Key:     models.ServiceItemParamNameServiceAreaDest,
				KeyType: models.ServiceItemParamTypeString,
				Value:   ddpTestServiceArea,
			},
		},
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

	domesticDestinationService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DDP",
				Name: "Dom. Destination Price",
			},
		})

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
