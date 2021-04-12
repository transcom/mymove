package movetaskorder

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"
)

type moveTaskOrderUpdater struct {
	db *pop.Connection
	moveTaskOrderFetcher
	builder            UpdateMoveTaskOrderQueryBuilder
	serviceItemCreator services.MTOServiceItemCreator
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(db *pop.Connection, builder UpdateMoveTaskOrderQueryBuilder, serviceItemCreator services.MTOServiceItemCreator) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{db, moveTaskOrderFetcher{db}, builder, serviceItemCreator}
}

// UpdateStatusServiceCounselingCompleted updates the status on the move (move task order) to service counseling completed
func (o moveTaskOrderUpdater) UpdateStatusServiceCounselingCompleted(moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	var err error
	var verrs *validate.Errors

	searchParams := services.FetchMoveTaskOrderParams{
		IncludeHidden: false,
	}
	move, err := o.FetchMoveTaskOrder(moveTaskOrderID, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	// check if status is in the right state
	// needs to be in MoveStatusNeedsServiceCounseling
	if move.Status != models.MoveStatusNeedsServiceCounseling {
		err = errors.Wrap(models.ErrInvalidTransition,
			fmt.Sprintf("Cannot move to Service Counseling Completed state when the Move is not in a Needs Service Counseling state for status: %s", move.Status))

		return &models.Move{}, services.NewConflictError(move.ID, err.Error())
	}

	// update field for move
	now := time.Now()
	move.ServiceCounselingCompletedAt = &now
	// set status to service counseling completed
	move.Status = models.MoveStatusServiceCounselingCompleted

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, services.NewPreconditionFailedError(move.ID, err)
	}

	verrs, err = o.db.ValidateAndSave(move)
	if verrs != nil && verrs.HasAny() {
		return &models.Move{}, services.NewInvalidInputError(move.ID, nil, verrs, "")
	}
	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, services.NewPreconditionFailedError(move.ID, err)
		default:
			return &models.Move{}, err
		}
	}

	return move, nil
}

