// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package serviceparamvaluelookups

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/route/mocks"
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

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithAllWeights(estimatedWeight *unit.Pound, originalWeight *unit.Pound, reweighWeight *unit.Pound, adjustedWeight *unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType, diversion bool, divertedFromShipmentID *uuid.UUID) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			EndDate: time.Now().Add(24 * time.Hour),
		},
	})
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: code,
				Name: string(code),
			},
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight:   estimatedWeight,
				PrimeActualWeight:      originalWeight,
				BillableWeightCap:      adjustedWeight,
				ShipmentType:           shipmentType,
				Diversion:              diversion,
				DivertedFromShipmentID: divertedFromShipmentID,
			},
		},
	}, nil)

	if reweighWeight != nil {
		var shipment models.MTOShipment
		suite.NoError(suite.DB().Find(&shipment, *mtoServiceItem.MTOShipmentID))

		_ = testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, shipment, *reweighWeight)
	}

	paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		},
	}, nil)

	paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	return mtoServiceItem, paymentRequest, paramLookup
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithEstimatedWeightForPPM(estimatedWeight *unit.Pound, originalWeight *unit.Pound, code models.ReServiceCode) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	move := factory.BuildMove(suite.DB(), nil, nil)
	pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destAddress := factory.BuildAddress(suite.DB(), nil, nil)
	mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: code,
				Name: string(code),
			},
		},
		{
			Model:    pickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    destAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: estimatedWeight,
				PrimeActualWeight:    originalWeight,
			},
		},
		{
			Model: models.MTOServiceItem{
				EstimatedWeight: estimatedWeight,
			},
		},
	}, nil)

	// BuildMTOShipment does not populate the addresses for PPM
	// shipments, so override ShipmentType after creation
	mtoShipment := mtoServiceItem.MTOShipment
	mtoShipment.ShipmentType = models.MTOShipmentTypePPM
	suite.MustSave(&mtoShipment)

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
	}

	serviceItemLookups := InitializeLookups(suite.AppContextForTest(), mtoShipment, mtoServiceItem)
	// i don't think this function gets called for PPMs, but need to verify
	//paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	//suite.FatalNoError(err)
	paramLookup := NewServiceItemParamKeyData(suite.planner, serviceItemLookups, mtoServiceItem, mtoShipment, testdatagen.DefaultContractCode)

	return mtoServiceItem, paymentRequest, &paramLookup
}

