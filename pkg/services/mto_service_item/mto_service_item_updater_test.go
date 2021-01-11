package mtoserviceitem

import (
	"testing"
	"time"

	"github.com/go-openapi/swag"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/services"
)

type testUpdateMTOServiceItemQueryBuilder struct {
	fakeFetchOne  func(model interface{}, filters []services.QueryFilter) error
	fakeUpdateOne func(models interface{}, eTag *string) (*validate.Errors, error)
}

func (t *testUpdateMTOServiceItemQueryBuilder) UpdateOne(model interface{}, eTag *string) (*validate.Errors, error) {
	return t.fakeUpdateOne(model, eTag)
}

func (t *testUpdateMTOServiceItemQueryBuilder) FetchOne(model interface{}, filters []services.QueryFilter) error {
	return t.fakeFetchOne(model, filters)
}

func (suite *MTOServiceItemServiceSuite) TestMTOServiceItemUpdater() {
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewMTOServiceItemUpdater(builder)

	serviceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	// Test not found error
	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundServiceItem := serviceItem
		notFoundServiceItem.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.DB(), &notFoundServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})

	// Test validation error
	suite.T().Run("Validation Error", func(t *testing.T) {
		invalidServiceItem := serviceItem
		invalidServiceItem.MoveTaskOrderID = serviceItem.ID // invalid Move ID

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.DB(), &invalidServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		invalidInputError := err.(services.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "moveTaskOrderID")
	})

	// Test precondition failed (stale eTag)
	suite.T().Run("Precondition Failed", func(t *testing.T) {
		newServiceItem := serviceItem
		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.DB(), &newServiceItem, "bloop")

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	// Test successful update
	suite.T().Run("Success", func(t *testing.T) {
		reason := "because we did this service"
		sitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)

		newServiceItem := serviceItem
		newServiceItem.Reason = &reason
		newServiceItem.SITEntryDate = &sitEntryDate
		newServiceItem.Status = "" // should keep the status from the original service item

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.DB(), &newServiceItem, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(updatedServiceItem.ID, serviceItem.ID)
		suite.Equal(updatedServiceItem.MTOShipmentID, serviceItem.MTOShipmentID)
		suite.Equal(updatedServiceItem.MoveTaskOrderID, serviceItem.MoveTaskOrderID)
		suite.Equal(updatedServiceItem.Reason, newServiceItem.Reason)
		suite.Equal(updatedServiceItem.SITEntryDate.Local(), newServiceItem.SITEntryDate.Local())
		suite.Equal(updatedServiceItem.Status, serviceItem.Status) // should not have been updated
		suite.NotEqual(updatedServiceItem.Status, newServiceItem.Status)
	})
}

