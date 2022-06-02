package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPackUnpack() {
	suite.Run("golden path with DNPK", func() {
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		isPPM := false
		priceCents, displayParams, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM)
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

	suite.Run("Invalid parameters to Price", func() {
		isPPM := false
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pack/unpack code")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, "", dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, time.Time{}, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, unit.Pound(250), dnpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, 0, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "Services schedule is required")
	})

	suite.Run("not finding domestic other price", func() {
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		badContractCode := "BOGUS"
		isPPM := false
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, badContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup domestic other price")
	})

	suite.Run("not finding contract year", func() {
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		twoYearsLaterPickupDate := dnpkTestRequestedPickupDate.AddDate(2, 0, 0)
		isPPM := false
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup contract year")
	})

	suite.Run("not finding shipment type price", func() {
		badMarket := models.MarketOconus
		isPPM := false
		suite.setupDomesticNTSPackPrices(dnpkTestServicesScheduleOrigin, dnpkTestIsPeakPeriod, dnpkTestBasePriceCents, badMarket, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDNPK, testdatagen.DefaultContractCode, dnpkTestRequestedPickupDate, dnpkTestWeight, dnpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup shipment type price")
	})

}

func (suite *GHCRateEngineServiceSuite) Test_domesticPackAndUnpackWithPPM() {
	suite.Run("golden path with DPK", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)

		isPPM := true
		priceCents, displayParams, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM)
		suite.NoError(err)
		suite.Equal(dpkTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dpkTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		isPPM := true
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)

		_, _, err := priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pack/unpack code")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, "", dnpkTestRequestedPickupDate, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, time.Time{}, dpkTestWeight, dpkTestServicesScheduleOrigin, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceDomesticPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDPK, testdatagen.DefaultContractCode, dpkTestRequestedPickupDate, dpkTestWeight, 0, isPPM)
		suite.Error(err)
		suite.Contains(err.Error(), "Services schedule is required")
	})
}
func (suite *GHCRateEngineServiceSuite) Test_priceDomesticFirstDaySIT() {
	suite.Run("destination golden path", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.NoError(err)
		suite.Equal(ddfsitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ddfsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ddfsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ddfsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ddfsitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeCS, DefaultContractCode, ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported first day sit code")
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		badWeight := unit.Pound(250)
		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddfsitTestRequestedPickupDate, badWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 250 less than the minimum")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, "BOGUS", ddfsitTestRequestedPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination first day SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDFSIT, ddfsitTestServiceArea, ddfsitTestIsPeakPeriod, ddfsitTestBasePriceCents, ddfsitTestContractYearName, ddfsitTestEscalationCompounded)

		twoYearsLaterPickupDate := ddfsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, DefaultContractCode, twoYearsLaterPickupDate, ddfsitTestWeight, ddfsitTestServiceArea)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticAdditionalDaysSIT() {

	suite.Run("destination golden path", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.NoError(err)
		suite.Equal(ddasitTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ddasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ddasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ddasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ddasitTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDFSIT, DefaultContractCode, ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported additional day sit code")
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		badWeight := unit.Pound(499)
		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, DefaultContractCode, ddasitTestRequestedPickupDate, badWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "weight of 499 less than the minimum")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, "BOGUS", ddasitTestRequestedPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination additional days SIT rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDDASIT, ddasitTestServiceArea, ddasitTestIsPeakPeriod, ddasitTestBasePriceCents, ddasitTestContractYearName, ddasitTestEscalationCompounded)

		twoYearsLaterPickupDate := ddasitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticAdditionalDaysSIT(suite.AppContextForTest(), models.ReServiceCodeDDASIT, DefaultContractCode, twoYearsLaterPickupDate, ddasitTestWeight, ddasitTestServiceArea, ddasitTestNumberOfDaysInSIT)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySITSameZip3s() {
	dshZipDest := "30907"
	dshZipSITDest := "30901" // same zip3
	dshDistance := unit.Miles(15)
	dshContractName := "dshTestYear"

	suite.Run("destination golden path for same zip3s", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDSH, dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestDomesticServiceAreaBasePriceCents, dshContractName, dddsitTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(53187) // dddsitTestDomesticServiceAreaBasePriceCents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dshContractName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dddsitTestDomesticServiceAreaBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pickup/delivery SIT code")
	})

	suite.Run("invalid weight", func() {
		badWeight := unit.Pound(250)
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, badWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		expectedError := fmt.Sprintf("weight of %d less than the minimum", badWeight)
		suite.Contains(err.Error(), expectedError)
	})

	suite.Run("bad destination zip", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, "309", dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid destination postal code")
	})

	suite.Run("bad SIT final destination zip", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, "1234", dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT final destination postal code")
	})

	suite.Run("error from shorthaul pricer", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dshZipDest, dshZipSITDest, dshDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price shorthaul")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50PlusMilesDiffZip3s() {
	dlhZipDest := "30907"
	dlhZipSITDest := "36106"       // different zip3
	dlhDistance := unit.Miles(305) // > 50 miles
	dlhContractName := "dhlTestYear"

	suite.Run("destination golden path for > 50 miles with different zip3s", func() {
		suite.setupDomesticLinehaulPrice(dddsitTestServiceArea, dddsitTestIsPeakPeriod, dddsitTestWeightLower, dddsitTestWeightUpper, dddsitTestMilesLower, dddsitTestMilesUpper, dddsitTestDomesticLinehaulBasePriceMillicents, dlhContractName, dddsitTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.NoError(err)
		expectedPriceMillicents := unit.Millicents(45944438) // dddsitTestDomesticLinehaulBasePriceMillicents * (dddsitTestWeight / 100) * distance * dddsitTestEscalationCompounded
		expectedPrice := expectedPriceMillicents.ToCents()

		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: dlhContractName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatFloat(dddsitTestDomesticLinehaulBasePriceMillicents.ToDollarFloatNoRound(), 3)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("error from linehaul pricer", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, dlhZipDest, dlhZipSITDest, dlhDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not price linehaul")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticPickupDeliverySIT50MilesOrLessDiffZip3s() {
	domOtherZipDest := "30907"
	domOtherZipSITDest := "29801"      // different zip3
	domOtherDistance := unit.Miles(37) // <= 50 miles
	domContractName := "domTestYear"

	suite.Run("destination golden path for <= 50 miles with different zip3s", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, domContractName, dddsitTestEscalationCompounded)
		priceCents, displayParams, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)
		expectedPrice := unit.Cents(58355) // dddsitTestDomesticOtherBasePriceCents * (dddsitTestWeight / 100) * dddsitTestEscalationCompounded
		suite.Equal(expectedPrice, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: domContractName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dddsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(dddsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dddsitTestDomesticOtherBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, "BOGUS", dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch domestic destination SIT delivery rate")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDDDSIT, dddsitTestSchedule, dddsitTestIsPeakPeriod, dddsitTestDomesticOtherBasePriceCents, domContractName, dddsitTestEscalationCompounded)
		_, _, err := priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, dddsitTestRequestedPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.NoError(err)

		twoYearsLaterPickupDate := dddsitTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err = priceDomesticPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, dddsitTestWeight, dddsitTestServiceArea, dddsitTestSchedule, domOtherZipDest, domOtherZipSITDest, domOtherDistance)
		suite.Error(err)
		suite.Contains(err.Error(), "could not fetch contract year")
	})
}

type pricerParamsSubtestData struct {
	params             services.PricingDisplayParams
	paymentServiceItem models.PaymentServiceItem
}

func (suite *GHCRateEngineServiceSuite) makePricerParamsSubtestData() (subtestData *pricerParamsSubtestData) {
	subtestData = &pricerParamsSubtestData{}
	subtestData.params = services.PricingDisplayParams{
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

	subtestData.paymentServiceItem = testdatagen.MakePaymentServiceItem(
		suite.DB(),
		testdatagen.Assertions{},
	)

	return subtestData
}

func (suite *GHCRateEngineServiceSuite) Test_createPricerGeneratedParams() {
	suite.Run("payment service item params created for the pricer", func() {
		subtestData := suite.makePricerParamsSubtestData()
		paymentServiceItemParams, err := createPricerGeneratedParams(suite.AppContextForTest(), subtestData.paymentServiceItem.ID, subtestData.params)
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

	suite.Run("errors if PaymentServiceItemID is invalid", func() {
		subtestData := suite.makePricerParamsSubtestData()
		invalidID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")

		_, err := createPricerGeneratedParams(suite.AppContextForTest(), invalidID, subtestData.params)
		suite.Error(err)
		suite.Contains(err.Error(), "validation error with creating payment service item param")
	})

	suite.Run("errors if PricingParm points to a serviceItem that doesnt originate from the Pricer", func() {
		subtestData := suite.makePricerParamsSubtestData()
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

		_, err := createPricerGeneratedParams(suite.AppContextForTest(), subtestData.paymentServiceItem.ID, invalidParam)
		suite.Error(err)
		suite.Contains(err.Error(), "Service item param key is not a pricer param")
	})

	suite.Run("errors if no PricingParms passed from the Pricer", func() {
		subtestData := suite.makePricerParamsSubtestData()
		emptyParams := services.PricingDisplayParams{}

		_, err := createPricerGeneratedParams(suite.AppContextForTest(), subtestData.paymentServiceItem.ID, emptyParams)
		suite.Error(err)
		suite.Contains(err.Error(), "PricingDisplayParams must not be empty")
	})
}

func (suite *GHCRateEngineServiceSuite) Test_priceDomesticShuttling() {
	suite.Run("destination golden path", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, DefaultContractCode, doshutTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDOSHUT, DefaultContractCode, doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(doshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(doshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(doshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, DefaultContractCode, doshutTestEscalationCompounded)
		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeCS, DefaultContractCode, doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported domestic shuttling code")
	})

	suite.Run("invalid weight", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, DefaultContractCode, doshutTestEscalationCompounded)

		badWeight := unit.Pound(250)
		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDOSHUT, DefaultContractCode, doshutTestRequestedPickupDate, badWeight, doshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDOSHUT, doshutTestServiceSchedule, doshutTestBasePriceCents, DefaultContractCode, doshutTestEscalationCompounded)

		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDOSHUT, "BOGUS", doshutTestRequestedPickupDate, doshutTestWeight, doshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, DefaultContractCode, ddshutTestEscalationCompounded)

		twoYearsLaterPickupDate := doshutTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticShuttling(suite.AppContextForTest(), models.ReServiceCodeDDSHUT, DefaultContractCode, twoYearsLaterPickupDate, ddshutTestWeight, ddshutTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "Could not lookup contract year")
	})
}
func (suite *GHCRateEngineServiceSuite) Test_priceDomesticCrating() {
	suite.Run("crating golden path", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, DefaultContractCode, dcrtTestEscalationCompounded)

		priceCents, displayParams, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, DefaultContractCode, dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(dcrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dcrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dcrtTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, DefaultContractCode, dcrtTestEscalationCompounded)
		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeCS, DefaultContractCode, dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported domestic crating code")
	})

	suite.Run("invalid crate size", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, DefaultContractCode, dcrtTestEscalationCompounded)

		badSize := unit.CubicFeet(1.0)
		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, DefaultContractCode, dcrtTestRequestedPickupDate, badSize, dcrtTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "crate must be billed for a minimum of 4 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, DefaultContractCode, dcrtTestEscalationCompounded)

		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, "BOGUS", dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDCRT, dcrtTestServiceSchedule, dcrtTestBasePriceCents, DefaultContractCode, dcrtTestEscalationCompounded)

		twoYearsLaterPickupDate := dcrtTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := priceDomesticCrating(suite.AppContextForTest(), models.ReServiceCodeDCRT, DefaultContractCode, twoYearsLaterPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}
