package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestZipDestAddressLookup() {
	key := models.ServiceItemParamNameZipDestAddress.String()

	suite.T().Run("zip destination address is present on MTO Shipment", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: "Domestic Line Haul",
				},
				MTOShipment: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("94535", valueStr)
	})

	suite.T().Run("nil DestinationAddressID", func(t *testing.T) {
		// Set the DestinationAddressID to nil
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: "Domestic Line Haul",
				},
			})

		mtoServiceItem.MTOShipment.DestinationAddressID = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		expected := fmt.Sprintf("looking for DestinationAddressID")
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})
}