func (suite *MTOServiceItemServiceSuite) TestValidateUpdateMTOServiceItem() {
	// Set up the data needed for updateMTOServiceItemData obj
	checker := movetaskorder.NewMoveTaskOrderChecker(suite.DB())
	oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
	oldServiceItemPrime := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: testdatagen.MakeAvailableMove(suite.DB()),
	})
	now := time.Now()

	// Test with bad string key
	suite.T().Run("bad validatorKey - failure", func(t *testing.T) {
		serviceItemData := updateMTOServiceItemData{}
		fakeKey := "FakeKey"
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, fakeKey)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.Contains(err.Error(), fakeKey)
	})

	// Test successful Basic validation
	suite.T().Run("UpdateMTOServiceItemBasicValidator - success", func(t *testing.T) {
		newServiceItem := models.MTOServiceItem{
			ID:              oldServiceItem.ID,
			MTOShipmentID:   oldServiceItem.MTOShipmentID,
			MoveTaskOrderID: oldServiceItem.MoveTaskOrderID,
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, UpdateMTOServiceItemBasicValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})

	// Test unsuccessful Basic validation
	suite.T().Run("UpdateMTOServiceItemBasicValidator - failure", func(t *testing.T) {
		newServiceItem := models.MTOServiceItem{
			ID:            oldServiceItem.ID,
			MTOShipmentID: &oldServiceItem.ID, // bad value
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, UpdateMTOServiceItemBasicValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	// Test successful Prime validation
	suite.T().Run("UpdateMTOServiceItemPrimeValidator - success", func(t *testing.T) {
		newServiceItemPrime := oldServiceItemPrime

		// Change something allowed by Prime:
		reason := "because"
		newServiceItemPrime.Reason = &reason

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
			db:                  suite.DB(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})

	// Test unsuccessful Prime validation - Not available to Prime
	suite.T().Run("UpdateMTOServiceItemPrimeValidator - not available failure", func(t *testing.T) {
		newServiceItemNotPrime := oldServiceItem // this service item should not be Prime-available

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemNotPrime,
			oldServiceItem:      oldServiceItem,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
			db:                  suite.DB(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	// Test unsuccessful Prime validation - Invalid input
	suite.T().Run("UpdateMTOServiceItemPrimeValidator - invalid input failure", func(t *testing.T) {
		newServiceItemPrime := oldServiceItemPrime

		// Change something unavailable to Prime:
		newServiceItemPrime.Status = models.MTOServiceItemStatusApproved
		newServiceItemPrime.ApprovedAt = &now

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
			db:                  suite.DB(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		invalidInputError := err.(services.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "status")
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "approvedAt")
	})

	// Test unsuccessful Prime validation - Payment requests
	suite.T().Run("UpdateMTOServiceItemPrimeValidator - payment request failure", func(t *testing.T) {
		newServiceItemPrime := oldServiceItemPrime

		// Create payment requests for service item:
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			PaymentRequest: paymentRequest,
			MTOServiceItem: oldServiceItemPrime,
		})

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
			db:                  suite.DB(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(services.ConflictError{}, err)
	})

	// Test with empty string key (successful Base validation)
	suite.T().Run("empty validatorKey - success", func(t *testing.T) {
		newServiceItem := oldServiceItem
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, "")

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})
}

func (suite *MTOServiceItemServiceSuite) createServiceItem() (string, models.MTOServiceItem, models.Move) {
	move := testdatagen.MakeApprovalsRequestedMove(suite.DB())

	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) createServiceItemForUnapprovedMove() (string, models.MTOServiceItem, models.Move) {
	move := testdatagen.MakeDefaultMove(suite.DB())

	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) TestUpdateMTOServiceItemStatus() {
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewMTOServiceItemUpdater(builder)

	rejectionReason := swag.String("")

	// Test that the move's status changes to Approved when the service item's
	// status is no longer SUBMITTED
	suite.T().Run("When TOO reviews move and approves service item", func(t *testing.T) {
		suite.SetupTest()
		eTag, serviceItem, move := suite.createServiceItem()

		updatedServiceItem, err := updater.UpdateMTOServiceItemStatus(
			serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)

		suite.DB().Find(&move, move.ID)
		suite.DB().Find(&serviceItem, serviceItem.ID)

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.Equal(models.MTOServiceItemStatusApproved, serviceItem.Status)
		suite.NotNil(updatedServiceItem)
		suite.NoError(err)
	})

	// Test that the move's status changes to Approvals Requested if any of its service
	// items' status is SUBMITTED
	suite.T().Run("When move is approved and service item is submitted", func(t *testing.T) {
		suite.SetupTest()
		eTag, serviceItem, move := suite.createServiceItem()
		move.Status = models.MoveStatusAPPROVED
		suite.MustSave(&move)

		updatedServiceItem, err := updater.UpdateMTOServiceItemStatus(
			serviceItem.ID, models.MTOServiceItemStatusSubmitted, rejectionReason, eTag)

		suite.DB().Find(&move, move.ID)
		suite.DB().Find(&serviceItem, serviceItem.ID)

		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItem.Status)
		suite.Nil(serviceItem.RejectionReason)
		suite.Nil(serviceItem.RejectedAt)
		suite.Nil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
		suite.NoError(err)
	})

	// Test that the move's status changes to Approved if the service item is
	// rejected
	suite.T().Run("When TOO reviews move and rejects service item", func(t *testing.T) {
		suite.SetupTest()
		eTag, serviceItem, move := suite.createServiceItem()
		rejectionReason = swag.String("incomplete")

		updatedServiceItem, err := updater.UpdateMTOServiceItemStatus(
			serviceItem.ID, models.MTOServiceItemStatusRejected, rejectionReason, eTag)

		suite.DB().Find(&move, move.ID)
		suite.DB().Find(&serviceItem, serviceItem.ID)

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.Equal(models.MTOServiceItemStatusRejected, serviceItem.Status)
		suite.Equal(rejectionReason, serviceItem.RejectionReason)
		suite.NotNil(serviceItem.RejectedAt)
		suite.Nil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
		suite.NoError(err)
	})

	// Test that a service item's status can only be updated if the Move's status
	// is either Approved or Approvals Requested. Neither the Move's status nor
	// the service item's status should be changed if the requirements aren't met.
	suite.T().Run("When the Move has not been approved yet", func(t *testing.T) {
		suite.SetupTest()
		eTag, serviceItem, move := suite.createServiceItemForUnapprovedMove()

		updatedServiceItem, err := updater.UpdateMTOServiceItemStatus(
			serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)

		suite.DB().Find(&move, move.ID)
		suite.DB().Find(&serviceItem, serviceItem.ID)

		suite.Equal(models.MoveStatusDRAFT, move.Status)
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItem.Status)
		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.Contains(err.Error(), "Cannot update a service item on a move that is not currently approved")
	})
}
