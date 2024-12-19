package shipment

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

func (suite *ShipmentSuite) TestCreateShipment() {

	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in makeSubtestData.

	createMTOShipmentMethodName := "CreateMTOShipment"
	createPPMShipmentMethodName := "CreatePPMShipmentWithDefaultCheck"
	updatePPMTypeMethodName := "UpdatePPMType"

	type subtestDataObjects struct {
		mockMTOShipmentCreator        *mocks.MTOShipmentCreator
		mockPPMShipmentCreator        *mocks.PPMShipmentCreator
		mockBoatShipmentCreator       *mocks.BoatShipmentCreator
		mockMobileHomeShipmentCreator *mocks.MobileHomeShipmentCreator
		mockMoveTaskOrderUpdater      *mocks.MoveTaskOrderUpdater
		shipmentCreatorOrchestrator   services.ShipmentCreator
		fakeError                     error
	}

	makeSubtestData := func(returnErrorForMTOShipment bool, returnErrorForPPMShipment bool) (subtestData subtestDataObjects) {
		mockMTOShipmentCreator := mocks.MTOShipmentCreator{}
		subtestData.mockMTOShipmentCreator = &mockMTOShipmentCreator

		mockPPMShipmentCreator := mocks.PPMShipmentCreator{}
		subtestData.mockPPMShipmentCreator = &mockPPMShipmentCreator

		mockBoatShipmentCreator := mocks.BoatShipmentCreator{}
		subtestData.mockBoatShipmentCreator = &mockBoatShipmentCreator

		mockMobileHomeShipmentCreator := mocks.MobileHomeShipmentCreator{}
		subtestData.mockMobileHomeShipmentCreator = &mockMobileHomeShipmentCreator

		mockMoveTaskOrderUpdater := mocks.MoveTaskOrderUpdater{}
		subtestData.mockMoveTaskOrderUpdater = &mockMoveTaskOrderUpdater

		subtestData.shipmentCreatorOrchestrator = NewShipmentCreator(subtestData.mockMTOShipmentCreator, subtestData.mockPPMShipmentCreator, subtestData.mockBoatShipmentCreator, subtestData.mockMobileHomeShipmentCreator, mtoshipment.NewShipmentRouter(), subtestData.mockMoveTaskOrderUpdater)

		if returnErrorForMTOShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Pickup date missing")

			subtestData.mockMTOShipmentCreator.
				On(
					createMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment")).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockMTOShipmentCreator.
				On(
					createMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment")).
				Return(
					func(_ appcontext.AppContext, ship *models.MTOShipment) *models.MTOShipment {
						ship.ID = uuid.Must(uuid.NewV4())

						return ship
					},
					func(_ appcontext.AppContext, _ *models.MTOShipment) error {
						return nil
					},
				)
		}

		if returnErrorForPPMShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Invalid input found while validating the PPM shipment.")

			subtestData.mockPPMShipmentCreator.
				On(
					createPPMShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.PPMShipment"),
				).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockPPMShipmentCreator.
				On(
					createPPMShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.PPMShipment"),
				).
				Return(
					func(_ appcontext.AppContext, ship *models.PPMShipment) *models.PPMShipment {
						ship.ID = uuid.Must(uuid.NewV4())

						return ship
					},
					func(_ appcontext.AppContext, _ *models.PPMShipment) error {
						return nil
					},
				)
		}
		subtestData.mockMoveTaskOrderUpdater.
			On(
				updatePPMTypeMethodName,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("uuid.UUID"),
			).
			Return(nil, nil)

		return subtestData
	}

	suite.Run("Returns an InvalidInputError if there is an error with the shipment info that was input", func() {
		subtestData := makeSubtestData(false, false)

		mtoShipment, err := subtestData.shipmentCreatorOrchestrator.CreateShipment(suite.AppContextForTest(), &models.MTOShipment{})

		suite.Nil(mtoShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "Invalid input found while validating the shipment.")
	})

	statusTestCases := map[string]struct {
		shipment       models.MTOShipment
		expectedStatus models.MTOShipmentStatus
	}{
		"PPM is set to Draft": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				PPMShipment:  &models.PPMShipment{},
			},
			models.MTOShipmentStatusDraft,
		},
		"PPM can be set to another status": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusSubmitted,
				PPMShipment:  &models.PPMShipment{},
			},
			models.MTOShipmentStatusSubmitted,
		},
		"HHG is set to Submitted": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
			models.MTOShipmentStatusSubmitted,
		},
		"NTS is set to Submitted": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
			},
			models.MTOShipmentStatusSubmitted,
		},
		"NTS-Release is set to Submitted": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
			models.MTOShipmentStatusSubmitted,
		},
	}

	for name, tc := range statusTestCases {
		name := name
		tc := tc

		suite.Run(fmt.Sprintf("Sets status as expected: %s", name), func() {
			subtestData := makeSubtestData(false, false)

			mtoShipment, err := subtestData.shipmentCreatorOrchestrator.CreateShipment(suite.AppContextForTest(), &tc.shipment)

			suite.Nil(err)

			suite.Equal(tc.expectedStatus, mtoShipment.Status)
		})
	}

	shipmentCreationTestCases := []models.MTOShipment{
		{
			ShipmentType: models.MTOShipmentTypePPM,
			PPMShipment:  &models.PPMShipment{},
		},
		{
			ShipmentType: models.MTOShipmentTypeHHG,
		},
		{
			ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
		},
		{
			ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
		},
	}

	for _, shipment := range shipmentCreationTestCases {
		shipment := shipment

		suite.Run(fmt.Sprintf("Calls necessary service objects for %s shipments", shipment.ShipmentType), func() {
			appCtx := suite.AppContextForTest()

			subtestData := makeSubtestData(false, false)

			// Need to start a transaction so we can assert the call with the correct appCtx
			err := appCtx.NewTransaction(func(txAppCtx appcontext.AppContext) error {
				mtoShipment, err := subtestData.shipmentCreatorOrchestrator.CreateShipment(txAppCtx, &shipment)

				suite.NoError(err)
				suite.NotNil(mtoShipment)

				subtestData.mockMTOShipmentCreator.AssertCalled(
					suite.T(),
					createMTOShipmentMethodName,
					txAppCtx,
					&shipment,
				)

				if shipment.ShipmentType == models.MTOShipmentTypePPM {
					subtestData.mockPPMShipmentCreator.AssertCalled(
						suite.T(),
						createPPMShipmentMethodName,
						txAppCtx,
						shipment.PPMShipment,
					)
				} else {
					subtestData.mockPPMShipmentCreator.AssertNotCalled(
						suite.T(),
						createPPMShipmentMethodName,
						mock.AnythingOfType("*appcontext.appContext"),
						mock.AnythingOfType("*models.PPMShipment"),
					)
				}

				return nil
			})

			suite.NoError(err) // just making golangci-lint happy
		})
	}

	suite.Run("Sets MTOShipment info on PPMShipment", func() {
		subtestData := makeSubtestData(false, false)

		shipment := &models.MTOShipment{
			ShipmentType: models.MTOShipmentTypePPM,
			PPMShipment:  &models.PPMShipment{},
		}

		mtoShipment, err := subtestData.shipmentCreatorOrchestrator.CreateShipment(suite.AppContextForTest(), shipment)

		suite.NoError(err)

		suite.NotNil(mtoShipment)

		suite.NotNil(mtoShipment.ID) // this one is mainly a sanity check to ensure we aren't comparing nil to nil in the next one
		suite.Equal(mtoShipment.ID, mtoShipment.PPMShipment.ShipmentID)
		suite.Equal(*mtoShipment, mtoShipment.PPMShipment.Shipment)
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
			subtestData := makeSubtestData(tc.returnErrorForMTOShipment, tc.returnErrorForPPMShipment)

			shipment := models.MTOShipment{
				ShipmentType: tc.shipmentType,
			}

			if tc.shipmentType == models.MTOShipmentTypePPM {
				shipment.PPMShipment = &models.PPMShipment{}
			}

			mtoShipment, err := subtestData.shipmentCreatorOrchestrator.CreateShipment(suite.AppContextForTest(), &shipment)

			suite.Nil(mtoShipment)

			suite.Error(err)
			suite.Equal(subtestData.fakeError, err)
		})
	}

	suite.Run("Returns error early if MTOShipment can't be created", func() {
		subtestData := makeSubtestData(true, false)

		mtoShipment, err := subtestData.shipmentCreatorOrchestrator.CreateShipment(suite.AppContextForTest(), &models.MTOShipment{
			ShipmentType: models.MTOShipmentTypePPM,
		})

		suite.Nil(mtoShipment)

		suite.Error(err)
		suite.Equal(subtestData.fakeError, err)

		subtestData.mockMTOShipmentCreator.AssertCalled(
			suite.T(),
			createMTOShipmentMethodName,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		)

		subtestData.mockPPMShipmentCreator.AssertNotCalled(
			suite.T(),
			createPPMShipmentMethodName,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
		)
	})
}
