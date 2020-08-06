package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServiceAreaDest() {
	key := models.ServiceItemParamNameServiceAreaDest.String()

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
	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)

	suite.T().Run("golden path", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("004", valueStr)
	})

	suite.T().Run("nil DestinationAddress ID", func(t *testing.T) {
		oldDestinationAddressID := mtoServiceItem.MTOShipment.DestinationAddressID

		mtoServiceItem.MTOShipment.DestinationAddress = nil
		mtoServiceItem.MTOShipment.DestinationAddressID = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find destination address for MTOShipment [%s]", mtoServiceItem.MTOShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		mtoServiceItem.MTOShipment.DestinationAddressID = oldDestinationAddressID
		suite.MustSave(&mtoServiceItem.MTOShipment)
	})

	suite.T().Run("nil MTOShipment ID", func(t *testing.T) {
		// Set the MTOShipmentID to nil
		oldMTOShipmentID := mtoServiceItem.MTOShipmentID
		mtoServiceItem.MTOShipmentID = nil
		suite.MustSave(&mtoServiceItem)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)

		mtoServiceItem.MTOShipmentID = oldMTOShipmentID
		suite.MustSave(&mtoServiceItem)
	})

	suite.T().Run("nil MTOServiceItem ID", func(t *testing.T) {
		// Pass in a non-existent MTOServiceItemID
		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
		badParamLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)

		valueStr, err := badParamLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)
	})
}
