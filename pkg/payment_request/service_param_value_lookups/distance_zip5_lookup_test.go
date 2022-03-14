package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestDistanceZip5Lookup() {
	key := models.ServiceItemParamNameDistanceZip5

	suite.Run("golden path", func() {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
				PickupAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "33607",
					},
				}),
				DestinationAddress: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						PostalCode: "33609",
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

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip5Distance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().
			Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		suite.NoError(err)
		suite.Equal(unit.Miles(defaultZip5Distance), *mtoShipment.Distance)
	})

	suite.Run("doesn't update mtoShipment distance when the pickup and destination zip3s are different", func() {
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
				Move: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip5Distance)
		suite.Equal(expected, distanceStr)

		var mtoShipment models.MTOShipment
		err = suite.DB().
			Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		suite.NoError(err)
		suite.Nil(mtoShipment.Distance)
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
