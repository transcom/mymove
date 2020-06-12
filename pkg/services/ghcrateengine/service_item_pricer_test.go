package ghcrateengine

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *GHCRateEngineServiceSuite) TestPriceServiceItem() {
	suite.setupPriceServiceItemData()
	paymentServiceItem := suite.setupPriceServiceItem()
	serviceItemPricer := NewServiceItemPricer(suite.DB())

	suite.T().Run("golden path", func(t *testing.T) {
		priceCents, err := serviceItemPricer.PriceServiceItem(paymentServiceItem)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.T().Run("not implemented pricer", func(t *testing.T) {
		badPaymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "BOGUS",
			},
		})

		_, err := serviceItemPricer.PriceServiceItem(badPaymentServiceItem)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupPriceServiceItemData() {
	contractYear := testdatagen.MakeDefaultReContractYear(suite.DB())

	counselingService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		})

	taskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     msPriceCents,
	}
	suite.MustSave(&taskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setupPriceServiceItem() models.PaymentServiceItem {
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
