package serviceparamvaluelookups

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestEIAFuelPriceLookup() {
	key := models.ServiceItemParamNameEIAFuelPrice.String()

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

	suite.T().Run("lookup GHC diesel fuel price successfully", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("90210", valueStr)
	})
}