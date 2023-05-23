// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package mtoserviceitem

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
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

		newAddress := factory.BuildAddress(nil, nil, nil)
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

	// Success for DDDSIT
	suite.Run("Success", func() {
		serviceItem, eTag := setupServiceItem()
		serviceItem.ReService.Code = models.ReServiceCodeDDDSIT
		reason := "because we did this service"
		sitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)

		newAddress := factory.BuildAddress(nil, nil, nil)
		newServiceItem := serviceItem
		newServiceItem.Reason = &reason
		newServiceItem.SITEntryDate = &sitEntryDate
		newServiceItem.Status = "" // should keep the status from the original service item
		newServiceItem.SITDestinationFinalAddress = &newAddress
		actualWeight := int64(4000)
		estimatedWeight := int64(4200)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&actualWeight)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&estimatedWeight)
		newServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC),
			},
		}
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
		suite.Equal(newServiceItem.CustomerContacts[0].TimeMilitary, updatedServiceItem.CustomerContacts[0].TimeMilitary)
		suite.Equal(newServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate, updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
		suite.NotEqual(newServiceItem.Status, updatedServiceItem.Status)
	})

	suite.Run("Successful Prime update - adding SITDestinationFinalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to add SITDestinationFinalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.NotNil(updatedServiceItem.SITDestinationFinalAddress)
		suite.Equal(newAddress.StreetAddress1, updatedServiceItem.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(newAddress.StreetAddress2, updatedServiceItem.SITDestinationFinalAddress.StreetAddress2)
		suite.Equal(newAddress.StreetAddress3, updatedServiceItem.SITDestinationFinalAddress.StreetAddress3)
		suite.Equal(newAddress.City, updatedServiceItem.SITDestinationFinalAddress.City)
		suite.Equal(newAddress.State, updatedServiceItem.SITDestinationFinalAddress.State)
		suite.Equal(newAddress.PostalCode, updatedServiceItem.SITDestinationFinalAddress.PostalCode)
	})

	suite.Run("Unsuccessful Prime update - updating existing SITDestinationFinalAddres", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to update SITDestinationFinalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "SITDestinationFinalAddress")
	})

	suite.Run("Unsuccessful basic update - adding SITDestinationOriginalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to update SITDestinationOriginalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationOriginalAddress = &newAddress

		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "SITDestinationOriginalAddress")
	})

	suite.Run("Unsuccessful prime update - adding SITDestinationOriginalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to update SITDestinationOriginalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationOriginalAddress = &newAddress

		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "SITDestinationOriginalAddress")
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
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Change something allowed by Prime:
		reason := "because"
		newServiceItemPrime.Reason = &reason
		newServiceItemPrime.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				TimeMilitary:               "1300Z",
				FirstAvailableDeliveryDate: time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC),
			},
		}
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
		suite.Equal(updatedServiceItem.CustomerContacts[0].TimeMilitary, newServiceItemPrime.CustomerContacts[0].TimeMilitary)
		suite.Equal(updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate, newServiceItemPrime.CustomerContacts[0].FirstAvailableDeliveryDate)
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
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
		}, nil)
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
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Create payment requests for service item:
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: paymentRequest, LinkOnly: true,
			}, {
				Model: oldServiceItemPrime, LinkOnly: true,
			},
		}, nil)

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
	move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)

	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) createServiceItemForUnapprovedMove() (string, models.MTOServiceItem, models.Move) {
	move := factory.BuildMove(suite.DB(), nil, nil)

	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) createServiceItemForMoveWithUnacknowledgedAmendedOrders() (string, models.MTOServiceItem, models.Move) {
	storer := storageTest.NewFakeS3Storage(true)
	userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
	suite.NoError(err)
	amendedDocument := factory.BuildDocument(suite.DB(), nil, nil)
	amendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
		{
			Model:    amendedDocument,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   suite.AppContextForTest(),
			},
		},
	}, nil)

	amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
	now := time.Now()
	move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				ExcessWeightQualifiedAt: &now,
			},
		},
		{
			Model:    amendedDocument,
			LinkOnly: true,
			Type:     &factory.Documents.UploadedAmendedOrders,
		},
		{
			Model:    amendedDocument.ServiceMember,
			LinkOnly: true,
		},
	}, nil)

	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) TestUpdateMTOServiceItemStatus() {
	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	updater := NewMTOServiceItemUpdater(builder, moveRouter)

	rejectionReason := models.StringPointer("")

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
		rejectionReason = models.StringPointer("incomplete")

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
		rejectionReason = models.StringPointer("incomplete")

		_, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusRejected, rejectionReason, "")

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), serviceItem.ID.String())
	})
}
