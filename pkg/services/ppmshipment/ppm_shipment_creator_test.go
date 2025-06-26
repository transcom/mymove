package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
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

		// adding existing HHG shipment to move
		factory.BuildMTOShipment(nil, []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
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

	suite.Run("Can successfully create a domestic PPMShipment", func() {
		// Under test:	CreatePPMShipment
		// Set up:		Established valid shipment and valid new PPM shipment
		// Expected:	New PPM shipment successfully created, market code is "d" on the parent shipment
		// Need a logged in user
		lgu := uuid.Must(uuid.NewV4()).String()
		user := models.User{
			OktaID:    lgu,
			OktaEmail: "email@example.com",
		}
		suite.MustSave(&user)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          user.ID,
			IDToken:         "fake token",
		}

		appCtx := suite.AppContextWithSessionForTest(session)

		// Set required fields for PPMShipment
		subtestData := createSubtestData(models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			SITExpected:           models.BoolPointer(false),
			PickupAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50308",
				County:         models.StringPointer("POLK"),
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 12345"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30183",
				County:         models.StringPointer("COLUMBIA"),
			},
		}, nil)

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil, nil).Once()

		ppmEstimator.On(
			"MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
		suite.Equal(createdPPMShipment.PPMType, models.PPMTypeIncentiveBased)
		suite.Equal(createdPPMShipment.Shipment.MarketCode, models.MarketCodeDomestic)
	})

	suite.Run("Can successfully create a domestic PPMShipment with applicable GCC multiplier", func() {
		lgu := uuid.Must(uuid.NewV4()).String()
		user := models.User{
			OktaID:    lgu,
			OktaEmail: "email@example.com",
		}
		suite.MustSave(&user)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          user.ID,
			IDToken:         "fake token",
			ActiveRole:      roles.Role{},
		}

		appCtx := suite.AppContextWithSessionForTest(session)
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")

		// Set required fields for PPMShipment
		subtestData := createSubtestData(models.PPMShipment{
			ExpectedDepartureDate: validGccMultiplierDate,
			SITExpected:           models.BoolPointer(false),
			PickupAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50308",
				County:         models.StringPointer("POLK"),
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 12345"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "30183",
				County:         models.StringPointer("COLUMBIA"),
			},
		}, nil)

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil, nil).Once()

		ppmEstimator.On(
			"MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
		suite.Equal(createdPPMShipment.PPMType, models.PPMTypeIncentiveBased)
		suite.Equal(createdPPMShipment.Shipment.MarketCode, models.MarketCodeDomestic)
		suite.NotNil(createdPPMShipment.GCCMultiplier)
		suite.NotNil(createdPPMShipment.GCCMultiplierID)
	})

	suite.Run("Can successfully create an international PPMShipment", func() {
		// Under test:	CreatePPMShipment
		// Set up:		Established valid shipment and valid new PPM shipment
		// Expected:	New PPM shipment successfully created, market code is "i" on the parent shipment
		// Need a logged in user
		lgu := uuid.Must(uuid.NewV4()).String()
		user := models.User{
			OktaID:    lgu,
			OktaEmail: "email@example.com",
		}
		suite.MustSave(&user)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          user.ID,
			IDToken:         "fake token",
			ActiveRole:      roles.Role{},
		}

		appCtx := suite.AppContextWithSessionForTest(session)

		// Set required fields for PPMShipment
		subtestData := createSubtestData(models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			SITExpected:           models.BoolPointer(false),
			PickupAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "ANCHORAGE",
				State:          "AK",
				PostalCode:     "99507",
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 12345"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Honolulu",
				State:          "HI",
				PostalCode:     "96821",
			},
		}, nil)

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil, nil).Once()

		ppmEstimator.On(
			"MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
		suite.Equal(createdPPMShipment.Shipment.MarketCode, models.MarketCodeInternational)
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
			// Need a logged in user
			lgu := uuid.Must(uuid.NewV4()).String()
			user := models.User{
				OktaID:    lgu,
				OktaEmail: "email@example.com",
			}
			suite.MustSave(&user)

			session := &auth.Session{
				ApplicationName: auth.OfficeApp,
				UserID:          user.ID,
				IDToken:         "fake token",
				ActiveRole:      roles.Role{},
			}

			appCtx := suite.AppContextWithSessionForTest(session)

			subtestData := createSubtestData(tt.ppmShipmentTemplate, tt.mtoShipmentTemplate)

			createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

			suite.Error(err)
			suite.Nil(createdPPMShipment)

			suite.IsType(apperror.InvalidInputError{}, err)

			suite.Equal(tt.expectedErrorMsg, err.Error())
		})
	}

	suite.Run("Can successfully create and auto-approve a PPMShipment as SC", func() {
		// Need a logged in user
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		identity, err := models.FetchUserIdentity(suite.DB(), scOfficeUser.User.OktaID)
		suite.NoError(err)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *scOfficeUser.UserID,
			IDToken:         "fake token",
		}
		defaultRole, err := identity.Roles.Default()
		suite.FatalNoError(err)
		session.ActiveRole = *defaultRole

		appCtx := suite.AppContextWithSessionForTest(session)

		// Set required fields for PPMShipment
		expectedDepartureDate := testdatagen.NextValidMoveDate
		sitExpected := false
		estimatedWeight := unit.Pound(2450)
		hasProGear := false
		estimatedIncentive := unit.Cents(123456)
		maxIncentive := unit.Cents(123456)

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

		tertiaryPickupAddress := models.Address{
			StreetAddress1: "123 Any Tertiary Pickup Street",
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
		tertiaryDestinationAddress := models.Address{
			StreetAddress1: "123 Any Tertiary Destination Street",
			City:           "SomeCity",
			State:          "CA",
			PostalCode:     "90210",
		}

		subtestData := createSubtestData(models.PPMShipment{
			Status:                      models.PPMShipmentStatusSubmitted,
			ExpectedDepartureDate:       expectedDepartureDate,
			SITExpected:                 &sitExpected,
			EstimatedWeight:             &estimatedWeight,
			HasProGear:                  &hasProGear,
			PickupAddress:               &pickupAddress,
			DestinationAddress:          &destinationAddress,
			SecondaryPickupAddress:      &secondaryPickupAddress,
			SecondaryDestinationAddress: &secondaryDestinationAddress,
			TertiaryPickupAddress:       &tertiaryPickupAddress,
			TertiaryDestinationAddress:  &tertiaryDestinationAddress,
		}, nil)

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(&estimatedIncentive, nil, nil).Once()

		ppmEstimator.On(
			"MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(&maxIncentive, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		if suite.NotNil(createdPPMShipment) {
			suite.NotZero(createdPPMShipment.ID)
			suite.NotEqual(uuid.Nil.String(), createdPPMShipment.ID.String())
			suite.Equal(expectedDepartureDate, createdPPMShipment.ExpectedDepartureDate)
			suite.Equal(&sitExpected, createdPPMShipment.SITExpected)
			suite.Equal(&estimatedWeight, createdPPMShipment.EstimatedWeight)
			suite.Equal(&hasProGear, createdPPMShipment.HasProGear)
			suite.Equal(models.PPMShipmentStatusDraft, createdPPMShipment.Status)
			suite.Equal(&estimatedIncentive, createdPPMShipment.EstimatedIncentive)
			suite.Equal(&maxIncentive, createdPPMShipment.MaxIncentive)
			suite.NotZero(createdPPMShipment.CreatedAt)
			suite.NotZero(createdPPMShipment.UpdatedAt)
			suite.Equal(pickupAddress.StreetAddress1, createdPPMShipment.PickupAddress.StreetAddress1)
			suite.Equal(secondaryPickupAddress.StreetAddress1, createdPPMShipment.SecondaryPickupAddress.StreetAddress1)
			suite.Equal(tertiaryPickupAddress.StreetAddress1, createdPPMShipment.TertiaryPickupAddress.StreetAddress1)
			suite.Equal(destinationAddress.StreetAddress1, createdPPMShipment.DestinationAddress.StreetAddress1)
			suite.Equal(secondaryDestinationAddress.StreetAddress1, createdPPMShipment.SecondaryDestinationAddress.StreetAddress1)
			suite.Equal(tertiaryDestinationAddress.StreetAddress1, createdPPMShipment.TertiaryDestinationAddress.StreetAddress1)
			suite.NotNil(createdPPMShipment.PickupAddressID)
			suite.NotNil(createdPPMShipment.DestinationAddressID)
			suite.NotNil(createdPPMShipment.SecondaryPickupAddressID)
			suite.NotNil(createdPPMShipment.SecondaryDestinationAddressID)
			suite.NotNil(createdPPMShipment.TertiaryPickupAddressID)
			suite.NotNil(createdPPMShipment.TertiaryDestinationAddressID)
			//ensure HasSecondaryPickupAddress/HasSecondaryDestinationAddress are set even if not initially provided
			suite.True(createdPPMShipment.HasSecondaryPickupAddress != nil)
			suite.True(createdPPMShipment.HasTertiaryPickupAddress != nil)
			suite.Equal(models.BoolPointer(true), createdPPMShipment.HasSecondaryPickupAddress)
			suite.Equal(models.BoolPointer(true), createdPPMShipment.HasTertiaryPickupAddress)
			suite.True(createdPPMShipment.HasSecondaryDestinationAddress != nil)
			suite.True(createdPPMShipment.HasTertiaryDestinationAddress != nil)
			suite.Equal(models.BoolPointer(true), createdPPMShipment.HasSecondaryDestinationAddress)
			suite.Equal(models.BoolPointer(true), createdPPMShipment.HasTertiaryDestinationAddress)
		}
	})

	suite.Run("Can successfully create an international PPM with incentives when existing iHHG shipment is on move", func() {
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		identity, err := models.FetchUserIdentity(suite.DB(), scOfficeUser.User.OktaID)
		suite.NoError(err)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          *scOfficeUser.UserID,
			IDToken:         "fake token",
			ActiveRole:      identity.Roles[0],
		}

		appCtx := suite.AppContextWithSessionForTest(session)

		subtestData := createSubtestData(models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			SITExpected:           models.BoolPointer(false),
			PickupAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				City:           "ANCHORAGE",
				State:          "AK",
				PostalCode:     "99507",
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				City:           "Tulsa",
				State:          "OK",
				PostalCode:     "74133",
			},
		}, nil)

		estimatedIncentive := 123456
		maxIncentive := 7890123
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), nil, nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(maxIncentive)), nil, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, subtestData.newPPMShipment)

		suite.Nil(err)
		suite.NotNil(createdPPMShipment)
		suite.Equal(createdPPMShipment.PPMType, models.PPMTypeIncentiveBased)
		suite.Equal(createdPPMShipment.Shipment.MarketCode, models.MarketCodeInternational)
		suite.NotNil(createdPPMShipment.EstimatedIncentive)
		suite.Equal(*createdPPMShipment.EstimatedIncentive, unit.Cents(estimatedIncentive))
		suite.NotNil(createdPPMShipment.MaxIncentive)
		suite.Equal(*createdPPMShipment.MaxIncentive, unit.Cents(maxIncentive))
	})

	suite.Run("Can update gun safe authorized when HasGunSafe value changes", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		// Set required fields for PPMShipment
		subtestData := createSubtestData(models.PPMShipment{
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			SITExpected:           models.BoolPointer(false),
			PickupAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50308",
				County:         models.StringPointer("POLK"),
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 12345"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "WALESKA",
				State:          "GA",
				PostalCode:     "30183",
				County:         models.StringPointer("COLUMBIA"),
			},
			HasGunSafe: models.BoolPointer(true),
		}, nil)

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil, nil).Once()

		ppmEstimator.On(
			"MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil).Once()

		newPPMShipment := subtestData.newPPMShipment

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, newPPMShipment)
		suite.NotNil(createdPPMShipment)
		suite.NilOrNoVerrs(err)

		var updatedEntitlement models.Entitlement
		err = appCtx.DB().Find(&updatedEntitlement, subtestData.move.Orders.EntitlementID)
		suite.NoError(err)

		suite.True(updatedEntitlement.GunSafe)
	})

	suite.Run("Returns QueryError When Entitlement Is Nil", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: uuid.Must(uuid.NewV4()),
		})

		// Create an order and link it to the entitlement
		orders := factory.BuildOrder(suite.DB(), nil, nil)

		// Manually set EntitlementID to a fake UUID to simulate broken preload
		orders.EntitlementID = nil
		orders.Entitlement = nil
		suite.MustSave(&orders)

		// Create a move and link it to the order
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    orders,
				LinkOnly: true,
			},
		}, nil)

		// Create an MTOShipment linked to the move
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypePPM,
					Status:       models.MTOShipmentStatusDraft,
				},
			},
		}, nil)

		// Initialize a valid PPMShipment properly with the MTOShipment
		newPPMShipment := &models.PPMShipment{
			ShipmentID:            mtoShipment.ID,
			Shipment:              mtoShipment,
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			SITExpected:           models.BoolPointer(false),
			PickupAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 1234"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "Des Moines",
				State:          "IA",
				PostalCode:     "50308",
				County:         models.StringPointer("POLK"),
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "987 Other Avenue",
				StreetAddress2: models.StringPointer("P.O. Box 12345"),
				StreetAddress3: models.StringPointer("c/o Another Person"),
				City:           "WALESKA",
				State:          "GA",
				PostalCode:     "30183",
				County:         models.StringPointer("COLUMBIA"),
			},
			HasGunSafe: models.BoolPointer(true),
		}

		ppmEstimator.On(
			"EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil, nil).Once()

		ppmEstimator.On(
			"MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment"),
		).Return(nil, nil).Once()

		createdPPMShipment, err := ppmShipmentCreator.CreatePPMShipmentWithDefaultCheck(appCtx, newPPMShipment)

		suite.Nil(createdPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.QueryError{}, err)
		suite.Contains(err.Error(), "Move is missing an associated entitlement.")

	})
}

