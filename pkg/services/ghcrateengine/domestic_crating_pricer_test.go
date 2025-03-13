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
	dcrtTestServiceSchedule      = 3
	dcrtTestBasePriceCents       = unit.Cents(2369)
	dcrtTestEscalationCompounded = 1.11
	dcrtTestBilledCubicFeet      = unit.CubicFeet(10)
	dcrtTestPriceCents           = unit.Cents(26300)
	dcrtTestStandaloneCrate      = false
	dcrtTestStandaloneCrateCap   = unit.Cents(1000000)
	dcrtTestUncappedRequestTotal = unit.Cents(26300)
)

var dcrtTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticCratingPricer() {
	pricer := NewDomesticCratingPricer()

	// var reContractYear = testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
	// 	ReContractYear: models.ReContractYear{
	// 		StartDate: testdatagen.ContractStartDate,
	// 		EndDate:   testdatagen.ContractEndDate,
	// 	},
	// })
	// var reService = testdatagen.FetchReService(suite.DB(), testdatagen.Assertions{
	// 	ReService: models.ReService{
	// 		Code: models.ReServiceCodeDCRT,
	// 	},
	// })

	// testdatagen.FetchOrMakeReDomesticAccessorialPrice(suite.DB(), testdatagen.Assertions{
	// 	ReDomesticAccessorialPrice: models.ReDomesticAccessorialPrice{
	// 		ContractID:       reContractYear.ContractID,
	// 		ServiceID:        reService.ID,
	// 		ServicesSchedule: dcrtTestServiceSchedule,
	// 	},
	// })

	suite.Run("success using PaymentServiceItemParams", func() {

		paymentServiceItem := suite.setupDomesticCratingServiceItem(dcrtTestBilledCubicFeet)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dcrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractYearName},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(dcrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(dcrtTestBasePriceCents)},
			{Key: models.ServiceItemParamNameUncappedRequestTotal, Value: FormatCents(dcrtTestUncappedRequestTotal)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
	suite.Run("success with truncating cubic feet", func() {

		paymentServiceItem := suite.setupDomesticCratingServiceItem(unit.CubicFeet(10.005))
		priceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dcrtTestPriceCents, priceCents)
	})

	suite.Run("success without PaymentServiceItemParams", func() {

		priceCents, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)
		suite.NoError(err)
		suite.Equal(dcrtTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid crating volume", func() {
		badVolume := unit.CubicFeet(3.0)
		_, _, err := pricer.Price(suite.AppContextForTest(), testdatagen.DefaultContractCode, dcrtTestRequestedPickupDate, badVolume, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)
		suite.Error(err)
		suite.Contains(err.Error(), "crate must be billed for a minimum of 4 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := pricer.Price(suite.AppContextForTest(), "BOGUS", dcrtTestRequestedPickupDate, dcrtTestBilledCubicFeet, dcrtTestServiceSchedule, dcrtTestStandaloneCrate, dcrtTestStandaloneCrateCap)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup Domestic Accessorial Area Price")
	})

}

func (suite *GHCRateEngineServiceSuite) setupDomesticCratingServiceItem(cubicFeet unit.CubicFeet) models.PaymentServiceItem {
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDCRT,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameCubicFeetBilled,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   cubicFeet.String(),
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   dcrtTestRequestedPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(dcrtTestServiceSchedule),
			},
			{
				Key:     models.ServiceItemParamNameStandaloneCrate,
				KeyType: models.ServiceItemParamTypeBoolean,
				Value:   strconv.FormatBool(false),
			},
			{
				Key:     models.ServiceItemParamNameStandaloneCrateCap,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.FormatInt(100000, 10),
			},
		}, nil, nil,
	)
}
