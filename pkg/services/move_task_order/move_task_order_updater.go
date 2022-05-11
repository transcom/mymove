package movetaskorder

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/services/query"
)

type moveTaskOrderUpdater struct {
	moveTaskOrderFetcher
	builder            UpdateMoveTaskOrderQueryBuilder
	serviceItemCreator services.MTOServiceItemCreator
	moveRouter         services.MoveRouter
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(builder UpdateMoveTaskOrderQueryBuilder, serviceItemCreator services.MTOServiceItemCreator, moveRouter services.MoveRouter) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{moveTaskOrderFetcher{}, builder, serviceItemCreator, moveRouter}
}

// UpdateStatusServiceCounselingCompleted updates the status on the move (move task order) to service counseling completed
func (o moveTaskOrderUpdater) UpdateStatusServiceCounselingCompleted(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	// Fetch the move and associations.
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, fetchErr := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if fetchErr != nil {
		return &models.Move{}, fetchErr
	}

	// Check the If-Match header against existing eTag before updating.
	encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return &models.Move{}, apperror.NewPreconditionFailedError(move.ID, nil)
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// Update move status, verifying that move/shipments are in expected state.
		err := o.moveRouter.CompleteServiceCounseling(appCtx, move)
		if err != nil {
			return err
		}

		// Save the move.
		var verrs *validate.Errors
		verrs, err = appCtx.DB().ValidateAndSave(move)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(move.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		ppmOnlyMove := true
		for _, s := range move.MTOShipments {
			if s.ShipmentType != models.MTOShipmentTypePPM {
				ppmOnlyMove = false
				break
			}
		}

		// If this is a PPM-only move, then we also need to adjust other statuses:
		//   - set MTO shipment status to APPROVED
		//   - set PPM shipment status to WAITING_ON_CUSTOMER
		// TODO: Perhaps this could be part of the shipment router. PPMs are a separate model/table,
		//   so would need to figure out how they factor in.
		if ppmOnlyMove {
			// Note: Avoiding the copy of the element in the range so we can preserve the changes to the
			// statuses when we return the entire move tree.
			for i := range move.MTOShipments { // We should only have PPM shipments if we get to here.
				move.MTOShipments[i].Status = models.MTOShipmentStatusApproved

				verrs, err = appCtx.DB().ValidateAndSave(&move.MTOShipments[i])
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(move.MTOShipments[i].ID, nil, verrs, "")
				}
				if err != nil {
					return err
				}

				if move.MTOShipments[i].PPMShipment != nil {
					move.MTOShipments[i].PPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

					verrs, err = appCtx.DB().ValidateAndSave(move.MTOShipments[i].PPMShipment)
					if verrs != nil && verrs.HasAny() {
						return apperror.NewInvalidInputError(move.MTOShipments[i].PPMShipment.ID, nil, verrs, "")
					}
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}

// UpdateReviewedBillableWeightsAt updates the BillableWeightsReviewedAt field on the move (move task order)
func (o moveTaskOrderUpdater) UpdateReviewedBillableWeightsAt(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	var err error

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// update field for move
		now := time.Now()
		move.BillableWeightsReviewedAt = &now

		// Check the If-Match header against existing eTag before updating
		encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return apperror.NewPreconditionFailedError(move.ID, err)
		}

		err = appCtx.DB().Update(move)
		return err
	})
	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}

// UpdateTIORemarks updates the TIORemarks field on the move (move task order)
func (o moveTaskOrderUpdater) UpdateTIORemarks(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string, remarks string) (*models.Move, error) {
	var err error

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// update field for move
		move.TIORemarks = &remarks

		// Check the If-Match header against existing eTag before updating
		encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return apperror.NewPreconditionFailedError(move.ID, err)
		}

		err = appCtx.DB().Update(move)
		return err
	})
	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}

// MakeAvailableToPrime approves a Move, makes it available to prime, and
// creates Move-level service items (counseling and move management) if the
// TOO selected them. If the move received service counseling, the counseling
// service item will automatically be created without the TOO having to select it.
func (o *moveTaskOrderUpdater) MakeAvailableToPrime(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string,
	includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error) {

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	existingETag := etag.GenerateEtag(move.UpdatedAt)
	if existingETag != eTag {
		return &models.Move{}, apperror.NewPreconditionFailedError(move.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	if move.AvailableToPrimeAt == nil {
		now := time.Now()
		move.AvailableToPrimeAt = &now

		err = o.moveRouter.Approve(appCtx, move)
		if err != nil {
			return &models.Move{}, apperror.NewConflictError(move.ID, err.Error())
		}

		transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			err = o.updateMove(txnAppCtx, move, order.CheckRequiredFields())
			if err != nil {
				return err
			}

			// When provided, this will create and approve these Move-level service items.
			if includeServiceCodeMS {
				err = o.createMoveLevelServiceItem(txnAppCtx, *move, models.ReServiceCodeMS)
			}

			if err != nil {
				return err
			}

			if includeServiceCodeCS {
				err = o.createMoveLevelServiceItem(txnAppCtx, *move, models.ReServiceCodeCS)
			}

			return err
		})

		if transactionError != nil {
			return &models.Move{}, transactionError
		}
	}

	return move, nil
}

