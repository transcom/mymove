package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_priceInternationalShuttling() {
	suite.Run("origin golden path", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		priceCents, displayParams, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)
		suite.NoError(err)
		suite.Equal(ioshutTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ioshutTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ioshutTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)
		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported international shuttling code")
	})

	suite.Run("invalid weight", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		badWeight := unit.Pound(250)
		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, testdatagen.DefaultContractCode, ioshutTestRequestedPickupDate, badWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "Weight must be a minimum of 500")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, "BOGUS", ioshutTestRequestedPickupDate, ioshutTestWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIOSHUT, ioshutTestMarket, ioshutTestBasePriceCents, testdatagen.DefaultContractCode, ioshutTestEscalationCompounded)

		tenYearsLaterPickupDate := ioshutTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := priceInternationalShuttling(suite.AppContextForTest(), models.ReServiceCodeIOSHUT, testdatagen.DefaultContractCode, tenYearsLaterPickupDate, ioshutTestWeight, ioshutTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not calculate escalated price: could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlPackUnpack() {
	suite.Run("success with IHPK", func() {
		suite.setupIntlPackServiceItem()
		totalCost, displayParams, err := priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, testdatagen.DefaultContractCode, ihpkTestRequestedPickupDate, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(ihpkTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ihpkTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ihpkTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ihpkTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ihpkTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		suite.setupIntlPackServiceItem()
		_, _, err := priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, ihpkTestRequestedPickupDate, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported pack/unpack code")

		_, _, err = priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, "", ihpkTestRequestedPickupDate, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, testdatagen.DefaultContractCode, time.Time{}, ihpkTestWeight, ihpkTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceIntlPackUnpack(suite.AppContextForTest(), models.ReServiceCodeIHPK, testdatagen.DefaultContractCode, ihpkTestRequestedPickupDate, ihpkTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "PerUnitCents is required")
	})

}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlFirstDaySIT() {
	suite.Run("success with IDFSIT", func() {
		suite.setupIntlDestinationFirstDayServiceItem()
		totalCost, displayParams, err := priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, testdatagen.DefaultContractCode, idfsitTestRequestedPickupDate, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(idfsitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: idfsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(idfsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(idfsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(idfsitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success with IOFSIT", func() {
		suite.setupIntlOriginFirstDayServiceItem()
		totalCost, displayParams, err := priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIOFSIT, testdatagen.DefaultContractCode, iofsitTestRequestedPickupDate, iofsitTestWeight, iofsitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(iofsitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: iofsitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(iofsitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(iofsitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(iofsitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		suite.setupIntlDestinationFirstDayServiceItem()
		_, _, err := priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, idfsitTestRequestedPickupDate, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported first day SIT code")

		_, _, err = priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, "", idfsitTestRequestedPickupDate, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, testdatagen.DefaultContractCode, time.Time{}, idfsitTestWeight, idfsitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceIntlFirstDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDFSIT, testdatagen.DefaultContractCode, idfsitTestRequestedPickupDate, idfsitTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "PerUnitCents is required")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlAdditionalDaySIT() {
	suite.Run("success with IDASIT", func() {
		suite.setupIntlDestinationAdditionalDayServiceItem()
		totalCost, displayParams, err := priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(idasitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: idasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(idasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(idasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(idasitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success with IOASIT", func() {
		suite.setupIntlOriginAdditionalDayServiceItem()
		totalCost, displayParams, err := priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIOASIT, testdatagen.DefaultContractCode, ioasitTestRequestedPickupDate, idasitNumerDaysInSIT, ioasitTestWeight, ioasitTestPerUnitCents.Int())
		suite.NoError(err)
		suite.Equal(ioasitTestTotalCost, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: ioasitTestContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ioasitTestEscalationCompounded)},
			{Key: models.ServiceItemParamNameIsPeak, Value: FormatBool(ioasitTestIsPeakPeriod)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ioasitTestPerUnitCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {
		suite.setupIntlDestinationAdditionalDayServiceItem()
		_, _, err := priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeDLH, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported additional day SIT code")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, "", idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, time.Time{}, idasitNumerDaysInSIT, idasitTestWeight, idasitTestPerUnitCents.Int())
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, idasitNumerDaysInSIT, idasitTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "PerUnitCents is required")

		_, _, err = priceIntlAdditionalDaySIT(suite.AppContextForTest(), models.ReServiceCodeIDASIT, testdatagen.DefaultContractCode, idasitTestRequestedPickupDate, 0, idasitTestWeight, 0)
		suite.Error(err)
		suite.Contains(err.Error(), "NumberDaysSIT is required")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlCratingUncrating() {
	suite.Run("crating golden path", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		priceCents, displayParams, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)
		suite.NoError(err)
		suite.Equal(icrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(icrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(icrtTestBasePriceCents)},
			{Key: models.ServiceItemParamNameUncappedRequestTotal, Value: FormatCents(dcrtTestUncappedRequestTotal)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("invalid service code", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)
		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeCS, testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "unsupported international crating/uncrating code")
	})

	suite.Run("invalid crate size - external crate", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		badSize := unit.CubicFeet(1.0)
		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, testdatagen.DefaultContractCode, icrtTestRequestedPickupDate, badSize, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, true, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "external crates must be billed for a minimum of 4.00 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, "BOGUS", icrtTestRequestedPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup International Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeICRT, icrtTestMarket, icrtTestBasePriceCents, testdatagen.DefaultContractCode, icrtTestEscalationCompounded)

		twoYearsLaterPickupDate := ioshutTestRequestedPickupDate.AddDate(10, 0, 0)
		_, _, err := priceIntlCratingUncrating(suite.AppContextForTest(), models.ReServiceCodeICRT, testdatagen.DefaultContractCode, twoYearsLaterPickupDate, icrtTestBilledCubicFeet, icrtTestStandaloneCrate, icrtTestStandaloneCrateCap, icrtTestExternalCrate, icrtTestMarket)

		suite.Error(err)
		suite.Contains(err.Error(), "could not calculate escalated price: could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlFuelSurchargeSIT() {
	fscPriceDifferenceInCents := (idsfscFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := idsfscWeightDistanceMultiplier * idsfscTestDistance.Float64()

	suite.Run("invalid service code", func() {
		invalidCode := models.ReServiceCodeIOSHUT
		_, _, err := priceIntlFuelSurchargeSIT(suite.AppContextForTest(), invalidCode, idsfscActualPickupDate, idsfscTestDistance, idsfscTestWeight, idsfscWeightDistanceMultiplier, idsfscFuelPrice)
		suite.NotNil(err)
	})

	suite.Run("success with IOSFSC", func() {
		totalCost, displayParams, err := priceIntlFuelSurchargeSIT(suite.AppContextForTest(), models.ReServiceCodeIOSFSC, iosfscActualPickupDate, iosfscTestDistance, iosfscTestWeight, iosfscWeightDistanceMultiplier, iosfscFuelPrice)
		suite.NoError(err)
		suite.Equal(iosfscPriceCents, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
		}

		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success with IDSFSC", func() {
		totalCost, displayParams, err := priceIntlFuelSurchargeSIT(suite.AppContextForTest(), models.ReServiceCodeIDSFSC, idsfscActualPickupDate, idsfscTestDistance, idsfscTestWeight, idsfscWeightDistanceMultiplier, idsfscFuelPrice)
		suite.NoError(err)
		suite.Equal(idsfscPriceCents, totalCost)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
		}

		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("Invalid parameters to Price", func() {

		invalidActualPickupDate := time.Time{}
		_, _, err := priceIntlFuelSurchargeSIT(suite.AppContextForTest(), models.ReServiceCodeIOSFSC, invalidActualPickupDate, idsfscTestDistance, idsfscTestWeight, idsfscWeightDistanceMultiplier, idsfscFuelPrice)
		suite.Error(err)
		suite.Contains(err.Error(), "ActualPickupDate is required")

		invalidDistance := unit.Miles(-1)
		_, _, err = priceIntlFuelSurchargeSIT(suite.AppContextForTest(), models.ReServiceCodeIOSFSC, idsfscActualPickupDate, invalidDistance, idsfscTestWeight, idsfscWeightDistanceMultiplier, idsfscFuelPrice)
		suite.Error(err)
		suite.Contains(err.Error(), "Distance cannot be less than 0")

		invalidWeight := unit.Pound(0)
		_, _, err = priceIntlFuelSurchargeSIT(suite.AppContextForTest(), models.ReServiceCodeIOSFSC, idsfscActualPickupDate, idsfscTestDistance, invalidWeight, idsfscWeightDistanceMultiplier, idsfscFuelPrice)
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("Weight must be a minimum of %d", minInternationalWeight))

		invalidWeightDistanceMultiplier := float64(0)
		_, _, err = priceIntlFuelSurchargeSIT(suite.AppContextForTest(), models.ReServiceCodeIOSFSC, idsfscActualPickupDate, idsfscTestDistance, idsfscTestWeight, invalidWeightDistanceMultiplier, idsfscFuelPrice)
		suite.Error(err)
		suite.Contains(err.Error(), "WeightBasedDistanceMultiplier is required")

		invalidFuelPrice := unit.Millicents(0)
		_, _, err = priceIntlFuelSurchargeSIT(suite.AppContextForTest(), models.ReServiceCodeIOSFSC, idsfscActualPickupDate, idsfscTestDistance, idsfscTestWeight, idsfscWeightDistanceMultiplier, invalidFuelPrice)
		suite.Error(err)
		suite.Contains(err.Error(), "EIAFuelPrice is required")
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceIntlPickupDeliverySIT() {
	suite.Run("invalid service code", func() {
		invalidCode := models.ReServiceCodeIOSHUT
		_, _, err := priceIntlPickupDeliverySIT(suite.AppContextForTest(), invalidCode, "test", iopsitTestRequestedPickupDate, unit.Pound(1000), 1000, 1)
		suite.Error(err)
		suite.Contains(err.Error(), "unsupported Intl PickupDeliverySIT code")

	})

	suite.Run("success  - distance less than 50 miles", func() {
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
					EscalationCompounded: iopsitTestEscalationCompounded,
				},
			})

		for _, code := range []models.ReServiceCode{models.ReServiceCodeIOPSIT, models.ReServiceCodeIDDSIT} {
			priceCents, displayParams, err := priceIntlPickupDeliverySIT(suite.AppContextForTest(), code, cy.Contract.Code, cy.StartDate.AddDate(0, 0, 1), iopsitTestWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceLessThan50Miles))
			suite.NoError(err)
			suite.Equal(expectIOPSITTestTotalCost, priceCents)

			expectedParams := services.PricingDisplayParams{
				{
					Key:   models.ServiceItemParamNamePriceRateOrFactor,
					Value: FormatCents(unit.Cents(iopsitTestPerUnitCents)),
				},
				{
					Key:   models.ServiceItemParamNameContractYearName,
					Value: cy.Name,
				},
				{
					Key:   models.ServiceItemParamNameIsPeak,
					Value: FormatBool(false),
				},
				{
					Key:   models.ServiceItemParamNameEscalationCompounded,
					Value: FormatEscalation(iopsitTestEscalationCompounded),
				},
			}
			suite.validatePricerCreatedParams(expectedParams, displayParams)
		}
	})

	suite.Run("success  - verify pricing when distance is over 50 miles", func() {
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
					EscalationCompounded: iopsitTestEscalationCompounded,
				},
			})

		for _, code := range []models.ReServiceCode{models.ReServiceCodeIOPSIT, models.ReServiceCodeIDDSIT} {
			priceCents, _, err := priceIntlPickupDeliverySIT(suite.AppContextForTest(), code, cy.Contract.Code, cy.StartDate.AddDate(0, 0, 1), iopsitTestWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceOver50Miles))
			suite.NoError(err)
			// If over 50 miles, mileage will be used as multiplier. Verify by dividing expected price(1 mile) by known distance.
			suite.Equal(expectIOPSITTestTotalCost, (priceCents / iopsitTestDistanceOver50Miles))
		}
	})

	suite.Run("failure - unable to retrieve contract by code", func() {
		_, _, err := priceIntlPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeIOPSIT, "UNKNOWN_CONTRACT_CODE", iopsitTestRequestedPickupDate, iopsitTestWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceLessThan50Miles))
		suite.Error(err)
		suite.Contains(err.Error(), "could not retrieve contract by code")
	})

	suite.Run("failure - could not calculate escalated price", func() {
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
					EscalationCompounded: iopsitTestEscalationCompounded,
				},
			})
		outOfBoundRequestTime := time.Date(2010, time.July, 5, 10, 22, 11, 456, time.UTC)
		_, _, err := priceIntlPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeIOPSIT, cy.Contract.Code, outOfBoundRequestTime, iopsitTestWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceLessThan50Miles))
		suite.Error(err)
		suite.Contains(err.Error(), "could not calculate escalated price")
	})

	suite.Run("Invalid parameters to Price", func() {

		_, _, err := priceIntlPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeIOPSIT, "", iopsitTestRequestedPickupDate, idsfscTestWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceLessThan50Miles))
		suite.Error(err)
		suite.Contains(err.Error(), "ContractCode is required")

		invalidActualPickupDate := time.Time{}
		_, _, err = priceIntlPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeIOPSIT, testdatagen.DefaultContractCode, invalidActualPickupDate, idsfscTestWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceLessThan50Miles))
		suite.Error(err)
		suite.Contains(err.Error(), "ReferenceDate is required")

		invalidWeight := unit.Pound(0)
		_, _, err = priceIntlPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeIOPSIT, testdatagen.DefaultContractCode, idsfscActualPickupDate, invalidWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceLessThan50Miles))
		suite.Error(err)
		suite.Contains(err.Error(), fmt.Sprintf("weight must be a minimum of %d", minInternationalWeight))

		_, _, err = priceIntlPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeIOPSIT, testdatagen.DefaultContractCode, idsfscActualPickupDate, idsfscTestWeight, 0, int(iopsitTestDistanceLessThan50Miles))
		suite.Error(err)
		suite.Contains(err.Error(), "perUnitCents is required")

		_, _, err = priceIntlPickupDeliverySIT(suite.AppContextForTest(), models.ReServiceCodeIOPSIT, testdatagen.DefaultContractCode, idsfscActualPickupDate, idsfscTestWeight, int(iopsitTestPerUnitCents), 0)
		suite.Error(err)
		suite.Contains(err.Error(), "distance is required")
	})
}
