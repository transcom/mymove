package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestZipDestAddress() {
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

		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

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

		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		expected := fmt.Sprintf("looking for DestinationAddressID")
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})

	suite.T().Run("nil MTOShipmentID", func(t *testing.T) {
		// Set the MTOShipmentID to nil
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: "Domestic Line Haul",
				},
			})

		mtoServiceItem.MTOShipmentID = nil
		suite.MustSave(&mtoServiceItem)

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		expected := fmt.Sprintf("looking for MTOShipmentID")
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})

	suite.T().Run("bogus MTOServiceItemID", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: "Domestic Line Haul",
				},
			})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		// Pass in a non-existent MTOServiceItemID
		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
		badParamLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

		valueStr, err := badParamLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		expected := fmt.Sprintf("looking for MTOServiceItemID")
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})
}
