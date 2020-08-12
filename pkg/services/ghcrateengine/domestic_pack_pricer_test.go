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
	servicesScheduleOrigin = 1
	weightBilledActual     = 3600
)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticPackWithServiceItemParamsBadData() {
	suite.setUpDomesticPackData(models.ReServiceCodeDPK)
	paymentServiceItem := suite.setupPaymentServiceItemWithParams(
		models.ReServiceCodeDPK,
		[]createParams{
			{
				models.ServiceItemParamNameContractCode,
				models.ServiceItemParamTypeString,
				testdatagen.DefaultContractCode,
			},
			{
				models.ServiceItemParamNameRequestedPickupDate,
				models.ServiceItemParamTypeDate,
				time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
			},
			{
				models.ServiceItemParamNameWeightBilledActual,
				models.ServiceItemParamTypeInteger,
				"0",
			},
			{
				models.ServiceItemParamNameServicesScheduleOrigin,
				models.ServiceItemParamTypeInteger,
				strconv.Itoa(servicesScheduleOrigin),
			},
		},
	)

	pricer := NewDomesticPackPricer(suite.DB())

	suite.T().Run("failure during pricing bubbles up", func(t *testing.T) {
		_, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticPackWithServiceItemParams() {
	suite.setUpDomesticPackData(models.ReServiceCodeDPK)
	paymentServiceItem := suite.setupDomesticPackServiceItems()

	pricer := NewDomesticPackPricer(suite.DB())

	suite.T().Run("success all params for domestic pack available", func(t *testing.T) {
		cost, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.T().Run("validation errors", func(t *testing.T) {
		// No contract code
		_, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
		suite.Equal("could not find param with key ContractCode", err.Error())

		// No requested pickup date
		missingRequestedPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameRequestedPickupDate)
		_, err = pricer.PriceUsingParams(missingRequestedPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key RequestedPickupDate", err.Error())

		// No weight
		missingBilledActualWeight := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilledActual)
		_, err = pricer.PriceUsingParams(missingBilledActualWeight)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilledActual", err.Error())

		// No service schedule origin
		missingServicesScheduleOrigin := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameServicesScheduleOrigin)
		_, err = pricer.PriceUsingParams(missingServicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("could not find param with key ServicesScheduleOrigin", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticPack() {
	suite.setUpDomesticPackData(models.ReServiceCodeDPK)

	pricer := NewDomesticPackPricer(suite.DB())

	suite.T().Run("success domestic pack cost within peak period", func(t *testing.T) {
		cost, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)
		expectedCost := unit.Cents(5470)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.T().Run("success domestic pack cost within non-peak period", func(t *testing.T) {
		nonPeakDate := peakStart.addDate(0, -1)
		cost, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, nonPeakDate.month, nonPeakDate.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)
		expectedCost := unit.Cents(4758)
		suite.NoError(err)
		suite.Equal(expectedCost, cost)
	})

	suite.T().Run("failure if contract code bogus", func(t *testing.T) {
		_, err := pricer.Price(
			"bogus_code",
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)

		suite.Error(err)
		suite.Equal("Could not lookup Domestic Other Price: sql: no rows in result set", err.Error())
	})

	suite.T().Run("failure if move date is outside of contract year", func(t *testing.T) {
		_, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear+1, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			weightBilledActual,
			servicesScheduleOrigin,
		)

		suite.Error(err)
		suite.Equal("Could not lookup contract year: sql: no rows in result set", err.Error())
	})

	suite.T().Run("weight below minimum", func(t *testing.T) {
		cost, err := pricer.Price(
			testdatagen.DefaultContractCode,
			time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC),
			unit.Pound(499),
			servicesScheduleOrigin,
		)
		suite.Equal(unit.Cents(0), cost)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())
	})

	suite.T().Run("validation errors", func(t *testing.T) {
		requestedPickupDate := time.Date(testdatagen.TestYear, time.July, 4, 0, 0, 0, 0, time.UTC)

		// No contract code
		_, err := pricer.Price("", requestedPickupDate, weightBilledActual, servicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("ContractCode is required", err.Error())

		// No requested pickup date
		_, err = pricer.Price(testdatagen.DefaultContractCode, time.Time{}, weightBilledActual, servicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("RequestedPickupDate is required", err.Error())

		// No weight
		_, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, 0, servicesScheduleOrigin)
		suite.Error(err)
		suite.Equal("Weight must be a minimum of 500", err.Error())

		// No service area
		_, err = pricer.Price(testdatagen.DefaultContractCode, requestedPickupDate, dshTestWeight, 0)
		suite.Error(err)
		suite.Equal("Service schedule is required", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticPackServiceItems() models.PaymentServiceItem {
	return suite.setupPaymentServiceItemWithParams(
		models.ReServiceCodeDPK,
		[]createParams{
			{
				models.ServiceItemParamNameContractCode,
				models.ServiceItemParamTypeString,
				testdatagen.DefaultContractCode,
			},
			{
				models.ServiceItemParamNameRequestedPickupDate,
				models.ServiceItemParamTypeDate,
				time.Date(testdatagen.TestYear, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC).Format(DateParamFormat),
			},
			{
				models.ServiceItemParamNameWeightBilledActual,
				models.ServiceItemParamTypeInteger,
				strconv.Itoa(weightBilledActual),
			},
			{
				models.ServiceItemParamNameServicesScheduleOrigin,
				models.ServiceItemParamTypeInteger,
				strconv.Itoa(servicesScheduleOrigin),
			},
		},
	)
}
