package serviceparamvaluelookups

import (
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type paramsCacheSubtestData struct {
	paymentRequest                models.PaymentRequest
	move                          models.Move
	mtoShipment1                  models.MTOShipment
	mtoShipment2                  models.MTOShipment
	mtoServiceItemShip1DLH        models.MTOServiceItem
	mtoServiceItemShip2DLH        models.MTOServiceItem
	mtoServiceItemMS              models.MTOServiceItem
	mtoServiceItemCrate1          models.MTOServiceItem
	mtoServiceItemCrate2          models.MTOServiceItem
	mtoServiceItemShuttle         models.MTOServiceItem
	paramKeyWeightEstimated       models.ServiceItemParamKey
	paramKeyRequestedPickupDate   models.ServiceItemParamKey
	paramKeyMTOAvailableToPrimeAt models.ServiceItemParamKey
	paramKeyCubicFeetBilled       models.ServiceItemParamKey
	paramKeyDimensionHeight       models.ServiceItemParamKey
	paramKeyDimensionWidth        models.ServiceItemParamKey
	paramKeyDimensionLength       models.ServiceItemParamKey
	estimatedWeight               unit.Pound
	shuttleEstimatedWeight        unit.Pound
	shuttleActualWeight           unit.Pound
}

func (suite *ServiceParamValueLookupsSuite) makeSubtestData() (subtestData *paramsCacheSubtestData) {
	subtestData = &paramsCacheSubtestData{}
	subtestData.move = testdatagen.MakeDefaultMove(suite.DB())

	subtestData.paymentRequest = testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: subtestData.move.ID,
			},
		})

	subtestData.estimatedWeight = unit.Pound(2048)
	subtestData.mtoShipment1 = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: subtestData.move})
	subtestData.mtoShipment1.PrimeEstimatedWeight = &subtestData.estimatedWeight
	suite.MustSave(&subtestData.mtoShipment1)

	reServiceDLH := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDLH,
		},
	})

	reServiceDOP := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOP,
		},
	})

	reServiceMS := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeMS,
		},
	})

	reServiceDCRT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDCRT,
		},
	})

	reServiceDOSHUT := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOSHUT,
		},
	})

	// DLH
	subtestData.mtoServiceItemShip1DLH = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        subtestData.move,
		ReService:   reServiceDLH,
		MTOShipment: subtestData.mtoShipment1,
	})

	// DOP
	mtoServiceItemShip1DOP := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        subtestData.move,
		ReService:   reServiceDOP,
		MTOShipment: subtestData.mtoShipment1,
	})

	subtestData.mtoShipment2 = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: subtestData.move})
	subtestData.mtoShipment2.PrimeEstimatedWeight = &subtestData.estimatedWeight
	suite.MustSave(&subtestData.mtoShipment2)

	// DLH
	subtestData.mtoServiceItemShip2DLH = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        subtestData.move,
		ReService:   reServiceDLH,
		MTOShipment: subtestData.mtoShipment2,
	})

	// MS
	subtestData.mtoServiceItemMS = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:      subtestData.move,
		ReService: reServiceMS,
	})

	subtestData.paramKeyWeightEstimated = testdatagen.FetchOrMakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameWeightEstimated,
			Description: "estimated weight",
			Type:        models.ServiceItemParamTypeInteger,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	subtestData.paramKeyRequestedPickupDate = testdatagen.FetchOrMakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameRequestedPickupDate,
			Description: "requested pickup date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginPrime,
		},
	})

	// DLH
	_ = testdatagen.FetchOrMakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItemShip1DLH.ReServiceID,
			ServiceItemParamKeyID: subtestData.paramKeyWeightEstimated.ID,
			ServiceItemParamKey:   subtestData.paramKeyWeightEstimated,
			IsOptional:            true,
		},
	})

	// DLH
	_ = testdatagen.FetchOrMakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItemShip1DLH.ReServiceID,
			ServiceItemParamKeyID: subtestData.paramKeyRequestedPickupDate.ID,
			ServiceItemParamKey:   subtestData.paramKeyRequestedPickupDate,
		},
	})

	// DOP
	_ = testdatagen.FetchOrMakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             mtoServiceItemShip1DOP.ReServiceID,
			ServiceItemParamKeyID: subtestData.paramKeyWeightEstimated.ID,
			ServiceItemParamKey:   subtestData.paramKeyWeightEstimated,
			IsOptional:            true,
		},
	})

	subtestData.paramKeyMTOAvailableToPrimeAt = testdatagen.FetchOrMakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameMTOAvailableToPrimeAt,
			Description: "prime mto made available date",
			Type:        models.ServiceItemParamTypeDate,
			Origin:      models.ServiceItemParamOriginSystem,
		},
	})

	_ = testdatagen.FetchOrMakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItemMS.ReServiceID,
			ServiceItemParamKeyID: subtestData.paramKeyMTOAvailableToPrimeAt.ID,
			ServiceItemParamKey:   subtestData.paramKeyMTOAvailableToPrimeAt,
		},
	})

	subtestData.shuttleEstimatedWeight = unit.Pound(400)
	subtestData.shuttleActualWeight = unit.Pound(450)
	subtestData.mtoServiceItemShuttle = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        subtestData.move,
		ReService:   reServiceDOSHUT,
		MTOShipment: subtestData.mtoShipment2,
		MTOServiceItem: models.MTOServiceItem{
			EstimatedWeight: &subtestData.shuttleEstimatedWeight,
			ActualWeight:    &subtestData.shuttleActualWeight,
		},
	})

	// DOSHUT estimated weight
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             subtestData.mtoServiceItemShuttle.ReServiceID,
			ServiceItemParamKeyID: subtestData.paramKeyWeightEstimated.ID, // estimated weight
			ServiceItemParamKey:   subtestData.paramKeyWeightEstimated,
			IsOptional:            true,
		},
	})
	subtestData.mtoServiceItemCrate1 = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        subtestData.move,
		ReService:   reServiceDCRT,
		MTOShipment: subtestData.mtoShipment2,
	})
	_ = testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
		MTOServiceItemDimension: models.MTOServiceItemDimension{
			MTOServiceItemID: subtestData.mtoServiceItemCrate1.ID,
			Type:             models.DimensionTypeCrate,
			// These dimensions are chosen to overflow 32bit ints if multiplied, and give a fractional result
			// when converted to cubic feet.
			Length:    16*12*1000 + 1000,
			Height:    8 * 12 * 1000,
			Width:     8 * 12 * 1000,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	})
	_ = testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
		MTOServiceItemDimension: models.MTOServiceItemDimension{
			MTOServiceItemID: subtestData.mtoServiceItemCrate1.ID,
			Type:             models.DimensionTypeItem,
			Length:           12000,
			Height:           12000,
			Width:            12000,
			CreatedAt:        time.Time{},
			UpdatedAt:        time.Time{},
		},
	})
	subtestData.mtoServiceItemCrate2 = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move:        subtestData.move,
		ReService:   reServiceDCRT,
		MTOShipment: subtestData.mtoShipment2,
	})
	_ = testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
		MTOServiceItemDimension: models.MTOServiceItemDimension{
			MTOServiceItemID: subtestData.mtoServiceItemCrate2.ID,
			Type:             models.DimensionTypeCrate,
			// These dimensions are chosen to overflow 32bit ints if multiplied, and give a fractional result
			// when converted to cubic feet.
			Length:    7000,
			Height:    7000,
			Width:     7000,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	})
	_ = testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
		MTOServiceItemDimension: models.MTOServiceItemDimension{
			MTOServiceItemID: subtestData.mtoServiceItemCrate2.ID,
			Type:             models.DimensionTypeItem,
			Length:           6000,
			Height:           6000,
			Width:            6000,
			CreatedAt:        time.Time{},
			UpdatedAt:        time.Time{},
		},
	})
	subtestData.paramKeyDimensionHeight = testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameDimensionHeight,
			Description: "height",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginSystem,
		},
	})
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             reServiceDCRT.ID,
			ServiceItemParamKeyID: subtestData.paramKeyDimensionHeight.ID,
			ServiceItemParamKey:   subtestData.paramKeyDimensionHeight,
		},
	})
	subtestData.paramKeyDimensionWidth = testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameDimensionWidth,
			Description: "width",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginSystem,
		},
	})
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             reServiceDCRT.ID,
			ServiceItemParamKeyID: subtestData.paramKeyDimensionWidth.ID,
			ServiceItemParamKey:   subtestData.paramKeyDimensionWidth,
		},
	})
	subtestData.paramKeyDimensionLength = testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameDimensionLength,
			Description: "length",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginSystem,
		},
	})
	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             reServiceDCRT.ID,
			ServiceItemParamKeyID: subtestData.paramKeyDimensionLength.ID,
			ServiceItemParamKey:   subtestData.paramKeyDimensionLength,
		},
	})
	subtestData.paramKeyCubicFeetBilled = testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameCubicFeetBilled,
			Description: "cubic feet billed",
			Type:        models.ServiceItemParamTypeDecimal,
			Origin:      models.ServiceItemParamOriginSystem,
		},
	})

	_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
		ServiceParam: models.ServiceParam{
			ServiceID:             reServiceDCRT.ID,
			ServiceItemParamKeyID: subtestData.paramKeyCubicFeetBilled.ID,
			ServiceItemParamKey:   subtestData.paramKeyCubicFeetBilled,
		},
	})

	return subtestData
}

