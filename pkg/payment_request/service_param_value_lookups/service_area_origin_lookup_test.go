package serviceparamvaluelookups

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServiceAreaOrigin() {
	key := models.ServiceItemParamNameServiceAreaOrigin.String()

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
		})

	testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{})

	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

	suite.T().Run("golden path", func(t *testing.T) {
		mtoServiceItem.MTOShipment.PickupAddress.PostalCode = "35007"
		suite.MustSave(&mtoServiceItem)

		fmt.Println("postalCode in test")
		fmt.Printf("%v", mtoServiceItem.MTOShipment.PickupAddress.PostalCode)

		valueStr, err := paramLookup.ServiceParamValue(key)

		suite.FatalNoError(err)
		suite.Equal("004", valueStr)
	})
}