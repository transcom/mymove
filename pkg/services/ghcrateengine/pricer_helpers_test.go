package ghcrateengine

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticFirstDaySIT() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestEscalationCompounded)

	suite.T().Run("destination golden path", func(t *testing.T) {
		priceCents, _, err := priceDomesticFirstDaySIT(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.NoError(err)
		suite.Equal(ddfsitTestPriceCents, priceCents)
	})

	suite.T().Run("invalid service code", func(t *testing.T) {
		_, _, err := priceDomesticFirstDaySIT(suite.DB(), models.ReServiceCodeCS, DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported first day sit code")
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(250)
		_, _, err := priceDomesticFirstDaySIT(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddfsitTestRequestedPickupDate, badWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, _, err := priceDomesticFirstDaySIT(suite.DB(), models.ReServiceCodeDDFSIT, "BOGUS", ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination first day SIT rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := ddfsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticFirstDaySIT(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, twoYearsLaterPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticAdditionalDaysSIT() {
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestEscalationCompounded)

	suite.T().Run("destination golden path", func(t *testing.T) {
		priceCents, _, err := priceDomesticAdditionalDaysSIT(suite.DB(), models.ReServiceCodeDDASIT, DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)
	})

	suite.T().Run("invalid service code", func(t *testing.T) {
		_, _, err := priceDomesticAdditionalDaysSIT(suite.DB(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported additional day sit code")
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(499)
		_, _, err := priceDomesticAdditionalDaysSIT(suite.DB(), models.ReServiceCodeDDASIT, DefaultContractCode, ddasitTestRequestedPickupDate, badWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 499 less than the minimum")
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, _, err := priceDomesticAdditionalDaysSIT(suite.DB(), models.ReServiceCodeDDASIT, "BOGUS", ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination additional days SIT rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := ddasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticAdditionalDaysSIT(suite.DB(), models.ReServiceCodeDDASIT, DefaultContractCode, twoYearsLaterPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySITSameZip3s() {
	dshZipDest := "30907"
	dshZipSITDest := "30901" // same zip3
	dshDistance := unit.Miles(15)

	suite.T().Run("destination golden path for same zip3s", func(t *testing.T) {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestDomesticServiceAreaBasePriceCents, dddsitTestEscalationCompounded)
		priceCents, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(53187) // dddsitTestDomesticServiceAreaBasePriceCents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("invalid service code", func(t *testing.T) {
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pickup/delivery SIT code")
	})

	suite.T().Run("invalid weight", func(t *testing.T) {
		badWeight := unit.Pound(250)
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, badWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.T().Run("bad destination zip", func(t *testing.T) {
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, "309", dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid destination postal code")
	})

	suite.T().Run("bad SIT final destination zip", func(t *testing.T) {
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, "1234", dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT final destination postal code")
	})

	suite.T().Run("error from shorthaul pricer", func(t *testing.T) {
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price shorthaul")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50PlusMilesDiffZip3s() {
	dlhZipDest := "30907"
	dlhZipSITDest := "36106"       // different zip3
	dlhDistance := unit.Miles(305) // > 50 miles

	suite.T().Run("destination golden path for > 50 miles with different zip3s", func(t *testing.T) {
		suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestWeightLower, dddsitTestWeightUpper, dddsitTestMilesLower, dddsitTestMilesUpper, dddsitTestDomesticLinehaulBasePriceMillicents, dddsitTestContractYearName, dddsitTestEscalationCompounded)
		priceCents, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.NoError(err)
		expectedPriceMillicents := unit.Millicents(45944438) // dddsitTestDomesticLinehaulBasePriceMillicents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
		expectedPrice := expectedPriceMillicents.ToCents()
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("error from linehaul pricer", func(t *testing.T) {
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50MilesOrLessDiffZip3s() {
	domOtherZipDest := "30907"
	domOtherZipSITDest := "29801"      // different zip3
	domOtherDistance := unit.Miles(37) // <= 50 miles

	suite.T().Run("destination golden path for <= 50 miles with different zip3s", func(t *testing.T) {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, dddsitTestEscalationCompounded)
		priceCents, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(58355) // dddsitTestDomesticOtherBasePriceCents * (dddsitTestWeight / 100) * dddsitTestEscalationCompounded
		suite.Equal(expectedPrice, priceCents)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.T().Run("not finding a contract year record", func(t *testing.T) {
		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticPickupDeliverySIT(suite.DB(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_createPricerGeneratedParams() {
	params := services.PricingDisplayParams{
		{
			Key:   models.ServiceItemParamNamePriceRateOrFactor,
			Value: "4000.90",
		}, {
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: "1.06",
		}, {
			Key:   models.ServiceItemParamNameIsPeak,
			Value: "True",
		}, {
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "TRUSS_TEST",
		},
	}

	testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNamePriceRateOrFactor,
			Description: "Price Rate Or Factor",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginPricer,
		},
	})
	testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameEscalationCompounded,
			Description: "Escalation compounded",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginPricer,
		},
	})
	testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameIsPeak,
			Description: "Is peak",
			Type:        models.ServiceItemParamTypeBoolean,
			Origin:      models.ServiceItemParamOriginPricer,
		},
	})
	testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameContractYearName,
			Description: "Contract year name",
			Type:        models.ServiceItemParamTypeString,
			Origin:      models.ServiceItemParamOriginPricer,
		},
	})

	paymentServiceItem := testdatagen.MakePaymentServiceItem(
		suite.DB(),
		testdatagen.Assertions{},
	)

	suite.T().Run("payment service item params created for the pricer", func(t *testing.T) {
		paymentServiceItemParams, err := createPricerGeneratedParams(suite.DB(), paymentServiceItem.ID, params)
		suite.NoError(err)
		expectedValues := [4]string{"4000.90", "1.06", "True", "TRUSS_TEST"}
		for _, paymentServiceItemParam := range paymentServiceItemParams {
			switch paymentServiceItemParam.ServiceItemParamKey.Key {
			case models.ServiceItemParamNamePriceRateOrFactor:
				suite.Equal(expectedValues[0], paymentServiceItemParam.Value)
			case models.ServiceItemParamNameEscalationCompounded:
				suite.Equal(expectedValues[1], paymentServiceItemParam.Value)
			case models.ServiceItemParamNameIsPeak:
				suite.Equal(expectedValues[2], paymentServiceItemParam.Value)
			case models.ServiceItemParamNameContractYearName:
				suite.Equal(expectedValues[3], paymentServiceItemParam.Value)
			}
		}
	})

	suite.T().Run("errors if PaymentServiceItemID is invalid", func(t *testing.T) {
		invalidID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

		_, err := createPricerGeneratedParams(suite.DB(), invalidID, params)
		suite.Error(err)
		suite.Contains(err.Error(), "validation error with creating payment service item param")
	})

	suite.T().Run("errors if PricingParm points to a serviceItem that doesnt originate from the Pricer", func(t *testing.T) {
		invalidParam := services.PricingDisplayParams{
			{
				Key:   models.ServiceItemParamNameServiceAreaOrigin,
				Value: "40000.9",
			},
		}

		testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameServiceAreaOrigin,
				Description: "service area actual",
				Type:        models.ServiceItemParamTypeString,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		})

		_, err := createPricerGeneratedParams(suite.DB(), paymentServiceItem.ID, invalidParam)
		suite.Error(err)
		suite.Contains(err.Error(), "Service item param key is not a pricer param")
	})

	suite.T().Run("errors if no PricingParms passed from the Pricer", func(t *testing.T) {
		emptyParams := services.PricingDisplayParams{}

		_, err := createPricerGeneratedParams(suite.DB(), paymentServiceItem.ID, emptyParams)
		suite.Error(err)
		suite.Contains(err.Error(), "PricingDisplayParams must not be empty")
	})
}