func (suite *ServiceParamValueLookupsSuite) TestServiceParamCache() {
	// Create some records we'll need to link to

	paramCache := NewServiceParamsCache()

	suite.Run("weight billed shuttling and DLH", func() {
		// A test to confirm that the cache returns correct values when having two lookups of the same param
		// on the same shipment.
		subtestData := suite.makeSubtestData()

		var paramLookupService1 *ServiceItemParamKeyData
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip2DLH.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyWeightEstimated.Key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(subtestData.estimatedWeight.Int())
		suite.Equal(expected, estimatedWeightStr)

		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShuttle.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		var shuttleEstimatedWeightStr string
		shuttleEstimatedWeightStr, err = paramLookupService2.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyWeightEstimated.Key)
		suite.FatalNoError(err)
		shuttleExpected := strconv.Itoa(subtestData.shuttleEstimatedWeight.Int())
		suite.Equal(shuttleExpected, shuttleEstimatedWeightStr)
	})

	// cubic feet billed
	suite.Run("cubic feet billed", func() {
		// Another test to confirm that the cache returns correct values when having two lookups of the same param
		// on the same shipment.
		subtestData := suite.makeSubtestData()
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemCrate1.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		var cubicFeet string
		cubicFeet, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyCubicFeetBilled.Key)
		suite.FatalNoError(err)
		suite.Equal("1029.33", cubicFeet)

		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemCrate2.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		cubicFeet, err = paramLookupService2.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyCubicFeetBilled.Key)
		suite.FatalNoError(err)
		suite.Equal("4.00", cubicFeet)
	})

	// Estimated Weight
	suite.Run("Shipment 1 estimated weight", func() {
		subtestData := suite.makeSubtestData()
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyWeightEstimated.Key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(subtestData.estimatedWeight.Int())
		suite.Equal(expected, estimatedWeightStr)
	})

	// Requested Pickup Date
	suite.Run("Shipment 1 requested pickup date", func() {
		subtestData := suite.makeSubtestData()
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)
		expectedRequestedPickupDate := subtestData.mtoShipment1.RequestedPickupDate.String()[:10]
		var requestedPickupDateStr string
		requestedPickupDateStr, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyRequestedPickupDate.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedRequestedPickupDate, requestedPickupDateStr)
	})

	// Estimated Weight changed on shipment1 but pulled from cache
	suite.Run("Shipment 1 estimated weight cache", func() {
		subtestData := suite.makeSubtestData()
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)
		expectedWeight := strconv.Itoa(subtestData.estimatedWeight.Int())
		changeExpectedEstimatedWeight := unit.Pound(3048)
		subtestData.mtoShipment1.PrimeEstimatedWeight = &changeExpectedEstimatedWeight
		suite.MustSave(&subtestData.mtoShipment1)
		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyWeightEstimated.Key)
		suite.FatalNoError(err)

		// EstimatedWeight hasn't changed from the cache
		suite.Equal(expectedWeight, estimatedWeightStr)
		// mtoShipment1 was changed to the new estimated weight
		suite.Equal(changeExpectedEstimatedWeight, *subtestData.mtoShipment1.PrimeEstimatedWeight)
	})

	// Requested Pickup Date changed on shipment1 but pulled from cache
	suite.Run("Shipment 1 requested pickup date changed", func() {
		subtestData := suite.makeSubtestData()
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)
		expectedRequestedPickupDate := subtestData.mtoShipment1.RequestedPickupDate.String()[:10]
		changeRequestedPickupDate := time.Date(testdatagen.GHCTestYear, time.April, 15, 0, 0, 0, 0, time.UTC)
		subtestData.mtoShipment1.RequestedPickupDate = &changeRequestedPickupDate
		suite.MustSave(&subtestData.mtoShipment1)

		var requestedPickupDateStr string
		requestedPickupDateStr, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyRequestedPickupDate.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedRequestedPickupDate, requestedPickupDateStr)
		// mtoShipment1 was changed to the new date
		suite.Equal(changeRequestedPickupDate.String()[:10], subtestData.mtoShipment1.RequestedPickupDate.String()[:10])
	})
	// DLH - for shipment 2
	// Estimated Weight
	suite.Run("Shipment 2 DLH estimated weight", func() {
		subtestData := suite.makeSubtestData()

		var paramLookupService2 *ServiceItemParamKeyData
		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip2DLH.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService2.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyWeightEstimated.Key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(subtestData.estimatedWeight.Int())
		suite.Equal(expected, estimatedWeightStr)
	})

	// Requested Pickup Date
	suite.Run("Shipment 2 Requested Pickup Date", func() {
		subtestData := suite.makeSubtestData()
		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip2DLH.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)
		expectedRequestedPickupDate := subtestData.mtoShipment2.RequestedPickupDate.String()[:10]
		var requestedPickupDateStr string
		requestedPickupDateStr, err = paramLookupService2.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyRequestedPickupDate.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedRequestedPickupDate, requestedPickupDateStr)
	})

	// MS - has no shipment
	// Prime MTO Made Available Date
	suite.Run("Task Order Service Prime MTO available", func() {
		subtestData := suite.makeSubtestData()

		subtestData.mtoServiceItemMS.MTOShipmentID = nil
		subtestData.mtoServiceItemMS.MTOShipment = models.MTOShipment{}
		suite.MustSave(&subtestData.mtoServiceItemMS)

		var paramLookupService3 *ServiceItemParamKeyData
		paramLookupService3, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemMS.ID, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		availToPrimeAt := time.Date(testdatagen.GHCTestYear, time.April, 15, 0, 0, 0, 0, time.UTC)
		subtestData.move.AvailableToPrimeAt = &availToPrimeAt
		suite.MustSave(&subtestData.move)
		expectedAvailToPrimeDate := subtestData.move.AvailableToPrimeAt.String()[:10]
		var availToPrimeDateStr string
		availToPrimeDateStr, err = paramLookupService3.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyMTOAvailableToPrimeAt.Key)
		suite.FatalNoError(err)
		suite.Equal(expectedAvailToPrimeDate, availToPrimeDateStr[:10])
	})
}
