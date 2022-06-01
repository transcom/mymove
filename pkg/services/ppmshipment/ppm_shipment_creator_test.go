package ppmshipment

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestPPMShipmentCreator() {

	// One-time test setup
	ppmEstimator := mocks.PPMEstimator{}
	ppmShipmentCreator := NewPPMShipmentCreator(&ppmEstimator)

	type createShipmentSubtestData struct {
		move           models.Move
		newPPMShipment *models.PPMShipment
	}

	// createSubtestData - Sets up objects/data that need to be set up on a per-test basis.
	createSubtestData := func(appCtx appcontext.AppContext, assertions testdatagen.Assertions) (subtestData *createShipmentSubtestData) {
		subtestData = &createShipmentSubtestData{}

		subtestData.move = testdatagen.MakeMove(appCtx.DB(), assertions)

		fullAssertions := testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				MoveTaskOrderID: subtestData.move.ID,
				ShipmentType:    models.MTOShipmentTypePPM,
				Status:          models.MTOShipmentStatusDraft,
			},
		}

		testdatagen.MergeModels(&fullAssertions, assertions)

		mtoShipment := testdatagen.MakeBaseMTOShipment(appCtx.DB(), fullAssertions)

		// Initialize a valid PPMShipment properly with the MTOShipment
		subtestData.newPPMShipment = &models.PPMShipment{
			ShipmentID: mtoShipment.ID,
			Shipment:   mtoShipment,
		}

		testdatagen.MergeModels(subtestData.newPPMShipment, assertions.PPMShipment)

		return subtestData
	}

	suite.Run("Can successfully create a PPMShipment", func() {
		// Under test:	CreatePPMShipment
		// Set up:		Established valid shipment and valid new PPM shipment
		// Expected:	New PPM shipment successfully created
		appCtx := suite.AppContextForTest()

		// Set required fields for PPMShipment
		subtestData := createSubtestData(appCtx, testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				ExpectedDepartureDate: testdatagen.NextValidMoveDate,
				PickupPostalCode:      "90909",
				DestinationPostalCode: "90905",
				SITExpected:           models.BoolPointer(false),
			},
		})

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
	})

	var invalidInputTests = []struct {
		name             string
		assertions       testdatagen.Assertions
		expectedErrorMsg string
	}{
		{
			"MTOShipment type is not PPM",
			testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
			"MTO shipment type must be PPM shipment",
		},
		{
			"MTOShipment is not a draft or submitted shipment",
			testdatagen.Assertions{
				MTOShipment: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			"Must have a DRAFT or SUBMITTED status associated with MTO shipment",
		},
		{
			"missing MTOShipment ID",
			testdatagen.Assertions{PPMShipment: models.PPMShipment{
				ShipmentID: uuid.Nil,
			}},
			"Invalid input found while validating the PPM shipment.",
		},
		{
			"already has a PPMShipment ID",
			testdatagen.Assertions{PPMShipment: models.PPMShipment{
				ID: uuid.Must(uuid.NewV4()),
			}},
			"Invalid input found while validating the PPM shipment.",
		},
		{
			"missing a required field",
			testdatagen.Assertions{}, // Passing in blank assertions, leaving out required fields.
			"Invalid input found while validating the PPM shipment.",
		},
	}

	for _, tt := range invalidInputTests {
		tt := tt

		suite.Run(fmt.Sprintf("Returns an InvalidInputError if %s", tt.name), func() {
			appCtx := suite.AppContextForTest()

			subtestData := createSubtestData(appCtx, tt.assertions)

			createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

			suite.Error(err)
			suite.Nil(createdPPMShipment)

			suite.IsType(apperror.InvalidInputError{}, err)

			suite.Equal(tt.expectedErrorMsg, err.Error())
		})
	}

	suite.Run("Can successfully create a PPMShipment as SC", func() {
		appCtx := suite.AppContextForTest()

		// Set required fields for PPMShipment
		expectedDepartureDate := testdatagen.NextValidMoveDate
		pickupPostalCode := "29212"
		destinationPostalCode := "78234"
		sitExpected := false
		estimatedWeight := unit.Pound(2450)
		hasProGear := false
		estimatedIncentive := unit.Cents(123456)
		subtestData := createSubtestData(appCtx, testdatagen.Assertions{
			PPMShipment: models.PPMShipment{
				Status:                models.PPMShipmentStatusSubmitted,
				ExpectedDepartureDate: expectedDepartureDate,
				PickupPostalCode:      pickupPostalCode,
				DestinationPostalCode: destinationPostalCode,
				SITExpected:           &sitExpected,
				EstimatedWeight:       &estimatedWeight,
				HasProGear:            &hasProGear,
			},
		})

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(&estimatedIncentive, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		if suite.NotNil(createdPPMShipment) {
			suite.NotZero(createdPPMShipment.ID)
			suite.NotEqual(uuid.Nil.String(), createdPPMShipment.ID.String())
			suite.Equal(expectedDepartureDate, createdPPMShipment.ExpectedDepartureDate)
			suite.Equal(pickupPostalCode, createdPPMShipment.PickupPostalCode)
			suite.Equal(destinationPostalCode, createdPPMShipment.DestinationPostalCode)
			suite.Equal(&sitExpected, createdPPMShipment.SITExpected)
			suite.Equal(&estimatedWeight, createdPPMShipment.EstimatedWeight)
			suite.Equal(&hasProGear, createdPPMShipment.HasProGear)
			suite.Equal(models.PPMShipmentStatusSubmitted, createdPPMShipment.Status)
			suite.Equal(&estimatedIncentive, createdPPMShipment.EstimatedIncentive)
			suite.NotZero(createdPPMShipment.CreatedAt)
			suite.NotZero(createdPPMShipment.UpdatedAt)
		}
	})
}