// Create a parent and child diverted shipment chain
func (suite *ServiceParamValueLookupsSuite) setupDivertedShipmentChainServiceItemsWithAllWeights(parentEstimatedWeight *unit.Pound, childEstimatedWeight *unit.Pound, parentActualWeight *unit.Pound, childActualWeight *unit.Pound, parentReweighWeight *unit.Pound, childReweighWeight *unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.MTOServiceItem, models.PaymentRequest, models.PaymentRequest, *ServiceItemParamKeyData, *ServiceItemParamKeyData) {
	// Create contract year
	testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			EndDate: time.Now().Add(24 * time.Hour),
		},
	})
	// Create the parent shipment
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{Model: move, LinkOnly: true},
		{Model: models.MTOShipment{
			PrimeEstimatedWeight:   parentEstimatedWeight,
			PrimeActualWeight:      parentActualWeight,
			ShipmentType:           shipmentType,
			Diversion:              true,
			DivertedFromShipmentID: nil,
		}},
	}, nil)

	// Create the child shipment
	childShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{Model: move, LinkOnly: true},
		{Model: models.MTOShipment{
			PrimeEstimatedWeight:   childEstimatedWeight,
			PrimeActualWeight:      childActualWeight,
			Diversion:              true,
			DivertedFromShipmentID: &parentShipment.ID,
		}},
	}, nil)

	// Create reweigh for parent shipment if provided
	if parentReweighWeight != nil {
		_ = testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, parentShipment, *parentReweighWeight)
	}

	// Create reweigh for child shipment if provided
	if childReweighWeight != nil {
		_ = testdatagen.MakeReweighForShipment(suite.DB(), testdatagen.Assertions{}, parentShipment, *childReweighWeight)
	}

	// Parent
	parentMtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{Model: parentShipment, LinkOnly: true},
		{Model: models.ReService{Code: code, Name: string(code)}},
	}, nil)

	parentPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{Model: move, LinkOnly: true},
		{Model: models.PaymentRequest{MoveTaskOrderID: parentMtoServiceItem.MoveTaskOrderID, SequenceNumber: 1}},
	}, nil)

	parentParamLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, parentMtoServiceItem, parentPaymentRequest.ID, parentPaymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	// Child
	childMtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{Model: childShipment, LinkOnly: true},
		{Model: models.ReService{Code: code, Name: string(code)}},
	}, nil)

	childPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{Model: move, LinkOnly: true},
		{Model: models.PaymentRequest{MoveTaskOrderID: childMtoServiceItem.MoveTaskOrderID, SequenceNumber: 2}},
	}, nil)

	childParamLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, childMtoServiceItem, childPaymentRequest.ID, childPaymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	return parentMtoServiceItem, childMtoServiceItem, parentPaymentRequest, childPaymentRequest, parentParamLookup, childParamLookup
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithOriginalWeightOnly(originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(nil, &originalWeight, nil, nil, code, shipmentType, false, nil)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithWeight(estimatedWeight unit.Pound, originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(&estimatedWeight, &originalWeight, nil, nil, code, shipmentType, false, nil)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithReweigh(reweighWeight unit.Pound, originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(nil, &originalWeight, &reweighWeight, nil, code, shipmentType, false, nil)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithAdjustedWeight(adjustedWeight *unit.Pound, originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(nil, &originalWeight, nil, adjustedWeight, code, shipmentType, false, nil)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithWeightOnDiversion(estimatedWeight unit.Pound, originalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType, diversion bool, divertedFromShipmentID *uuid.UUID) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	return suite.setupTestMTOServiceItemWithAllWeights(&estimatedWeight, &originalWeight, nil, nil, code, shipmentType, diversion, divertedFromShipmentID)
}

func (suite *ServiceParamValueLookupsSuite) setupTestDivertedShipmentChain(parentEstimatedWeight *unit.Pound, childEstimatedWeight *unit.Pound, parentOriginalWeight *unit.Pound, childOriginalWeight *unit.Pound, parentReweighWeight *unit.Pound, childReweighWeight *unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.MTOServiceItem, models.PaymentRequest, models.PaymentRequest, *ServiceItemParamKeyData, *ServiceItemParamKeyData) {
	return suite.setupDivertedShipmentChainServiceItemsWithAllWeights(parentEstimatedWeight, childEstimatedWeight, parentOriginalWeight, childOriginalWeight, parentReweighWeight, childReweighWeight, code, shipmentType)
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithShuttleWeight(itemEstimatedWeight unit.Pound, itemOriginalWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			EndDate: time.Now().Add(24 * time.Hour),
		},
	})
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: code,
				Name: string(code),
			},
		},
		{
			Model: models.MTOServiceItem{
				EstimatedWeight: &itemEstimatedWeight,
				ActualWeight:    &itemOriginalWeight,
			},
		},
		{
			Model: models.MTOShipment{
				ShipmentType: shipmentType,
			},
		},
	}, nil)

	paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		},
	}, nil)

	paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	return mtoServiceItem, paymentRequest, paramLookup
}

