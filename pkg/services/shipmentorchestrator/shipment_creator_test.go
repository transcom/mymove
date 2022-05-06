package shipmentorchestrator

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite ShipmentSuite) TestCreateShipment() {

	// Setup in this area should only be for objects that can be created once for all the sub-tests. Any model data,
	// mocks, or objects that can be modified in subtests should instead be set up in makeSubtestData.

	createMTOShipmentMethodName := "CreateMTOShipment"
	createPPMShipmentMethodName := "CreatePPMShipmentWithDefaultCheck"

	type subtestDataObjects struct {
		mockMTOShipmentCreator      *mocks.MTOShipmentCreator
		mockPPMShipmentCreator      *mocks.PPMShipmentCreator
		shipmentCreatorOrchestrator services.ShipmentCreator
		fakeError                   error
	}

	makeSubtestData := func(returnErrorForMTOShipment bool) (subtestData subtestDataObjects) {
		mockMTOShipmentCreator := mocks.MTOShipmentCreator{}
		subtestData.mockMTOShipmentCreator = &mockMTOShipmentCreator

		mockPPMShipmentCreator := mocks.PPMShipmentCreator{}
		subtestData.mockPPMShipmentCreator = &mockPPMShipmentCreator

		subtestData.shipmentCreatorOrchestrator = NewShipmentCreator(subtestData.mockMTOShipmentCreator, subtestData.mockPPMShipmentCreator)

		if returnErrorForMTOShipment {
			subtestData.fakeError = apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Pickup date missing")

			subtestData.mockMTOShipmentCreator.
				On(
					createMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment"),
					mock.AnythingOfType("models.MTOServiceItems")).
				Return(nil, subtestData.fakeError)
		} else {
			subtestData.mockMTOShipmentCreator.
				On(
					createMTOShipmentMethodName,
					mock.AnythingOfType("*appcontext.appContext"),
					mock.AnythingOfType("*models.MTOShipment"),
					mock.AnythingOfType("models.MTOServiceItems")).
				Return(
					func(_ appcontext.AppContext, ship *models.MTOShipment, _ models.MTOServiceItems) *models.MTOShipment {
						ship.ID = uuid.Must(uuid.NewV4())

						return ship
					},
					func(_ appcontext.AppContext, ship *models.MTOShipment, _ models.MTOServiceItems) error {
						return nil
					},
				)
		}

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
				func(_ appcontext.AppContext, ship *models.PPMShipment) error {
					return nil
				},
			)

		return
	}

	suite.Run("Returns an InvalidInputError if there is an error with the shipment info that was input", func() {
		subtestData := makeSubtestData(false)

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
		"HHG is set to Submitted": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
			models.MTOShipmentStatusSubmitted,
		},
		"HHG_LONGHAUL_DOMESTIC is set to Submitted": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGLongHaulDom,
			},
			models.MTOShipmentStatusSubmitted,
		},
		"HHG_SHORTHAUL_DOMESTIC is set to Submitted": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGShortHaulDom,
			},
			models.MTOShipmentStatusSubmitted,
		},
		"NTS is set to Submitted": {
			models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
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
			subtestData := makeSubtestData(false)

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
			ShipmentType: models.MTOShipmentTypeHHGLongHaulDom,
		},
		{
			ShipmentType: models.MTOShipmentTypeHHGShortHaulDom,
		},
		{
			ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
		},
		{
			ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
		},
	}

	for _, shipment := range shipmentCreationTestCases {
		shipment := shipment

		suite.Run(fmt.Sprintf("Calls necessary service objects for %s shipments", shipment.ShipmentType), func() {
			appCtx := suite.AppContextForTest()

			subtestData := makeSubtestData(false)

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
					mock.AnythingOfType("models.MTOServiceItems"),
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
		subtestData := makeSubtestData(false)

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

	suite.Run("Returns transaction error if one is raised", func() {
		subtestData := makeSubtestData(true)

		mtoShipment, err := subtestData.shipmentCreatorOrchestrator.CreateShipment(suite.AppContextForTest(), &models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHG,
		})

		suite.Nil(mtoShipment)

		suite.Error(err)
		suite.Equal(subtestData.fakeError, err)
	})

	suite.Run("Returns error early if MTOShipment can't be created", func() {
		subtestData := makeSubtestData(true)

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
			mock.AnythingOfType("models.MTOServiceItems"),
		)

		subtestData.mockPPMShipmentCreator.AssertNotCalled(
			suite.T(),
			createPPMShipmentMethodName,
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PPMShipment"),
		)
	})
}
