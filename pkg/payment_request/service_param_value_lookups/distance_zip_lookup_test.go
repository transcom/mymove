package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/route/mocks"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestDistanceLookup() {
	key := models.ServiceItemParamNameDistanceZip

	suite.Run("Calculate transit zip distance", func() {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				PickupAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "33607",
					},
				}),
				DestinationAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "90210",
					},
				}),
			}),
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

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
		ppmShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())

		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      models.Address{PostalCode: ppmShipment.PickupPostalCode},
			DestinationAddress: models.Address{PostalCode: ppmShipment.DestinationPostalCode},
		}

		appContext := suite.AppContextForTest()
		distance, err := distanceZipLookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &ppmShipment.ShipmentID,
		})
		suite.NoError(err)

		planner := suite.planner.(*mocks.Planner)
		planner.AssertCalled(suite.T(), "ZipTransitDistance", appContext, ppmShipment.PickupPostalCode, ppmShipment.DestinationPostalCode)

		err = suite.DB().Reload(&ppmShipment.Shipment)
		suite.NoError(err)

		suite.Equal(fmt.Sprintf("%d", defaultZipDistance), distance)
		suite.Equal(unit.Miles(defaultZipDistance), *ppmShipment.Shipment.Distance)
	})

	suite.Run("Call ZipTransitDistance on PPMs with shipments that have a distance", func() {
		miles := unit.Miles(defaultZipDistance)
		ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Distance:     &miles,
				ShipmentType: models.MTOShipmentTypePPM,
			},
		})

		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      models.Address{PostalCode: ppmShipment.PickupPostalCode},
			DestinationAddress: models.Address{PostalCode: ppmShipment.DestinationPostalCode},
		}

		appContext := suite.AppContextForTest()
		distance, err := distanceZipLookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &ppmShipment.ShipmentID,
		})
		suite.NoError(err)

		planner := suite.planner.(*mocks.Planner)
		planner.AssertCalled(suite.T(), "ZipTransitDistance", appContext, ppmShipment.PickupPostalCode, ppmShipment.DestinationPostalCode)

		err = suite.DB().Reload(&ppmShipment.Shipment)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZipDistance), *ppmShipment.Shipment.Distance)
		suite.Equal(fmt.Sprintf("%d", defaultZipDistance), distance)
	})

	suite.Run("Do not call ZipTransitDistance on PPMs with shipments that have a distance", func() {
		miles := unit.Miles(defaultZipDistance)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Distance:     &miles,
				ShipmentType: models.MTOShipmentTypeHHG,
			},
		})

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
		planner.AssertNotCalled(suite.T(), "ZipTransitDistance", appContext, shipment.PickupAddress.PostalCode, shipment.DestinationAddress.PostalCode)

		err = suite.DB().Reload(&shipment)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZipDistance), *shipment.Distance)
		suite.Equal(fmt.Sprintf("%d", defaultZipDistance), distance)
	})

	suite.Run("Sucessfully updates mtoShipment distance when the pickup and destination zips are the same", func() {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				PickupAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "90211",
					},
				}),
				DestinationAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "90210",
					},
				}),
			}),
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

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
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				PickupAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "33607",
					},
				}),
				DestinationAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "90210",
					},
				}),
			}),
		})
		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		// DLH
		reServiceDLH := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
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
		serviceItemParamKey1 := testdatagen.FetchOrMakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameDistanceZip,
				Description: "zip distance",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		})

		_ = testdatagen.FetchOrMakeServiceParam(suite.DB(), testdatagen.Assertions{
			ServiceParam: models.ServiceParam{
				ServiceID:             mtoServiceItemDLH.ReServiceID,
				ServiceItemParamKeyID: serviceItemParamKey1.ID,
				ServiceItemParamKey:   serviceItemParamKey1,
			},
		})

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
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				PickupAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "33",
					},
				}),
				DestinationAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "90103",
					},
				}),
			}),
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "Shipment must have valid pickup zipcode")
	})

	suite.Run("returns error if the destination zipcode isn't at least 5 digits", func() {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				PickupAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "33607",
					},
				}),
				DestinationAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "901",
					},
				}),
			}),
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "Shipment must have valid destination zipcode")
	})

	suite.Run("returns a not found error if the service item shipment id doesn't exist", func() {
		distanceZipLookup := DistanceZipLookup{
			PickupAddress:      testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{Stub: true}),
			DestinationAddress: testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true}),
		}

		mtoShipmentID := uuid.Must(uuid.NewV4())

		_, err := distanceZipLookup.lookup(suite.AppContextForTest(), &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &mtoShipmentID,
		})

		suite.Error(err)
		suite.Equal(fmt.Sprintf("ID: %s not found looking for MTOShipmentID", mtoShipmentID), err.Error())
	})
}
