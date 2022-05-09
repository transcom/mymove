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

func (suite *ServiceParamValueLookupsSuite) TestDistanceZip5Lookup() {
	key := models.ServiceItemParamNameDistanceZip5

	suite.Run("Calculate zip5 distance", func() {
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

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
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

	suite.Run("Calculate zip5 distance lookup without a saved service item", func() {
		ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				PickupPostalCode:      "33607",
				DestinationPostalCode: "33609",
			},
		})

		distanceZip5Lookup := DistanceZip5Lookup{
			PickupAddress:      models.Address{PostalCode: ppmShipment.PickupPostalCode},
			DestinationAddress: models.Address{PostalCode: ppmShipment.DestinationPostalCode},
		}

		appContext := suite.AppContextForTest()
		distance, err := distanceZip5Lookup.lookup(appContext, &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &ppmShipment.ShipmentID,
		})
		suite.NoError(err)

		planner := suite.planner.(*mocks.Planner)
		planner.AssertCalled(suite.T(), "Zip5TransitDistance", appContext, ppmShipment.PickupPostalCode, ppmShipment.DestinationPostalCode)

		err = suite.DB().Reload(&ppmShipment.Shipment)
		suite.NoError(err)

		suite.Equal(fmt.Sprintf("%d", defaultZip5Distance), distance)
		suite.Equal(unit.Miles(defaultZip5Distance), *ppmShipment.Shipment.Distance)
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

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
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
		distanceZip3Lookup := DistanceZip5Lookup{
			PickupAddress:      testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{Stub: true}),
			DestinationAddress: testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{Stub: true}),
		}

		mtoShipmentID := uuid.Must(uuid.NewV4())

		_, err := distanceZip3Lookup.lookup(suite.AppContextForTest(), &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &mtoShipmentID,
		})

		suite.Error(err)
		suite.Equal(fmt.Sprintf("ID: %s not found looking for MTOShipmentID", mtoShipmentID), err.Error())
	})
}
