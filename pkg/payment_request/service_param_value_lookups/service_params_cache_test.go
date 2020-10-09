package serviceparamvaluelookups

import (
	"strconv"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestServiceParamCache() {
	// Create some records we'll need to link to

	move := testdatagen.MakeDefaultMove(suite.DB())

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: move.ID,
			},
		})

	estimatedWeight := unit.Pound(2048)
	mtoShipment1 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move})
	mtoShipment1.PrimeEstimatedWeight = &estimatedWeight
	suite.MustSave(&mtoShipment1)

	reService1 := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DLH",
		},
	})

	reService2 := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "DOP",
		},
	})

	reService3 := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: "MS",
		},
	})

	// DLH
	mtoServiceItem1 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        move,
		ReService:   reService1,
		MTOShipment: mtoShipment1,
	})

	// DOP
	mtoServiceItem2 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        move,
		ReService:   reService2,
		MTOShipment: mtoShipment1,
	})

	mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: move})
	mtoShipment2.PrimeEstimatedWeight = &estimatedWeight
	suite.MustSave(&mtoShipment2)

	// DLH
	mtoServiceItem3 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        move,
		ReService:   reService1,
		MTOShipment: mtoShipment2,
	})

	// MS
	mtoServiceItem4 := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:      move,
		ReService: reService3,
	})

	serviceItemParamKey1 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameWeightEstimated,
			Description: "estimated weight",
			Type:        models.ServiceItemParamTypeInteger,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})
	serviceItemParamKey2 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameRequestedPickupDate,
			Description: "requested pickup date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	// DLH
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
		},
	})

	// DLH
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem1.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey2.ID,
			ServiceItemParamKey:   serviceItemParamKey2,
		},
	})

	// DOP
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem2.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey1.ID,
			ServiceItemParamKey:   serviceItemParamKey1,
		},
	})

	serviceItemParamKey3 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameMTOAvailableToPrimeAt,
			Description: "prime mto made available date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginSystem,
		},
	})

	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItem4.ReServiceID,
			ServiceItemParamKeyID: serviceItemParamKey3.ID,
			ServiceItemParamKey:   serviceItemParamKey3,
		},
	})

	paramCache := ServiceParamsCache{}
	paramCache.Initialize(suite.DB())

	paramLookupService1, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem1.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
	suite.NoError(err)

	// Estimated Weight
	suite.T().Run("Shipment 1 "+serviceItemParamKey1.Key.String(), func(t *testing.T) {
		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService1.ServiceParamValue(serviceItemParamKey1.Key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(estimatedWeight.Int())
		suite.Equal(expected, estimatedWeightStr)
	})

	// Requested Pickup Date
	suite.T().Run("Shipment 1 "+serviceItemParamKey2.Key.String(), func(t *testing.T) {
		expectedRequestedPickupDate := mtoShipment1.RequestedPickupDate.String()[:10]
		var requestedPickupDateStr string
		requestedPickupDateStr, err = paramLookupService1.ServiceParamValue(serviceItemParamKey2.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedRequestedPickupDate, requestedPickupDateStr)
	})

	// Estimated Weight changed on shipment1 but pulled from cache
	suite.T().Run("Shipment 1 "+serviceItemParamKey1.Key.String(), func(t *testing.T) {
		expectedWeight := strconv.Itoa(estimatedWeight.Int())
		changeExpectedEstimatedWeight := unit.Pound(3048)
		mtoShipment1.PrimeEstimatedWeight = &changeExpectedEstimatedWeight
		suite.MustSave(&mtoShipment1)
		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService1.ServiceParamValue(serviceItemParamKey1.Key)
		suite.FatalNoError(err)

		// EstimatedWeight hasn't changed from the cache
		suite.Equal(expectedWeight, estimatedWeightStr)
		// mtoShipment1 was changed to the new estimated weight
		suite.Equal(changeExpectedEstimatedWeight, *mtoShipment1.PrimeEstimatedWeight)
	})

	// Requested Pickup Date changed on shipment1 but pulled from cache
	suite.T().Run("Shipment 1 "+serviceItemParamKey2.Key.String(), func(t *testing.T) {
		expectedRequestedPickupDate := mtoShipment1.RequestedPickupDate.String()[:10]
		changeRequestedPickupDate := time.Date(testdatagen.GHCTestYear, time.April, 15, 0, 0, 0, 0, time.UTC)
		mtoShipment1.RequestedPickupDate = &changeRequestedPickupDate
		suite.MustSave(&mtoShipment1)

		var requestedPickupDateStr string
		requestedPickupDateStr, err = paramLookupService1.ServiceParamValue(serviceItemParamKey2.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedRequestedPickupDate, requestedPickupDateStr)
		// mtoShipment1 was changed to the new date
		suite.Equal(changeRequestedPickupDate.String()[:10], mtoShipment1.RequestedPickupDate.String()[:10])
	})

	var paramLookupService2 *ServiceItemParamKeyData
	paramLookupService2, err = ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
	suite.NoError(err)

	// DLH - for shipment 2
	// Estimated Weight
	suite.T().Run("Shipment 2 "+serviceItemParamKey1.Key.String(), func(t *testing.T) {
		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService2.ServiceParamValue(serviceItemParamKey1.Key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(estimatedWeight.Int())
		suite.Equal(expected, estimatedWeightStr)
	})

	// Requested Pickup Date
	suite.T().Run("Shipment 2 "+serviceItemParamKey2.Key.String(), func(t *testing.T) {
		expectedRequestedPickupDate := mtoShipment2.RequestedPickupDate.String()[:10]
		var requestedPickupDateStr string
		requestedPickupDateStr, err = paramLookupService2.ServiceParamValue(serviceItemParamKey2.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedRequestedPickupDate, requestedPickupDateStr)
	})

	mtoServiceItem4.MTOShipmentID = nil
	mtoServiceItem4.MTOShipment = models.MTOShipment{}
	suite.MustSave(&mtoServiceItem4)

	var paramLookupService3 *ServiceItemParamKeyData
	paramLookupService3, err = ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem4.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
	suite.NoError(err)

	// MS - has no shipment
	// Prime MTO Made Available Date
	suite.T().Run("Task Order Service "+serviceItemParamKey3.Key.String(), func(t *testing.T) {
		availToPrimeAt := time.Date(testdatagen.GHCTestYear, time.April, 15, 0, 0, 0, 0, time.UTC)
		move.AvailableToPrimeAt = &availToPrimeAt
		suite.MustSave(&move)
		expectedAvailToPrimeDate := move.AvailableToPrimeAt.String()[:10]
		var availToPrimeDateStr string
		availToPrimeDateStr, err = paramLookupService3.ServiceParamValue(serviceItemParamKey3.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedAvailToPrimeDate, availToPrimeDateStr[:10])
	})
}