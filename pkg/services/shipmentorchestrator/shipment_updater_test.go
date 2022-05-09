package shipmentorchestrator

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite ShipmentSuite) TestUpdateShipment() {

	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in makeSubtestData.

	updateMTOShipmentMethodName := "UpdateMTOShipmentCustomer"
	updatePPMShipmentMethodName := "UpdatePPMShipmentWithDefaultCheck"

	type subtestDataObjects struct {
		mockMTOShipmentUpdater      *mocks.MTOShipmentUpdater
		mockPPMShipmentUpdater      *mocks.PPMShipmentUpdater
		shipmentUpdaterOrchestrator services.ShipmentUpdater

		fakeError error
	}

	makeSubtestData := func(returnErrorForMTOShipment bool, returnErrorForPPMShipment bool) (subtestData subtestDataObjects) {
		mockMTOShipmentUpdater := mocks.MTOShipmentUpdater{}
		subtestData.mockMTOShipmentUpdater = &mockMTOShipmentUpdater

		mockPPMShipmentUpdater := mocks.PPMShipmentUpdater{}
		subtestData.mockPPMShipmentUpdater = &mockPPMShipmentUpdater

		subtestData.shipmentUpdaterOrchestrator = NewShipmentUpdater(subtestData.mockMTOShipmentUpdater, subtestData.mockPPMShipmentUpdater)

		if returnErrorForMTOShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Pickup date missing")

			subtestData.mockMTOShipmentUpdater.
				On(
					updateMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment"),
					mock.AnythingOfType("string")).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockMTOShipmentUpdater.
				On(
					updateMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment"),
					mock.AnythingOfType("string")).
				Return(
					func(_ appcontext.AppContext, ship *models.MTOShipment, _ string) *models.MTOShipment {
						// Mimicking how the MTOShipment updater actually returns a new pointer so that we can test
						// a bit more realistically while still using mocks.
						updatedShip := *ship
						updatedShip.PPMShipment = nil // Currently returns an MTOShipment without PPMShipment info

						return &updatedShip
					},
					func(_ appcontext.AppContext, ship *models.MTOShipment, _ string) error {
						return nil
					},
				)
		}

		if returnErrorForPPMShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Invalid input found while validating the PPM shipment.")

			subtestData.mockPPMShipmentUpdater.
				On(
					updatePPMShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.PPMShipment"),
					mock.AnythingOfType("uuid.UUID"),
				).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockPPMShipmentUpdater.
				On(
					updatePPMShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.PPMShipment"),
					mock.AnythingOfType("uuid.UUID"),
				).
				Return(
					func(_ appcontext.AppContext, ship *models.PPMShipment, _ uuid.UUID) *models.PPMShipment {
						// Mimicking how the PPMShipment updater actually returns a new pointer so that we can test
						// a bit more realistically while still using mocks.
						updatedShip := *ship

						return &updatedShip
					},
					func(_ appcontext.AppContext, ship *models.PPMShipment, _ uuid.UUID) error {
						return nil
					},
				)
		}

		return subtestData
	}

	suite.Run("Returns an InvalidInputError if there is an error with the shipment info that was input", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(false, false)

		shipment := testdatagen.MakeDefaultMTOShipment(appCtx.DB())

		// Set invalid data, can't pass in blank to the generator above (it'll default to HHG if blank) so we're setting it afterward.
		shipment.ShipmentType = ""

		updatedShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt))

		suite.Nil(updatedShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the shipment.")
	})

	shipmentTypeTestCases := []models.MTOShipmentType{

		models.MTOShipmentTypeHHG,
		models.MTOShipmentTypeHHGLongHaulDom,
		models.MTOShipmentTypeHHGShortHaulDom,
		models.MTOShipmentTypeHHGIntoNTSDom,
		models.MTOShipmentTypeHHGOutOfNTSDom,
		models.MTOShipmentTypePPM,
	}

	for _, shipmentType := range shipmentTypeTestCases {
		shipmentType := shipmentType

		suite.Run(fmt.Sprintf("Calls necessary service objects for %s shipments", shipmentType), func() {
			appCtx := suite.AppContextForTest()

			subtestData := makeSubtestData(false, false)

			var shipment models.MTOShipment

			isPPMShipment := shipmentType == models.MTOShipmentTypePPM

			if isPPMShipment {
				ppmShipment := testdatagen.MakeDefaultPPMShipment(appCtx.DB())

				shipment = ppmShipment.Shipment
			} else {
				shipment = testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
					MTOShipment: models.MTOShipment{
						ShipmentType: shipmentType,
					},
				})
			}

			eTag := etag.GenerateEtag(shipment.UpdatedAt)

			// Need to start a transaction so we can assert the call with the correct appCtx
			err := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
				mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(txAppCtx, &shipment, eTag)

				suite.NoError(err)
				suite.NotNil(mtoShipment)

				subtestData.mockMTOShipmentUpdater.AssertCalled(
					suite.T(),
					updateMTOShipmentMethodName,
					txAppCtx,
					&shipment,
					eTag,
				)

				if isPPMShipment {
					subtestData.mockPPMShipmentUpdater.AssertCalled(
						suite.T(),
						updatePPMShipmentMethodName,
						txAppCtx,
						shipment.PPMShipment,
						shipment.ID,
					)
				} else {
					subtestData.mockPPMShipmentUpdater.AssertNotCalled(
						suite.T(),
						updatePPMShipmentMethodName,
						mock.AnythingOfType("*appcontext.appContext"),
						mock.AnythingOfType("*models.PPMShipment"),
						mock.AnythingOfType("uuid.UUID"),
					)
				}

				return nil
			})

			suite.NoError(err) // just making golangci-lint happy
		})
	}

	suite.Run("Sets MTOShipment info on PPMShipment and updated PPMShipment back on MTOShipment", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(false, false)

		ppmShipment := testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				HasProGear:          models.BoolPointer(true),
				ProGearWeight:       models.PoundPointer(unit.Pound(1900)),
				SpouseProGearWeight: models.PoundPointer(unit.Pound(300)),
			},
		})

		shipment := ppmShipment.Shipment

		// set new field to update
		shipment.PPMShipment.HasProGear = models.BoolPointer(false)

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt))

		suite.NoError(err)

		suite.NotNil(mtoShipment)

		// check that the PPMShipment has the expected MTOShipment fields
		suite.Equal(mtoShipment.ID, mtoShipment.PPMShipment.ShipmentID)
		suite.Equal(*mtoShipment, mtoShipment.PPMShipment.Shipment)

		// check we got the latest version back
		suite.NotEqual(&ppmShipment, mtoShipment.PPMShipment)
		suite.False(*mtoShipment.PPMShipment.HasProGear)
	})

	serviceObjectErrorTestCases := map[string]struct {
		shipmentType              models.MTOShipmentType
		returnErrorForMTOShipment bool
		returnErrorForPPMShipment bool
	}{
		"error updating MTOShipment": {
			shipmentType:              models.MTOShipmentTypeHHG,
			returnErrorForMTOShipment: true,
			returnErrorForPPMShipment: false,
		},
		"error updating PPMShipment": {
			shipmentType:              models.MTOShipmentTypePPM,
			returnErrorForMTOShipment: false,
			returnErrorForPPMShipment: true,
		},
	}

	for name, tc := range serviceObjectErrorTestCases {
		name := name
		tc := tc

		suite.Run(fmt.Sprintf("Returns transaction error if there is an %s", name), func() {
			appCtx := suite.AppContextForTest()

			subtestData := makeSubtestData(tc.returnErrorForMTOShipment, tc.returnErrorForPPMShipment)

			var shipment models.MTOShipment

			if tc.shipmentType == models.MTOShipmentTypePPM {
				ppmShipment := testdatagen.MakeDefaultPPMShipment(appCtx.DB())

				shipment = ppmShipment.Shipment
			} else {
				shipment = testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
					MTOShipment: models.MTOShipment{
						ShipmentType: tc.shipmentType,
					},
				})
			}

			mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, etag.GenerateEtag(shipment.UpdatedAt))

			suite.Nil(mtoShipment)

			suite.Error(err)
			suite.Equal(subtestData.fakeError, err)
		})
	}

	suite.Run("Returns error early if MTOShipment can't be updated", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeSubtestData(true, false)

		shipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
		})

		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		mtoShipment, err := subtestData.shipmentUpdaterOrchestrator.UpdateShipment(appCtx, &shipment, eTag)

		suite.Nil(mtoShipment)

		suite.Error(err)
		suite.Equal(subtestData.fakeError, err)

		subtestData.mockMTOShipmentUpdater.AssertCalled(
			suite.T(),
			updateMTOShipmentMethodName,
			mock.AnythingOfType("*appcontext.appContext"),
			&shipment,
			eTag,
		)

		subtestData.mockPPMShipmentUpdater.AssertNotCalled(
			suite.T(),
			updatePPMShipmentMethodName,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
			mock.AnythingOfType("uuid.UUID"),
		)
	})
}
