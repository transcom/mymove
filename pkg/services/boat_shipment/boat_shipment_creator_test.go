package boatshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *BoatShipmentSuite) TestBoatShipmentCreator() {

	boatShipmentCreator := NewBoatShipmentCreator()

	type createShipmentSubtestData struct {
		move            models.Move
		newBoatShipment *models.BoatShipment
	}

	// createSubtestData - Sets up objects/data that need to be set up on a per-test basis.
	createSubtestData := func(boatShipmentTemplate models.BoatShipment, mtoShipmentTemplate *models.MTOShipment) (subtestData *createShipmentSubtestData) {
		subtestData = &createShipmentSubtestData{}

		// TODO: pass customs through once we refactor this function to take in []factory.Customization instead of assertions
		subtestData.move = factory.BuildMove(suite.DB(), nil, nil)

		customMTOShipment := models.MTOShipment{
			MoveTaskOrderID: subtestData.move.ID,
			ShipmentType:    models.MTOShipmentTypeBoatHaulAway,
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

		// Initialize a valid BoatShipment properly with the MTOShipment
		subtestData.newBoatShipment = &models.BoatShipment{
			ShipmentID: mtoShipment.ID,
			Shipment:   mtoShipment,
		}

		testdatagen.MergeModels(subtestData.newBoatShipment, boatShipmentTemplate)

		return subtestData
	}

	suite.Run("Can successfully create a BoatShipment", func() {
		// Under test:	CreateBoatShipment
		// Set up:		Established valid shipment and valid new Boat shipment
		// Expected:	New Boat shipment successfully created
		appCtx := suite.AppContextForTest()

		// Set required fields for BoatShipment
		subtestData := createSubtestData(models.BoatShipment{
			Type:           models.BoatShipmentTypeHaulAway,
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Boat Make"),
			Model:          models.StringPointer("Boat Model"),
			LengthInInches: models.IntPointer(300),
			WidthInInches:  models.IntPointer(108),
			HeightInInches: models.IntPointer(72),
			HasTrailer:     models.BoolPointer(true),
			IsRoadworthy:   models.BoolPointer(false)}, nil)

		createdBoatShipment, err := boatShipmentCreator.CreateBoatShipmentWithDefaultCheck(appCtx, subtestData.newBoatShipment)

		suite.Nil(err)
		suite.NotNil(createdBoatShipment)
	})

	var invalidInputTests = []struct {
		name                 string
		mtoShipmentTemplate  *models.MTOShipment
		boatShipmentTemplate models.BoatShipment
		expectedErrorMsg     string
	}{
		{
			"MTOShipment type is not Boat",
			&models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
			models.BoatShipment{
				Type: models.BoatShipmentTypeHaulAway,
			},
			"MTO shipment type must be Boat shipment",
		},
		{
			"missing MTOShipment ID",
			nil,
			models.BoatShipment{
				ShipmentID: uuid.Nil,
				Type:       models.BoatShipmentTypeHaulAway,
			},
			"Invalid input found while validating the Boat shipment.",
		},
		{
			"already has a BoatShipment ID",
			nil,
			models.BoatShipment{
				ID:   uuid.Must(uuid.NewV4()),
				Type: models.BoatShipmentTypeHaulAway,
			},
			"Invalid input found while validating the Boat shipment.",
		},
		{
			"missing a required field",
			// Passing in blank assertions, leaving out required
			// fields.
			nil,
			models.BoatShipment{},
			"Must have a HAUL_AWAY or TOW_AWAY type associated with Boat shipment",
		},
	}

	for _, tt := range invalidInputTests {
		tt := tt

		suite.Run(fmt.Sprintf("Returns an InvalidInputError if %s", tt.name), func() {
			appCtx := suite.AppContextForTest()

			subtestData := createSubtestData(tt.boatShipmentTemplate, tt.mtoShipmentTemplate)

			createdBoatShipment, err := boatShipmentCreator.CreateBoatShipmentWithDefaultCheck(appCtx, subtestData.newBoatShipment)

			suite.Error(err)
			suite.Nil(createdBoatShipment)

			suite.IsType(apperror.InvalidInputError{}, err)

			suite.Equal(tt.expectedErrorMsg, err.Error())
		})
	}

	suite.Run("Can successfully create a BoatShipment as SC", func() {
		appCtx := suite.AppContextForTest()

		subtestData := createSubtestData(models.BoatShipment{
			Type:           models.BoatShipmentTypeHaulAway,
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Boat Make"),
			Model:          models.StringPointer("Boat Model"),
			LengthInInches: models.IntPointer(300),
			WidthInInches:  models.IntPointer(108),
			HeightInInches: models.IntPointer(72),
			HasTrailer:     models.BoolPointer(true),
			IsRoadworthy:   models.BoolPointer(false),
		}, nil)

		createdBoatShipment, err := boatShipmentCreator.CreateBoatShipmentWithDefaultCheck(appCtx, subtestData.newBoatShipment)

		suite.Nil(err)
		if suite.NotNil(createdBoatShipment) {
			suite.NotZero(createdBoatShipment.ID)
			suite.NotEqual(uuid.Nil.String(), createdBoatShipment.ID.String())
		}
	})
}
