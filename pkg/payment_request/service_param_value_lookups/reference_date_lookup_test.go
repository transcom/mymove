package serviceparamvaluelookups

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestReferenceDateLookup() {
	key := models.ServiceItemParamNameReferenceDate

	requestedPickupDate := time.Date(testdatagen.TestYear, time.May, 10, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(testdatagen.TestYear, time.May, 14, 0, 0, 0, 0, time.UTC)

	setupTestData := func(shipmentType models.MTOShipmentType) models.MTOServiceItem {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					ActualPickupDate:    &actualPickupDate,
					RequestedPickupDate: &requestedPickupDate,
					ShipmentType:        shipmentType,
				},
			})

		// Don't need a payment request for this test.

		return mtoServiceItem
	}

	suite.Run("golden path for HHG", func() {
		mtoServiceItem := setupTestData(models.MTOShipmentTypeHHG)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := requestedPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.Run("golden path for NTS-Release", func() {
		mtoServiceItem := setupTestData(models.MTOShipmentTypeHHGOutOfNTSDom)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := actualPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.Run("error path for HHG", func() {
		mtoServiceItem := setupTestData(models.MTOShipmentTypeHHG)

		// Set the RequestedPickupDate to nil
		mtoServiceItem.MTOShipment.RequestedPickupDate = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find a valid requested pickup date for MTOShipmentID [%s]", mtoServiceItem.MTOShipment.ID)
		suite.Contains(err.Error(), expected)
	})

	suite.Run("error path for NTS-Release", func() {
		mtoServiceItem := setupTestData(models.MTOShipmentTypeHHGOutOfNTSDom)

		// Set the ActualPickupDate to nil
		mtoServiceItem.MTOShipment.ActualPickupDate = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find a valid actual pickup date for MTOShipmentID [%s]", mtoServiceItem.MTOShipment.ID)
		suite.Contains(err.Error(), expected)
	})
}
