package movetaskorder

import (
	"database/sql"
	"fmt"
	"github.com/gobuffalo/validate"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveTaskOrderFetcher struct {
	db *pop.Connection
}

func (f moveTaskOrderFetcher) ListMoveTaskOrders(moveOrderID uuid.UUID) ([]models.MoveTaskOrder, error) {
	var moveTaskOrders []models.MoveTaskOrder
	err := f.db.Where("move_order_id = $1", moveOrderID).Eager().All(&moveTaskOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.MoveTaskOrder{}, services.NotFoundError{}
		default:
			return []models.MoveTaskOrder{}, err
		}
	}
	return moveTaskOrders, nil
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.MoveTaskOrderFetcher {
	return &moveTaskOrderFetcher{db}
}

//FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderFetcher) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	mto := &models.MoveTaskOrder{}
	if err := f.db.Eager().Find(mto, moveTaskOrderID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MoveTaskOrder{}, services.NewNotFoundError(moveTaskOrderID)
		default:
			return &models.MoveTaskOrder{}, err
		}
	}

	f.createDefaultServiceItems(mto)

	return mto, nil
}

func (f moveTaskOrderFetcher) createDefaultServiceItems(mto *models.MoveTaskOrder) error {
	var reServices []models.ReService
	err := f.db.Where("code in (?)", []string{"MS", "CS"}).All(&reServices)

	if err != nil {
		return err
	}

	defaultServiceItems := make(map[uuid.UUID]models.MTOServiceItem)
	for _, reService := range reServices {
		defaultServiceItems[reService.ID] = models.MTOServiceItem{
			ReServiceID:     reService.ID,
			MoveTaskOrderID: mto.ID,
		}
	}

	// Remove the ones that exist on the mto
	for _, item := range mto.MTOServiceItems {
		for _, reService := range reServices {
			if item.ReServiceID == reService.ID {
				delete(defaultServiceItems, reService.ID)
			}
		}
	}

	for _, serviceItem := range defaultServiceItems {
		_, err := f.db.ValidateAndCreate(&serviceItem)

		if err != nil {
			return err
		}
	}

	return nil
}

type moveTaskOrderUpdater struct {
	db *pop.Connection
	moveTaskOrderFetcher
	builder UpdateMoveTaskOrderQueryBuilder
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(db *pop.Connection, builder UpdateMoveTaskOrderQueryBuilder) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{db, moveTaskOrderFetcher{db}, builder}
}

//MakeAvailableToPrime updates the status of a MoveTaskOrder for a given UUID to make it available to prime
func (f moveTaskOrderFetcher) MakeAvailableToPrime(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	mto, err := f.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	mto.IsAvailableToPrime = true
	vErrors, err := f.db.ValidateAndUpdate(mto)
	if vErrors.HasAny() {
		return &models.MoveTaskOrder{}, services.InvalidInputError{}
	}
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}
	return mto, nil
}

// UpdateMoveTaskOrderQueryBuilder is the query builder for updating MTO
type UpdateMoveTaskOrderQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

func (o *moveTaskOrderUpdater) UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.MoveTaskOrder, error) {
	var moveTaskOrder models.MoveTaskOrder

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveTaskOrderID),
	}

	err := o.builder.FetchOne(&moveTaskOrder, queryFilters)

	if err != nil {
		return nil, NotFoundError{id: moveTaskOrderID}
	}

	moveTaskOrder.PPMType = body.PpmType
	moveTaskOrder.PPMEstimatedWeight = unit.Pound(body.PpmEstimatedWeight)

	verrs, _ := o.builder.UpdateOne(&moveTaskOrder, &eTag)

	if verrs != nil && verrs.HasAny() {
		return nil, ValidationError{
			id:    moveTaskOrder.ID,
			Verrs: verrs,
		}
	}

	return &moveTaskOrder, nil
}

// NotFoundError is the not found error
type NotFoundError struct {
	id uuid.UUID
}

// Error is the string representation of the error
func (e NotFoundError) Error() string {
	return fmt.Sprintf("move_task_order with id '%s' not found", e.id.String())
}

// ValidationError is the validation error
type ValidationError struct {
	id    uuid.UUID
	Verrs *validate.Errors
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("move_task_order with id: '#{e.id.String()} could not be updated due to a validation error")
}

// PreconditionFailedError is the precondition failed error
type PreconditionFailedError struct {
	id  uuid.UUID
	Err error
}

// Error is the string representation of the precondition failed error
func (e PreconditionFailedError) Error() string {
	return fmt.Sprintf("move_task_order with id: '%s' could not be updated due to the record being stale", e.id.String())
}