func (suite *PPMShipmentSuite) TestPPMShipmentCreator_StatusMapping() {
	ppmEstimator := &mocks.PPMEstimator{}

	creator := NewPPMShipmentCreator(ppmEstimator, address.NewAddressCreator())

	sc := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *sc.UserID,
		IDToken:         "fake token",
		ActiveRole:      sc.User.Roles[0],
	}
	appCtx := suite.AppContextWithSessionForTest(session)

	makePPM := func(ms models.MoveStatus) *models.PPMShipment {
		move := factory.BuildMove(suite.DB(), []factory.Customization{{
			Model: models.Move{Status: ms},
		}}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypePPM,
				Status:       models.MTOShipmentStatusDraft,
			}},
			{Model: move, LinkOnly: true},
		}, nil)
		return &models.PPMShipment{
			ShipmentID:            shipment.ID,
			Shipment:              shipment,
			ExpectedDepartureDate: testdatagen.NextValidMoveDate,
			SITExpected:           models.BoolPointer(false),
			PPMType:               models.PPMTypeIncentiveBased,
			PickupAddress: &models.Address{
				StreetAddress1: "123 Test St",
				City:           "Tulsa",
				State:          "OK",
				PostalCode:     "74133",
			},
			DestinationAddress: &models.Address{
				StreetAddress1: "456 Test Ave",
				City:           "Beverly Hills",
				State:          "CA",
				PostalCode:     "90210",
			},
		}
	}

	cases := []struct {
		name          string
		moveStatus    models.MoveStatus
		wantPPMStatus models.PPMShipmentStatus
	}{
		{"draft → draft", models.MoveStatusDRAFT, models.PPMShipmentStatusDraft},
		{"needs service counseling → submitted", models.MoveStatusNeedsServiceCounseling, models.PPMShipmentStatusSubmitted},
		{"submitted → waiting on customer", models.MoveStatusSUBMITTED, models.PPMShipmentStatusWaitingOnCustomer},
		{"approvals requested → waiting on customer", models.MoveStatusAPPROVALSREQUESTED, models.PPMShipmentStatusWaitingOnCustomer},
		{"approved → waiting on customer", models.MoveStatusAPPROVED, models.PPMShipmentStatusWaitingOnCustomer},
		{"service counseling completed → waiting on customer", models.MoveStatusServiceCounselingCompleted, models.PPMShipmentStatusWaitingOnCustomer},
	}

	for _, tc := range cases {
		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil, nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil, nil).Once()
		ppm := makePPM(tc.moveStatus)
		created, err := creator.CreatePPMShipmentWithDefaultCheck(appCtx, ppm)
		suite.NoError(err, tc.name)
		suite.Equal(tc.wantPPMStatus, created.Status, tc.name)
	}
}
