package serviceparamvaluelookups

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestDistanceLookup() {
	key := models.ServiceItemParamNameDistanceZip

	suite.Run("Calculate transit zip distance", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "33607",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZipDistance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZipDistance), *mtoShipment.Distance)
	})

	suite.Run("Calculate zip distance lookup without a saved service item", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      models.Address{PostalCode: ppmShipment.PickupAddress.PostalCode},
			DestinationAddress: models.Address{PostalCode: ppmShipment.DestinationAddress.PostalCode},
		}

		appContext := suite.AppContextForTest()
		distance, err := distanceZipLookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &ppmShipment.ShipmentID,
		})
		suite.NoError(err)

		planner := suite.planner.(*mocks.Planner)
		planner.AssertCalled(suite.T(), "ZipTransitDistance", appContext, ppmShipment.PickupAddress.PostalCode, ppmShipment.DestinationAddress.PostalCode, false, false)

		err = suite.DB().Reload(&ppmShipment.Shipment)
		suite.NoError(err)

		suite.Equal(fmt.Sprintf("%d", defaultZipDistance), distance)
		suite.Equal(unit.Miles(defaultZipDistance), *ppmShipment.Shipment.Distance)
	})

	suite.Run("Call ZipTransitDistance on non-PPMs with shipments that have a distance", func() {
		miles := unit.Miles(defaultZipDistance)
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Distance: &miles,
				},
			},
		}, nil)
		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      models.Address{PostalCode: ppmShipment.PickupAddress.PostalCode},
			DestinationAddress: models.Address{PostalCode: ppmShipment.DestinationAddress.PostalCode},
		}

		appContext := suite.AppContextForTest()
		distance, err := distanceZipLookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &ppmShipment.ShipmentID,
		})
		suite.NoError(err)

		planner := suite.planner.(*mocks.Planner)
		planner.AssertCalled(suite.T(), "ZipTransitDistance", appContext, ppmShipment.PickupAddress.PostalCode, ppmShipment.DestinationAddress.PostalCode, false, false)

		err = suite.DB().Reload(&ppmShipment.Shipment)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZipDistance), *ppmShipment.Shipment.Distance)
		suite.Equal(fmt.Sprintf("%d", defaultZipDistance), distance)
	})

	suite.Run("Do not call ZipTransitDistance on PPMs with shipments that have a distance", func() {
		miles := unit.Miles(defaultZipDistance)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Distance:     &miles,
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
		}, nil)

		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      models.Address{PostalCode: shipment.PickupAddress.PostalCode},
			DestinationAddress: models.Address{PostalCode: shipment.DestinationAddress.PostalCode},
		}

		appContext := suite.AppContextForTest()
		distance, err := distanceZipLookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &shipment.ID,
		})
		suite.NoError(err)

		planner := suite.planner.(*mocks.Planner)
		planner.AssertNotCalled(suite.T(), "ZipTransitDistance", appContext, shipment.PickupAddress.PostalCode, shipment.DestinationAddress.PostalCode, false, false)

		err = suite.DB().Reload(&shipment)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZipDistance), *shipment.Distance)
		suite.Equal(fmt.Sprintf("%d", defaultZipDistance), distance)
	})

	suite.Run("Sucessfully updates mtoShipment distance when the pickup and destination zips are the same", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "90211",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZipDistance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZipDistance), *mtoShipment.Distance)
	})

	suite.Run("Calculate zip distance with param cache", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "33607",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "90210",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		// DLH
		reServiceDLH := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)

		estimatedWeight := unit.Pound(2048)

		// DLH
		mtoServiceItemDLH := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    reServiceDLH,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
		}, nil)
		mtoServiceItemDLH.MoveTaskOrderID = paymentRequest.MoveTaskOrderID
		mtoServiceItemDLH.MoveTaskOrder = paymentRequest.MoveTaskOrder
		suite.MustSave(&mtoServiceItemDLH)

		// ServiceItemParamNameActualPickupDate
		serviceItemParamKey1 := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameDistanceZip,
					Description: "zip distance",
					Type:        models.ServiceItemParamTypeInteger,
					Origin:      models.ServiceItemParamOriginSystem,
				},
			},
		}, nil)

		factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItemDLH.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey1,
				LinkOnly: true,
			},
		}, nil)

		paramCache := NewServiceParamsCache()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDLH, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZipDistance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().Find(&mtoShipment, mtoServiceItemDLH.MTOShipmentID)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZipDistance), *mtoShipment.Distance)

		// Verify value from paramCache
		paramCacheValue := paramCache.ParamValue(*mtoServiceItemDLH.MTOShipmentID, key)
		suite.Equal(expected, *paramCacheValue)
	})

	suite.Run("returns error if the pickup zipcode isn't at least 5 digits", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "33",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "90103",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "Shipment must have valid pickup zipcode")
	})

	suite.Run("returns error if the destination zipcode isn't at least 5 digits", func() {

		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "33607",
				},
				Type: &factory.Addresses.PickupAddress,
			},
			{
				Model: models.Address{
					PostalCode: "901",
				},
				Type: &factory.Addresses.DeliveryAddress,
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "Shipment must have valid destination zipcode")
	})

	suite.Run("returns a not found error if the service item shipment id doesn't exist", func() {
		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      factory.BuildAddress(nil, nil, nil),
			DestinationAddress: factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2}),
		}

		mtoShipmentID := uuid.Must(uuid.NewV4())

		_, err := distanceZipLookup.lookup(suite.AppContextForTest(), &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &mtoShipmentID,
		})

		suite.Error(err)
		suite.Equal(fmt.Sprintf("ID: %s not found looking for MTOShipmentID", mtoShipmentID), err.Error())
	})

	suite.Run("sets distance to one when origin and destination postal codes are the same", func() {
		MTOShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: factory.BuildAddress(suite.DB(), []factory.Customization{
					{
						Model: models.Address{
							PostalCode: "90211",
						},
					},
				}, nil),
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model: factory.BuildAddress(suite.DB(), []factory.Customization{
					{
						Model: models.Address{
							PostalCode: "90211",
						},
					},
				}, nil),
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
		}, nil)

		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      models.Address{PostalCode: MTOShipment.PickupAddress.PostalCode},
			DestinationAddress: models.Address{PostalCode: MTOShipment.DestinationAddress.PostalCode},
		}

		distance, err := distanceZipLookup.lookup(suite.AppContextForTest(), &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &MTOShipment.ID,
		})

		//Check if distance equal 1
		suite.Equal("1", distance)
		suite.FatalNoError(err)

	})
}
