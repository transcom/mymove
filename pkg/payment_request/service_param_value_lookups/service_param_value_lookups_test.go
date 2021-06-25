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
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

const defaultZip3Distance = 1234
const defaultZip5Distance = 48

type ServiceParamValueLookupsSuite struct {
	testingsuite.PopTestSuite
	logger  Logger
	planner route.Planner
}

func TestServiceParamValueLookupsSuite(t *testing.T) {
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.Anything,
		mock.Anything,
	).Return(defaultZip5Distance, nil)
	planner.On("Zip3TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(defaultZip3Distance, nil)
	planner.On("Zip5TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(defaultZip5Distance, nil)

	ts := &ServiceParamValueLookupsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
		planner:      planner,
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ServiceParamValueLookupsSuite) setupTestMTOServiceItemWithWeight(estimatedWeight unit.Pound, actualWeight unit.Pound, code models.ReServiceCode, shipmentType models.MTOShipmentType) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: code,
				Name: string(code),
			},
			MTOShipment: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         shipmentType,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
	suite.FatalNoError(err)

	return mtoServiceItem, paymentRequest, paramLookup
}

func (suite *ServiceParamValueLookupsSuite) TestServiceParamValueLookup() {
	suite.T().Run("contract passed in", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

		suite.FatalNoError(err)
		suite.Equal(ghcrateengine.DefaultContractCode, paramLookup.ContractCode)
	})

	suite.T().Run("MTOServiceItem passed in", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

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
		suite.T().Run(fmt.Sprintf("MTOShipment not looked up for %s", code), func(t *testing.T) {
			mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: code,
					Name: string(code),
				},
			})

			mtoServiceItem.MTOShipmentID = nil
			mtoServiceItem.MTOShipment = models.MTOShipment{}
			suite.MustSave(&mtoServiceItem)

			paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
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

	suite.T().Run("MTOShipment is looked up for other serivce items", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDLH,
				Name: "DLH",
			},
		})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		suite.NotNil(paramLookup.MTOServiceItem)
		if rpdl, ok := paramLookup.lookups[models.ServiceItemParamNameRequestedPickupDate].(RequestedPickupDateLookup); ok {
			suite.Equal(*mtoServiceItem.MTOShipmentID, rpdl.MTOShipment.ID)
		} else {
			suite.Fail("lookup not RequestedPickupDateLookup type")
		}
	})

	suite.T().Run("DestinationAddress is looked up for other serivce items", func(t *testing.T) {
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: "DLH",
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDUPK,
					Name: "DUPK",
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipDestAddress].(ZipAddressLookup); ok {
				suite.Equal(mtoServiceItem.MTOShipment.DestinationAddress.PostalCode, zdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
		}
	})

	suite.T().Run("DestinationAddress is not required for service items like domestic pack", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDPK,
				Name: "DPK",
			},
		})

		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.DestinationAddressID = nil
		suite.DB().Save(&mtoShipment)

		_, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)
	})

	suite.T().Run("PickupAddress is looked up for other serivce items", func(t *testing.T) {
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: "DLH",
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDPK,
					Name: "DPK",
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zpal, ok := paramLookup.lookups[models.ServiceItemParamNameZipPickupAddress].(ZipAddressLookup); ok {
				suite.Equal(mtoServiceItem.MTOShipment.PickupAddress.PostalCode, zpal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipAddressLookup type")
			}
		}
	})

	suite.T().Run("PickupAddress is not required for service items like domestic unpack", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDUPK,
				Name: "DUPK",
			},
		})

		mtoShipment := mtoServiceItem.MTOShipment
		mtoShipment.PickupAddressID = nil
		suite.DB().Save(&mtoShipment)

		_, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)
	})

	suite.T().Run("SITDestinationAddress is looked up for destination sit", func(t *testing.T) {
		sitFinalDestAddress := testdatagen.MakeAddress3(suite.DB(), testdatagen.Assertions{})
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDDASIT,
					Name: "DDASIT",
				},
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &sitFinalDestAddress.ID,
					SITDestinationFinalAddress:   &sitFinalDestAddress,
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
					Name: "DDDSIT",
				},
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &sitFinalDestAddress.ID,
					SITDestinationFinalAddress:   &sitFinalDestAddress,
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
					Name: "DDFSIT",
				},
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &sitFinalDestAddress.ID,
					SITDestinationFinalAddress:   &sitFinalDestAddress,
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipSITDestHHGFinalAddress].(ZipAddressLookup); ok {
				suite.Equal(mtoServiceItem.SITDestinationFinalAddress.PostalCode, zdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipSitAddress destination type")
			}
		}
	})

	suite.T().Run("SITDestinationAddress is not loaded non sit", func(t *testing.T) {
		testData := []models.MTOServiceItem{
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDLH,
					Name: "DLH",
				},
			}),
			testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDSH,
					Name: "DSH",
				},
			}),
		}

		for _, mtoServiceItem := range testData {
			paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
			suite.FatalNoError(err)

			suite.NotNil(paramLookup.MTOServiceItem)
			if zdal, ok := paramLookup.lookups[models.ServiceItemParamNameZipSITDestHHGFinalAddress].(ZipAddressLookup); ok {
				suite.Equal("", zdal.Address.PostalCode)
			} else {
				suite.Fail("lookup not ZipSitAddress destination type")
			}
		}
	})

	suite.T().Run("nil MTOServiceItemID", func(t *testing.T) {
		badMTOServiceItemID := uuid.Must(uuid.NewV4())
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, badMTOServiceItemID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), fmt.Sprintf("id: %s not found looking for MTOServiceItemID", badMTOServiceItemID))
		var expected *ServiceItemParamKeyData = nil
		suite.Equal(expected, paramLookup)
	})
}
