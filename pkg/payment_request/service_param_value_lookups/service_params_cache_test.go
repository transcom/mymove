package serviceparamvaluelookups

import (
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/factory"
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
	testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			EndDate: time.Now().Add(24 * time.Hour),
		},
	})
	subtestData = &paramsCacheSubtestData{}
	subtestData.move = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

	subtestData.paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
	}, nil)

	subtestData.estimatedWeight = unit.Pound(2048)
	subtestData.mtoShipment1 = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		}}, nil)
	subtestData.mtoShipment1.PrimeEstimatedWeight = &subtestData.estimatedWeight
	suite.MustSave(&subtestData.mtoShipment1)

	reServiceDLH := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
	reServiceDOP := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDOP)
	reServiceMS := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeMS)
	reServiceDCRT := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDCRT)
	reServiceDOSHUT := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDOSHUT)

	// DLH
	subtestData.mtoServiceItemShip1DLH = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model:    reServiceDLH,
			LinkOnly: true,
		},
		{
			Model:    subtestData.mtoShipment1,
			LinkOnly: true,
		},
	}, nil)

	// DOP
	mtoServiceItemShip1DOP := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model:    reServiceDOP,
			LinkOnly: true,
		},
		{
			Model:    subtestData.mtoShipment1,
			LinkOnly: true,
		},
	}, nil)

	subtestData.mtoShipment2 = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		}}, nil)
	subtestData.mtoShipment2.PrimeEstimatedWeight = &subtestData.estimatedWeight
	suite.MustSave(&subtestData.mtoShipment2)

	// DLH
	subtestData.mtoServiceItemShip2DLH = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model:    reServiceDLH,
			LinkOnly: true,
		},
		{
			Model:    subtestData.mtoShipment2,
			LinkOnly: true,
		},
	}, nil)

	// MS
	subtestData.mtoServiceItemMS = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model:    reServiceMS,
			LinkOnly: true,
		},
	}, nil)

	subtestData.paramKeyWeightEstimated = factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameWeightEstimated,
				Description: "estimated weight",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		},
	}, nil)

	subtestData.paramKeyRequestedPickupDate = factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameRequestedPickupDate,
				Description: "requested pickup date",
				Type:        models.ServiceItemParamTypeDate,
				Origin:      models.ServiceItemParamOriginPrime,
			},
		},
	}, nil)

	// DLH
	factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItemShip1DLH.ReService,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyWeightEstimated,
			LinkOnly: true,
		},
		{
			Model: models.ServiceParam{
				IsOptional: true,
			},
		},
	}, nil)

	// DLH
	factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItemShip1DLH.ReService,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyRequestedPickupDate,
			LinkOnly: true,
		},
	}, nil)

	// DOP
	factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    mtoServiceItemShip1DOP.ReService,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyWeightEstimated,
			LinkOnly: true,
		},
		{
			Model: models.ServiceParam{
				IsOptional: true,
			},
		},
	}, nil)

	subtestData.paramKeyMTOAvailableToPrimeAt = factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameMTOAvailableToPrimeAt,
				Description: "prime mto made available date",
				Type:        models.ServiceItemParamTypeDate,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)

	factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItemMS.ReService,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyMTOAvailableToPrimeAt,
			LinkOnly: true,
		},
	}, nil)

	subtestData.shuttleEstimatedWeight = unit.Pound(400)
	subtestData.shuttleActualWeight = unit.Pound(450)
	subtestData.mtoServiceItemShuttle = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model:    reServiceDOSHUT,
			LinkOnly: true,
		},
		{
			Model:    subtestData.mtoShipment2,
			LinkOnly: true,
		},
		{
			Model: models.MTOServiceItem{
				EstimatedWeight: &subtestData.shuttleEstimatedWeight,
				ActualWeight:    &subtestData.shuttleActualWeight,
			},
		},
	}, nil)

	// DOSHUT estimated weight
	factory.BuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.mtoServiceItemShuttle.ReService,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyWeightEstimated,
			LinkOnly: true,
		},
		{
			Model: models.ServiceParam{
				IsOptional: true,
			},
		},
	}, nil)
	subtestData.mtoServiceItemCrate1 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model:    reServiceDCRT,
			LinkOnly: true,
		},
		{
			Model:    subtestData.mtoShipment2,
			LinkOnly: true,
		},
	}, nil)
	_ = factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItemDimension{
				Type: models.DimensionTypeCrate,
				// These dimensions are chosen to overflow 32bit ints if multiplied, and give a fractional result
				// when converted to cubic feet.
				Length:    16*12*1000 + 1000,
				Height:    8 * 12 * 1000,
				Width:     8 * 12 * 1000,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		{
			Model:    subtestData.mtoServiceItemCrate1,
			LinkOnly: true,
		},
	}, nil)
	_ = factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItemDimension{
				Type:      models.DimensionTypeItem,
				Length:    12000,
				Height:    12000,
				Width:     12000,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		{
			Model:    subtestData.mtoServiceItemCrate1,
			LinkOnly: true,
		},
	}, nil)
	subtestData.mtoServiceItemCrate2 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.move,
			LinkOnly: true,
		},
		{
			Model:    reServiceDCRT,
			LinkOnly: true,
		},
		{
			Model:    subtestData.mtoShipment2,
			LinkOnly: true,
		},
	}, nil)
	_ = factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItemDimension{
				Type: models.DimensionTypeCrate,
				// These dimensions are chosen to overflow 32bit ints if multiplied, and give a fractional result
				// when converted to cubic feet.
				Length:    7000,
				Height:    7000,
				Width:     7000,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		{
			Model:    subtestData.mtoServiceItemCrate2,
			LinkOnly: true,
		},
	}, nil)
	_ = factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItemDimension{
				Type:      models.DimensionTypeItem,
				Length:    6000,
				Height:    6000,
				Width:     6000,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		},
		{
			Model:    subtestData.mtoServiceItemCrate2,
			LinkOnly: true,
		},
	}, nil)
	subtestData.paramKeyDimensionHeight = factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameDimensionHeight,
				Description: "height",
				Type:        models.ServiceItemParamTypeDecimal,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)
	factory.BuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    reServiceDCRT,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyDimensionHeight,
			LinkOnly: true,
		},
	}, nil)
	subtestData.paramKeyDimensionWidth = factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameDimensionWidth,
				Description: "width",
				Type:        models.ServiceItemParamTypeDecimal,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)
	factory.BuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    reServiceDCRT,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyDimensionWidth,
			LinkOnly: true,
		},
	}, nil)
	subtestData.paramKeyDimensionLength = factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameDimensionLength,
				Description: "length",
				Type:        models.ServiceItemParamTypeDecimal,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)
	factory.BuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    reServiceDCRT,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyDimensionLength,
			LinkOnly: true,
		},
	}, nil)
	subtestData.paramKeyCubicFeetBilled = factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameCubicFeetBilled,
				Description: "cubic feet billed",
				Type:        models.ServiceItemParamTypeDecimal,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		},
	}, nil)

	factory.BuildServiceParam(suite.DB(), []factory.Customization{
		{
			Model:    reServiceDCRT,
			LinkOnly: true,
		},
		{
			Model:    subtestData.paramKeyCubicFeetBilled,
			LinkOnly: true,
		},
	}, nil)

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
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip2DLH, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		var estimatedWeightStr string
		estimatedWeightStr, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyWeightEstimated.Key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(subtestData.estimatedWeight.Int())
		suite.Equal(expected, estimatedWeightStr)

		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShuttle, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemCrate1, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		var cubicFeet string
		cubicFeet, err = paramLookupService1.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyCubicFeetBilled.Key)
		suite.FatalNoError(err)
		suite.Equal("1029.33", cubicFeet)

		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemCrate2, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
		suite.NoError(err)

		cubicFeet, err = paramLookupService2.ServiceParamValue(suite.AppContextForTest(), subtestData.paramKeyCubicFeetBilled.Key)
		suite.FatalNoError(err)
		suite.Equal("4.00", cubicFeet)
	})

	// Estimated Weight
	suite.Run("Shipment 1 estimated weight", func() {
		subtestData := suite.makeSubtestData()
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
		paramLookupService1, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip1DLH, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip2DLH, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
		paramLookupService2, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemShip2DLH, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
		paramLookupService3, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, subtestData.mtoServiceItemMS, subtestData.paymentRequest.ID, subtestData.paymentRequest.MoveTaskOrderID, &paramCache)
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
