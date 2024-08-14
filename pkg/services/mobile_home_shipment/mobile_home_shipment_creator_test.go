package mobilehomeshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MobileHomeShipmentSuite) TestMobileHomeShipmentCreator() {
	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)

	mobileHomeShipmentCreator := NewMobileHomeShipmentCreator()

	type createShipmentSubtestData struct {
		move                  models.Move
		newMobileHomeShipment *models.MobileHome
	}

	// createSubtestData - Sets up objects/data that need to be set up on a per-test basis.
	createSubtestData := func(mobileHomeShipmentTemplate models.MobileHome, mtoShipmentTemplate *models.MTOShipment) (subtestData *createShipmentSubtestData) {
		subtestData = &createShipmentSubtestData{}

		// TODO: pass customs through once we refactor this function to take in []factory.Customization instead of assertions
		subtestData.move = factory.BuildMove(suite.DB(), nil, nil)

		customMTOShipment := models.MTOShipment{
			MoveTaskOrderID: subtestData.move.ID,
			ShipmentType:    models.MTOShipmentTypeMobileHome,
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

		// Initialize a valid MobileHomeShipment properly with the MTOShipment
		subtestData.newMobileHomeShipment = &models.MobileHome{
			ShipmentID: mtoShipment.ID,
			Shipment:   mtoShipment,
		}

		testdatagen.MergeModels(subtestData.newMobileHomeShipment, mobileHomeShipmentTemplate)

		return subtestData
	}

	suite.Run("Can successfully create a MobileHomeShipment", func() {
		// Under test:	CreateMobileHomeShipment
		// Set up:		Established valid shipment and valid new Mobile Home shipment
		// Expected:	New Mobile Home shipment successfully created
		appCtx := suite.AppContextForTest()

		// Set required fields for MobileHomeShipment
		subtestData := createSubtestData(models.MobileHome{
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Mobile Home Make"),
			Model:          models.StringPointer("Mobile Home Model"),
			LengthInInches: models.IntPointer(300),
			HeightInInches: models.IntPointer(72),
			WidthInInches:  models.IntPointer(108),
			CreatedAt:      date}, nil)

		createdMobileHomeShipment, err := mobileHomeShipmentCreator.CreateMobileHomeShipmentWithDefaultCheck(appCtx, subtestData.newMobileHomeShipment)

		suite.Nil(err)
		suite.NotNil(createdMobileHomeShipment)
	})

	var invalidInputTests = []struct {
		name                       string
		mtoShipmentTemplate        *models.MTOShipment
		mobileHomeShipmentTemplate models.MobileHome
		expectedErrorMsg           string
	}{
		{
			"missing MTOShipment ID",
			nil,
			models.MobileHome{
				ShipmentID: uuid.Nil,
			},
			"Invalid input found while validating the Mobile Home shipment.",
		},
		{
			"already has a MobileHomeShipment ID",
			nil,
			models.MobileHome{
				ID: uuid.Must(uuid.NewV4()),
			},
			"Invalid input found while validating the Mobile Home shipment.",
		},
		{
			"missing a required field",
			// Passing in blank assertions, leaving out required
			// fields.
			nil,
			models.MobileHome{},
			"Invalid input found while validating the Mobile Home shipment.",
		},
	}

	for _, tt := range invalidInputTests {
		tt := tt

		suite.Run(fmt.Sprintf("Returns an InvalidInputError if %s", tt.name), func() {
			appCtx := suite.AppContextForTest()

			subtestData := createSubtestData(tt.mobileHomeShipmentTemplate, tt.mtoShipmentTemplate)

			createdMobileHomeShipment, err := mobileHomeShipmentCreator.CreateMobileHomeShipmentWithDefaultCheck(appCtx, subtestData.newMobileHomeShipment)

			suite.Error(err)
			suite.Nil(createdMobileHomeShipment)

			suite.IsType(apperror.InvalidInputError{}, err)

			suite.Equal(tt.expectedErrorMsg, err.Error())
		})
	}

	suite.Run("Can successfully create a MobileHomeShipment as SC", func() {
		appCtx := suite.AppContextForTest()

		subtestData := createSubtestData(models.MobileHome{
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Mobile Home Make"),
			Model:          models.StringPointer("Mobile Home Model"),
			LengthInInches: models.IntPointer(300),
			WidthInInches:  models.IntPointer(108),
			HeightInInches: models.IntPointer(72),
			CreatedAt:      date,
		}, nil)

		createdMobileHomeShipment, err := mobileHomeShipmentCreator.CreateMobileHomeShipmentWithDefaultCheck(appCtx, subtestData.newMobileHomeShipment)

		suite.Nil(err)
		if suite.NotNil(createdMobileHomeShipment) {
			suite.NotZero(createdMobileHomeShipment.ID)
			suite.NotEqual(uuid.Nil.String(), createdMobileHomeShipment.ID.String())
		}
	})
}