package movetaskorder_test

import (
	"encoding/base64"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/apperror"
	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/etag"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_UpdateStatusServiceCounselingCompleted() {
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())
	expectedMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusNeedsServiceCounseling,
		},
		Order: expectedOrder,
	})
	moveRouter := moverouter.NewMoveRouter()

	queryBuilder := query.NewQueryBuilder()
	mtoUpdater := NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
		moveRouter,
	)

	suite.RunWithRollback("MTO status is updated succesfully", func() {
		eTag := etag.GenerateEtag(expectedMTO.UpdatedAt)

		actualMTO, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(suite.AppContextForTest(), expectedMTO.ID, eTag)

		suite.NoError(err)
		suite.NotZero(actualMTO.ID)
		suite.NotNil(actualMTO.ServiceCounselingCompletedAt)
		suite.Equal(actualMTO.Status, models.MoveStatusServiceCounselingCompleted)
	})

	suite.RunWithRollback("MTO status is updated succesfully with facility info", func() {
		storageFacility := testdatagen.MakeStorageFacility(suite.DB(), testdatagen.Assertions{
			StorageFacility: models.StorageFacility{
				Address: testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
					Address: models.Address{
						StreetAddress1: "1234 Over Here Street",
						City:           "Houston",
						State:          "TX",
						PostalCode:     "77083",
						Country:        swag.String("US"),
					},
				}),
				Email: swag.String("old@email.com"),
			},
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				StorageFacility: &storageFacility,
			},
		})
		var mtoShipments []models.MTOShipment
		mtoShipments = append(mtoShipments, shipment)
		expectedMTOWithFacility := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status:       models.MoveStatusNeedsServiceCounseling,
				MTOShipments: mtoShipments,
			},
			Order: expectedOrder,
		})
		eTag := etag.GenerateEtag(expectedMTOWithFacility.UpdatedAt)

		actualMTO, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(suite.AppContextForTest(), expectedMTOWithFacility.ID, eTag)

		suite.NoError(err)
		suite.NotZero(actualMTO.ID)
		suite.NotNil(actualMTO.ServiceCounselingCompletedAt)
		suite.Equal(actualMTO.Status, models.MoveStatusServiceCounselingCompleted)
	})

	suite.RunWithRollback("Invalid input error when there is no facility information on NTS-r shipment", func() {
		noFacilityInfoMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
			Order: expectedOrder,
		})
		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
			},
			Move: noFacilityInfoMove,
		})
		eTag := etag.GenerateEtag(noFacilityInfoMove.UpdatedAt)

		_, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(suite.AppContextForTest(), noFacilityInfoMove.ID, eTag)

		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "NTS-release shipment must include facility info")
	})

	suite.RunWithRollback("MTO status is in a conflicted state", func() {
		expectedMTO = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusDRAFT,
			},
			Order: expectedOrder,
		})
		eTag := etag.GenerateEtag(expectedMTO.UpdatedAt)

		_, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(suite.AppContextForTest(), expectedMTO.ID, eTag)

		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.RunWithRollback("Etag is stale", func() {
		expectedMTO = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
			Order: expectedOrder,
		})
		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(suite.AppContextForTest(), expectedMTO.ID, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_UpdatePostCounselingInfo() {
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())
	expectedMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})

	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	mtoUpdater := NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
		moveRouter,
	)
	body := movetaskorderops.UpdateMTOPostCounselingInformationBody{
		PpmType:            "FULL",
		PpmEstimatedWeight: 3000,
		PointOfContact:     "user@prime.com",
	}

	suite.RunWithRollback("MTO post counseling information is updated succesfully", func() {
		eTag := base64.StdEncoding.EncodeToString([]byte(expectedMTO.UpdatedAt.Format(time.RFC3339Nano)))

		actualMTO, err := mtoUpdater.UpdatePostCounselingInfo(suite.AppContextForTest(), expectedMTO.ID, body, eTag)

		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		suite.NotNil(expectedMTO.Orders.ServiceMember.FirstName)
		suite.NotNil(expectedMTO.Orders.ServiceMember.LastName)
		suite.NotNil(expectedMTO.Orders.NewDutyStation.Address.City)
		suite.NotNil(expectedMTO.Orders.NewDutyStation.Address.State)
	})

	suite.RunWithRollback("Etag is stale", func() {
		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.UpdatePostCounselingInfo(suite.AppContextForTest(), expectedMTO.ID, body, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_ShowHide() {
	// Set up a default move:
	show := true
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Show: &show,
		},
	})

	// Set up the necessary updater objects:
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	updater := NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
		moveRouter,
	)

	// Case: Move successfully deactivated
	suite.RunWithRollback("Success - Set show field to false", func() {
		show = false
		updatedMove, err := updater.ShowHide(suite.AppContextForTest(), move.ID, &show)

		suite.NotNil(updatedMove)
		suite.NoError(err)
		suite.Equal(updatedMove.ID, move.ID)
		suite.Equal(*updatedMove.Show, show)
	})

	// Case: Move successfully activated
	suite.RunWithRollback("Success - Set show field to true", func() {
		show = true
		updatedMove, err := updater.ShowHide(suite.AppContextForTest(), move.ID, &show)

		suite.NotNil(updatedMove)
		suite.NoError(err)
		suite.Equal(updatedMove.ID, move.ID)
		suite.Equal(*updatedMove.Show, show)
	})

	// Case: Move UUID not found in DB
	suite.Run("Fail - Move not found", func() {
		badMoveID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		updatedMove, err := updater.ShowHide(suite.AppContextForTest(), badMoveID, &show)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), badMoveID.String())
	})

	// Case: Show input value is nil, not True or False
	suite.RunWithRollback("Fail - Nil value in show field", func() {
		updatedMove, err := updater.ShowHide(suite.AppContextForTest(), move.ID, nil)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Contains(err.Error(), "The 'show' field must be either True or False - it cannot be empty")
	})

	// Case: Invalid input found while updating the move
	// TODO: Is there a way to mock ValidateUpdate so that these tests actually mean something?
	suite.RunWithRollback("Fail - Invalid input found on move", func() {
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mockUpdater.On("ShowHide",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, // our arguments aren't important here because there's no specific way to trigger this error
			mock.Anything,
		).Return(nil, apperror.InvalidInputError{})

		updatedMove, err := mockUpdater.ShowHide(suite.AppContextForTest(), move.ID, &show)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	// Case: Query error encountered while updating the move
	suite.RunWithRollback("Fail - Query error", func() {
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mockUpdater.On("ShowHide",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, // our arguments aren't important here because there's no specific way to trigger this error
			mock.Anything,
		).Return(nil, apperror.QueryError{})

		updatedMove, err := mockUpdater.ShowHide(suite.AppContextForTest(), move.ID, &show)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(apperror.QueryError{}, err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_MakeAvailableToPrime() {
	suite.Run("Service item creator is not called if move fails to get approved", func() {
		mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, mockserviceItemCreator, moveRouter)
		// Create move in DRAFT status, which should fail to get approved
		move := testdatagen.MakeDefaultMove(suite.DB())
		eTag := etag.GenerateEtag(move.UpdatedAt)
		fetchedMove := models.Move{}

		_, err := mtoUpdater.MakeAvailableToPrime(suite.AppContextForTest(), move.ID, eTag, true, true)

		mockserviceItemCreator.AssertNumberOfCalls(suite.T(), "CreateMTOServiceItem", 0)
		suite.Error(err)
		err = suite.DB().Find(&fetchedMove, move.ID)
		suite.NoError(err)
		suite.Nil(fetchedMove.AvailableToPrimeAt)
	})

	suite.Run("When ETag is stale", func() {
		mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, mockserviceItemCreator, moveRouter)

		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		})

		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.MakeAvailableToPrime(suite.AppContextForTest(), move.ID, eTag, true, true)

		mockserviceItemCreator.AssertNumberOfCalls(suite.T(), "CreateMTOServiceItem", 0)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Makes move available to Prime and creates Move management and Service counseling service items when both are specified", func() {
		suite.createMSAndCSReServices()

		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, serviceItemCreator, moveRouter)

		move := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
		eTag := etag.GenerateEtag(move.UpdatedAt)
		fetchedMove := models.Move{}
		var serviceItems models.MTOServiceItems

		suite.Nil(move.AvailableToPrimeAt)

		updatedMove, err := mtoUpdater.MakeAvailableToPrime(suite.AppContextForTest(), move.ID, eTag, true, true)

		suite.NoError(err)
		suite.NotNil(updatedMove.AvailableToPrimeAt)
		suite.Equal(models.MoveStatusAPPROVED, updatedMove.Status)
		err = suite.DB().Eager("ReService").Where("move_id = ?", move.ID).All(&serviceItems)
		suite.NoError(err)
		suite.Len(serviceItems, 2, "Expected to find at most 2 service items")
		suite.True(suite.containsServiceCode(serviceItems, models.ReServiceCodeMS), "Expected to find reServiceCode, MS, in array.")
		suite.True(suite.containsServiceCode(serviceItems, models.ReServiceCodeCS), "Expected to find reServiceCode, CS, in array.")
		err = suite.DB().Find(&fetchedMove, move.ID)
		suite.NoError(err)
		suite.NotNil(fetchedMove.AvailableToPrimeAt)
		suite.Equal(models.MoveStatusAPPROVED, fetchedMove.Status)
	})

	suite.Run("Makes move available to Prime and only creates Move management when it's the only one specified", func() {
		suite.createMSAndCSReServices()

		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, serviceItemCreator, moveRouter)

		move := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
		eTag := etag.GenerateEtag(move.UpdatedAt)
		fetchedMove := models.Move{}
		var serviceItems models.MTOServiceItems

		suite.Nil(move.AvailableToPrimeAt)

		_, err := mtoUpdater.MakeAvailableToPrime(suite.AppContextForTest(), move.ID, eTag, true, false)

		suite.NoError(err)
		err = suite.DB().Find(&fetchedMove, move.ID)
		suite.NoError(err)
		suite.NotNil(fetchedMove.AvailableToPrimeAt)
		err = suite.DB().Eager("ReService").Where("move_id = ?", move.ID).All(&serviceItems)
		suite.NoError(err)
		suite.Len(serviceItems, 1, "Expected to find at most 1 service item")
		suite.True(suite.containsServiceCode(serviceItems, models.ReServiceCodeMS), "Expected to find reServiceCode, MS, in array.")
		suite.False(suite.containsServiceCode(serviceItems, models.ReServiceCodeCS), "Expected to find reServiceCode, CS, in array.")
	})

	suite.Run("Makes move available to Prime and only creates CS service item when it's the only one specified", func() {
		suite.createMSAndCSReServices()

		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, serviceItemCreator, moveRouter)

		move := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
		eTag := etag.GenerateEtag(move.UpdatedAt)
		fetchedMove := models.Move{}
		var serviceItems models.MTOServiceItems

		suite.Nil(move.AvailableToPrimeAt)

		_, err := mtoUpdater.MakeAvailableToPrime(suite.AppContextForTest(), move.ID, eTag, false, true)

		suite.NoError(err)
		err = suite.DB().Find(&fetchedMove, move.ID)
		suite.NoError(err)
		suite.NotNil(fetchedMove.AvailableToPrimeAt)
		err = suite.DB().Eager("ReService").Where("move_id = ?", move.ID).All(&serviceItems)
		suite.NoError(err)
		suite.Len(serviceItems, 1, "Expected to find at most 1 service item")
		suite.False(suite.containsServiceCode(serviceItems, models.ReServiceCodeMS), "Expected to find reServiceCode, MS, in array.")
		suite.True(suite.containsServiceCode(serviceItems, models.ReServiceCodeCS), "Expected to find reServiceCode, CS, in array.")
	})

	suite.Run("Does not create service items if neither CS nor MS are requested", func() {
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, mockserviceItemCreator, moveRouter)

		move := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
		eTag := etag.GenerateEtag(move.UpdatedAt)
		fetchedMove := models.Move{}

		suite.Nil(move.AvailableToPrimeAt)

		_, err := mtoUpdater.MakeAvailableToPrime(suite.AppContextForTest(), move.ID, eTag, false, false)

		mockserviceItemCreator.AssertNumberOfCalls(suite.T(), "CreateMTOServiceItem", 0)
		suite.NoError(err)
		err = suite.DB().Find(&fetchedMove, move.ID)
		suite.NoError(err)
		suite.NotNil(fetchedMove.AvailableToPrimeAt)
	})

	suite.Run("Does not make move available to prime if Order is missing required fields", func() {
		mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, mockserviceItemCreator, moveRouter)

		orderWithoutDefaults := testdatagen.MakeOrderWithoutDefaults(suite.DB(), testdatagen.Assertions{})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusServiceCounselingCompleted,
			},
			Order: orderWithoutDefaults,
		})
		eTag := etag.GenerateEtag(move.UpdatedAt)
		fetchedMove := models.Move{}

		_, err := mtoUpdater.MakeAvailableToPrime(suite.AppContextForTest(), move.ID, eTag, true, true)

		mockserviceItemCreator.AssertNumberOfCalls(suite.T(), "CreateMTOServiceItem", 0)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		err = suite.DB().Find(&fetchedMove, move.ID)
		suite.NoError(err)
		suite.Nil(fetchedMove.AvailableToPrimeAt)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_BillableWeightsReviewedAt() {
	suite.Run("Service item creator is not called if move fails to get approved", func() {
		mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, mockserviceItemCreator, moveRouter)
		move := testdatagen.MakeDefaultMove(suite.DB())
		eTag := etag.GenerateEtag(move.UpdatedAt)

		updatedMove, err := mtoUpdater.UpdateReviewedBillableWeightsAt(suite.AppContextForTest(), move.ID, eTag)

		suite.NoError(err)
		suite.NotNil(updatedMove.BillableWeightsReviewedAt)
	})

	suite.Run("When ETag is stale", func() {
		mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, mockserviceItemCreator, moveRouter)

		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		})

		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.UpdateReviewedBillableWeightsAt(suite.AppContextForTest(), move.ID, eTag)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_TIORemarks() {
	remarks := "Reweigh requested"
	mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	mtoUpdater := NewMoveTaskOrderUpdater(queryBuilder, mockserviceItemCreator, moveRouter)
	suite.Run("Service item creator is not called if move fails to get approved", func() {
		move := testdatagen.MakeDefaultMove(suite.DB())
		eTag := etag.GenerateEtag(move.UpdatedAt)

		updatedMove, err := mtoUpdater.UpdateTIORemarks(suite.AppContextForTest(), move.ID, eTag, remarks)

		suite.NoError(err)
		suite.NotNil(updatedMove.TIORemarks)
	})

	suite.Run("When ETag is stale", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		})

		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.UpdateTIORemarks(suite.AppContextForTest(), move.ID, eTag, remarks)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Fail - Move not found", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		})
		eTag := etag.GenerateEtag(move.UpdatedAt)

		badMoveID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		_, err := mtoUpdater.UpdateTIORemarks(suite.AppContextForTest(), badMoveID, eTag, remarks)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), badMoveID.String())
	})
}

func (suite *MoveTaskOrderServiceSuite) containsServiceCode(items models.MTOServiceItems, target models.ReServiceCode) bool {
	for _, si := range items {
		if si.ReService.Code == target {
			return true
		}
	}

	return false
}

func (suite *MoveTaskOrderServiceSuite) createMSAndCSReServices() {
	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"),
			Code: "MS",
		},
	})

	testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
		ReService: models.ReService{
			ID:   uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"),
			Code: "CS",
		},
	})
}
