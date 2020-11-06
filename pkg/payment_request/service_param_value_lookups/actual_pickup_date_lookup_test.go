package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestActualPickupDateLookup() {
	key := models.ServiceItemParamNameActualPickupDate

	actualPickupDate := time.Date(testdatagen.TestYear, time.May, 18, 0, 0, 0, 0, time.UTC)
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ActualPickupDate: &actualPickupDate,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
		})

	suite.T().Run("golden path with param cache", func(t *testing.T) {
		// FSC
		mtoServiceItemFSC := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "FSC",
			},
			MTOShipment: models.MTOShipment{
				ActualPickupDate: &actualPickupDate,
			},
		})
		mtoServiceItemFSC.MoveTaskOrderID = paymentRequest.MoveTaskOrderID
		mtoServiceItemFSC.MoveTaskOrder = paymentRequest.MoveTaskOrder
		suite.MustSave(&mtoServiceItemFSC)

		// ServiceItemParamNameActualPickupDate
		serviceItemParamKey1 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameActualPickupDate,
				Description: "actual pickup date",
				Type:        models.ServiceItemParamTypeDate,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		})

		_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
			ServiceParam: models.ServiceParam{
				ServiceID:             mtoServiceItemFSC.ReServiceID,
				ServiceItemParamKeyID: serviceItemParamKey1.ID,
				ServiceItemParamKey:   serviceItemParamKey1,
			},
		})

		paramCache := ServiceParamsCache{}
		paramCache.Initialize(suite.DB())
		paramLookupWithCache, _ := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemFSC.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)

		valueStr, err := paramLookupWithCache.ServiceParamValue(serviceItemParamKey1.Key)
		suite.FatalNoError(err)
		expected := actualPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)

		// Verify value from paramCache
		paramCacheValue := paramCache.ParamValue(*mtoServiceItemFSC.MTOShipmentID, key)
		suite.Equal(expected, *paramCacheValue)
	})

	paramLookup, _ := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)

	suite.T().Run("golden path", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := actualPickupDate.Format(ghcrateengine.DateParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.T().Run("nil actual pickup date", func(t *testing.T) {
		// Set the actual pickup date to nil
		mtoShipment := mtoServiceItem.MTOShipment
		oldActualPickupDate := mtoShipment.ActualPickupDate
		mtoShipment.ActualPickupDate = nil
		suite.MustSave(&mtoShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find an actual pickup date for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		mtoShipment.ActualPickupDate = oldActualPickupDate
		suite.MustSave(&mtoShipment)
	})

	suite.T().Run("nil MTOShipmentID", func(t *testing.T) {
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

	suite.T().Run("bogus MTOServiceItemID", func(t *testing.T) {
		// Pass in a non-existent MTOServiceItemID
		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
		_, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})
}
