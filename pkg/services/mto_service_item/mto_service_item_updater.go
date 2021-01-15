package mtoserviceitem

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type mtoServiceItemQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoServiceItemUpdater struct {
	builder mtoServiceItemQueryBuilder
}

// NewMTOServiceItemUpdater returns a new mto service item updater
func NewMTOServiceItemUpdater(builder mtoServiceItemQueryBuilder) services.MTOServiceItemUpdater {
	return &mtoServiceItemUpdater{builder: builder}
}

func (p *mtoServiceItemUpdater) UpdateMTOServiceItemStatus(mtoServiceItemID uuid.UUID, status models.MTOServiceItemStatus, rejectionReason *string, eTag string) (*models.MTOServiceItem, error) {
	var mtoServiceItem models.MTOServiceItem

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoServiceItemID),
	}
	err := p.builder.FetchOne(&mtoServiceItem, queryFilters)

	if err != nil {
		return nil, services.NewNotFoundError(mtoServiceItemID, "MTOServiceItemID")
	}

	var move models.Move
	moveFilter := []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoServiceItem.MoveTaskOrderID),
	}
	err = p.builder.FetchOne(&move, moveFilter)
	moveID := mtoServiceItem.MoveTaskOrderID
	if err != nil {
		return nil, services.NewNotFoundError(moveID, "MoveID")
	}

	// Service item status can only be updated if a Move's status is either Approved
	// or Approvals Requested, so check and fail early.
	if move.Status != models.MoveStatusAPPROVED && move.Status != models.MoveStatusAPPROVALSREQUESTED {
		return nil, services.NewConflictError(
			move.ID,
			fmt.Sprintf("Cannot update a service item on a move that is not currently approved. The current status for the move with ID %s is %s", move.ID, move.Status),
		)
	}

	mtoServiceItem.Status = status
	updatedAt := time.Now()
	mtoServiceItem.UpdatedAt = updatedAt

	if status == models.MTOServiceItemStatusRejected {
		if rejectionReason == nil {
			verrs := validate.NewErrors()
			verrs.Add("rejectionReason", "field must be provided when status is set to REJECTED")
			err := services.NewInvalidInputError(mtoServiceItemID, nil, verrs, "Invalid input found in the request.")
			return nil, err
		}
		mtoServiceItem.RejectionReason = rejectionReason
		mtoServiceItem.RejectedAt = &updatedAt
		// clear field if previously accepted
		mtoServiceItem.ApprovedAt = nil
	} else if status == models.MTOServiceItemStatusApproved {
		// clear fields if previously rejected
		mtoServiceItem.RejectionReason = nil
		mtoServiceItem.RejectedAt = nil
		mtoServiceItem.ApprovedAt = &updatedAt
	}

	verrs, err := p.builder.UpdateOne(&mtoServiceItem, &eTag)

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(mtoServiceItemID, err, verrs, "")
	}

	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(mtoServiceItemID, "")
		}

		switch err.(type) {
		case query.StaleIdentifierError:
			return &models.MTOServiceItem{}, services.NewPreconditionFailedError(mtoServiceItemID, err)
		}
	}

	// If there are no service items that are SUBMITTED then we need to change
	// the move status to APPROVED
	moveShouldBeMoveApproved := true

	if status == models.MTOServiceItemStatusSubmitted {
		moveShouldBeMoveApproved = false
	} else {
		// Fetch move again since other service items could have been updated
		err = p.builder.FetchOne(&move, moveFilter)
		if err != nil {
			return nil, services.NewNotFoundError(mtoServiceItemID, "MTOServiceItemID")
		}

		for _, mtoServiceItem := range move.MTOServiceItems {
			if mtoServiceItem.Status == models.MTOServiceItemStatusSubmitted {
				moveShouldBeMoveApproved = false
				break
			}
		}
	}

	if moveShouldBeMoveApproved {
		err = move.Approve()
		if err != nil {
			return nil, err
		}
		verrs, err = p.builder.UpdateOne(&move, nil)
		if verrs != nil && verrs.HasAny() {
			return nil, services.NewInvalidInputError(move.ID, err, verrs, "")
		}

		if err != nil {
			if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
				return nil, services.NewNotFoundError(move.ID, "")
			}
		}
	}
	// If the move should not be APPROVED (due to one or more service items
	// being in SUBMITTED status) then it needs to be set to APPROVALS REQUESTED
	// so the TOO can review it.
	if !moveShouldBeMoveApproved {
		err = move.SetApprovalsRequested()
		if err != nil {
			return nil, err
		}
		verrs, err = p.builder.UpdateOne(&move, nil)
		if verrs != nil && verrs.HasAny() {
			return nil, services.NewInvalidInputError(move.ID, err, verrs, "")
		}

		if err != nil {
			if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
				return nil, services.NewNotFoundError(move.ID, "")
			}
		}
	}

	return &mtoServiceItem, err
}

