package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestDistanceZip3Lookup() {
	key := models.ServiceItemParamNameDistanceZip3

	suite.Run("Calculate zip3 distance", func() {
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

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip3Distance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZip3Distance), *mtoShipment.Distance)
	})

	suite.Run("Doesn't update mtoShipment distance when the pickup and destination zip3s are the same", func() {
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

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip3Distance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		suite.NoError(err)
		suite.Nil(mtoShipment.Distance)
	})

	suite.Run("Calculate zip3 distance with param cache", func() {
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
		reServiceDLH := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
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

		paramCache := NewServiceParamsCache()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDLH.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip3Distance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().Find(&mtoShipment, mtoServiceItemDLH.MTOShipmentID)
		suite.NoError(err)

		suite.Equal(unit.Miles(defaultZip3Distance), *mtoShipment.Distance)

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

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
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

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "Shipment must have valid destination zipcode")
	})
}
