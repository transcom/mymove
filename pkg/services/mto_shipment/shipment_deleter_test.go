package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *MTOShipmentServiceSuite) TestShipmentDeleter() {
	builder := query.NewQueryBuilder()
	moveRouter, err := moveservices.NewMoveRouter()
	suite.FatalNoError(err)
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}

	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		builder,
		mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(mockFeatureFlagFetcher), ghcrateengine.NewDomesticPackPricer(mockFeatureFlagFetcher), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(mockFeatureFlagFetcher), ghcrateengine.NewDomesticDestinationPricer(mockFeatureFlagFetcher), ghcrateengine.NewFuelSurchargePricer(), mockFeatureFlagFetcher),
		moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil),
	)
	suite.Run("Returns an error when shipment is not found", func() {
		shipmentDeleter := NewShipmentDeleter(moveTaskOrderUpdater, moveRouter)
		id := uuid.Must(uuid.NewV4())
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := shipmentDeleter.DeleteShipment(session, id)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the Move is neither in Draft nor in NeedsServiceCounseling status", func() {
		moveRouter, err := moveservices.NewMoveRouter()
		suite.FatalNoError(err)
		shipmentDeleter := NewShipmentDeleter(moveTaskOrderUpdater, moveRouter)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		move := shipment.MoveTaskOrder
		move.Status = models.MoveStatusServiceCounselingCompleted
		suite.MustSave(&move)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
			Roles:           factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor}).User.Roles,
		})

		_, err = shipmentDeleter.DeleteShipment(session, shipment.ID)

		suite.Error(err)
		suite.IsType(apperror.ForbiddenError{}, err)
	})

	suite.Run("Soft deletes the shipment when it is found", func() {
		moveRouter, err := moveservices.NewMoveRouter()
		suite.FatalNoError(err)
		shipmentDeleter := NewShipmentDeleter(moveTaskOrderUpdater, moveRouter)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)

		validStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Draft", models.MoveStatusDRAFT},
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
		}
		for _, validStatus := range validStatuses {
			move := shipment.MoveTaskOrder
			move.Status = validStatus.status
			suite.MustSave(&move)
			session := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.OfficeApp,
				OfficeUserID:    uuid.Must(uuid.NewV4()),
			})
			moveID, err := shipmentDeleter.DeleteShipment(session, shipment.ID)
			suite.NoError(err)
			// Verify that the shipment's Move ID is returned because the
			// handler needs it to generate the TriggerEvent.
			suite.Equal(shipment.MoveTaskOrderID, moveID)

			// Verify the shipment still exists in the DB
			var shipmentInDB models.MTOShipment
			err = suite.DB().Find(&shipmentInDB, shipment.ID)
			suite.NoError(err)

			actualDeletedAt := shipmentInDB.DeletedAt
			suite.WithinDuration(time.Now(), *actualDeletedAt, 10*time.Second)

			// Reset the deleted_at field to nil to allow the shipment to be
			// deleted a second time when testing the other move status (a
			// shipment can only be deleted once)
			shipmentInDB.DeletedAt = nil
			suite.MustSave(&shipment)
		}
	})

	suite.Run("Soft deletes the shipment when it is found and check if shipment_seq_num changed", func() {
		moveRouter, err := moveservices.NewMoveRouter()
		suite.FatalNoError(err)
		move := factory.BuildMove(suite.DB(), nil, nil)
		shipmentDeleter := NewShipmentDeleter(moveTaskOrderUpdater, moveRouter)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		validStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Draft", models.MoveStatusDRAFT},
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
		}
		for _, validStatus := range validStatuses {
			move := shipment.MoveTaskOrder
			move.Status = validStatus.status

			var moveInDB models.Move
			err := suite.DB().Find(&moveInDB, move.ID)
			suite.NoError(err)

			move.ShipmentSeqNum = moveInDB.ShipmentSeqNum

			suite.MustSave(&move)

			shipmentSeqNum := move.ShipmentSeqNum
			session := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.OfficeApp,
				OfficeUserID:    uuid.Must(uuid.NewV4()),
			})
			moveID, err := shipmentDeleter.DeleteShipment(session, shipment.ID)
			suite.NoError(err)
			// Verify that the shipment's Move ID is returned because the
			// handler needs it to generate the TriggerEvent.
			suite.Equal(shipment.MoveTaskOrderID, moveID)

			// Verify the shipment still exists in the DB
			var shipmentInDB models.MTOShipment
			err = suite.DB().Find(&shipmentInDB, shipment.ID)
			suite.NoError(err)

			actualDeletedAt := shipmentInDB.DeletedAt
			suite.WithinDuration(time.Now(), *actualDeletedAt, 10*time.Second)

			// Reset the deleted_at field to nil to allow the shipment to be
			// deleted a second time when testing the other move status (a
			// shipment can only be deleted once)
			shipmentInDB.DeletedAt = nil
			suite.MustSave(&shipment)

			// Get updated Move in DB
			err = suite.DB().Find(&moveInDB, move.ID)
			suite.NoError(err)

			suite.Equal(shipmentSeqNum, moveInDB.ShipmentSeqNum)
		}
	})

	suite.Run("Returns not found error when the shipment is already deleted", func() {
		moveRouter, err := moveservices.NewMoveRouter()
		suite.FatalNoError(err)
		shipmentDeleter := NewShipmentDeleter(moveTaskOrderUpdater, moveRouter)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err = shipmentDeleter.DeleteShipment(session, shipment.ID)

		suite.NoError(err)

		// Try to delete the shipment a second time
		_, err = shipmentDeleter.DeleteShipment(session, shipment.ID)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Soft deletes the associated PPM shipment", func() {
		moveRouter, err := moveservices.NewMoveRouter()
		suite.FatalNoError(err)
		shipmentDeleter := NewShipmentDeleter(moveTaskOrderUpdater, moveRouter)
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		moveID, err := shipmentDeleter.DeleteShipment(session, ppmShipment.ShipmentID)
		suite.NoError(err)
		// Verify that the shipment's Move ID is returned because the
		// handler needs it to generate the TriggerEvent.
		suite.Equal(ppmShipment.Shipment.MoveTaskOrderID, moveID)

		// Verify the shipment still exists in the DB
		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("PPMShipment").Find(&shipmentInDB, ppmShipment.ShipmentID)
		suite.NoError(err)

		actualDeletedAt := shipmentInDB.DeletedAt
		suite.WithinDuration(time.Now(), *actualDeletedAt, 10*time.Second)

		actualDeletedAt = shipmentInDB.PPMShipment.DeletedAt
		suite.WithinDuration(time.Now(), *actualDeletedAt, 10*time.Second)
	})
}

