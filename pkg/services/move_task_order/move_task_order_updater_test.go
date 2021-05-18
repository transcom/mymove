package movetaskorder_test

import (
	"encoding/base64"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/transcom/mymove/pkg/etag"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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

	queryBuilder := query.NewQueryBuilder(suite.DB())
	mtoUpdater := NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder))

	suite.RunWithRollback("MTO status is updated succesfully", func() {
		eTag := etag.GenerateEtag(expectedMTO.UpdatedAt)

		actualMTO, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(expectedMTO.ID, eTag)

		suite.NoError(err)
		suite.NotZero(actualMTO.ID)
		suite.NotNil(actualMTO.ServiceCounselingCompletedAt)
		suite.Equal(actualMTO.Status, models.MoveStatusServiceCounselingCompleted)
	})

	suite.RunWithRollback("MTO status is in a conflicted state", func() {
		expectedMTO = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusDRAFT,
			},
			Order: expectedOrder,
		})
		eTag := etag.GenerateEtag(expectedMTO.UpdatedAt)

		_, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(expectedMTO.ID, eTag)

		suite.IsType(services.ConflictError{}, err)
	})

	suite.RunWithRollback("Etag is stale", func() {
		expectedMTO = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
			Order: expectedOrder,
		})
		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.UpdateStatusServiceCounselingCompleted(expectedMTO.ID, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_UpdatePostCounselingInfo() {
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())
	expectedMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})

	queryBuilder := query.NewQueryBuilder(suite.DB())
	mtoUpdater := NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder))
	body := movetaskorderops.UpdateMTOPostCounselingInformationBody{
		PpmType:            "FULL",
		PpmEstimatedWeight: 3000,
		PointOfContact:     "user@prime.com",
	}

	suite.RunWithRollback("MTO post counseling information is updated succesfully", func() {
		eTag := base64.StdEncoding.EncodeToString([]byte(expectedMTO.UpdatedAt.Format(time.RFC3339Nano)))

		actualMTO, err := mtoUpdater.UpdatePostCounselingInfo(expectedMTO.ID, body, eTag)

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
		_, err := mtoUpdater.UpdatePostCounselingInfo(expectedMTO.ID, body, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
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
	queryBuilder := query.NewQueryBuilder(suite.DB())
	updater := NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder))

	// Case: Move successfully deactivated
	suite.RunWithRollback("Success - Set show field to false", func() {
		show = false
		updatedMove, err := updater.ShowHide(move.ID, &show)

		suite.NotNil(updatedMove)
		suite.NoError(err)
		suite.Equal(updatedMove.ID, move.ID)
		suite.Equal(*updatedMove.Show, show)
	})

	// Case: Move successfully activated
	suite.RunWithRollback("Success - Set show field to true", func() {
		show = true
		updatedMove, err := updater.ShowHide(move.ID, &show)

		suite.NotNil(updatedMove)
		suite.NoError(err)
		suite.Equal(updatedMove.ID, move.ID)
		suite.Equal(*updatedMove.Show, show)
	})

	// Case: Move UUID not found in DB
	suite.Run("Fail - Move not found", func() {
		badMoveID := uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")
		updatedMove, err := updater.ShowHide(badMoveID, &show)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), badMoveID.String())
	})

	// Case: Show input value is nil, not True or False
	suite.RunWithRollback("Fail - Nil value in show field", func() {
		updatedMove, err := updater.ShowHide(move.ID, nil)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
		suite.Contains(err.Error(), "The 'show' field must be either True or False - it cannot be empty")
	})

	// Case: Invalid input found while updating the move
	// TODO: Is there a way to mock ValidateUpdate so that these tests actually mean something?
	suite.RunWithRollback("Fail - Invalid input found on move", func() {
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mockUpdater.On("ShowHide",
			mock.Anything, // our arguments aren't important here because there's no specific way to trigger this error
			mock.Anything,
		).Return(nil, services.InvalidInputError{})

		updatedMove, err := mockUpdater.ShowHide(move.ID, &show)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	// Case: Query error encountered while updating the move
	suite.RunWithRollback("Fail - Query error", func() {
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mockUpdater.On("ShowHide",
			mock.Anything, // our arguments aren't important here because there's no specific way to trigger this error
			mock.Anything,
		).Return(nil, services.QueryError{})

		updatedMove, err := mockUpdater.ShowHide(move.ID, &show)

		suite.Nil(updatedMove)
		suite.Error(err)
		suite.IsType(services.QueryError{}, err)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_MakeAvailableToPrime() {
	mockserviceItemCreator := &mocks.MTOServiceItemCreator{}
	queryBuilder := query.NewQueryBuilder(suite.DB())
	mtoUpdater := NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, mockserviceItemCreator)

	suite.RunWithRollback("Service item creator is not called if move fails to get approved", func() {
		// Create move in DRAFT status, which should fail to get approved
		move := testdatagen.MakeDefaultMove(suite.DB())
		eTag := etag.GenerateEtag(move.UpdatedAt)
		_, err := mtoUpdater.MakeAvailableToPrime(move.ID, eTag, true, true)

		mockserviceItemCreator.AssertNumberOfCalls(suite.T(), "CreateMTOServiceItem", 0)
		suite.Error(err)
	})

	suite.RunWithRollback("When ETag is stale", func() {
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusSUBMITTED,
			},
		})

		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.MakeAvailableToPrime(move.ID, eTag, false, false)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})
}