//MakeAvailableToPrime updates the status of a MoveTaskOrder for a given UUID to make it available to prime
func (o moveTaskOrderUpdater) MakeAvailableToPrime(moveTaskOrderID uuid.UUID, eTag string,
	includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error) {
	var err error
	var verrs *validate.Errors

	searchParams := services.FetchMoveTaskOrderParams{
		IncludeHidden: false,
	}
	move, err := o.FetchMoveTaskOrder(moveTaskOrderID, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	// Fail early if the Order is invalid due to missing required fields
	order := move.Orders
	err = o.db.Load(&order, "Moves")
	if err != nil {
		return &models.Move{}, err
	}
	if verrs, err = o.db.ValidateAndUpdate(&order); verrs.HasAny() || err != nil {
		return &models.Move{}, services.NewInvalidInputError(move.ID, nil, verrs, "")
	}

	if move.AvailableToPrimeAt == nil {
		// update field for move
		now := time.Now()
		move.AvailableToPrimeAt = &now

		if move.Status == models.MoveStatusSUBMITTED {
			err = move.Approve()
			if err != nil {
				return &models.Move{}, services.NewConflictError(move.ID, err.Error())
			}
		}

		verrs, err = o.builder.UpdateOne(move, &eTag)
		if verrs != nil && verrs.HasAny() {
			return &models.Move{}, services.NewInvalidInputError(move.ID, nil, verrs, "")
		}
		if err != nil {
			switch err.(type) {
			case query.StaleIdentifierError:
				return nil, services.NewPreconditionFailedError(move.ID, err)
			default:
				return &models.Move{}, err
			}
		}

		// When provided, this will auto create and approve MTO level service items. This is going to typically happen
		// from the ghc api via the office app. The handler in question is this one: UpdateMoveTaskOrderStatusHandlerFunc
		// in ghcapi/move_task_order.go
		if includeServiceCodeMS {
			// create if doesn't exist
			_, verrs, err = o.serviceItemCreator.CreateMTOServiceItem(&models.MTOServiceItem{
				MoveTaskOrderID: moveTaskOrderID,
				MTOShipmentID:   nil,
				ReService:       models.ReService{Code: models.ReServiceCodeMS},
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &now,
			})
		}

		if err != nil {
			if errors.Is(err, models.ErrInvalidTransition) {
				return &models.Move{}, services.NewConflictError(move.ID, err.Error())
			}
			return &models.Move{}, err
		}
		if verrs != nil {
			return &models.Move{}, verrs
		}

		if includeServiceCodeCS {
			// create if doesn't exist
			_, verrs, err = o.serviceItemCreator.CreateMTOServiceItem(&models.MTOServiceItem{
				MoveTaskOrderID: moveTaskOrderID,
				MTOShipmentID:   nil,
				ReService:       models.ReService{Code: models.ReServiceCodeCS},
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &now,
			})
		}

		if err != nil {
			if errors.Is(err, models.ErrInvalidTransition) {
				return &models.Move{}, services.NewConflictError(move.ID, err.Error())
			}
			return &models.Move{}, err
		}
		if verrs != nil {
			return &models.Move{}, verrs
		}

		// CreateMTOServiceItem may have updated the move status so refetch as to not return incorrect status
		// TODO: Modify CreateMTOServiceItem to return the updated move or refactor to operate on the passed in reference
		move, err = o.FetchMoveTaskOrder(moveTaskOrderID, nil)
		if err != nil {
			return &models.Move{}, err
		}
	}

	return move, nil
}

// UpdateMoveTaskOrderQueryBuilder is the query builder for updating MTO
type UpdateMoveTaskOrderQueryBuilder interface {
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

// UpdatePostCounselingInfo updates the counseling info
func (o *moveTaskOrderUpdater) UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.Move, error) {
	var moveTaskOrder models.Move

	err := o.db.Q().EagerPreload(
		"Orders.NewDutyStation.Address",
		"Orders.ServiceMember",
		"Orders.Entitlement",
		"MTOShipments",
		"PaymentRequests",
	).Find(&moveTaskOrder, moveTaskOrderID)

	if err != nil {
		return nil, services.NewNotFoundError(moveTaskOrderID, "while looking for moveTaskOrder.")
	}

	estimatedWeight := unit.Pound(body.PpmEstimatedWeight)
	moveTaskOrder.PPMType = &body.PpmType
	moveTaskOrder.PPMEstimatedWeight = &estimatedWeight
	verrs, err := o.builder.UpdateOne(&moveTaskOrder, &eTag)

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(moveTaskOrder.ID, err, verrs, "")
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, services.NewPreconditionFailedError(moveTaskOrder.ID, err)
		default:
			return nil, err
		}
	}

	return &moveTaskOrder, nil
}

// ShowHide changes the value in the "Show" field for a Move. This can be either True or False and indicates if the move has been deactivated or not.
func (o *moveTaskOrderUpdater) ShowHide(moveID uuid.UUID, show *bool) (*models.Move, error) {
	searchParams := services.FetchMoveTaskOrderParams{
		IncludeHidden: true, // We need to search every move to change its status
	}
	move, err := o.FetchMoveTaskOrder(moveID, &searchParams)
	if err != nil {
		return nil, services.NewNotFoundError(moveID, "while fetching the Move")
	}

	if show == nil {
		return nil, services.NewInvalidInputError(moveID, nil, nil, "The 'show' field must be either True or False - it cannot be empty")
	}

	move.Show = show
	verrs, err := o.db.ValidateAndSave(move)
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(move.ID, err, verrs, "Invalid input found while updating the Move")
	} else if err != nil {
		return nil, services.NewQueryError("Move", err, "")
	}

	// Get the updated Move and return
	updatedMove, err := o.FetchMoveTaskOrder(move.ID, &searchParams)
	if err != nil {
		return nil, services.NewQueryError("Move", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}

	return updatedMove, nil
}