func (suite *ServiceParamValueLookupsSuite) TestServiceParamValueLookup() {
	suite.Run("contract passed in", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)

		suite.FatalNoError(err)
		suite.Equal(testdatagen.DefaultContractCode, paramLookup.ContractCode)
	})

	suite.Run("MTOServiceItem passed in", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)

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
			testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EndDate: time.Now().Add(24 * time.Hour),
				},
			})
			mtoServiceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: code,
						Name: string(code),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			})

			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: models.ReServiceCodeDLH.String(),
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		suite.NotNil(paramLookup.MTOServiceItem)
		if rpdl, ok := paramLookup.lookups[models.ServiceItemParamNameRequestedPickupDate].(RequestedPickupDateLookup); ok {
			suite.Equal(*mtoServiceItem.MTOShipmentID, rpdl.MTOShipment.ID)
		} else {
			suite.Fail("lookup not RequestedPickupDateLookup type")
		}
	})

	suite.Run("DestinationAddress is looked up for other service items", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		testData := []models.MTOServiceItem{
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
						Name: models.ReServiceCodeDLH.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDUPK,
						Name: models.ReServiceCodeDUPK.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipDestAddress].(ZipAddressLookup); ok {
				suite.Equal(mtoServiceItem.MTOShipment.DestinationAddress.PostalCode, zdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
		}
	})

	suite.Run("DestinationAddress will not change from when SIT Destination service items were approved", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		testData := []models.MTOServiceItem{
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
						Name: models.ReServiceCodeDLH.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeFSC,
						Name: models.ReServiceCodeFSC.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDDDSIT,
						Name: models.ReServiceCodeDDDSIT.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)

			originalAddress, err := getDestinationAddressForService(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, mtoServiceItem.MTOShipment)
			suite.FatalNoError(err)

			if sdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipDestAddress].(ZipAddressLookup); ok {
				suite.Equal(originalAddress.PostalCode, sdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
		}
	})

	suite.Run("DestinationAddress is new address when there's a DeliverAddressUpdate and partially approved SIT", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		testData := []models.MTOServiceItem{
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDDASIT,
						Name: models.ReServiceCodeDDASIT.String(),
					},
				},
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusSubmitted,
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDDDSIT,
						Name: models.ReServiceCodeDDDSIT.String(),
					},
				},
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusSubmitted,
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),

			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDDSFSC,
						Name: models.ReServiceCodeDDSFSC.String(),
					},
				},
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusApproved,
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDDFSIT,
						Name: models.ReServiceCodeDDFSIT.String(),
					},
				},
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusSubmitted,
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
		}

		addressUpdate := factory.BuildShipmentAddressUpdate(nil, nil, nil)

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
			suite.FatalNoError(err)

			mtoServiceItem.MTOShipment.DeliveryAddressUpdate = &addressUpdate

			suite.NotNil(paramLookup.MTOServiceItem)

			originalAddress, err := getDestinationAddressForService(suite.AppContextForTest(), models.ReServiceCodeDDDSIT, mtoServiceItem.MTOShipment)
			suite.FatalNoError(err)

			if sdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipDestAddress].(ZipAddressLookup); ok {
				suite.Equal(originalAddress.PostalCode, sdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}

			suite.Equal(originalAddress.PostalCode, addressUpdate.OriginalAddress.PostalCode)
		}
	})

	suite.Run("DestinationAddress is not required for service items like domestic pack", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		servicesToTest := []models.ReServiceCode{models.ReServiceCodeDPK, models.ReServiceCodeDNPK}
		for _, service := range servicesToTest {
			mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: service,
						Name: service.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			})

			mtoShipment := mtoServiceItem.MTOShipment
			mtoShipment.DestinationAddressID = nil
			suite.DB().Save(&mtoShipment)

			_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
			suite.FatalNoError(err)
		}
	})

	suite.Run("PickupAddress is looked up for other service items", func() {
		testData := []models.MTOServiceItem{
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
						Name: models.ReServiceCodeDLH.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDPK,
						Name: models.ReServiceCodeDPK.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
		}

		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDUPK,
					Name: models.ReServiceCodeDUPK.String(),
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.PickupAddressID = nil
		suite.DB().Save(&mtoShipment)

		_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
	})

	suite.Run("Correct addresses are used for NTS and NTS-release shipments", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		// Make a move and service for reuse.
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)

		// NTS should have a pickup address and storage facility address.
		pickupPostalCode := "29212"
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: pickupPostalCode,
				},
			},
		}, nil)
		storageFacilityPostalCode := "30907"
		storageFacility := factory.BuildStorageFacility(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: storageFacilityPostalCode,
				},
			},
		}, nil)
		ntsServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
				},
			},
		}, nil)

		// Check to see if the distance lookup got the expected NTS addresses.
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, ntsServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)
		suite.FatalNoError(err)
		if dz3l, ok := paramLookup.lookups[models.ServiceItemParamNameDistanceZip].(DistanceZipLookup); ok {
			suite.Equal(pickupPostalCode, dz3l.PickupAddress.PostalCode)
			suite.Equal(storageFacilityPostalCode, dz3l.DestinationAddress.PostalCode)
		} else {
			suite.Fail("lookup not DistanceZipLookup type")
		}

		// NTS-Release should have a storage facility address and destination address.
		destinationPostalCode := "29440"
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: destinationPostalCode,
				},
			},
		}, nil)
		ntsrServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				},
			},
		}, nil)

		// Check to see if the distance lookup got the expected NTS-Release addresses.
		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, ntsrServiceItem, uuid.Must(uuid.NewV4()), ntsrServiceItem.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		if dz3l, ok := paramLookup.lookups[models.ServiceItemParamNameDistanceZip].(DistanceZipLookup); ok {
			suite.Equal(storageFacilityPostalCode, dz3l.PickupAddress.PostalCode)
			suite.Equal(destinationPostalCode, dz3l.DestinationAddress.PostalCode)
		} else {
			suite.Fail("lookup not DistanceZipLookup type")
		}
	})

	suite.Run("SITDestinationAddress is looked up for destination sit", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		sitFinalDestAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		testData := []models.MTOServiceItem{
			factory.BuildMTOServiceItem(suite.DB(),
				[]factory.Customization{
					{
						Model: models.ReService{
							Code: models.ReServiceCodeDDASIT,
							Name: models.ReServiceCodeDDASIT.String(),
						},
					},
					{
						Model:    sitFinalDestAddress,
						LinkOnly: true,
						Type:     &factory.Addresses.SITDestinationFinalAddress,
					},
				}, []factory.Trait{
					factory.GetTraitAvailableToPrimeMove,
				}),
			factory.BuildMTOServiceItem(suite.DB(),
				[]factory.Customization{
					{
						Model: models.ReService{
							Code: models.ReServiceCodeDDDSIT,
							Name: models.ReServiceCodeDDDSIT.String(),
						},
					},
					{
						Model:    sitFinalDestAddress,
						LinkOnly: true,
						Type:     &factory.Addresses.SITDestinationFinalAddress,
					},
				}, []factory.Trait{
					factory.GetTraitAvailableToPrimeMove,
				}),
			factory.BuildMTOServiceItem(suite.DB(),
				[]factory.Customization{
					{
						Model: models.ReService{
							Code: models.ReServiceCodeDDFSIT,
							Name: models.ReServiceCodeDDFSIT.String(),
						},
					},
					{
						Model:    sitFinalDestAddress,
						LinkOnly: true,
						Type:     &factory.Addresses.SITDestinationFinalAddress,
					},
				}, []factory.Trait{
					factory.GetTraitAvailableToPrimeMove,
				}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		testData := []models.MTOServiceItem{
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDLH,
						Name: models.ReServiceCodeDLH.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
			factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.ReService{
						Code: models.ReServiceCodeDSH,
						Name: models.ReServiceCodeDSH.String(),
					},
				},
			}, []factory.Trait{
				factory.GetTraitAvailableToPrimeMove,
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, uuid.Must(uuid.NewV4()), mtoServiceItem.MoveTaskOrderID, nil)
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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		badMTOServiceItem := models.MTOServiceItem{ID: uuid.Must(uuid.NewV4()), MoveTaskOrderID: move.ID, MoveTaskOrder: move}
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, badMTOServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "Not found the shipment service item is missing a MTOShipmentID")
		var expected *ServiceItemParamKeyData
		suite.Equal(expected, paramLookup)
	})

	suite.Run("Non-basic MTOServiceItem has a MTOShipmentID that is not found", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		badMTOServiceItem := models.MTOServiceItem{ID: uuid.Must(uuid.NewV4()), MTOShipmentID: models.UUIDPointer(uuid.Must(uuid.NewV4()))}
		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, badMTOServiceItem, uuid.Must(uuid.NewV4()), move.ID, nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found looking for MTOShipment")
		var expected *ServiceItemParamKeyData
		suite.Equal(expected, paramLookup)
	})
}