func (suite *MTOShipmentServiceSuite) TestPrimeShipmentDeleter() {
	builder := query.NewQueryBuilder()
	moveRouter, err := moveservices.NewMoveRouter()
	suite.FatalNoError(err)
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}

	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		builder,
		mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(mockFeatureFlagFetcher), ghcrateengine.NewDomesticPackPricer(mockFeatureFlagFetcher), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(mockFeatureFlagFetcher), ghcrateengine.NewDomesticDestinationPricer(mockFeatureFlagFetcher), ghcrateengine.NewFuelSurchargePricer(), mockFeatureFlagFetcher),
		moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil),
	)
	suite.Run("Doesn't return an error when allowed to delete a shipment", func() {
		shipmentDeleter := NewPrimeShipmentDeleter(moveTaskOrderUpdater)
		now := time.Now()
		shipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
			{
				Model: models.PPMShipment{
					Status: models.PPMShipmentStatusSubmitted,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := shipmentDeleter.DeleteShipment(session, shipment.ID)
		suite.Error(err)
	})

	suite.Run("Returns an error when a shipment is not available to prime", func() {
		shipmentDeleter := NewPrimeShipmentDeleter(moveTaskOrderUpdater)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: nil,
					ApprovedAt:         nil,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := shipmentDeleter.DeleteShipment(session, shipment.ID)
		suite.Error(err)
	})

	suite.Run("Returns an error when a shipment is not a PPM", func() {
		shipmentDeleter := NewPrimeShipmentDeleter(moveTaskOrderUpdater)
		now := time.Now()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := shipmentDeleter.DeleteShipment(session, shipment.ID)
		suite.Error(err)
	})

	suite.Run("Returns an error when PPM status is WAITING_ON_CUSTOMER", func() {
		shipmentDeleter := NewPrimeShipmentDeleter(moveTaskOrderUpdater)
		now := time.Now()
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := shipmentDeleter.DeleteShipment(session, shipment.ID)
		suite.Error(err)
	})
}
