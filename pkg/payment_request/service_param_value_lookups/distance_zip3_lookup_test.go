//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package serviceparamvaluelookups

import (
	"strconv"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestDistanceZip3Lookup() {
	key := models.ServiceItemParamNameDistanceZip3

	suite.T().Run("Calculate zip3 distance", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip3Distance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		suite.DB().Find(&mtoShipment, mtoServiceItem.MTOShipmentID)

		suite.Equal(unit.Miles(defaultZip3Distance), *mtoShipment.Distance)
	})

	suite.T().Run("Calculate zip3 distance with param cache", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		// DLH
		reServiceDLH := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DLH",
			},
		})

		estimatedWeight := unit.Pound(2048)

		// DLH
		mtoServiceItemDLH := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: reServiceDLH,
			MTOShipment: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
			},
		})
		mtoServiceItemDLH.MoveTaskOrderID = paymentRequest.MoveTaskOrderID
		mtoServiceItemDLH.MoveTaskOrder = paymentRequest.MoveTaskOrder
		suite.MustSave(&mtoServiceItemDLH)

		// ServiceItemParamNameActualPickupDate
		serviceItemParamKey1 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameDistanceZip3,
				Description: "zip 3 distance",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		})

		_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
			ServiceParam: models.ServiceParam{
				ServiceID:             mtoServiceItemDLH.ReServiceID,
				ServiceItemParamKeyID: serviceItemParamKey1.ID,
				ServiceItemParamKey:   serviceItemParamKey1,
			},
		})

		paramCache := ServiceParamsCache{}
		paramCache.Initialize(suite.DB())

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemDLH.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip3Distance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		suite.DB().Find(&mtoShipment, mtoServiceItemDLH.MTOShipmentID)

		suite.Equal(unit.Miles(defaultZip3Distance), *mtoShipment.Distance)

		// Verify value from paramCache
		paramCacheValue := paramCache.ParamValue(*mtoServiceItemDLH.MTOShipmentID, key)
		suite.Equal(expected, *paramCacheValue)
	})
}