func (suite *ServiceParamValueLookupsSuite) TestFetchContract() {
	setupTestData := func() {
		firstContract := testdatagen.MakeReContract(suite.DB(), testdatagen.Assertions{
			ReContract: models.ReContract{
				Code: "first",
			},
		})
		secondContract := testdatagen.MakeReContract(suite.DB(), testdatagen.Assertions{
			ReContract: models.ReContract{
				Code: "second",
			},
		})
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2021, 8, 31, 0, 0, 0, 0, time.UTC),
			},
			ReContract: firstContract,
		})
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC),
			},
			ReContract: firstContract,
		})
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 8, 31, 0, 0, 0, 0, time.UTC),
			},
			ReContract: secondContract,
		})
	}
	type testCase struct {
		date                 time.Time
		expectedContractCode string
		expectedError        error
		description          string
	}
	testCases := []testCase{
		{
			date:          time.Date(2020, 8, 31, 0, 0, 0, 0, time.UTC),
			expectedError: apperror.NotFoundError{},
			description:   "before first contract year",
		},
		{
			date:                 time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
			expectedContractCode: "first",
			expectedError:        nil,
			description:          "first day of first contract year",
		},
		{
			date:                 time.Date(2021, 8, 31, 23, 0, 0, 0, time.UTC),
			expectedContractCode: "first",
			expectedError:        nil,
			description:          "last day of first contract year, after time 0",
		},
		{
			date:                 time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
			expectedContractCode: "first",
			expectedError:        nil,
			description:          "second year of first contract",
		},
		{
			date:          time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC),
			expectedError: apperror.NotFoundError{},
			description:   "after all contract years",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			setupTestData()
			contract, err := FetchContract(suite.AppContextForTest(), tc.date)
			if tc.expectedError != nil {
				suite.Error(err)
				suite.IsType(tc.expectedError, err)
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedContractCode, contract.Code)
			}
		})
	}
}

func (suite *ServiceParamValueLookupsSuite) TestFetchContractForMove() {
	suite.Run("should return error for nonexistent move", func() {
		moveID := uuid.Must(uuid.NewV4())
		_, err := fetchContractForMove(suite.AppContextForTest(), moveID)
		suite.IsType(apperror.NotFoundError{}, err)
	})
	suite.Run("should return error for move that is not available to prime", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)
		_, err := fetchContractForMove(suite.AppContextForTest(), move.ID)
		suite.IsType(apperror.ConflictError{}, err)
	})
	suite.Run("should find contract for move that is available to prime", func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		contract, err := fetchContractForMove(suite.AppContextForTest(), move.ID)
		suite.NoError(err)
		suite.Equal(testdatagen.DefaultContractCode, contract.Code)
	})
}