// UpdateMTOServiceItemBasic updates the MTO Service Item using base validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemBasic(db *pop.Connection, mtoServiceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error) {
	return p.UpdateMTOServiceItem(db, mtoServiceItem, eTag, UpdateMTOServiceItemBasicValidator)
}

// UpdateMTOServiceItemPrime updates the MTO Service Item using Prime API validators
func (p *mtoServiceItemUpdater) UpdateMTOServiceItemPrime(db *pop.Connection, mtoServiceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error) {
	return p.UpdateMTOServiceItem(db, mtoServiceItem, eTag, UpdateMTOServiceItemPrimeValidator)
}

// UpdateMTOServiceItem updates the given service item
func (p *mtoServiceItemUpdater) UpdateMTOServiceItem(db *pop.Connection, mtoServiceItem *models.MTOServiceItem, eTag string, validatorKey string) (*models.MTOServiceItem, error) {
	oldServiceItem := models.MTOServiceItem{}

	// Find the service item, return error if not found
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoServiceItem.ID),
	}
	err := p.builder.FetchOne(&oldServiceItem, queryFilters)
	if err != nil {
		return nil, services.NewNotFoundError(mtoServiceItem.ID, "while looking for MTOServiceItem")
	}

	checker := movetaskorder.NewMoveTaskOrderChecker(db)
	serviceItemData := updateMTOServiceItemData{
		updatedServiceItem:  *mtoServiceItem,
		oldServiceItem:      oldServiceItem,
		availabilityChecker: checker,
		db:                  db,
		verrs:               validate.NewErrors(),
	}

	validServiceItem, err := ValidateUpdateMTOServiceItem(&serviceItemData, validatorKey)
	if err != nil {
		return nil, err
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldServiceItem.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, services.NewPreconditionFailedError(validServiceItem.ID, nil)
	}

	if validServiceItem.SITDestinationFinalAddress != nil {
		verrs, createErr := p.builder.CreateOne(validServiceItem.SITDestinationFinalAddress)
		if verrs != nil || createErr != nil {
			return nil, fmt.Errorf("%#v %e", verrs, createErr)
		}
		validServiceItem.SITDestinationFinalAddressID = &validServiceItem.SITDestinationFinalAddress.ID
	}

	// Make the update and create a InvalidInputError if there were validation issues
	verrs, err := p.builder.UpdateOne(validServiceItem, &eTag)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(validServiceItem.ID, err, verrs, "Invalid input found while updating the service item.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("MTOServiceItem", err, "")
	}

	// Get the updated address and return
	updatedServiceItem := models.MTOServiceItem{}
	err = p.builder.FetchOne(&updatedServiceItem, queryFilters) // using the same queryFilters set at the beginning
	if err != nil {
		return nil, services.NewQueryError("MTOServiceItem", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedServiceItem, nil
}

// ValidateUpdateMTOServiceItem checks the provided serviceItemData struct against the validator indicated by validatorKey.
// Defaults to base validation if the empty string is entered as the key.
// Returns an MTOServiceItem that has been set up for update.
func ValidateUpdateMTOServiceItem(serviceItemData *updateMTOServiceItemData, validatorKey string) (*models.MTOServiceItem, error) {
	if validatorKey == "" {
		validatorKey = UpdateMTOServiceItemBasicValidator
	}
	validator, ok := UpdateMTOServiceItemValidators[validatorKey]
	if !ok {
		err := fmt.Errorf("validator key %s was not found in update MTO Service Item validators", validatorKey)
		return nil, err
	}
	err := validator.validate(serviceItemData)
	if err != nil {
		return nil, err
	}

	newServiceItem := serviceItemData.setNewMTOServiceItem()

	return newServiceItem, nil
}
