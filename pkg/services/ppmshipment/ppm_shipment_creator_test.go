package ppmshipment

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestPPMShipmentCreator() {

	// One-time test setup
	ppmEstimator := mocks.PPMEstimator{}
	addressCreator := address.NewAddressCreator()
	ppmShipmentCreator := NewPPMShipmentCreator(&ppmEstimator, addressCreator)

	type createShipmentSubtestData struct {
		move           models.Move
		newPPMShipment *models.PPMShipment
	}

	// createSubtestData - Sets up objects/data that need to be set up on a per-test basis.
	createSubtestData := func(ppmShipmentTemplate models.PPMShipment, mtoShipmentTemplate *models.MTOShipment) (subtestData *createShipmentSubtestData) {
		subtestData = &createShipmentSubtestData{}

		// TODO: pass customs through once we refactor this function to take in []factory.Customization instead of assertions
		subtestData.move = factory.BuildMove(suite.DB(), nil, nil)

		customMTOShipment := models.MTOShipment{
			MoveTaskOrderID: subtestData.move.ID,
			ShipmentType:    models.MTOShipmentTypePPM,
			Status:          models.MTOShipmentStatusDraft,
		}

		if mtoShipmentTemplate != nil {
			testdatagen.MergeModels(&customMTOShipment, *mtoShipmentTemplate)
		}

		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: customMTOShipment,
			},
		}, nil)

		// Initialize a valid PPMShipment properly with the MTOShipment
		subtestData.newPPMShipment = &models.PPMShipment{
			ShipmentID: mtoShipment.ID,
			Shipment:   mtoShipment,
		}

		testdatagen.MergeModels(subtestData.newPPMShipment, ppmShipmentTemplate)

		return subtestData
	}

	suite.Run("Can successfully create a PPMShipment", func() {
		// Under test:	CreatePPMShipment
		// Set up:		Established valid shipment and valid new PPM shipment
		// Expected:	New PPM shipment successfully created
		appCtx := suite.AppContextForTest()

		// Set required fields for PPMShipment
		subtestData := createSubtestData(models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			PickupPostalCode:      "90909",
			DestinationPostalCode: "90905",
			SITExpected:           models.BoolPointer(false),
		}, nil)

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
	})

	var invalidInputTests = []struct {
		name                string
		mtoShipmentTemplate *models.MTOShipment
		ppmShipmentTemplate models.PPMShipment
		expectedErrorMsg    string
	}{
		{
			"MTOShipment type is not PPM",
			&models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
			models.PPMShipment{},
			"MTO shipment type must be PPM shipment",
		},
		{
			"MTOShipment is not a draft or submitted shipment",
			&models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
			models.PPMShipment{},
			"Must have a DRAFT or SUBMITTED status associated with MTO shipment",
		},
		{
			"missing MTOShipment ID",
			nil,
			models.PPMShipment{
				ShipmentID: uuid.Nil,
			},
			"Invalid input found while validating the PPM shipment.",
		},
		{
			"already has a PPMShipment ID",
			nil,
			models.PPMShipment{
				ID: uuid.Must(uuid.NewV4()),
			},
			"Invalid input found while validating the PPM shipment.",
		},
		{
			"missing a required field",
			// Passing in blank assertions, leaving out required
			// fields.
			nil,
			models.PPMShipment{},
			"Invalid input found while validating the PPM shipment.",
		},
	}

	for _, tt := range invalidInputTests {
		tt := tt

		suite.Run(fmt.Sprintf("Returns an InvalidInputError if %s", tt.name), func() {
			appCtx := suite.AppContextForTest()

			subtestData := createSubtestData(tt.ppmShipmentTemplate, tt.mtoShipmentTemplate)

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

		pickupAddress := models.Address{
			StreetAddress1: "123 Any Pickup Street",
			City:           "SomeCity",
			State:          "CA",
			PostalCode:     "90210",
		}

		secondaryPickupAddress := models.Address{
			StreetAddress1: "123 Any Secondary Pickup Street",
			City:           "SomeCity",
			State:          "CA",
			PostalCode:     "90210",
		}

		destinationAddress := models.Address{
			StreetAddress1: "123 Any Destination Street",
			City:           "SomeCity",
			State:          "CA",
			PostalCode:     "90210",
		}

		secondaryDestinationAddress := models.Address{
			StreetAddress1: "123 Any Secondary Destination Street",
			City:           "SomeCity",
			State:          "CA",
			PostalCode:     "90210",
		}

		subtestData := createSubtestData(models.PPMShipment{
			Status:                      models.PPMShipmentStatusSubmitted,
			ExpectedDepartureDate:       expectedDepartureDate,
			PickupPostalCode:            pickupPostalCode,
			DestinationPostalCode:       destinationPostalCode,
			SITExpected:                 &sitExpected,
			EstimatedWeight:             &estimatedWeight,
			HasProGear:                  &hasProGear,
			PickupAddress:               &pickupAddress,
			DestinationAddress:          &destinationAddress,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			SecondaryDestinationAddress: &secondaryDestinationAddress,
		}, nil)

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(&estimatedIncentive, nil, nil).Once()

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
			suite.Equal(pickupAddress.StreetAddress1, createdPPMShipment.PickupAddress.StreetAddress1)
			suite.Equal(secondaryPickupAddress.StreetAddress1, createdPPMShipment.SecondaryPickupAddress.StreetAddress1)
			suite.Equal(destinationAddress.StreetAddress1, createdPPMShipment.DestinationAddress.StreetAddress1)
			suite.Equal(secondaryDestinationAddress.StreetAddress1, createdPPMShipment.SecondaryDestinationAddress.StreetAddress1)
			suite.NotNil(createdPPMShipment.PickupAddressID)
			suite.NotNil(createdPPMShipment.DestinationAddressID)
			suite.NotNil(createdPPMShipment.SecondaryPickupAddressID)
			suite.NotNil(createdPPMShipment.SecondaryDestinationAddressID)
			//ensure HasSecondaryPickupAddress/HasSecondaryDestinationAddress are set even if not initially provided
			suite.True(createdPPMShipment.HasSecondaryPickupAddress != nil)
			suite.Equal(models.BoolPointer(true), createdPPMShipment.HasSecondaryPickupAddress)
			suite.True(createdPPMShipment.HasSecondaryDestinationAddress != nil)
			suite.Equal(models.BoolPointer(true), createdPPMShipment.HasSecondaryDestinationAddress)
		}
	})
}
