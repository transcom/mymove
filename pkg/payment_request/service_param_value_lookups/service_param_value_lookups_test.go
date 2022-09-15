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
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

const defaultZipDistance = 1234

type ServiceParamValueLookupsSuite struct {
	*testingsuite.PopTestSuite
	planner route.Planner
}

func TestServiceParamValueLookupsSuite(t *testing.T) {
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(defaultZipDistance, nil)
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(defaultZipDistance, nil)

	ts := &ServiceParamValueLookupsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		planner:      planner,
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithAllWeights(estimatedWeight *unit.Pound, originalWeight *unit.Pound, reweighWeight *unit.Pound, adjustedWeight *unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			Move: move,
			ReService: models.ReService{
				Code: code,
				Name: string(code),
			},
			MTOShipment: models.MTOShipment{
				PrimeEstimatedWeight: estimatedWeight,
				PrimeActualWeight:    originalWeight,
				BillableWeightCap:    adjustedWeight,
				ShipmentType:         shipmentType,
			},
		})

	if reweighWeight != nil {
		var shipment models.MTOShipment
		suite.NoError(suite.DB().Find(&shipment, *mtoServiceItem.MTOShipmentID))

		_ = testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, *reweighWeight)
	}

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	return mtoServiceItem, paymentRequest, paramLookup
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithEstimatedWeightForPPM(estimatedWeight *unit.Pound, originalWeight *unit.Pound, code models.ReServiceCode) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	destAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			Move: move,
			ReService: models.ReService{
				Code: code,
				Name: string(code),
			},
			MTOShipment: models.MTOShipment{
				PrimeEstimatedWeight: estimatedWeight,
				PrimeActualWeight:    originalWeight,
				ShipmentType:         models.MTOShipmentTypePPM,
				PickupAddress:        &pickupAddress,
				DestinationAddress:   &destAddress,
				PPMShipment:          &models.PPMShipment{},
			},
			Address: pickupAddress,
			MTOServiceItem: models.MTOServiceItem{
				EstimatedWeight: estimatedWeight,
			},
		})

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
	}

	paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	return mtoServiceItem, paymentRequest, paramLookup
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithOriginalWeightOnly(originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(nil, &originalWeight, nil, nil, code, shipmentType)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithWeight(estimatedWeight unit.Pound, originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(&estimatedWeight, &originalWeight, nil, nil, code, shipmentType)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithReweigh(reweighWeight unit.Pound, originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(nil, &originalWeight, &reweighWeight, nil, code, shipmentType)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithAdjustedWeight(adjustedWeight *unit.Pound, originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(nil, &originalWeight, nil, adjustedWeight, code, shipmentType)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithShuttleWeight(itemEstimatedWeight unit.Pound, itemOriginalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			Move: move,
			ReService: models.ReService{
				Code: code,
				Name: string(code),
			},
			MTOServiceItem: models.MTOServiceItem{
				EstimatedWeight: &itemEstimatedWeight,
				ActualWeight:    &itemOriginalWeight,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: shipmentType,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: move,
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	return mtoServiceItem, paymentRequest, paramLookup
}

func (suite *ServiceParamValueLookupsSuite) TestServiceParamValueLookup() {
	suite.Run("contract passed in", func() {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

		suite.FatalNoError(err)
		suite.Equal(ghcrateengine.DefaultContractCode, paramLookup.ContractCode)
	})

	suite.Run("MTOServiceItem passed in", func() {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

		suite.FatalNoError(err)
		suite.Equal(mtoServiceItem.ID, paramLookup.MTOServiceItemID)
		suite.NotNil(paramLookup.MTOServiceItem)
		suite.Equal(mtoServiceItem.MoveTaskOrderID, paramLookup.MTOServiceItem.MoveTaskOrderID)
	})

	// Setup data for testing service items not dependent on the shipment
	serviceCodesWithoutShipment := []models.ReServiceCode{
		models.ReServiceCodeCS,
		models.ReServiceCodeMS,
	}

	for _, code := range serviceCodesWithoutShipment {
		suite.Run(fmt.Sprintf("MTOShipment not looked up for %s", code), func() {
			mtoServiceItem := testdatagen.MakeMTOServiceItemBasic(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: code,
					Name: string(code),
				},
			})

			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if rpdl, ok := paramLookup.lookups[models.ServiceItemParamNameRequestedPickupDate].(RequestedPickupDateLookup); ok {
				suite.Equal(uuid.Nil, rpdl.MTOShipment.ID)
			} else {
				suite.Fail("lookup not RequestedPickupDateLookup type")
			}
			if zpal, ok := paramLookup.lookups[models.ServiceItemParamNameZipPickupAddress].(ZipAddressLookup); ok {
				suite.Equal(uuid.Nil, zpal.Address.ID)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipDestAddress].(ZipAddressLookup); ok {
				suite.Equal(uuid.Nil, zdal.Address.ID)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
		})
	}

	suite.Run("MTOShipment is looked up for other service items", func() {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
				Name: models.ReServiceCodeDLH.String(),
			},
		})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		suite.NotNil(paramLookup.MTOServiceItem)
		if rpdl, ok := paramLookup.lookups[models.ServiceItemParamNameRequestedPickupDate].(RequestedPickupDateLookup); ok {
			suite.Equal(*mtoServiceItem.MTOShipmentID, rpdl.MTOShipment.ID)
		} else {
			suite.Fail("lookup not RequestedPickupDateLookup type")
		}
	})

	suite.Run("DestinationAddress is looked up for other service items", func() {
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: models.ReServiceCodeDLH.String(),
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDUPK,
					Name: models.ReServiceCodeDUPK.String(),
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipDestAddress].(ZipAddressLookup); ok {
				suite.Equal(mtoServiceItem.MTOShipment.DestinationAddress.PostalCode, zdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
		}
	})

	suite.Run("DestinationAddress is not required for service items like domestic pack", func() {
		servicesToTest := []models.ReServiceCode{models.ReServiceCodeDPK, models.ReServiceCodeDNPK}
		for _, service := range servicesToTest {
			mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: service,
					Name: service.String(),
				},
			})

			mtoShipment := mtoServiceItem.MTOShipment
			mtoShipment.DestinationAddressID = nil
			suite.DB().Save(&mtoShipment)

			_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)
		}
	})

	suite.Run("PickupAddress is looked up for other service items", func() {
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: models.ReServiceCodeDLH.String(),
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDPK,
					Name: models.ReServiceCodeDPK.String(),
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zpal, ok := paramLookup.lookups[models.ServiceItemParamNameZipPickupAddress].(ZipAddressLookup); ok {
				suite.Equal(mtoServiceItem.MTOShipment.PickupAddress.PostalCode, zpal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
		}
	})

	suite.Run("PickupAddress is not required for service items like domestic unpack", func() {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDUPK,
				Name: models.ReServiceCodeDUPK.String(),
			},
		})

		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.PickupAddressID = nil
		suite.DB().Save(&mtoShipment)

		_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)
	})

	suite.Run("Correct addresses are used for NTS and NTS-release shipments", func() {
		// Make a move and service for reuse.
		move := testdatagen.MakeDefaultMove(suite.DB())
		reService := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			Move: move,
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
		})

		// NTS should have a pickup address and storage facility address.
		pickupPostalCode := "29212"
		pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: pickupPostalCode,
			},
		})
		storageFacilityPostalCode := "30907"
		storageFacility := testdatagen.MakeStorageFacility(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: storageFacilityPostalCode,
			},
		})
		ntsServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move:      move,
			ReService: reService,
			MTOShipment: models.MTOShipment{
				ShipmentType:    models.MTOShipmentTypeHHGIntoNTSDom,
				PickupAddress:   &pickupAddress,
				StorageFacility: &storageFacility,
			},
		})

		// Check to see if the distance lookup got the expected NTS addresses.
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, ntsServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)
		if dz3l, ok := paramLookup.lookups[models.ServiceItemParamNameDistanceZip].(DistanceZipLookup); ok {
			suite.Equal(pickupPostalCode, dz3l.PickupAddress.PostalCode)
			suite.Equal(storageFacilityPostalCode, dz3l.DestinationAddress.PostalCode)
		} else {
			suite.Fail("lookup not DistanceZipLookup type")
		}

		// NTS-Release should have a storage facility address and destination address.
		destinationPostalCode := "29440"
		destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: destinationPostalCode,
			},
		})
		ntsrServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move:      move,
			ReService: reService,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				DestinationAddress: &destinationAddress,
				StorageFacility:    &storageFacility,
			},
		})

		// Check to see if the distance lookup got the expected NTS-Release addresses.
		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, ntsrServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)
		if dz3l, ok := paramLookup.lookups[models.ServiceItemParamNameDistanceZip].(DistanceZipLookup); ok {
			suite.Equal(storageFacilityPostalCode, dz3l.PickupAddress.PostalCode)
			suite.Equal(destinationPostalCode, dz3l.DestinationAddress.PostalCode)
		} else {
			suite.Fail("lookup not DistanceZipLookup type")
		}
	})

	suite.Run("SITDestinationAddress is looked up for destination sit", func() {
		sitFinalDestAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDDASIT,
					Name: models.ReServiceCodeDDASIT.String(),
				},
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &sitFinalDestAddress.ID,
					SITDestinationFinalAddress:   &sitFinalDestAddress,
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
					Name: models.ReServiceCodeDDDSIT.String(),
				},
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &sitFinalDestAddress.ID,
					SITDestinationFinalAddress:   &sitFinalDestAddress,
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
					Name: models.ReServiceCodeDDFSIT.String(),
				},
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &sitFinalDestAddress.ID,
					SITDestinationFinalAddress:   &sitFinalDestAddress,
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipSITDestHHGFinalAddress].(ZipAddressLookup); ok {
				suite.Equal(mtoServiceItem.SITDestinationFinalAddress.PostalCode, zdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipSitAddress destination type")
			}
		}
	})

	suite.Run("SITDestinationAddress is not loaded non sit", func() {
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: models.ReServiceCodeDLH.String(),
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDSH,
					Name: models.ReServiceCodeDSH.String(),
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipSITDestHHGFinalAddress].(ZipAddressLookup); ok {
				suite.Equal("", zdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipSitAddress destination type")
			}
		}
	})

	suite.Run("Non-basic MTOServiceItem is missing a MTOShipmentID", func() {
		badMTOServiceItem := models.MTOServiceItem{ID: uuid.Must(uuid.NewV4())}
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, badMTOServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Not found the shipment service item is missing a MTOShipmentID")
		var expected *ServiceItemParamKeyData
		suite.Equal(expected, paramLookup)
	})

	suite.Run("Non-basic MTOServiceItem has a MTOShipmentID that is not found", func() {
		badMTOServiceItem := models.MTOServiceItem{ID: uuid.Must(uuid.NewV4()), MTOShipmentID: models.UUIDPointer(uuid.Must(uuid.NewV4()))}
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, badMTOServiceItem, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found looking for MTOShipment")
		var expected *ServiceItemParamKeyData
		suite.Equal(expected, paramLookup)
	})
}
