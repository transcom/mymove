package movetaskorder

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"

	"github.com/gobuffalo/validate/v3"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"
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
	var err error
	var verrs *validate.Errors

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	// check if status is in the right state
	// needs to be in MoveStatusNeedsServiceCounseling
	if move.Status != models.MoveStatusNeedsServiceCounseling {
		err = errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Cannot move to Service Counseling Completed state when the Move is not in a Needs Service Counseling state for status: %s", move.Status))

		return &models.Move{}, apperror.NewConflictError(move.ID, err.Error())
	}

	for _, s := range move.MTOShipments {
		if s.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom && s.StorageFacilityID == nil {
			return &models.Move{}, apperror.NewConflictError(
				s.ID, "NTS-release shipment must include facility info")
		}
	}

	// update field for move
	now := time.Now()
	move.ServiceCounselingCompletedAt = &now
	// set status to service counseling completed
	move.Status = models.MoveStatusServiceCounselingCompleted

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, apperror.NewPreconditionFailedError(move.ID, err)
	}

	verrs, err = appCtx.DB().ValidateAndSave(move)
	if verrs != nil && verrs.HasAny() {
		return &models.Move{}, apperror.NewInvalidInputError(move.ID, nil, verrs, "")
	}
	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, apperror.NewPreconditionFailedError(move.ID, err)
		default:
			return &models.Move{}, err
		}
	}

	return move, nil
}

func (o moveTaskOrderUpdater) UpdateStatusServiceCounselingPPMApproved(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	var err error
	var verrs *validate.Errors

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// check if status is in the right state
		// needs to be in MoveStatusNeedsServiceCounseling
		if move.Status != models.MoveStatusNeedsServiceCounseling {
			err = errors.Wrap(models.ErrInvalidTransition,
				fmt.Sprintf("Cannot move to Approved state when the move is not in a Needs Service Counseling state for status: %s", move.Status))

			return apperror.NewConflictError(move.ID, err.Error())
		}

		if len(move.MTOShipments) == 0 {
			return apperror.NewConflictError(move.ID, "No shipments associated with move")
		}

		ppmOnlyMove := true
		for _, s := range move.MTOShipments {
			if s.ShipmentType != models.MTOShipmentTypePPM {
				ppmOnlyMove = false
				break
			}
		}
		if !ppmOnlyMove {
			return apperror.NewConflictError(move.ID, "Move should only contain PPM shipments")
		}

		// Set move status to APPROVED
		move.Status = models.MoveStatusAPPROVED

		// Check the If-Match header against existing eTag before updating
		encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return apperror.NewPreconditionFailedError(move.ID, err)
		}

		verrs, err = appCtx.DB().ValidateAndSave(move)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(move.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		// Set MTO shipment status to APPROVED and set PPM shipment status to WAITING_ON_CUSTOMER
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

			// Due to a Pop bug, we cannot EagerPreload "PPMShipment" likely because it is a pointer and
			// a "has_one" field.  This seems similar to other EagerPreload issues we've found (and
			// sometimes fixed): https://github.com/gobuffalo/pop/issues?q=author%3Areggieriser
			loadErr := appCtx.DB().Load(&move.MTOShipments[i], "PPMShipment")
			if loadErr != nil {
				return apperror.NewQueryError("PPMShipment", err, "")
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
func (o *moveTaskOrderUpdater) UpdatePostCounselingInfo(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.Move, error) {
	var moveTaskOrder models.Move

	err := appCtx.DB().Q().EagerPreload(
		"Orders.NewDutyLocation.Address",
		"Orders.ServiceMember",
		"Orders.Entitlement",
		"MTOShipments",
		"PaymentRequests",
	).Find(&moveTaskOrder, moveTaskOrderID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveTaskOrderID, "while looking for moveTaskOrder.")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	estimatedWeight := unit.Pound(body.PpmEstimatedWeight)
	moveTaskOrder.PPMType = &body.PpmType
	moveTaskOrder.PPMEstimatedWeight = &estimatedWeight
	verrs, err := o.builder.UpdateOne(appCtx, &moveTaskOrder, &eTag)

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(moveTaskOrder.ID, err, verrs, "")
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, apperror.NewPreconditionFailedError(moveTaskOrder.ID, err)
		default:
			return nil, err
		}
	}

	// Filtering external vendor shipments (if requested) in code since we can't do it easily in Pop
	// without a raw query (which could be painful since we'd have to populate all the associations).
	var filteredShipments models.MTOShipments
	if moveTaskOrder.MTOShipments != nil {
		filteredShipments = models.MTOShipments{}
	}
	for _, shipment := range moveTaskOrder.MTOShipments {
		if !shipment.UsesExternalVendor {
			filteredShipments = append(filteredShipments, shipment)
		}
	}
	moveTaskOrder.MTOShipments = filteredShipments

	return &moveTaskOrder, nil
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
