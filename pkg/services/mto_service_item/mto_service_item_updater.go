package mtoserviceitem

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type mtoServiceItemQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoServiceItemUpdater struct {
	builder mtoServiceItemQueryBuilder
}

// NewMTOServiceItemUpdater returns a new mto service item updater
func NewMTOServiceItemUpdater(builder mtoServiceItemQueryBuilder) services.MTOServiceItemUpdater {
	return &mtoServiceItemUpdater{builder}
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

	mtoServiceItem.Status = status
	updatedAt := time.Now()
	mtoServiceItem.UpdatedAt = updatedAt

	if status == models.MTOServiceItemStatusRejected {
		if rejectionReason == nil {
			return nil, services.NewConflictError(mtoServiceItemID, "Rejecting an MTO Service item requires a rejection reason")
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

	var move models.Move
	moveFilter := []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoServiceItem.MoveTaskOrderID),
	}
	err = p.builder.FetchOne(&move, moveFilter)
	if err != nil {
		return nil, services.NewNotFoundError(mtoServiceItemID, "MTOServiceItemID")
	}

	// If there are no service items that are SUBMITTED then we need to change the move status to MOVE APPROVED
	moveShouldBeMoveApproved := true
	for _, mtoServiceItem := range move.MTOServiceItems {
		if mtoServiceItem.Status == models.MTOServiceItemStatusSubmitted {
			moveShouldBeMoveApproved = false
			break
		}
	}
	// Doing the change
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
	// If we didn't set to MOVE APPROVED and we aren't already at APPROVALS REQUESTED we need to get there
	if move.Status != models.MoveStatusAPPROVALSREQUESTED && move.Status != models.MoveStatusAPPROVED {
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
