package ghcrateengine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	ducrtTestServiceSchedule      = 3
	ducrtTestBasePriceCents       = unit.Cents(595)
	ducrtTestEscalationCompounded = 1.125
	ducrtTestBilledCubicFeet      = 10
	ducrtTestPriceCents           = unit.Cents(6694) // ducrtTestBasePriceCents * ducrtTestBilledCubicFeet * ducrtTestEscalationCompounded
)

var ducrtTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestDomesticUncratingPricer() {
	suite.setupDomesticAccessorialPrice(models.ReServiceCodeDUCRT, ducrtTestServiceSchedule, ducrtTestBasePriceCents, testdatagen.DefaultContractCode, ducrtTestEscalationCompounded)

	paymentServiceItem := suite.setupDomesticUncratingServiceItem()
	pricer := NewDomesticUncratingPricer(suite.DB())

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := pricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ducrtTestPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameContractYearName, Value: testdatagen.DefaultContractCode},
			{Key: models.ServiceItemParamNameEscalationCompounded, Value: FormatEscalation(ducrtTestEscalationCompounded)},
			{Key: models.ServiceItemParamNamePriceRateOrFactor, Value: FormatCents(ducrtTestBasePriceCents)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := pricer.Price(testdatagen.DefaultContractCode, ducrtTestRequestedPickupDate, ducrtTestBilledCubicFeet, ducrtTestServiceSchedule)
		suite.NoError(err)
		suite.Equal(ducrtTestPriceCents, priceCents)
	})

	suite.Run("PriceUsingParams but sending empty params", func() {
		_, _, err := pricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("invalid crating volume", func() {
		badVolume := unit.CubicFeet(-50.0)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, ducrtTestRequestedPickupDate, badVolume, ducrtTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "crate must be billed for a minimum of 4 cubic feet")
	})

	suite.Run("not finding a rate record", func() {
		_, _, err := pricer.Price("BOGUS", ducrtTestRequestedPickupDate, ducrtTestBilledCubicFeet, ducrtTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup Domestic Accessorial Area Price")
	})

	suite.Run("not finding a contract year record", func() {
		twoYearsLaterPickupDate := ducrtTestRequestedPickupDate.AddDate(2, 0, 0)
		_, _, err := pricer.Price(testdatagen.DefaultContractCode, twoYearsLaterPickupDate, ducrtTestBilledCubicFeet, ducrtTestServiceSchedule)
		suite.Error(err)
		suite.Contains(err.Error(), "could not lookup contract year")
	})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticUncratingServiceItem() models.PaymentServiceItem {
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDUCRT,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameServicesScheduleDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   strconv.Itoa(ducrtTestServiceSchedule),
			},
			{
				Key:     models.ServiceItemParamNameCubicFeetBilled,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   fmt.Sprintf("%d", int(ducrtTestBilledCubicFeet)),
			},
			{
				Key:     models.ServiceItemParamNameRequestedPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   ducrtTestRequestedPickupDate.Format(DateParamFormat),
			},
		},
	)
}
