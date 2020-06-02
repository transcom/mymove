package movetaskorder

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

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
	builder UpdateMoveTaskOrderQueryBuilder
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(db *pop.Connection, builder UpdateMoveTaskOrderQueryBuilder) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{db, moveTaskOrderFetcher{db}, builder}
}

//MakeAvailableToPrime updates the status of a MoveTaskOrder for a given UUID to make it available to prime
func (o moveTaskOrderUpdater) MakeAvailableToPrime(moveTaskOrderID uuid.UUID, eTag string) (*models.MoveTaskOrder, error) {
	mto, err := o.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.MoveTaskOrder{}, err
	}

	if mto.AvailableToPrimeAt == nil {
		now := time.Now()
		mto.AvailableToPrimeAt = &now
	}

	verrs, err := o.builder.UpdateOne(mto, &eTag)
	if verrs != nil && verrs.HasAny() {
		return &models.MoveTaskOrder{}, services.InvalidInputError{}
	}
	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, services.NewPreconditionFailedError(mto.ID, err)
		default:
			return &models.MoveTaskOrder{}, err
		}
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
		return nil, services.NewNotFoundError(moveTaskOrder.ID, "")
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