func (o *moveTaskOrderUpdater) updateMove(appCtx appcontext.AppContext, move *models.Move, checks ...order.Validator) error {
	if verr := order.ValidateOrder(&move.Orders, checks...); verr != nil {
		return verr
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(move)

	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(move.ID, nil, verrs, "")
	}

	return err
}

func (o *moveTaskOrderUpdater) createMoveLevelServiceItem(appCtx appcontext.AppContext, move models.Move, code models.ReServiceCode) error {
	now := time.Now()

	siCreator := o.serviceItemCreator

	_, verrs, err := siCreator.CreateMTOServiceItem(appCtx, &models.MTOServiceItem{
		MoveTaskOrderID: move.ID,
		MTOShipmentID:   nil,
		ReService:       models.ReService{Code: code},
		Status:          models.MTOServiceItemStatusApproved,
		ApprovedAt:      &now,
	})

	if err != nil {
		if errors.Is(err, models.ErrInvalidTransition) {
			return apperror.NewConflictError(move.ID, err.Error())
		}
		return err
	}

	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(move.ID, nil, verrs, "")
	}

	return nil
}

// UpdateMoveTaskOrderQueryBuilder is the query builder for updating MTO
type UpdateMoveTaskOrderQueryBuilder interface {
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

// UpdatePostCounselingInfo updates the counseling info
func (o *moveTaskOrderUpdater) UpdatePostCounselingInfo(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	// Fetch the move and associations.
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:            false,
		MoveTaskOrderID:          moveTaskOrderID,
		ExcludeExternalShipments: true,
	}
	moveTaskOrder, fetchErr := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if fetchErr != nil {
		return &models.Move{}, fetchErr
	}

	approvedForPrimeCounseling := false
	for _, serviceItem := range moveTaskOrder.MTOServiceItems {
		if serviceItem.ReService.Code == models.ReServiceCodeCS && serviceItem.Status == models.MTOServiceItemStatusApproved {
			approvedForPrimeCounseling = true
			break
		}
	}
	if !approvedForPrimeCounseling {
		return &models.Move{}, apperror.NewConflictError(moveTaskOrderID, "Counseling is not an approved service item")
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// Check the If-Match header against existing eTag before updating.
		encodedUpdatedAt := etag.GenerateEtag(moveTaskOrder.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return apperror.NewPreconditionFailedError(moveTaskOrderID, nil)
		}

		now := time.Now()
		moveTaskOrder.PrimeCounselingCompletedAt = &now

		verrs, err := appCtx.DB().ValidateAndSave(moveTaskOrder)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(moveTaskOrderID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		// Note: Avoiding the copy of the element in the range so we can preserve the changes to the
		// statuses when we return the entire move tree.
		for i := range moveTaskOrder.MTOShipments {
			if moveTaskOrder.MTOShipments[i].PPMShipment != nil {
				moveTaskOrder.MTOShipments[i].PPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

				verrs, err = appCtx.DB().ValidateAndSave(moveTaskOrder.MTOShipments[i].PPMShipment)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(moveTaskOrder.MTOShipments[i].PPMShipment.ID, nil, verrs, "")
				}
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return moveTaskOrder, nil
}

// ShowHide changes the value in the "Show" field for a Move. This can be either True or False and indicates if the move has been deactivated or not.
func (o *moveTaskOrderUpdater) ShowHide(appCtx appcontext.AppContext, moveID uuid.UUID, show *bool) (*models.Move, error) {
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   true, // We need to search every move to change its status
		MoveTaskOrderID: moveID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return nil, err
	}

	if show == nil {
		return nil, apperror.NewInvalidInputError(moveID, nil, nil, "The 'show' field must be either True or False - it cannot be empty")
	}

	move.Show = show
	verrs, err := appCtx.DB().ValidateAndSave(move)
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "Invalid input found while updating the Move")
	} else if err != nil {
		return nil, apperror.NewQueryError("Move", err, "")
	}

	// Get the updated Move and return
	updatedMove, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return nil, apperror.NewQueryError("Move", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}

	return updatedMove, nil
}
