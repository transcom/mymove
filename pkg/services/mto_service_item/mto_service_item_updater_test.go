//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package mtoserviceitem

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *MTOServiceItemServiceSuite) TestMTOServiceItemUpdater() {

	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	updater := NewMTOServiceItemUpdater(builder, moveRouter)

	setupServiceItem := func() (models.MTOServiceItem, string) {
		serviceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		return serviceItem, eTag
	}

	// Test not found error
	suite.Run("Not Found Error", func() {
		serviceItem, eTag := setupServiceItem()
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundServiceItem := serviceItem
		notFoundServiceItem.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &notFoundServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})

	// Test validation error
	suite.Run("Validation Error", func() {
		serviceItem, eTag := setupServiceItem()
		invalidServiceItem := serviceItem
		invalidServiceItem.MoveTaskOrderID = serviceItem.ID // invalid Move ID

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &invalidServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "moveTaskOrderID")
	})

	// Test precondition failed (stale eTag)
	suite.Run("Precondition Failed", func() {
		serviceItem, _ := setupServiceItem()
		newServiceItem := serviceItem
		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &newServiceItem, "bloop")

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	// Test successful update
	suite.Run("Success", func() {
		serviceItem, eTag := setupServiceItem()
		reason := "because we did this service"
		sitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)

		newAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{Stub: true})
		newServiceItem := serviceItem
		newServiceItem.Reason = &reason
		newServiceItem.SITEntryDate = &sitEntryDate
		newServiceItem.Status = "" // should keep the status from the original service item
		newServiceItem.SITDestinationFinalAddress = &newAddress
		actualWeight := int64(4000)
		estimatedWeight := int64(4200)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&actualWeight)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&estimatedWeight)

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &newServiceItem, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(serviceItem.ID, updatedServiceItem.ID)
		suite.Equal(serviceItem.MTOShipmentID, updatedServiceItem.MTOShipmentID)
		suite.Equal(serviceItem.MoveTaskOrderID, updatedServiceItem.MoveTaskOrderID)
		suite.Equal(newServiceItem.Reason, updatedServiceItem.Reason)
		suite.Equal(newServiceItem.SITEntryDate.Local(), updatedServiceItem.SITEntryDate.Local())
		suite.Equal(serviceItem.Status, updatedServiceItem.Status) // should not have been updated
		suite.Equal(newAddress.StreetAddress1, updatedServiceItem.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(newAddress.City, updatedServiceItem.SITDestinationFinalAddress.City)
		suite.Equal(newAddress.State, updatedServiceItem.SITDestinationFinalAddress.State)
		suite.Equal(newAddress.Country, updatedServiceItem.SITDestinationFinalAddress.Country)
		suite.Equal(newAddress.PostalCode, updatedServiceItem.SITDestinationFinalAddress.PostalCode)
		suite.Equal(newServiceItem.ActualWeight, updatedServiceItem.ActualWeight)
		suite.Equal(newServiceItem.EstimatedWeight, updatedServiceItem.EstimatedWeight)
		suite.NotEqual(newServiceItem.Status, updatedServiceItem.Status)
	})
}

func (suite *MTOServiceItemServiceSuite) TestValidateUpdateMTOServiceItem() {
	// Set up the data needed for updateMTOServiceItemData obj
	checker := movetaskorder.NewMoveTaskOrderChecker()
	now := time.Now()

	// Test with bad string key
	suite.Run("bad validatorKey - failure", func() {
		serviceItemData := updateMTOServiceItemData{}
		fakeKey := "FakeKey"
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, fakeKey)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.Contains(err.Error(), fakeKey)
	})

	// Test successful Basic validation
	suite.Run("UpdateMTOServiceItemBasicValidator - success", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
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
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemBasicValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})

	// Test unsuccessful Basic validation
	suite.Run("UpdateMTOServiceItemBasicValidator - failure", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		newServiceItem := models.MTOServiceItem{
			ID:            oldServiceItem.ID,
			MTOShipmentID: &oldServiceItem.ID, // bad value
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemBasicValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	// Test successful Prime validation
	suite.Run("UpdateMTOServiceItemPrimeValidator - success", func() {
		oldServiceItemPrime := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
		})
		newServiceItemPrime := oldServiceItemPrime

		// Change something allowed by Prime:
		reason := "because"
		newServiceItemPrime.Reason = &reason

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})

	// Test unsuccessful Prime validation - Not available to Prime
	suite.Run("UpdateMTOServiceItemPrimeValidator - not available failure", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		newServiceItemNotPrime := oldServiceItem // this service item should not be Prime-available

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemNotPrime,
			oldServiceItem:      oldServiceItem,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	// Test unsuccessful Prime validation - Invalid input
	suite.Run("UpdateMTOServiceItemPrimeValidator - invalid input failure", func() {
		oldServiceItemPrime := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
		})
		newServiceItemPrime := oldServiceItemPrime

		// Change something unavailable to Prime:
		newServiceItemPrime.Status = models.MTOServiceItemStatusApproved
		newServiceItemPrime.ApprovedAt = &now

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "status")
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "approvedAt")
	})

	// Test unsuccessful Prime validation - Payment requests
	suite.Run("UpdateMTOServiceItemPrimeValidator - payment request failure", func() {
		oldServiceItemPrime := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
		})
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
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	// Test with empty string key (successful Base validation)
	suite.Run("empty validatorKey - success", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		newServiceItem := oldServiceItem
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, "")

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})
}

