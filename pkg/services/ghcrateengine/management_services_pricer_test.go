package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	msPriceCents = unit.Cents(12303)
)

var msAvailableToPrimeAt = time.Date(testdatagen.TestYear, time.June, 3, 12, 57, 33, 123, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceManagementServices() {
	suite.setupTaskOrderFeeData(models.ReServiceCodeMS, msPriceCents)
	paymentServiceItem := suite.setupManagementServicesItem()
	counselingServicesPricer := NewManagementServicesPricer(suite.DB())

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := counselingServicesPricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := counselingServicesPricer.Price(testdatagen.DefaultContractCode, msAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		_, err := counselingServicesPricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := counselingServicesPricer.Price("BOGUS", msAvailableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupManagementServicesItem() models.PaymentServiceItem {
	return suite.setupPaymentServiceItemWithParams(
		models.ReServiceCodeMS,
		[]createParams{
			{
				models.ServiceItemParamNameContractCode,
				models.ServiceItemParamTypeString,
				testdatagen.DefaultContractCode,
			},
			{
				models.ServiceItemParamNameMTOAvailableToPrimeAt,
				models.ServiceItemParamTypeTimestamp,
				msAvailableToPrimeAt.Format(TimestampParamFormat),
			},
		},
	)
}
