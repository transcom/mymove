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
	"github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"
)

type moveTaskOrderUpdater struct {
	db *pop.Connection
	moveTaskOrderFetcher
	builder            UpdateMoveTaskOrderQueryBuilder
	serviceItemCreator services.MTOServiceItemCreator
	moveRouter         services.MoveRouter
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(db *pop.Connection, builder UpdateMoveTaskOrderQueryBuilder, serviceItemCreator services.MTOServiceItemCreator, moveRouter services.MoveRouter) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{db, moveTaskOrderFetcher{db}, builder, serviceItemCreator, moveRouter}
}

// UpdateStatusServiceCounselingCompleted updates the status on the move (move task order) to service counseling completed
func (o moveTaskOrderUpdater) UpdateStatusServiceCounselingCompleted(moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	var err error
	var verrs *validate.Errors

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(&searchParams)
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

// MakeAvailableToPrime approves a Move, makes it available to prime, and
// creates Move-level service items (counseling and move management) if the
// TOO selected them. If the move received service counseling, the counseling
// service item will automatically be created without the TOO having to select it.
func (o *moveTaskOrderUpdater) MakeAvailableToPrime(moveTaskOrderID uuid.UUID, eTag string,
	includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error) {

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(&searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	existingETag := etag.GenerateEtag(move.UpdatedAt)
	if existingETag != eTag {
		return &models.Move{}, services.NewPreconditionFailedError(move.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	if move.AvailableToPrimeAt == nil {
		now := time.Now()
		move.AvailableToPrimeAt = &now

		err = o.moveRouter.Approve(move)
		if err != nil {
			return &models.Move{}, services.NewConflictError(move.ID, err.Error())
		}

		transactionError := o.db.Transaction(func(tx *pop.Connection) error {
			err = o.updateMove(tx, *move, order.CheckRequiredFields())
			if err != nil {
				return err
			}

			// When provided, this will create and approve these Move-level service items.
			if includeServiceCodeMS {
				err = o.createMoveLevelServiceItem(tx, *move, models.ReServiceCodeMS)
			}

			if err != nil {
				return err
			}

			if includeServiceCodeCS {
				err = o.createMoveLevelServiceItem(tx, *move, models.ReServiceCodeCS)
			}

			return err
		})

		if transactionError != nil {
			return &models.Move{}, transactionError
		}
	}

	return move, nil
}

func (o *moveTaskOrderUpdater) updateMove(tx *pop.Connection, move models.Move, checks ...order.Validator) error {
	if verr := order.ValidateOrder(&move.Orders, checks...); verr != nil {
		return verr
	}

	verrs, err := tx.ValidateAndUpdate(&move)

	if verrs != nil && verrs.HasAny() {
		return services.NewInvalidInputError(move.ID, nil, verrs, "")
	}

	return err
}

func (o *moveTaskOrderUpdater) createMoveLevelServiceItem(tx *pop.Connection, move models.Move, code models.ReServiceCode) error {
	now := time.Now()

	siCreator := o.serviceItemCreator
	siCreator.SetConnection(tx)

	_, verrs, err := siCreator.CreateMTOServiceItem(&models.MTOServiceItem{
		MoveTaskOrderID: move.ID,
		MTOShipmentID:   nil,
		ReService:       models.ReService{Code: code},
		Status:          models.MTOServiceItemStatusApproved,
		ApprovedAt:      &now,
	})

	if err != nil {
		if errors.Is(err, models.ErrInvalidTransition) {
			return services.NewConflictError(move.ID, err.Error())
		}
		return err
	}

	if verrs != nil && verrs.HasAny() {
		return services.NewInvalidInputError(move.ID, nil, verrs, "")
	}

	return nil
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
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   true, // We need to search every move to change its status
		MoveTaskOrderID: moveID,
	}
	move, err := o.FetchMoveTaskOrder(&searchParams)
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
	updatedMove, err := o.FetchMoveTaskOrder(&searchParams)
	if err != nil {
		return nil, services.NewQueryError("Move", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}

	return updatedMove, nil
}

func (o *moveTaskOrderUpdater) UpdateApprovedAmendedOrders(move models.Move) error {
	eTag := etag.GenerateEtag(move.UpdatedAt)
	verrs, err := o.builder.UpdateOne(&move, &eTag)

	if verrs != nil && verrs.HasAny() {
		return services.NewInvalidInputError(move.ID, err, verrs, "")
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return services.NewPreconditionFailedError(move.ID, err)
		default:
			return err
		}
	}

	return nil
}