func (suite *MTOServiceItemServiceSuite) createServiceItem() (string, models.MTOServiceItem, models.Move) {
	move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})

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

func (suite *MTOServiceItemServiceSuite) createServiceItemForMoveWithUnacknowledgedAmendedOrders() (string, models.MTOServiceItem, models.Move) {
	storer := storageTest.NewFakeS3Storage(true)
	userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
	suite.NoError(err)
	amendedDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
	amendedUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
		UserUpload: models.UserUpload{
			DocumentID: &amendedDocument.ID,
			Document:   amendedDocument,
			UploaderID: amendedDocument.ServiceMember.UserID,
		},
		UserUploader: userUploader,
	})

	amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
	now := time.Now()
	move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			UploadedAmendedOrders:   &amendedDocument,
			UploadedAmendedOrdersID: &amendedDocument.ID,
			ServiceMember:           amendedDocument.ServiceMember,
			ServiceMemberID:         amendedDocument.ServiceMemberID,
		},
		Move: models.Move{ExcessWeightQualifiedAt: &now},
	})

	serviceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: move,
	})

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) TestUpdateMTOServiceItemStatus() {
	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	updater := NewMTOServiceItemUpdater(builder, moveRouter)

	rejectionReason := swag.String("")

	// Test that the move's status changes to Approved when the service item's
	// status is no longer SUBMITTED
	suite.Run("When TOO reviews move and approves service item", func() {
		eTag, serviceItem, move := suite.createServiceItem()

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.Equal(models.MTOServiceItemStatusApproved, serviceItem.Status)
		suite.NotNil(serviceItem.ApprovedAt)
		suite.Nil(serviceItem.RejectionReason)
		suite.Nil(serviceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
	})

	// Test that the move's status changes to Approvals Requested if any of its service
	// items' status is SUBMITTED
	suite.Run("When move is approved and service item is submitted", func() {
		eTag, serviceItem, move := suite.createServiceItem()
		move.Status = models.MoveStatusAPPROVED
		suite.MustSave(&move)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusSubmitted, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItem.Status)
		suite.Nil(serviceItem.RejectionReason)
		suite.Nil(serviceItem.RejectedAt)
		suite.Nil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
	})

	// Test that the move's status changes to Approved if the service item is
	// rejected
	suite.Run("When TOO reviews move and rejects service item", func() {
		eTag, serviceItem, move := suite.createServiceItem()
		rejectionReason = swag.String("incomplete")

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusRejected, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.Equal(models.MTOServiceItemStatusRejected, serviceItem.Status)
		suite.Equal(rejectionReason, serviceItem.RejectionReason)
		suite.NotNil(serviceItem.RejectedAt)
		suite.Nil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
	})

	// Test that a service item's status can only be updated if the Move's status
	// is either Approved or Approvals Requested. Neither the Move's status nor
	// the service item's status should be changed if the requirements aren't met.
	suite.Run("When the Move has not been approved yet", func() {
		eTag, serviceItem, move := suite.createServiceItemForUnapprovedMove()

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)

		suite.Error(err)
		suite.Contains(err.Error(), "Cannot approve or reject a service item if the move's status is neither Approved nor Approvals Requested.")

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusDRAFT, move.Status)
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItem.Status)
		suite.Nil(updatedServiceItem)
	})

	suite.Run("does not approve the move if unacknowledged amended orders exist", func() {

		eTag, serviceItem, move := suite.createServiceItemForMoveWithUnacknowledgedAmendedOrders()
		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.Equal(models.MTOServiceItemStatusApproved, serviceItem.Status)
		suite.Nil(serviceItem.RejectionReason)
		suite.Nil(serviceItem.RejectedAt)
		suite.NotNil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
	})

	suite.Run("Returns an error when eTag is stale", func() {
		_, serviceItem, _ := suite.createServiceItem()
		rejectionReason = swag.String("incomplete")

		_, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusRejected, rejectionReason, "")

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), serviceItem.ID.String())
	})
}
