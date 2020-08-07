package serviceparamvaluelookups

import (
	"strconv"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type createParams struct {
	key     models.ServiceItemParamName
	keyType models.ServiceItemParamType
	value   string
}

const (
	csPriceCents = unit.Cents(8327)
)

var csAvailableToPrimeAt = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)


func (suite *ServiceParamValueLookupsSuite) TestParamCacheDistanceZip3Lookup() {
	key := models.ServiceItemParamNameDistanceZip3.String()

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	paramCache := ServiceParamsCache{}
	paramCache.Initialize(suite.DB())

	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)

	suite.T().Run("Calculate zip3 distance", func(t *testing.T) {
		distanceStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultDistance)
		suite.Equal(expected, distanceStr)
	})

	suite.T().Run("nil PickupAddressID", func(t *testing.T) {
		oldPickupAddressID := mtoServiceItem.MTOShipment.PickupAddressID

		mtoServiceItem.MTOShipment.PickupAddress = nil
		mtoServiceItem.MTOShipment.PickupAddressID = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)

		// Address value has changed in the system but the cache will be the same
		// value as the first run wih the default distance
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultDistance)
		suite.Equal(expected, valueStr)

		mtoServiceItem.MTOShipment.PickupAddressID = oldPickupAddressID
		suite.MustSave(&mtoServiceItem.MTOShipment)
	})

	suite.T().Run("nil DestinationAddressID", func(t *testing.T) {
		oldDestinationAddressID := mtoServiceItem.MTOShipment.PickupAddressID

		mtoServiceItem.MTOShipment.DestinationAddress = nil
		mtoServiceItem.MTOShipment.DestinationAddressID = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)

		// Address value has changed in the system but the cache will be the same
		// value as the first run wih the default distance
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultDistance)
		suite.Equal(expected, valueStr)

		mtoServiceItem.MTOShipment.PickupAddressID = oldDestinationAddressID
		suite.MustSave(&mtoServiceItem.MTOShipment)
	})

}


func (suite *ServiceParamValueLookupsSuite) setupCounselingServicesItem() models.PaymentServiceItem {
	return suite.setupPaymentServiceItemWithParams(
		models.ReServiceCodeCS,
		[]createParams{
			{
				models.ServiceItemParamNameContractCode,
				models.ServiceItemParamTypeString,
				testdatagen.DefaultContractCode,
			},
			{
				models.ServiceItemParamNameMTOAvailableToPrimeAt,
				models.ServiceItemParamTypeTimestamp,
				csAvailableToPrimeAt.Format(TimestampParamFormat),
			},
		},
	)
}


func (suite *ServiceParamValueLookupsSuite) setupPaymentServiceItemWithParams(serviceCode models.ReServiceCode, paramsToCreate []createParams) models.PaymentServiceItem {
	var params models.PaymentServiceItemParams

	paymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: serviceCode,
		},
	})

	for _, param := range paramsToCreate {
		serviceItemParamKey := testdatagen.MakeServiceItemParamKey(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: models.ServiceItemParamKey{
					Key:  param.key,
					Type: param.keyType,
				},
			})

		serviceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem:  paymentServiceItem,
				ServiceItemParamKey: serviceItemParamKey,
				PaymentServiceItemParam: models.PaymentServiceItemParam{
					Value: param.value,
				},
			})
		params = append(params, serviceItemParam)
	}

	paymentServiceItem.PaymentServiceItemParams = params

	return paymentServiceItem
}