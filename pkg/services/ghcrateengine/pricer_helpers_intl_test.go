package ghcrateengine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_fetchContractFromParams() {
	suite.Run("can fetch contract for the move in the params", func() {
		// Setup contract
		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().Add(time.Hour * 12),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		contract, err := models.FetchContractForMove(suite.AppContextForTest(), move.ID)
		suite.NoError(err)
		suite.Equal(contract.ID, reContract.ID)

		// Setup INPK params
		availableToPrimeAt := time.Now().Add(time.Hour * 6)
		paymentServiceItems := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeINPK,
			[]factory.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameContractCode,
					KeyType: models.ServiceItemParamTypeString,
					Value:   contract.Code,
				},
				{
					Key:     models.ServiceItemParamNameReferenceDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   ihpkTestRequestedPickupDate.Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNamePerUnitCents,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ihpkTestPerUnitCents)),
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   strconv.Itoa(ihpkTestWeight.Int()),
				},
			}, []factory.Customization{{
				// Available to prime is used to fetch market factors for moves
				// The market factor can only be fetched if it's available to the Prime
				// And if it isn't available to the Prime, then we shouldn't be processing the creation
				// of this payment request
				Model: models.Move{
					AvailableToPrimeAt: &availableToPrimeAt,
				},
			},
			}, nil,
		)
		suite.FatalTrue(suite.NotEmpty(paymentServiceItems))

		fetchedContract, err := fetchContractFromParams(suite.AppContextForTest(), paymentServiceItems.PaymentServiceItemParams)
		suite.FatalNoError(err)
		suite.Equal(fetchedContract.ID, contract.ID)
	})
}

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
		suite.setupIntlPackServiceItem(models.ReServiceCodeIHPK)
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
		suite.setupIntlPackServiceItem(models.ReServiceCodeIHPK)
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
