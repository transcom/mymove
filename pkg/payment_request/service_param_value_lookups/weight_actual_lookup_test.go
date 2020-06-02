package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestWeightActualLookup() {
	key := "WeightActual"

	suite.T().Run("actual weight is present on MTO Shipment", func(t *testing.T) {
		_, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.T().Run("nil PrimeActualWeight", func(t *testing.T) {
		// Set the actual weight to nil
		mtoServiceItem, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.PrimeActualWeight = nil
		suite.MustSave(&mtoShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find actual weight for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})

	suite.T().Run("nil MTOShipmentID", func(t *testing.T) {
		// Set the MTOShipmentID to nil
		mtoServiceItem, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
		mtoServiceItem.MTOShipmentID = nil
		suite.MustSave(&mtoServiceItem)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		expected := fmt.Sprintf("looking for MTOShipmentID")
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})

	suite.T().Run("bogus MTOServiceItemID", func(t *testing.T) {
		// Pass in a non-existent MTOServiceItemID
		_, paymentRequest, _ := suite.setupTest(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH, models.MTOShipmentTypeHHG)
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
