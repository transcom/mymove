package serviceparamvaluelookups

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServiceAreaOrigin() {
	key := models.ServiceItemParamNameServiceAreaOrigin.String()

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			PostalCode: "35007",
		},
	})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
		})

	domesticServiceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea: "004",
		},
	})

	testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            domesticServiceArea.Contract,
			DomesticServiceArea: domesticServiceArea,
			Zip3:                "350",
		},
	})

	suite.T().Run("golden path", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("004", valueStr)
	})

	suite.T().Run("nil PickupAddress ID", func(t *testing.T) {
		oldPickupAddressID := mtoServiceItem.MTOShipment.PickupAddressID

		mtoServiceItem.MTOShipment.PickupAddress = nil
		mtoServiceItem.MTOShipment.PickupAddressID = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.Contains(err.Error(), "looking for PickupAddressID")
		suite.Equal("", valueStr)

		mtoServiceItem.MTOShipment.PickupAddressID = oldPickupAddressID
		suite.MustSave(&mtoServiceItem.MTOShipment)
	})
}
